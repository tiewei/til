/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if err := run(os.Args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %s\n", err)
		os.Exit(1)
	}
}

// run executes the command.
func run(args []string, stdout, stderr io.Writer) error {
	cmdName := filepath.Base(args[0])

	if len(args) == 1 {
		return fmt.Errorf("no subcommand provided.\n\n%s", usage(cmdName))
	}

	flagSet := flag.NewFlagSet(cmdName, flag.ExitOnError)
	flagSet.SetOutput(stderr)
	setUsageFn(flagSet, usage)

	_ = flagSet.Parse(args[1:]) // ignore err; the FlagSet uses ExitOnError

	common := GenericCommand{
		stdout:  stdout,
		flagSet: flagSet,
	}

	switch subcommand := args[1]; subcommand {
	case cmdGenerate:
		cmd := &GenerateCommand{
			GenericCommand: common,
		}
		return cmd.Run(args[2:]...)

	case cmdValidate:
		cmd := &ValidateCommand{
			GenericCommand: common,
		}
		return cmd.Run(args[2:]...)

	case cmdGraph:
		cmd := &GraphCommand{
			GenericCommand: common,
		}
		return cmd.Run(args[2:]...)

	default:
		return fmt.Errorf("unknow subcommand %q.\n\n%s", subcommand, usage(cmdName))
	}
}
