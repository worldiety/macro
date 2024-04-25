package main

import (
	"flag"
	"fmt"
	"github.com/worldiety/macro/ast/golang"
	"github.com/worldiety/macro/macro"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current dir: %w", err)
	}

	dir := flag.String("dir", wd, "the directory to expand the files")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		return fmt.Errorf("cannot resolve absolute path: %w", err)
	}

	modDir, ok := golang.ModRootDir(absDir)
	if !ok {
		return fmt.Errorf("not within a Go module directory: %s", absDir)
	}

	slog.Info("expanding in module", slog.String("dir", modDir))

	prog, err := golang.Parse(modDir)
	if err != nil {
		return fmt.Errorf("cannot parse and resolve Go module source: %w", err)
	}

	engine := macro.NewEngine(prog)
	return engine.Exec()
}
