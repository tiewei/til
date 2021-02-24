# This file contains a Bridge description with duplicated block identifiers.

source "some_source" "MySource" {
  some_block { }

  some_attribute = "xyz"

  to = target.SomeTarget
}

#! this block uses the same identifier as the block above
source "other_source" "MySource" {
  some_block { }

  some_attribute = "xyz"

  to = target.OtherTarget
}
