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

// Secret class "gcloud_service_account"
const (
	secrClassGCloudSvcAccountKey = "key.json"
)

// Secret class "github"
const (
	secrClassGitHubAccessToken   = "access_token"
	secrClassGitHubWebhookSecret = "webhook_secret"
)

// Secret class "kafka"
const (
	secrClassKafkaSASLMechanism = "sasl.mechanism"
	secrClassKafkaSASLUser      = "user"
	secrClassKafkaSASLPasswd    = "password"
	secrClassKafkaTLSCACert     = "ca.crt"
	secrClassKafkaTLSCert       = "user.crt"
	secrClassKafkaTLSKey        = "user.key"
)

// Secret class "logz"
const (
	secrClassLogzAPIToken = "token"
)

// Secret class "salesforce_oauth_jwt"
const (
	secrClassSalesforceOAuthJWTKey = "secret_key"
)

// Secret class "slack"
const (
	secrClassSlackAPIToken = "token"
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

// Secret class "zendesk"
const (
	secrClassZendeskToken = "token"
)
