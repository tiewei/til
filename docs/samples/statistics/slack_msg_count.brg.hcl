/*
  This Bridge describes an integration between Slack and Datadog.
  Its purpose is to record a metric in Datadog whenever a user sends a message
  in Slack.

  Events originating from Slack are sent to a message transformer, which is a
  functions that converts incoming data into a metric payload that is
  consumable by the Datadog event target.

  Whenever a metric fails to be created in Datadog, an error response is
  returned to a Content-Based router, which dispatches it to Zendesk via
  another message transformer.

  Additionally, the router sends a copy of each incoming event to a CloudEvent
  receiver called Sockeye for monitoring purposes. The Sockeye service embeds a
  web UI which allows received events to be visualized in real time.

  Visual representation:

                                 transformers
                             ┌──────────────────┐     ┌─────────┐
                   ┌─────────► Error to Zendesk ├─────► Zendesk │
                   │         └──────────────────┘     └─────────┘
                   │
  ┌───────┐    ┌───┴────┐    ┌───────────────────┐    ┌─────────┐
  │ Slack ├────► Router ├────► Slack to metric   ├────► Datadog │
  └───────┘    └───┬──▲─┘    └───────────────────┘    └─────────┘
                   │  ╵                                       ╵
                   │  └ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ╶ ┘
                   │                          error responses
                   │
                   │   ┌─────────┐
                   └───► Sockeye │
                       └─────────┘
                        monitoring
*/

bridge "slack_msg_count" {}

// ---- Event Sources ----

source "slack" "general_channel" {
  signing_secret = secret_name("my-slack-app")
  app_id = "A12345"

  to = router.dispatch
}

// ---- Event Targets ----

target "datadog" "slack_stats" {
  metric_prefix = "slackgeneral"
  auth = secret_name("my-datadog-credentials")

  reply_to = router.dispatch
}

target "zendesk" "datadog_errors" {
  subdomain = "example-corp"

  email = "admin@example.com"
  api_auth = secret_name("zendesktarget-secret")

  subject = "Ticket from TriggerMesh"
}

target "container" "sockeye" {
  image = "docker.io/n3wscott/sockeye:v0.7.0"
  public = true
}

// ---- Event Routing ----

router "content_based" "dispatch" {
  
  route {
    attributes = {
      type: "com.slack.events"
    }
    to = transformer.slack_datadog_metric
  }

  route {
    attributes = {
      type: "io.triggermesh.datadog.response"
    }
    to = transformer.handle_datadog_responses
  }

  route {
    to = target.sockeye
  }

}

// ---- Event Transformers ----

transformer "function" "slack_datadog_metric" {
  runtime = "js"

  code = <<-EOF
  function handle(input) {
    input.type = "io.triggermesh.datadog.metric.submit";
    input.data.series = [{"metric": input.data.user, "points": [[Date.now(),"1"]]}];
    return input;
  }
  EOF

  ce_context {
    type = "io.triggermesh.datadog.metric.submit"
  }

  to = target.slack_stats
}

transformer "function" "handle_datadog_responses" {
  runtime = "js"

  code = <<-EOF
  function handle(input) {
    input.type = "com.zendesk.ticket.create";
    input.data.subject = "An error occured while creating a Datadog metric";
    input.data.body = input.data;
    return input;
  }
  EOF

  ce_context {
    type = "com.zendesk.ticket.create"
  }

  to = target.datadog_errors
}
