# This file contains a Bridge description with duplicated global settings.

bridge "some_bridge_id" {}

#! this block is the second occurent of a "bridge" block
bridge "other_bridge_id" {}

source "some_source" "MySource" {
  some_block { }

  some_attribute = "xyz"

  to = target.SomeTarget
}
