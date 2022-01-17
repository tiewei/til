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

package targets

// All includes all "target" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"aws_dynamodb":     (*AWSDynamoDB)(nil),
	"aws_kinesis":      (*AWSKinesis)(nil),
	"aws_lambda":       (*AWSLambda)(nil),
	"aws_s3":           (*AWSS3)(nil),
	"aws_sns":          (*AWSSNS)(nil),
	"aws_sqs":          (*AWSSQS)(nil),
	"container":        (*Container)(nil),
	"datadog":          (*Datadog)(nil),
	"event_display":    (*EventDisplay)(nil),
	"function":         (*Function)(nil),
	"ibmmq":            (*IBMMQ)(nil),
	"gcloud_storage":   (*GCloudStorage)(nil),
	"gcloud_firestore": (*GCloudFirestore)(nil),
	"kafka":            (*Kafka)(nil),
	"logz":             (*Logz)(nil),
	"sendgrid":         (*Sendgrid)(nil),
	"slack":            (*Slack)(nil),
	"sockeye":          (*Sockeye)(nil),
	"splunk":           (*Splunk)(nil),
	"twilio":           (*Twilio)(nil),
	"zendesk":          (*Zendesk)(nil),
}
