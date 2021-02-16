# This file contains a Bridge definition with blocks that are missing required
# attributes.

source "some_source" "MySource" {
  some_block { }

  some_attribute = "xyz"

  #! the required "to" attribute is missing
}

channel "some_channel" "MyChannel" {
  some_block { }

  some_attribute = "xyz"

  #! the required "to" attribute is missing
}

transformer "some_transformer" "MyTransformer" {
  some_block { }

  some_attribute = "xyz"

  #! the required "to" attribute is missing
}
