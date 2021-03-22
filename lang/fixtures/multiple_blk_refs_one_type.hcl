# This file contains a configuration body with multiple references to the same
# type of block.

some_attr = "some_value"

some_ref = channel.some_channel

errors {
  nested_attr = false
  nested_ref  = channel.other_channel
}
