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

package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// CLI implements a command-line interface with subcommands.
type CLI struct {
	// Name of the CLI (the invoked program)
	name string
	// Rest of the command line arguments
	args []string

	// Main usage text, printed when the user asks for help or when an
	// error occurs
	usage string

	// Supported subcommands
	commands map[string]Command

	// UI for reporting messages/errors to the terminal.
	// Defaulted to stdout/stderr during initialization.
	ui UI
}

// New returns a new instance of a CLI.
func New(osArgs []string, usage usageFn, opts ...Option) *CLI {
	// passing empty OS arguments here would be inconsiderate from the
	// caller, so we make the error loud
	if len(osArgs) == 0 {
		panic(errors.New("CLI constructor invoked without OS argument. " + //nolint:stylecheck
			"At least a program name/path is required."))
	}

	c := &CLI{
		name:     osArgs[0],
		args:     osArgs[1:],
		usage:    usage(osArgs[0]),
		commands: make(map[string]Command),
		ui: UI{
			StdWriter: os.Stdout,
			ErrWriter: os.Stderr,
		},
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// Run executes the command, passing OS arguments provided during the CLI
// initialization to the subcommand.
func (c *CLI) Run() error {
	// initializing a FlagSet here ensures users can display the general
	// usage using Go's special '-h, --help' flag
	f := flag.NewFlagSet(c.name, flag.ExitOnError)
	f.SetOutput(c.ui.ErrWriter)
	f.Usage = func() {
		fmt.Fprint(f.Output(), c.usage)
	}
	_ = f.Parse(c.args) // ignore err; the FlagSet uses ExitOnError

	if len(c.args) == 0 {
		return errors.New("no subcommand provided.\n\n" + c.usage)
	}

	cmd, exists := c.commands[c.args[0]]
	if !exists {
		return fmt.Errorf("unknown subcommand %q.\n\n%s", c.args[0], c.usage)
	}

	ctx := context.Background()
	ctx = withFlagSet(ctx, initCmdFlagSet(c.name, c.args[0], c.ui.ErrWriter))
	ctx = withUI(ctx, &c.ui)

	return cmd.Run(ctx, c.args[1:])
}

// usageFn returns the usage text for a command, based on the given CLI name.
type usageFn func(cliName string) string

// Option is a functional option for a CLI.
type Option func(*CLI)

// Subcommand defines a CLI subcommand.
func Subcommand(name string, cmd Command) Option {
	return func(c *CLI) {
		if _, exists := c.commands[name]; exists {
			// this would indicate a programmer mistake, so we make
			// the error loud
			panic(fmt.Errorf("command %q is defined more than once", name))
		}

		c.commands[name] = cmd
	}
}

// StdWriter sets the writer to which regular messages are written.
func StdWriter(w io.Writer) Option {
	return func(c *CLI) {
		c.ui.StdWriter = w
	}
}

// ErrWriter sets the writer to which error messages are written.
func ErrWriter(w io.Writer) Option {
	return func(c *CLI) {
		c.ui.ErrWriter = w
	}
}

// UI wraps io.Writer interfaces used by commands to report messages/errors to
// the terminal.
type UI struct {
	StdWriter io.Writer
	ErrWriter io.Writer
}

// Command can execute a CLI subcommand.
//
// The given context.Context is infused with data and abstractions a subcommand
// typically needs from the main CLI:
//  - an initialized flag.FlagSet, which can be augmented with vars and usage text
//  - a UI instance, for writing messages and errors to the caller's terminal
type Command interface {
	Run(ctx context.Context, args []string) error
}

// initCmdFlagSet returns a flag.FlagSet initialized with a standard name,
// error handling property and output for a CLI subcommand.
func initCmdFlagSet(cliName, cmd string, output io.Writer) *flag.FlagSet {
	flagSet := flag.NewFlagSet(cliName+" "+cmd, flag.ExitOnError)
	flagSet.SetOutput(output)
	return flagSet
}
