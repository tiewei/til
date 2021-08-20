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

package cli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"til/cli"
)

// cmdName is the name of the program passed in all faked OS arguments in this
// test suite.
const cmdName = "clitest"

func TestCLI(t *testing.T) {
	t.Run("run existing subcommand", func(t *testing.T) {
		var stdout strings.Builder
		var stderr strings.Builder

		c := cli.New([]string{cmdName, "mycommand", "-A", "-B"}, usageFn,
			cli.StdWriter(&stdout),
			cli.ErrWriter(&stderr),
			cli.Subcommand("mycommand", new(writeAndSucceedCmd)),
		)

		// The command should run to completion and write its output
		// to the standard CLI writer.
		// Flags are ignored due to the absence of flag parsing in our
		// test's subcommand.
		err := c.Run()

		if err != nil {
			t.Fatal("Unexpected error:", err)
		}
		if stderr.String() != "" {
			t.Errorf("Expected no write to error writer. Got:\n%s", &stderr)
		}

		if stdout.String() != stdCmdMsg {
			t.Errorf("Unexpected command output:\n%s", &stdout)
		}
	})

	t.Run("global help requested", func(t *testing.T) {
		var stdout strings.Builder
		var stderr strings.Builder

		c := cli.New([]string{cmdName, "-help"}, usageFn,
			cli.StdWriter(&stdout),
			cli.ErrWriter(&stderr),
		)

		var err error

		func() {
			defer func() {
				// tests need to recover from os.Exit()
				if r := recover(); r != nil {
					t.Log("Recovered:", r)
				} else {
					t.Fatal("Expected program to exit")
				}
			}()

			// the command should exit (code 0) because the special
			// flag '-h, --help' is encountered
			err = c.Run()
		}()

		if err != nil {
			t.Fatal("Unexpected error:", err)
		}
		if stdout.String() != "" {
			t.Errorf("Expected no write to standard writer. Got:\n%s", &stdout)
		}

		if stderr.String() != expectedUsageTxt {
			t.Errorf("Expected only usage text to be written to error writer. Got:\n%s", &stderr)
		}
	})

	t.Run("no subcommand given", func(t *testing.T) {
		var stdout strings.Builder
		var stderr strings.Builder

		c := cli.New([]string{cmdName}, usageFn,
			cli.StdWriter(&stdout),
			cli.ErrWriter(&stderr),
		)

		err := c.Run()

		if err == nil {
			t.Fatal("Expected an error to occur")
		}

		assertNoWrite(t, &stdout, &stderr)
		assertContainsUsage(t, err.Error())

		if !strings.Contains(err.Error(), "no subcommand provided") {
			t.Errorf("Unexpected error message: %q", err)
		}
	})

	t.Run("subcommand given but not defined", func(t *testing.T) {
		var stdout strings.Builder
		var stderr strings.Builder

		c := cli.New([]string{cmdName, "mycommand", "-A", "-B"}, usageFn,
			cli.StdWriter(&stdout),
			cli.ErrWriter(&stderr),
		)

		err := c.Run()

		if err == nil {
			t.Fatal("Expected an error to occur")
		}

		assertNoWrite(t, &stdout, &stderr)
		assertContainsUsage(t, err.Error())

		if !strings.Contains(err.Error(), `unknown subcommand "mycommand"`) {
			t.Errorf("Unexpected error message: %q", err)
		}
	})
}

// assertNoWrite checks that no write occurred to the CLI writers (stdout/err).
func assertNoWrite(t *testing.T, stdout, stderr fmt.Stringer) {
	t.Helper()

	if stdout.String() != "" {
		t.Errorf("Expected no write to standard writer. Got:\n%s", stdout)
	}
	if stderr.String() != "" {
		t.Errorf("Expected no write to error writer. Got:\n%s", stderr)
	}
}

// assertContainsUsage checks that the given message contains the sample usage text.
func assertContainsUsage(t *testing.T, msg string) {
	t.Helper()

	if !strings.Contains(msg, expectedUsageTxt) {
		t.Errorf("Expected given message to contain command usage:\n%s", msg)
	}
}

// usageFn returns a short sample usage text for use in tests.
func usageFn(cliName string) string {
	return `usage for ` + cliName + ` a.k.a. "HELP!"`
}

const expectedUsageTxt = `usage for ` + cmdName + ` a.k.a. "HELP!"`

// writeAndSucceedCmd is a minimal implementation of cli.Command for tests,
// which write a succinct message to the standard writer and returns no error.
type writeAndSucceedCmd struct{}

var _ cli.Command = (*writeAndSucceedCmd)(nil)

// stdCmdMsg is the message written to test CLIs' standard writer.
const stdCmdMsg = "Hello, World!"

func (*writeAndSucceedCmd) Run(ctx context.Context, _ []string) error {
	ui := cli.UIFromContext(ctx)
	_, _ = ui.StdWriter.Write([]byte(stdCmdMsg))
	return nil
}
