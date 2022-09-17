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
