/*
  This Bridge describes an integration between Amazon CloudWatch Logs and Logz.io.
  Its purpose is to perform a one-way synchronization of log entries between
  these two services.

  Events originating from CloudWatch are sent to a message transformer, which
  is a function that converts incoming data into a format that is consumable by
  the Logz.io event target.

  Additionally, a Content-Based router sends a copy of each incoming event to a
  CloudEvent receiver called Sockeye for monitoring purposes. The Sockeye service
  embeds a web UI which allows received events to be visualized in real time.

  Visual representation:

  ┌─────────────────┐    ┌────────┐    ┌────────────────┐    ┌─────────┐
  │ CloudWatch Logs ├────► Router ├────► Transformation ├────► Logz.io │
  └─────────────────┘    └─────┬──┘    └────────────────┘    └─────────┘
                               │
                               │
                               │ ┌─────────┐
                               └─► Sockeye │
                                 └─────────┘
                                  monitoring
*/

bridge "cloudwatchlogs_to_logz" {}

// ---- Event Sources ----

source "aws_cloudwatch_logs" "my_logs" {
  arn = "arn:aws:logs:us-east-2:123456789012:log-group:/my/log/group:*"
  polling_interval = "1m"

  credentials = secret_name("my-aws-access-keys")

  to = router.dispatch
}

// ---- Event Targets ----

target "logz" "my_logs" {
  logs_listener_url = "listener.logz.io"

  auth = secret_name("my-logz-credentials")
}

target "container" "sockeye" {
  image = "docker.io/n3wscott/sockeye:v0.7.0"
  public = true
}

// ---- Event Routing ----

router "content_based" "dispatch" {

  route {
    attributes = {
      type: "com.amazon.logs.log"
    }
    to = transformer.cloudwatch_logz
  }

  route {
    to = target.sockeye
  }

}

// ---- Event Transformers ----

transformer "function" "cloudwatch_logz" {
  runtime = "js-otto"

  code = <<-EOF
  function handle(input) {
    input.type = "io.triggermesh.logz.ship";
    return input;
  }
  EOF

  ce_context {
    type = "io.triggermesh.logz.ship"
  }

  to = target.my_logs
}
