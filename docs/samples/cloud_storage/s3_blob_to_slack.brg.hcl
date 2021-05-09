/*
  This Bridge describes an integration between Amazon S3 / Azure Storage and Slack.

  Events originating from Amazon S3 and Azure Storage are dispatched to their
  respective message transformer by a Content-Based router. Message transformers
  are functions which convert incoming data into a format that is consumable by
  the Slack event target.

  Additionally, the router sends a copy of each incoming event to a CloudEvent
  receiver called Sockeye for monitoring purposes. The Sockeye service embeds a
  web UI which allows received events to be visualized in real time.

  Visual representation:

                                        transformers
  ┌───────────────┐                   ┌───────────────┐
  │ Amazon S3     ├───┐             ┌─► S3 to Slack   ├───┐
  └───────────────┘   │ ┌────────┐  │ └───────────────┘   │ ┌───────────┐
                      ├─► Router ├──┤                     ├─► Slack App │
  ┌───────────────┐   │ └─────┬──┘  │ ┌───────────────┐   │ └───────────┘
  │ Azure Storage ├───┘       │     └─► Blob to Slack ├───┘
  └───────────────┘           │       └───────────────┘
                              │
                              │
                              │ ┌─────────┐
                              └─► Sockeye │
                                └─────────┘
                                 monitoring
*/

// ---- Event Sources ----

source "aws_s3" "my_bucket" {
  arn = "arn:aws:s3:us-east-2:123456789012:my-bucket"

  event_types = [
    "s3:ObjectCreated:*",
    "s3:ObjectRemoved:*"
  ]

  credentials = secret_name("my-aws-access-keys")

  to = router.dispatch
}

source "azure_blob_storage" "my_files" {
  storage_account_id = "/subscriptions/1234/resourceGroups/my-group/providers/Microsoft.Storage/storageAccounts/myfiles"
  event_hub_id = "/subscriptions/1234/resourceGroups/my-group/providers/Microsoft.EventHub/namespaces/myevents"

  event_types = [
    "Microsoft.Storage.BlobCreated",
    "Microsoft.Storage.BlobDeleted",
  ]

  auth = secret_name("my-azure-service-principal")

  to = router.dispatch
}

// ---- Event Targets ----

target "slack" "chat_notifications" {
  auth = secret_name("my-slack-app")
}

target "container" "sockeye" {
  image = "docker.io/n3wscott/sockeye:v0.7.0"
  public = true
}

// ---- Event Routing ----

router "content_based" "dispatch" {

  route {
    attributes = {
      source = "arn:aws:s3:::my-bucket"
    }
    to = transformer.s3_slack_message
  }

  route {
    attributes = {
      source = "/subscriptions/1234/resourceGroups/my-group/providers/Microsoft.Storage/storageAccounts/myfiles"
    }
    to = transformer.blob_slack_message
  }

  route {
    to = target.sockeye
  }

}

// ---- Event Transformers ----

transformer "function" "s3_slack_message" {
  runtime = "python"

  code = <<-EOF
  def main(event, context):
    return {
      "channel": "C0000000000",
      "text": f"Event from S3: `{event['eventName']}`\n"
              f"Bucket: `{event['s3']['bucket']['name']}`\n"
              f"Object: `{event['s3']['object']['key']}`"
    }
  EOF

  entrypoint = "main"

  to = target.chat_notifications
}

transformer "function" "blob_slack_message" {
  runtime = "python"

  code = <<-EOF
  def main(event, context):
    return {
      "channel": "C0000000000",
      "text": f"Event from Azure Storage: `{event['api']}`"
    }
  EOF

  entrypoint = "main"

  to = target.chat_notifications
}
