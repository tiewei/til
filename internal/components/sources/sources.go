package sources

// All includes all "source" component types supported by TriggerMesh.
var All = map[string]interface{}{
	"aws_codecommit":       (*AWSCodeCommit)(nil),
	"aws_cognito_userpool": (*AWSCognitoUserPool)(nil),
	"aws_dynamodb":         (*AWSDynamoDB)(nil),
	"aws_kinesis":          (*AWSKinesis)(nil),
	"aws_s3":               (*AWSS3)(nil),
	"aws_sqs":              (*AWSSQS)(nil),
}
