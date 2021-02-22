# This file contains a Bridge description with malformed block references.

source "some_source" "MySource1" {
  some_block { }

  some_attribute = "xyz"

  #! this reference is not static (template expression)
  to = "${foo.bar}"
}

source "some_source" "MySource2" {
  some_block { }

  some_attribute = "xyz"

  #! this reference is not static (calculation)
  to = foo+bar
}

source "some_source" "MySource3" {
  some_block { }

  some_attribute = "xyz"

  #! this reference is not static (function call)
  to = foo(bar)
}
