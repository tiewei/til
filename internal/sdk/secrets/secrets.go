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

// SecretKeyRefsZendesk returns secret key selectors for the "zendesk" secret class.
func SecretKeyRefsZendesk(secretName string) (token secretKeySelector) {
	return newSecretKeySelector(secretName, secrClassZendeskToken)
}
