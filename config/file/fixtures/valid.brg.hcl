# This file contains a syntactically valid Bridge description.
#
# It is used as a "golden file" and should therefore always contain at least
# one valid occurence of each supported block type and top-level attribute.

bridge "some_bridge" {
  delivery {
    retries = 2
    dead_letter_sink = channel.foo
  }
}

source some_source "MySource" {
  some_block { }

  some_attribute = "xyz"

  to = router.MyRouter
}

router some_router "MyRouter" {
  some_block { }

  some_attribute = "xyz"
}

transformer some_transformer "MyTransformer" {
  some_block { }

  some_attribute = "xyz"

  to = channel.MyChannel
}

channel some_channel "MyChannel" {
  some_block { }

  some_attribute = "xyz"
}

target sometarget "MyTarget" {
  some_block { }

  some_attribute = "xyz"
}
