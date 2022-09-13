# Documentation

## Overview

* [Vaults](/docs/vaults.md) describes valid vault types and their configuration
* [Transformations](/docs/transformations.md) is about transformation steps.
* [Sinks](/docs/sinks.md)describe the use of sinks.

## Running

```bash
$ go-secretshelper 
Usage: go-secretshelper [-v] [-e] <command>
where commands are
  version               print out version
  run                   run specified config
```

Global flags are:
* -v: be more verbose
* -e: substitute environment variables when processing configuration files

```bash
$ go-secretshelper run -h
Usage of run:
  -c string
        configuration file
```

The `run` command takes the name of a yaml-based configuration file in `-c`, and starts processing the file.