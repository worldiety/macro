package render

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// Write emits the given artifact into the destination.
func Write(dir string, artifact Artifact) error {
	switch t := artifact.(type) {
	case *File:
		dst := filepath.Join(dir, t.FileName)
		if err := ioutil.WriteFile(dst, t.Buf, 0600); err != nil {
			return fmt.Errorf("unable to emit file: %w", err)
		}
	case *Dir:
		dst := filepath.Join(dir, t.DirName)
		if err := os.MkdirAll(dst, 0700); err != nil {
			return fmt.Errorf("unable to create directory: %s: %w", dst, err)
		}

		for _, file := range t.Files {
			if err := Write(dst, file); err != nil {
				return fmt.Errorf("unable to emit file: %s: %w", dst, err)
			}
		}

		for _, d := range t.Dirs {
			if err := Write(dst, d); err != nil {
				return fmt.Errorf("unable to emit dir: %w", err)
			}
		}
	default:
		return fmt.Errorf("invalid artifact type: %s", reflect.TypeOf(t).String())
	}

	return nil
}

// Clean takes the given magic bytes and searches in the very first bytes of each file in dir recursively, if
// it contains one of the magic sequences and deletes it. It ignores any hidden (prefixed with .) folders and files.
func Clean(dir string, magic ...[]byte) error {
	var files []string

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, fname := range files {
		ok, err := fileHasMagic(fname, magic...)
		if err != nil {
			return fmt.Errorf("unable to check %s for magic: %w", fname, err)
		}

		if ok {
			if err := os.Remove(fname); err != nil {
				return fmt.Errorf("unable to delete file: %s: %w", fname, err)
			}
		}
	}

	return nil
}

func fileHasMagic(fname string, magic ...[]byte) (bool, error) {
	var buf [1024]byte

	file, err := os.Open(fname)
	if err != nil {
		return false, fmt.Errorf("cannot open file: %w", err)
	}

	defer file.Close() // suppress error for this read-only case

	n, err := file.Read(buf[:])
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false, nil
		}

		return false, fmt.Errorf("unable to read buffer: %w", err)
	}

	for _, m := range magic {
		if len(m) == 0 {
			return false, fmt.Errorf("magic sequence is not allowed to be empty")
		}

		if bytes.Contains(buf[:n], m) {
			return true, nil
		}
	}

	return false, nil
}
