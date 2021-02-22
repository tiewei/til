# This file contains a Bridge description with a block of an unknown type.

source "some_source" "MySource" {
  some_block { }

  some_attribute = "xyz"

  to = target.MyTarget
}

#! the "foo" block type is not part of the spec
foo "some_foo" "MyFoo" {
  some_block { }

  some_attribute = "xyz"

  to = target.MyTarget
}
