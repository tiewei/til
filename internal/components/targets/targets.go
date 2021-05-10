package targets

// All includes all "target" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"aws_kinesis":    (*AWSKinesis)(nil),
	"aws_lambda":     (*AWSLambda)(nil),
	"aws_s3":         (*AWSS3)(nil),
	"aws_sns":        (*AWSSNS)(nil),
	"aws_sqs":        (*AWSSQS)(nil),
	"container":      (*Container)(nil),
	"datadog":        (*Datadog)(nil),
	"function":       (*Function)(nil),
	"gcloud_storage": (*GCloudStorage)(nil),
	"kafka":          (*Kafka)(nil),
	"logz":           (*Logz)(nil),
	"slack":          (*Slack)(nil),
	"twilio":         (*Twilio)(nil),
	"zendesk":        (*Zendesk)(nil),
}
