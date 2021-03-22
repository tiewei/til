# This file contains a configuration body with multiple block references
# referencing different categories of Bridge components.

some_attr = "some_value"

some_ref = router.some_router

errors {
  nested_attr = false
  nested_ref  = channel.some_channel
}
