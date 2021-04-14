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
