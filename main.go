package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"bridgedl/bridge"
	"bridgedl/config/file"
)

const defaultFilePath = "config.brg.hcl"

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %s\n", err)
		os.Exit(1)
	}
}

// run executes the command.
func run(args []string, stderr io.Writer) error {
	opts := parseFlags(args, stderr)

	brg, diags := file.NewParser().LoadBridge(opts.filePath)
	if diags.HasErrors() {
		return diags
	}

	ctx := bridge.Context{
		Bridge: brg,
	}

	g, diags := ctx.Graph()
	if diags.HasErrors() {
		return diags
	}

	_ = g

	return nil
}

// cmdOpts are the options that can be passed to the command.
type cmdOpts struct {
	filePath string
}

// parseFlags parses the given command line arguments and returns the values
// associated with the supported flags.
func parseFlags(args []string, output io.Writer) *cmdOpts {
	cmdName := filepath.Base(args[0])
	flags := flag.NewFlagSet(cmdName, flag.ExitOnError)
	flags.SetOutput(output)

	opts := &cmdOpts{}

	flags.StringVar(&opts.filePath, "f", defaultFilePath, "Path of the config file to parse")

	_ = flags.Parse(args[1:]) // ignore err; the FlagSet uses ExitOnError

	return opts
}
