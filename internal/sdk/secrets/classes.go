package secrets

// All secret classes in this file match the secret classes used in the
// TriggerMesh web frontend.

// Secret class "aws"
const (
	secrClassAWSAccessKeyID     = "access_key_id"
	secrClassAWSSecretAccessKey = "secret_access_key"
)

// Secret class "kafka_sasl"
const (
	secrClassKafkaSASLUser      = "username"
	secrClassKafkaSASLPasswd    = "password"
	secrClassKafkaSASLMechanism = "mechanism"
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
