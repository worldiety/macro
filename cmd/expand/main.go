package main

import (
	"flag"
	"fmt"
	"github.com/worldiety/macro/ast/golang"
	"log"
	"log/slog"
	"os"
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

	modDir, ok := golang.ModRootDir(*dir)
	if !ok {
		return fmt.Errorf("not within a Go module directory: %s", *dir)
	}

	slog.Info("expanding in module", slog.String("dir", modDir))

	pkgs, err := golang.Load(modDir)
	if err != nil {
		return fmt.Errorf("cannot parse and resolve Go module source: %w", err)
	}

	typeDeclr := golang.Macros(pkgs)

	fmt.Printf("%+v", typeDeclr)
	return nil
}
