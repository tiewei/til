# TriggerMesh Integration Language Specification

The TriggerMesh Integration Language is a configuration language based on the [HCL syntax][hcl-spec] which purpose is to
offer a user-friendly interface for describing [TriggerMesh Bridges][tm-brg].

## Contents

1. [Configuration Files](#configuration-files)
1. [Attributes and Blocks](#attributes-and-blocks)
1. [Blocks Labels](#block-labels)
1. [Component Identifiers](#component-identifiers)
1. [Block References](#block-references)
1. [Global Configurations](#global-configurations)
1. [Component Categories](#component-categories)
   * [channel](#channel)
   * [router](#router)
   * [transformer](#transformer)
   * [source](#source)
   * [target](#target)

## Configuration Files

A Bridge Description File contains the description of a _single_ Bridge.

The description of a Bridge can only span a _single_ Bridge Description File.

We suggest using the two extensions `.brg.hcl` with Bridge Description Files for the following reasons:

* The `.brg` part allows IDEs to understand that the file describes a TriggerMesh Bridge and enable, if supported,
  _optional_ language-specific integrations, such as autocompletion and live validation.
* The `.hcl` part allows text editors and IDEs to fall back to a _generic_ syntax highlighting for HCL files, since
  `HCL` is a widely supported file format.

## Attributes and Blocks

The language does not support any top-level [attribute][hcl-elems].

The following [block][hcl-elems] type must appear _at most once_ in a configuration file. It contains configurations
which pertain to an entire Bridge. Details are presented in the [Global Configurations](#global-configurations) section.

* `bridge`

The following [block][hcl-elems] types can appear in a configuration file in any order and number of occurrences. Each
of them represents a different _component category_. Details are presented in the [Component
Categories](#component-categories) section.

* `channel`
* `router`
* `transformer`
* `source`
* `target`

## Block Labels

[Labels][hcl-elems] that appear in top-level blocks must be valid [HCL identifiers][hcl-ident], and can therefore be
written either as quoted literal strings or naked identifiers. 

* Both `"foo_bar"` and `foo_bar` are acceptable labels (valid HCL identifiers, quoted and unquoted are equivalent)
* Neither `foo/bar` nor `00foo_bar` are acceptable labels (invalid HCL identifiers)

## Component Identifiers

Labels that represent _component identifiers_ must be unique in a given _component category_.

* There can be both a `channel` and a `router` blocks with the same `foo` identifier.
* There _cannot_ be two `channel` blocks with the same `foo` identifier.

## Block References

Certain types of blocks can contain `to` and/or `reply_to` attributes which are references to other blocks. Whether
those attributes are present at the root of a block or nested inside sub-blocks, their value must be a [variable
expression][hcl-varexpr] composed of a _component identifier_ separated from a _block type_ by an [attribute access
operator][hcl-attrop].

* `channel.my_channel` is a syntactically valid block reference.

## Global Configurations

```hcl
bridge <BRIDGE IDENTIFIER> {
    delivery {
      retries = <integer> // optional
      dead_letter_sink = <block reference> // optional
    }
}
```

A `bridge` block has exactly one label, which represents its _identifier_. This identifier is used to set the Bridge's
components apart from other resources in the destination environment.

A `delivery` block may be set inside a `bridge` block. Its attributes control global aspects of message deliveries:

- `retries`: the minimum number of retries a sender should attempt when sending an event.
- `dead_letter_sink`: component where events that fail to get delivered are moved to.

## Component Categories

Unless otherwise specified, each documented top-level attribute is _required_.

### `channel`

```hcl
channel <CHANNEL TYPE> <CHANNEL IDENTIFIER> {
    # component-type-specific configuration
}
```

### `router`

```hcl
router <ROUTER TYPE> <ROUTER IDENTIFIER> {
    # component-type-specific configuration
}
```

### `transformer`

```hcl
transformer <TRANSFORMER TYPE> <TRANSFORMER IDENTIFIER> {
    # component-type-specific configuration
}
```

### `source`

```hcl
source <SOURCE TYPE> <SOURCE IDENTIFIER> {
    to = <block reference>

    # component-type-specific configuration
}
```

### `target`

```hcl
target <TARGET TYPE> <TARGET IDENTIFIER> {
    reply_to = <block reference> // optional

    # component-type-specific configuration
}
```

[tm-brg]: https://www.triggermesh.com/integrations

[hcl-spec]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md
[hcl-elems]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#structural-elements
[hcl-ident]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#identifiers
[hcl-varexpr]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#variables-and-variable-expressions
[hcl-attrop]: https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#attribute-access-operator
