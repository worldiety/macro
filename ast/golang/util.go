package golang

import (
	"os"
	"path/filepath"
)

func ModRootDir(dir string) (string, bool) {
	for len(dir) > 0 {
		gomod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(gomod); err != nil {
			dir = filepath.Dir(dir)
		} else {
			return dir, true
		}
	}

	return dir, false
}
