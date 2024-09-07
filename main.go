//    Copyright 2021-2023 Olivier Mengué
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/moby/patternmatcher"
	"github.com/moby/patternmatcher/ignorefile"
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
	var dockerFile string
	flag.BoolVar(&verbose, "v", false, "verbose mode: show ignored files on stderr")
	flag.StringVar(&dockerFile, "f", "", "name of the `Dockerfile` to derive the dockerignore file")
	flag.StringVar(&dockerFile, "file", "", "name of the `Dockerfile` to derive the dockerignore file (also --file)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: %s [-v] [PATH]\n\n", os.Args[0])
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	switch flag.NArg() {
	case 0:
		// If the path is not explicitely given, check there is a Dockerfile
		if _, err := os.Stat(dockerFile); dockerFile != "" && errors.Is(err, os.ErrNotExist) {
			return errors.New("no Dockerfile here")
		}
	case 1:
		if dockerFile != "" {
			var err error
			dockerFile, err = filepath.Abs(dockerFile)
			if err != nil {
				return err
			}
		}
		if err := os.Chdir(flag.Arg(0)); err != nil {
			return err
		}
	default:
		flag.Usage()
		os.Exit(2)
	}

	if dockerFile == "" {
		dockerFile = "Dockerfile"
	}

	var dockerIgnore string
	var ignorePatterns []string
	// Handle .dockerignore attached to a Dockerfile
	// https://github.com/moby/buildkit/releases/tag/dockerfile%2F1.1.0
	for _, dockerIgnore = range []string{dockerFile + ".dockerignore", ".dockerignore"} {
		if f, err := os.Open(dockerIgnore); err == nil {
			if ignorePatterns, err = ignorefile.ReadAll(f); err != nil {
				return fmt.Errorf(dockerIgnore+": %w", err)
			}
			break
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	ignore, err := patternmatcher.New(ignorePatterns)
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
		if m, _ := ignore.MatchesOrParentMatches(relFilePath); m {
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
