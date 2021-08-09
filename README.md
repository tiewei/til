# TriggerMesh Integration Language

Specification and interpreter for TriggerMesh's Integration Language.

The TriggerMesh Integration Language (TIL) is a configuration language based on the [HCL syntax][hcl-spec] which purpose
is to provide a user-friendly interface for describing [TriggerMesh Bridges][tm-brg].

Using the `til` CLI tool, it is possible to turn Bridge definitions into deployment manifests which can run
complete messaging systems on the TriggerMesh platform.

## Documentation

For instructions about the usage of the language and its tooling, please refer to the [Wiki][wiki] _(temporary
location)_.

For a catalog of sample Bridge descriptions, please refer to the [`docs/samples/`](./docs/samples/) directory.

For details about the language specification and other technical documents about the interpreter, please refer to the
[`docs/`](./docs) directory.

## Getting started

The interpreter is written in the [Go programming language][go] and leverages the [HCL toolkit][hcl].

After installing the [Go toolchain][go-dl] (version 1.16 or above), the interpreter can be compiled for the current OS
and architecture by executing the following command inside the root of the repository (the `main` Go package):

```console
$ go build .
```

The above command creates an executable called `til` inside the current directory.

The `-h` or `--help` flag can be appended to any command or subcommand to print some usage instructions about that
command:

```console
$ ./til --help
```

## Contributions and support

We would love to hear your feedback. Please don't hesitate to submit bug reports and suggestions by
[filing issues][gh-issue], or contribute by [submitting pull-requests][gh-pr].

## Commercial Support

TriggerMesh Inc. supports TIL commercially. Email us at <info@triggermesh.com> to get more details.

## Code of Conduct

Although this project is not part of the [CNCF][cncf], we abide by its [code of conduct][cncf-conduct].

[gh-issue]: https://github.com/triggermesh/til/issues
[gh-pr]: https://github.com/triggermesh/til/pulls

[cncf]: https://www.cncf.io/
[cncf-conduct]: https://github.com/cncf/foundation/blob/master/code-of-conduct.md

[tm-brg]: https://www.triggermesh.com/integrations
[wiki]: https://github.com/triggermesh/til/wiki

[go]: https://golang.org/
[go-dl]: https://golang.org/dl/

[hcl]: https://github.com/hashicorp/hcl
[hcl-spec]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
