package secrets

// All secret classes in this file match the secret classes used in the
// TriggerMesh web frontend.

// Secret class "aws"
const (
	secrClassAWSAccessKeyID     = "access_key_id"
	secrClassAWSSecretAccessKey = "secret_access_key"
)

// Secret class "basic_auth"
// https://kubernetes.io/docs/concepts/configuration/secret/#basic-authentication-secret
const (
	secrClassBasicAuthUser   = "username"
	secrClassBasicAuthPasswd = "password"
)

// Secret class "kafka_sasl"
const (
	secrClassKafkaSASLUser      = "username"
	secrClassKafkaSASLPasswd    = "password"
	secrClassKafkaSASLMechanism = "mechanism"
)

// Secret class "salesforce_oauth_jwt"
const (
	secrClassSalesforceOAuthJWTKey = "secret_key"
)

// Secret class "slack_app"
const (
	secrClassSlackAppSignSecr = "signing_secret"
)

// Secret class "tls"
const (
	secrClassTLSCert   = "certificate"
	secrClassTLSKey    = "key"
	secrClassTLSCACert = "ca_certificate"
)

// Secret class "github"
const (
	secrClassGithubAccessToken = "accessToken"
	secrClassGithubSecretToken = "secretToken"
)

// Secret class "zendesk"
const (
	secrClassZendeskToken = "token"
)
