# Software Architecture

This document describes the architecture of the `bridgedl` interpreter, and explains the design choices behind it.

## Contents

1. [Overview](#overview)
   * [The Component Graph](#the-component-graph)
   * [Loose Comparisons with Terraform](#loose-comparisons-with-terraform)
1. [Request Flow](#request-flow)
1. [Graph Builder](#graph-builder)
1. [Components Implementations](#components-implementations)
1. [Bridge Translator](#bridge-translator)

## Overview

### The Component Graph

The architecture of `bridgedl` borrows a few ideas from [HashiCorp Terraform][tf-arch]. For instance, we reason about a
Bridge description as a [directed graph][wiki-graph] where each component of a messaging system is a _vertex_ (_"node"_)
connected to one or more other components by an _edge_ (_"link"_) that represents the flow of events in that system
(Terraform calls it the ["Resource Graph"][tf-graph]).

A graph is also a good fit for HCL's evaluation model, where the dependencies of each [expression][hcl-expr] should be
analyzed in order to infer the evaluation context. For example, the evaluation of a `router` block containing a
reference to a `transformer` block can only be performed _after_ the `transformer` block itself is evaluated. The HCL
guide refers to this pattern as [Interdependent Blocks][hcl-idpblk].

However, this is where the comparison with Terraform ends due to fundamental differences in what both of these tools try
to achieve.

### Loose Comparisons with Terraform

Like Terraform, `bridgedl`...

* Interprets configuration files written in a language based on the [HCL syntax][hcl-spec] (the [_Bridge Description
  Language_][bdl-spec]).
* Supports cross-references between HCL configuration blocks using traversal expressions.
* Represents configuration blocks and the relationships between them as a directed graph.
* Uses internal schemas ([`hcldec.Spec`][hcldec-spec]) for decoding the contents of HCL blocks into concrete
  types/values.

Unlike Terraform, `bridgedl`...

* Is not concerned about the deployment aspects of the components it interprets. The interpreter generates deployment
  manifests which users are free to deploy using the tool of their choice (kubectl, Helm, Terraform, ...).
* Does not depend on side-effects during runtime to evaluate certain parts of a configuration (such as creating a
  certain resource in order to be able to determine the value of a variable). A _Bridge Description File_ is interpreted
  statically.
* Leverages graphs to represent data flows, not dependencies.

## Request Flow

Below is a high level representation of the internal request flow within the interpreter when a user generates
deployment manifests from a Bridge description. Each subsystem is described in more details in the rest of this
document.

![Architecture - Request flow][arch-graph]

## Graph Builder

The role of the [`core.GraphBuilder`][pkgcore-graphb] is to place each component of a Bridge onto a graph and connect it
to other components, based on the information contained in its HCL block. 

This is achieved by successive transformations performed by implementers of the
[`core.GraphTransformer`][pkgcore-grapht] interface. In the process, the graph builder attaches to each graph vertex any
information susceptible to be useful to translators (e.g. a schema to decode the contents of the corresponding HCL
block, looked up in a map of supported component implementations). The inspiration for this iterative approach comes
from Terraform's design, and aims at making the process of building the graph more testable, even though the complexity
of `bridgedl` is not comparable to Terraform's.

In order to avoid leaking the "graph" domain into the "configuration" domain, and having to implement large interfaces
for each type in the [`config`][pkgconfig] package, we leverage Go's behavioral polymorphism by wrapping `config` types
(such as `config.Channel`, `config.Source`, etc.) into structs that represent graph vertices. Each vertex type (e.g.
[`core.ChannelVertex`][pkgcore-chvrtx]) implements methods that are only relevant to the [`core`][pkgcore] package, and
generally to a single `core.GraphTransformer`. For instance, if a type of component can be the destination of events
within the Bridge, and by definition exposes a Knative [Addressable][kn-addr] type, its vertex type will implement
[`core.AddressableVertex`][pkgcore-addrvrtx].

## Components Implementations

For a component type to be able to appear in a Bridge description, it needs to have a corresponding _implementation_
that allows the interpreter to decode its HCL configuration and translate it into one or more Kubernetes API objects.

Such implementation combines three Go interfaces defined in the [`translation`][pkgtransl] package: `Decodable`,
`Translatable` and `Addressable`. The semantics of those interfaces maps directly into the types of operations performed
by the interpreter:

* `Decodable` provides the [`hcldec.Spec`][hcldec-spec] that is used by HCL APIs (such as [`hcldec.Decode`][hcldec-dec])
  for decoding component-type-specific segments of HCL configuration into data that can be interpreted in the Go
  language. As long as a component requires some kind of non-generic configuration, it must implement this interface.
* `Translatable` bridges the gap between a HCL configuration block and its representation in the Kubernetes space by
  generating values that represent Kubernetes objects, from configurations decoded from the Bridge description. All
  components implement this interface.
* `Addressable` allows determining the address at which a component accepts events. Components that can ingest events
  implement this interface so that HCL reference expressions in other components can be evaluated to actual addresses.

As a means of comparisons, Terraform has a concept of "provider" which also interacts with custom ["resource"
types][tf-res] by using decode schemas and calling CRUD functions on the decoded configuration data.

Currently, all components implementations reside in-tree, inside sub-packages of [`internal/components`][pkgcomps].

## Bridge Translator

The [`core.BridgeTranslator`][pkgcore-brgtrsl] is responsible for evaluating the simple connected graph built by
`core.GraphBuilder`, and translating it into a list of Kubernetes manifests which can be deployed to a target
environment (cluster) to materialize the Bridge.

It does so by creating a topological ordering of the graph which ensures, whenever possible, that the address (event
destination) of a given Bridge component is readily available when components which depend on that address are
evaluated/translated. This ordering allows the translation to be accurate without having to visit a graph vertex
multiple times (e.g. once to store its event address in the evaluation context, once to translate it after all addresses
within the Bridge have been determined).

Event addresses are stored as variables inside a [`core.Evaluator`][pkgcore-eval], which accumulates them as the graph
is being walked.

![Evaluation - Topological sort][topo-sort-graph]

<!-- internal links -->
[arch-graph]: .assets/arch-request-flow.svg
[topo-sort-graph]: .assets/topological-sort.svg
[bdl-spec]: language-spec.md
[pkgcore-graphb]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/graph.go#L10-L15
[pkgcore-grapht]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/graph.go#L52-L55
[pkgcore-brgtrsl]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/translate.go#L13-L18
[pkgcore-eval]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/eval_context.go#L12-L27
[pkgcore]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/doc.go
[pkgcore-chvrtx]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/vertex_channel.go#L16-L26
[pkgcore-addrvrtx]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/core/transform_connect_refs.go#L11-L25
[pkgconfig]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/config/doc.go
[pkgcomps]: ../internal/components
[pkgtransl]: https://github.com/triggermesh/bridgedl/blob/bf0d78bd/translation/doc.go

<!-- HashiCorp links -->
[tf-arch]: https://github.com/hashicorp/terraform/blob/main/docs/architecture.md
[tf-graph]: https://www.terraform.io/docs/internals/graph.html
[tf-res]: https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema#Resource
[hcl-spec]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
[hcl-expr]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#expressions
[hcl-idpblk]: https://hcl.readthedocs.io/en/latest/go_patterns.html#interdependent-blocks
[hcldec-spec]: https://pkg.go.dev/github.com/hashicorp/hcl/v2/hcldec#Spec
[hcldec-dec]: https://pkg.go.dev/github.com/hashicorp/hcl/v2/hcldec#Decode

<!-- misc links -->
[kn-addr]: https://pkg.go.dev/github.com/knative/pkg/apis/duck/v1#Addressable
[wiki-graph]: https://en.wikipedia.org/wiki/Directed_graph
