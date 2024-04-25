package java

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/worldiety/macro/pkg/src"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const googleJavaFormatUrl = "https://github.com/google/google-java-format/releases/download/google-java-format-1.9/google-java-format-1.9-all-deps.jar"
const googleJavaFormatShaSum = "3214b30794e15bf501ae66805ca40475de71e831"

// Format tries to apply the google-java-format rules to the given text
// (see https://google.github.io/styleguide/javaguide.html).
// If it fails, the error is returned and the string contains the text with line enumeration.
// The formatter is automatically downloaded into the tmp folder and executed from there.
func Format(source []byte) ([]byte, error) {
	dir := filepath.Join(os.TempDir(), "golangee", "google-java-format")
	jarFile := filepath.Join(dir, "google-java-format")
	if _, err := os.Stat(jarFile); err != nil {
		fmt.Printf("downloading google-java-format...\n")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("unable to create temporary folder: %w", err)
		}

		// download
		res, err := http.Get(googleJavaFormatUrl)
		if err != nil {
			return nil, fmt.Errorf("cannot start download %s: %w", googleJavaFormatUrl, err)
		}

		defer res.Body.Close()

		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("cannot download %s: %w", googleJavaFormatUrl, err)
		}

		shaSum := sha1.Sum(buf)
		if hex.EncodeToString(shaSum[:]) != googleJavaFormatShaSum {
			return nil, fmt.Errorf("download corruption detected %s. Expected: %s", googleJavaFormatUrl, googleJavaFormatShaSum)
		}

		if err := ioutil.WriteFile(jarFile, buf, os.ModePerm); err != nil {
			return nil, fmt.Errorf("cannot save jar file %s: %w", jarFile, err)
		}
	}

	cmd := exec.Command("java", "-jar", jarFile, "-")
	cmd.Env = os.Environ()
	cmd.Dir = dir
	w, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("unable to open stdinPipe: %w", err)
	}

	if _, err := w.Write(source); err != nil {
		return nil, fmt.Errorf("unable to pipe to stdin: %w", err)
	}

	_ = w.Close()

	res, err := cmd.CombinedOutput()
	if err != nil {
		return []byte(src.WithLineNumbers(string(source))), fmt.Errorf("cannot format: %s: %w", string(res), err)
	}

	return res, nil
}
