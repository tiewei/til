package secrets

// secretKeySelector represents a corev1.SecretKeySelector in a format that can
// be passed to unstructured.SetNestedMap.
type secretKeySelector map[string]interface{}

// SecretKeyRefsAWS returns secret key selectors for the "aws" secret class.
func SecretKeyRefsAWS(secretName string) (accessKeyID, secretAccessKey secretKeySelector) {
	return secretKeySelector{
			"name": secretName,
			"key":  secrClassAWSAccessKeyID,
		},
		secretKeySelector{
			"name": secretName,
			"key":  secrClassAWSSecretAccessKey,
		}
}

// SecretKeyRefsConfluent returns secret key selectors for the "confluent" secret class.
func SecretKeyRefsConfluent(secretName string) (passwd secretKeySelector) {
	return secretKeySelector{
		"name": secretName,
		"key":  secrClassConfluentPasswd,
	}
}

// SecretKeyRefsKafkaSASL returns secret key selectors for the "kafka_sasl" secret class.
func SecretKeyRefsKafkaSASL(secretName string) (usr, passwd, typ secretKeySelector) {
	return secretKeySelector{
			"name": secretName,
			"key":  secrClassKafkaSASLUser,
		},
		secretKeySelector{
			"name": secretName,
			"key":  secrClassKafkaSASLPasswd,
		},
		secretKeySelector{
			"name": secretName,
			"key":  secrClassKafkaSASLType,
		}

}

// SecretKeyRefsTLS returns secret key selectors for the "tls" secret class.
func SecretKeyRefsTLS(secretName string) (cert, key, caCert secretKeySelector) {
	return secretKeySelector{
			"name": secretName,
			"key":  secrClassTLSCert,
		},
		secretKeySelector{
			"name": secretName,
			"key":  secrClassTLSKey,
		},
		secretKeySelector{
			"name": secretName,
			"key":  secrClassTLSCACert,
		}
}

// SecretKeyRefsZendesk returns secret key selectors for the "zendesk" secret class.
func SecretKeyRefsZendesk(secretName string) (token secretKeySelector) {
	return secretKeySelector{
		"name": secretName,
		"key":  secrClassZendeskToken,
	}
}
