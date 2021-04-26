package secrets

// secretKeySelector represents a corev1.SecretKeySelector in a format that can
// be passed to unstructured.SetNestedMap.
type secretKeySelector map[string]interface{}

// newSecretKeySelector returns a secretKeySelector with the given secret name
// and data key.
func newSecretKeySelector(name, key string) secretKeySelector {
	return secretKeySelector{
		"name": name,
		"key":  key,
	}
}

// SecretKeyRefsAWS returns secret key selectors for the "aws" secret class.
func SecretKeyRefsAWS(secretName string) (accessKeyID, secretAccessKey secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassAWSAccessKeyID),
		newSecretKeySelector(secretName, secrClassAWSSecretAccessKey)
}

// SecretKeyRefsBasicAuth returns secret key selectors for the "basic_auth" secret class.
func SecretKeyRefsBasicAuth(secretName string) (usr, passwd secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassBasicAuthUser),
		newSecretKeySelector(secretName, secrClassBasicAuthPasswd)
}

// SecretKeyRefsSalesforceOAuthJWT returns secret key selectors for the "salesforce_oauth_jwt" secret class.
func SecretKeyRefsSalesforceOAuthJWT(secretName string) (key secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassSalesforceOAuthJWTKey)
}

// SecretKeyRefsSlackApp returns secret key selectors for the "slack_app" secret class.
func SecretKeyRefsSlackApp(secretName string) (signSecr secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassSlackAppSignSecr)
}

// SecretKeyRefsKafkaSASL returns secret key selectors for the "kafka_sasl" secret class.
func SecretKeyRefsKafkaSASL(secretName string) (usr, passwd, mechanism secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassKafkaSASLUser),
		newSecretKeySelector(secretName, secrClassKafkaSASLPasswd),
		newSecretKeySelector(secretName, secrClassKafkaSASLMechanism)
}

// SecretKeyRefsTLS returns secret key selectors for the "tls" secret class.
func SecretKeyRefsTLS(secretName string) (cert, key, caCert secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassTLSCert),
		newSecretKeySelector(secretName, secrClassTLSKey),
		newSecretKeySelector(secretName, secrClassTLSCACert)
}

// SecretKeyRefsGitHub returns secret key selectors for the "github" secret class.
func SecretKeyRefsGitHub(secretName string) (accessToken, secretToken secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassGitHubAccessToken),
		newSecretKeySelector(secretName, secrClassGitHubSecretToken)
}

// SecretKeyRefsZendesk returns secret key selectors for the "zendesk" secret class.
func SecretKeyRefsZendesk(secretName string) (token secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassZendeskToken)
}
