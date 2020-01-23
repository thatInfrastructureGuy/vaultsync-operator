package secretsmanager

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/providers/checks"
)

type SecretsManager struct {
	result                 *secretsmanager.GetSecretValueOutput
	DestinationLastUpdated time.Time
	VaultName              string
}

func (s *SecretsManager) listSecrets() (err error) {
	//Create a Secrets Manager client
	svc := secretsmanager.New(session.New())
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(s.VaultName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html
	s.result, err = svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			}
		}
		return err
	}
	return nil
}

func (s *SecretsManager) GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error) {
	secretList := make(map[string]data.SecretAttribute)
	err := s.listSecrets()
	if err != nil {
		return nil, err
	}
	secrets, err := s.getSecretValuesJson()
	if err != nil {
		return nil, err
	}
	var keyvalue map[string]interface{}
	err = json.Unmarshal([]byte(secrets), &keyvalue)
	if err != nil {
		return nil, err
	}
	for originalSecretName, valueInterface := range keyvalue {
		//Checks against key metadata
		dateUpdated := *s.result.CreatedDate
		secretName, skipUpdate := checks.CommonProviderChecks(env, originalSecretName, dateUpdated, s.DestinationLastUpdated)
		if skipUpdate {
			break
		}
		value := fmt.Sprintf("%v", valueInterface)
		secretList[secretName] = data.SecretAttribute{
			Value:       value,
			DateUpdated: dateUpdated,
		}
	}
	return secretList, nil
}

func (s *SecretsManager) getSecretValuesJson() (secretsJson string, err error) {
	// Decrypts secret using the associated KMS CMK.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString, decodedBinarySecret string
	if s.result.SecretString != nil {
		secretString = *s.result.SecretString
		return secretString, nil
	}
	decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(s.result.SecretBinary)))
	length, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, s.result.SecretBinary)
	if err != nil {
		fmt.Println("Base64 Decode Error:", err)
		return "", err
	}
	decodedBinarySecret = string(decodedBinarySecretBytes[:length])
	return decodedBinarySecret, nil
}
