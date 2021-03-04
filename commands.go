package main

import (
	"flag"
	"fmt"
	"io"

	"bridgedl/config/file"
	"bridgedl/core"
	"bridgedl/graph/dot"
	"bridgedl/translate"
	"bridgedl/translate/router"
)

// CLI subcommands
const (
	cmdGenerate = "generate"
	cmdValidate = "validate"
	cmdGraph    = "graph"
)

// usage is a usageFn for the top level command.
func usage(cmdName string) string {
	return "Usage: " + cmdName + " <command>\n" +
		"\n" +
		"Commands:\n" +
		"  " + cmdGenerate + "     generate Kubernetes manifests for deploying a Bridge\n" +
		"  " + cmdValidate + "     validate the syntax of a Bridge Description File\n" +
		"  " + cmdGraph + "        represent a Bridge as a directed graph in DOT format\n"
}

// usageGenerate is a usageFn for the "generate" subcommand.
func usageGenerate(cmdName string) string {
	return "Usage: " + cmdName + " " + cmdGenerate + " FILE\n\n" +
		"Generates the Kubernetes manifests that allow the Bridge to be deployed " +
		"to TriggerMesh, and writes them to standard output.\n"
}

// usageValidate is a usageFn for the "validate" subcommand.
func usageValidate(cmdName string) string {
	return "Usage: " + cmdName + " " + cmdValidate + " FILE\n\n" +
		"Returns with an exit code of 0 if FILE is a syntactically valid Bridge " +
		"Description File, with an exit code of 1 otherwise.\n"
}

// usageGraph is a usageFn for the "usage" subcommand.
func usageGraph(cmdName string) string {
	return "Usage: " + cmdName + " " + cmdGraph + " FILE\n\n" +
		"Generates a DOT representation of the Bridge parsed from FILE and writes " +
		"it to standard output.\n"
}

type usageFn func(cmdName string) string

// setUsageFn uses the given usageFn to set the Usage function of the provided
// flag.FlagSet.
func setUsageFn(f *flag.FlagSet, u usageFn) {
	f.Usage = func() {
		fmt.Fprint(f.Output(), u(f.Name()))
	}
}

type Command interface {
	Run(args ...string) error
}

var (
	_ Command = (*ValidateCommand)(nil)
	_ Command = (*GraphCommand)(nil)
)

type GenericCommand struct {
	stdout  io.Writer
	flagSet *flag.FlagSet
}

type GenerateCommand struct {
	GenericCommand
}

// Run implements Command.
func (c *GenerateCommand) Run(args ...string) error {
	setUsageFn(c.flagSet, usageGenerate)

	pos, flags := splitArgs(1, args)
	_ = c.flagSet.Parse(flags) // ignore err; the FlagSet uses ExitOnError

	if len(pos) != 1 {
		return fmt.Errorf("unexpected number of positional arguments.\n\n%s", usageGenerate(c.flagSet.Name()))
	}
	filePath := pos[0]

	_, diags := file.NewParser().LoadBridge(filePath)
	if diags.HasErrors() {
		return diags
	}

	fmt.Fprintln(c.stdout, "Not implemented")

	return nil
}

type ValidateCommand struct {
	GenericCommand
}

// Run implements Command.
func (c *ValidateCommand) Run(args ...string) error {
	setUsageFn(c.flagSet, usageValidate)

	pos, flags := splitArgs(1, args)
	_ = c.flagSet.Parse(flags) // ignore err; the FlagSet uses ExitOnError

	if len(pos) != 1 {
		return fmt.Errorf("unexpected number of positional arguments.\n\n%s", usageValidate(c.flagSet.Name()))
	}
	filePath := pos[0]

	brg, diags := file.NewParser().LoadBridge(filePath)
	if diags.HasErrors() {
		return diags
	}

	tp := &translate.TranslatorProviders{
		Routers: router.AllRouters,
	}

	ctx := core.Context{
		Bridge:      brg,
		Translators: tp,
	}

	if _, diags := ctx.Graph(); diags.HasErrors() {
		return diags
	}

	return nil
}

type GraphCommand struct {
	GenericCommand
}

// Run implements Command.
func (c *GraphCommand) Run(args ...string) error {
	setUsageFn(c.flagSet, usageGraph)

	pos, flags := splitArgs(1, args)
	_ = c.flagSet.Parse(flags) // ignore err; the FlagSet uses ExitOnError

	if len(pos) != 1 {
		return fmt.Errorf("unexpected number of positional arguments.\n\n%s", usageGraph(c.flagSet.Name()))
	}
	filePath := pos[0]

	brg, diags := file.NewParser().LoadBridge(filePath)
	if diags.HasErrors() {
		return diags
	}

	tp := &translate.TranslatorProviders{
		Routers: router.AllRouters,
	}

	ctx := core.Context{
		Bridge:      brg,
		Translators: tp,
	}

	g, diags := ctx.Graph()
	if diags.HasErrors() {
		return diags
	}

	dg, err := dot.Marshal(g)
	if err != nil {
		return fmt.Errorf("marshaling graph to DOT: %w", err)
	}

	if _, err := c.stdout.Write(dg); err != nil {
		return fmt.Errorf("writing generated DOT graph: %w", err)
	}

	return nil
}

// splitArgs attempts to separate n positional arguments from the rest of the
// given arguments list. The caller is responsible for ensuring that the
// correct number of positional arguments could be extracted.
//
// It is meant as a helper to implement CLI commands of the shape:
//   cmd ARG1 ARG2 [flags]
func splitArgs(n int, args []string) ( /*positional*/ []string /*flags*/, []string) {
	if len(args) == 0 {
		return nil, nil
	}

	// no positional, or user passed only flags (e.g. "cmd -h")
	if n == 0 || (args[0] != "" && args[0][0] == '-') {
		return nil, args
	}

	if len(args) <= n {
		return args, nil
	}

	return args[:n], args[n:]
}
