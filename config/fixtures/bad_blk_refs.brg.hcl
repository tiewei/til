# This file contains a Bridge definition with malformed block references.

source "some_source" "MySource1" {
  some_block { }

  some_attribute = "xyz"

  #! this reference contains more than 2 attributes
  to = target.foo.bar
}

source "some_source" "MySource2" {
  some_block { }

  some_attribute = "xyz"

  #! this reference contains 1 attribute instead of 2
  to = target
}

source "some_source" "MySource3" {
  some_block { }

  some_attribute = "xyz"

  #! this reference refers to a block of an unknown type
  to = foo.bar
}
