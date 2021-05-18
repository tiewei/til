/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sources

// All includes all "source" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"aws_cloudwatch":           (*AWSCloudWatch)(nil),
	"aws_cloudwatch_logs":      (*AWSCloudWatchLogs)(nil),
	"aws_codecommit":           (*AWSCodeCommit)(nil),
	"aws_cognito_userpool":     (*AWSCognitoUserPool)(nil),
	"aws_dynamodb":             (*AWSDynamoDB)(nil),
	"aws_kinesis":              (*AWSKinesis)(nil),
	"aws_performance_insights": (*AWSPerformanceInsights)(nil),
	"aws_s3":                   (*AWSS3)(nil),
	"aws_sns":                  (*AWSSNS)(nil),
	"aws_sqs":                  (*AWSSQS)(nil),
	"azure_activity_logs":      (*AzureActivityLogs)(nil),
	"azure_blob_storage":       (*AzureBlobStorage)(nil),
	"azure_event_hubs":         (*AzureEventHubs)(nil),
	"github":                   (*GitHub)(nil),
	"httppoller":               (*HTTPPoller)(nil),
	"kafka":                    (*Kafka)(nil),
	"ping":                     (*Ping)(nil),
	"slack":                    (*Slack)(nil),
	"salesforce":               (*Salesforce)(nil),
	"webhook":                  (*Webhook)(nil),
	"zendesk":                  (*Zendesk)(nil),
}
