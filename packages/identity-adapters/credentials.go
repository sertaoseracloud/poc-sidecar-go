package adapters

// Credentials represents cloud auth credentials.
type Credentials struct {
    AccessKeyId     string `json:"accessKeyId,omitempty"`
    SecretAccessKey string `json:"secretAccessKey,omitempty"`
    SessionToken    string `json:"sessionToken,omitempty"`
    Token           string `json:"token,omitempty"`
}

// GetAwsCredentials returns placeholder AWS credentials.
func GetAwsCredentials() (Credentials, error) {
    // TODO: integrate AWS SDK to fetch real credentials
    return Credentials{
        AccessKeyId:     "fake",
        SecretAccessKey: "fake",
        SessionToken:    "fake",
    }, nil
}

// GetAzureToken returns a placeholder Azure token.
func GetAzureToken() (Credentials, error) {
    // TODO: integrate Azure SDK to fetch real token
    return Credentials{Token: "fake"}, nil
}
