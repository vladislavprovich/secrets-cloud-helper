// Package main contains the main command
//
// MIT License
//
// Copyright (c) 2021 Andreas Schmidt
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"io/ioutil"
	"log"
	"os"
)

var (
	commit = "none"
	date   = "unknown"
)

const (
	// ExitCodeOk is ok
	ExitCodeOk = 0

	// ExitCodeNoOrUnknownCommand we're not able to run the command
	ExitCodeNoOrUnknownCommand = 1

	// ExitCodeInvalidConfig something wrong with config
	ExitCodeInvalidConfig = 2
)

func usage() {
	fmt.Println("Usage: go-secretshelper [-v] [-e] [-c config] <command>")
	fmt.Println("where commands are")
	fmt.Println("  version		print out version")
	fmt.Println("  run			run specified config")
}

func main() {

	verboseFlag := flag.Bool("v", false, "Enables verbose output")
	envFlag := flag.Bool("e", false, "Enables environment variable substitution")
	flag.Parse()

	var l *log.Logger
	if *verboseFlag {
		l = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		l = log.New(ioutil.Discard, "", 0)
	}

	values := flag.Args()
	if len(values) == 0 {
		usage()
		os.Exit(ExitCodeNoOrUnknownCommand)
	}

	switch values[0] {
	case "version":
		fmt.Printf("%s (%s)\n", commit, date)
		os.Exit(ExitCodeOk)

	case "run":

		fs := flag.NewFlagSet("run", flag.ExitOnError)
		configFlag := fs.String("c", "", "configuration file")

		if err := fs.Parse(values[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "error parsind commands: %s\n", err)
			os.Exit(ExitCodeNoOrUnknownCommand)
		}

		// read config
		config, err := core.NewConfigFromFile(*configFlag, *envFlag)
		if err != nil {
			fmt.Printf("Unable to read config from file %s: %s\n", *configFlag, err)
			os.Exit(ExitCodeInvalidConfig)
		}

		// validate
		f := adapters.NewBuiltinFactory(l, afero.NewOsFs())
		if err := config.Validate(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error validating configuration: %s\n", err)
			os.Exit(ExitCodeInvalidConfig)
		}

		// run
		cmd := core.NewMainUseCaseImpl(l)

		err = cmd.Process(context.Background(), f,
			&config.Defaults,
			&config.Vaults,
			&config.Secrets,
			&config.Transformations,
			&config.Sinks)

		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		os.Exit(ExitCodeOk)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", values[0])
		usage()
		os.Exit(ExitCodeNoOrUnknownCommand)
	}

}
