# This file contains a syntactically valid Bridge definition.
#
# It is used as a "golden file" and should therefore always contain at least
# one valid occurence of each supported block type and top-level attribute.

source "some_source" "MySource" {
  some_block { }

  some_attribute = "xyz"

  to = router.MyRouter
}

router "some_router" "MyRouter" {
  some_block { }

  some_attribute = "xyz"
}

transformer "some_transformer" "MyTransformer" {
  some_block { }

  some_attribute = "xyz"

  to = channel.MyChannel
}

channel "some_channel" "MyChannel" {
  some_block { }

  some_attribute = "xyz"

  to = transformer.MyTransformer
}

function "MyFunction" {
  reply_to = router.MyRouter
}

target "sometarget" "MyTarget" {
  some_block { }

  some_attribute = "xyz"
}
