# Sample Bridge Description File

source "aws_sqs" "my_queue" {
  arn = "arn:aws:sqs:us-east-2:123456789012:triggermesh"

  access_key = "AKIA0000000000000000"
  secret_key = "TG9yZW0gaXBzdW0gZG9sb3Igc2l0IHZpdmVycmEu"

  to = router.my_router
}

router "content_based" "my_router" {

  route {
    attributes = {
      type = "com.amazon.sqs.message"
    }
    to = transformer.my_transformation
  }

  route {
    attributes = {
      type = "corp.acme.my.processing"
    }
    to = router.even_uid
  }

}

router "data_expression_filter" "even_uid" {
  condition = "$user.id.(int64) % 2 == 0"
  to = target.custom_logic
}

transformer "bumblebee" "my_transformation" {

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
        value = "$${person}-$${id}"
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

  to = target.custom_logic
}

target "kafka" "my_kafka_topic" {
  topic = "myapp"
}

target "function" "custom_logic" {
  runtime = "python"
  entrypoint = "foo"
  public = true
  code =<<EOF
import urllib.request

def foo(event, context):
  resp = urllib.request.urlopen(event['url'])
  page = resp.read()

  response = {
    "statusCode": resp.status,
    "body": str(page)
  }

return response
  EOF
}
