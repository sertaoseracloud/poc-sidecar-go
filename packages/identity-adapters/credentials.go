package adapters

import (
       "context"
       "errors"
       "os"

       "github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
       "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
       "github.com/aws/aws-sdk-go-v2/credentials"
)

// Credentials represents cloud auth credentials.
type Credentials struct {
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
	SessionToken    string `json:"sessionToken,omitempty"`
	Token           string `json:"token,omitempty"`
}

// GetAwsCredentials retrieves AWS credentials using the default credential chain.
func GetAwsCredentials() (Credentials, error) {
       accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
       secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
       sessionToken := os.Getenv("AWS_SESSION_TOKEN")

       provider := credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken)
       creds, err := provider.Retrieve(context.Background())
       if err != nil {
               return Credentials{}, err
       }

       if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
               return Credentials{}, errors.New("AWS credentials not set")
       }

       return Credentials{
               AccessKeyId:     creds.AccessKeyID,
               SecretAccessKey: creds.SecretAccessKey,
               SessionToken:    creds.SessionToken,
       }, nil
}

// GetAzureToken retrieves an Azure access token using a service principal.
func GetAzureToken() (Credentials, error) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return Credentials{}, err
	}
	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return Credentials{}, err
	}
	return Credentials{Token: token.Token}, nil
}