package sources

// All includes all "source" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"aws_codecommit":       (*AWSCodeCommit)(nil),
	"aws_cognito_userpool": (*AWSCognitoUserPool)(nil),
	"aws_dynamodb":         (*AWSDynamoDB)(nil),
	"aws_kinesis":          (*AWSKinesis)(nil),
	"aws_s3":               (*AWSS3)(nil),
	"aws_sqs":              (*AWSSQS)(nil),
	"azure_activity_logs":  (*AzureActivityLogs)(nil),
	"azure_blob_storage":   (*AzureBlobStorage)(nil),
	"azure_event_hubs":     (*AzureEventHubs)(nil),
	"github":               (*GitHub)(nil),
	"httppoller":           (*HTTPPoller)(nil),
	"kafka":                (*Kafka)(nil),
	"ping":                 (*Ping)(nil),
	"slack":                (*Slack)(nil),
	"salesforce":           (*Salesforce)(nil),
	"webhook":              (*Webhook)(nil),
	"zendesk":              (*Zendesk)(nil),
}
