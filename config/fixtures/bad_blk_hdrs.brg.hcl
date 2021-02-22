# This file contains a Bridge description with malformed block headers.

#! this block contains too many labels, it won't be evaluated
source "some_source" "MySource1" "another_label" {
  some_block { }

  some_attribute = "xyz"

  to = target.foo
}

#! this block contains too few labels, it won't be evaluated
source "MySource2" {
  some_block { }

  some_attribute = "xyz"

  to = target.foo
}

#! this block's labels are invalid identifiers
source "some source" "My/Source3" {
  some_block { }

  some_attribute = "xyz"

  to = target.foo
}
