# This file contains a Bridge description with malformed global settings.

bridge "some_bridge" {
  delivery {
    #! this attribute value is not a number
    retries = "two"
    #! this attribute value is not a traversal expression
    dead_letter_sink = 0
  }
}
