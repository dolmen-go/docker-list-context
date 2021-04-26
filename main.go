package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/moby/buildkit/frontend/dockerfile/dockerignore"
)

func main() {
	exitCode := 0
	if err := _main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		exitCode = 1
	}
	os.Exit(exitCode)
}

func _main() error {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "verbose mode: show ignored files on stderr")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [-v] [PATH]\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	switch flag.NArg() {
	case 0:
		// If the path is not explicitely given, check there is a Dockerfile
		if _, err := os.Stat("Dockerfile"); errors.Is(err, os.ErrNotExist) {
			return errors.New("no Dockerfile here")
		}
	case 1:
		if err := os.Chdir(flag.Arg(0)); err != nil {
			return err
		}
	default:
		flag.Usage()
		os.Exit(2)
	}

	const dockerIgnore = ".dockerignore"

	var ignorePatterns []string
	if f, err := os.Open(dockerIgnore); err == nil {
		if ignorePatterns, err = dockerignore.ReadAll(f); err != nil {
			return fmt.Errorf(dockerIgnore+": %w", err)
		}
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	ignore, err := fileutils.NewPatternMatcher(ignorePatterns)
	if err != nil {
		return fmt.Errorf(dockerIgnore+": %w", err)
	}

	err = filepath.Walk(".", func(filePath string, f os.FileInfo, err error) error {
		if filePath == "." {
			return nil
		}
		relFilePath, err := filepath.Rel(".", filePath)
		if err != nil {
			return nil
		}
		if m, _ := ignore.Matches(relFilePath); m {
			if verbose {
				fmt.Fprintln(os.Stderr, "IGNORE", relFilePath)
			}
			return nil
		}
		fmt.Println(relFilePath)
		return nil
	})
	return err
}
