source "awssqs" "MyQueue" {
  arn = "arn:aws:sqs:us-east-2:123456789012:triggermesh"
  to = router.MyRouter
}

function "MyFunction" {
  repo = "github.com/acme/my-function"
  reply_to = router.MyRouter
}

router "content-based" "MyRouter" {

  route "ToMyFunction" {
    attributes {
      type = "com.amazon.sqs.message"
    }
    to = transformer.MyTransformation
  }

  route "ToKafka" {
    attributes {
      type = "corp.acme.my.processing"
    }
    to = target.MyKafkaTopic
  }

}

transformer "bumblebee" "MyTransformation" {

  context {
    operation "store" {
      path {
        key = "$id"
        value = "id"
      }
    }
    operation "add" {
      path {
        key = "id"
        value = "${person}-${id}"
      }
    }
  }

  data {
    operation "store" {
      path {
        key = "$person"
        value = "Alice"
      }
    }
    operation "add" {
      path {
        key = "event.ID"
        value = "$id"
      }
    }
  }

  to = function.MyFunction
}

target "kafka" "MyKafkaTopic" {
  topic = "myapp"
}
