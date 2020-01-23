package keyvault

import (
	"context"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/providers/checks"
)

// listSecrets Get all the secrets from specified keyvault
func (k *Keyvault) listSecrets(env *data.Env) (secretList map[string]data.SecretAttribute, err error) {
	ctx := context.Background()
	secretItr, err := k.basicClient.GetSecrets(ctx, "https://"+k.VaultName+".vault.azure.net", nil)
	if err != nil {
		return nil, err
	}

	secretList = make(map[string]data.SecretAttribute)

	for {
		if secretItr.Values() == nil {
			break
		}

		for _, secretProperties := range secretItr.Values() {
			originalSecretName := path.Base(*secretProperties.ID)
			dateUpdated := time.Time(*secretProperties.Attributes.Updated)

			//Checks against key metadata
			updatedSecretName, skipUpdate := checks.CommonProviderChecks(env, originalSecretName, dateUpdated, k.DestinationLastUpdated)
			markForDeletion := customProviderChecks(secretProperties)

			//Get Secret Values
			var secretValue string
			if !markForDeletion {
				secretValue, err = k.getSecretValue(originalSecretName)
				if err != nil {
					return nil, err
				}
			}

			//Create Key-Value map
			if !skipUpdate {
				secretList[updatedSecretName] = data.SecretAttribute{
					DateUpdated:       dateUpdated,
					Value:             secretValue,
					MarkedForDeletion: markForDeletion,
				}
			}
		}

		err = secretItr.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	return secretList, nil
}

func customProviderChecks(secretProperties keyvault.SecretItem) (markForDeletion bool) {
	currentTimeUTC := time.Now().UTC()
	// Check Activation date
	if secretProperties.Attributes.NotBefore != nil {
		activationDate := time.Time(*secretProperties.Attributes.NotBefore)
		if activationDate.After(currentTimeUTC) {
			markForDeletion = true
		}
	}

	// Check Expiry date
	if secretProperties.Attributes.Expires != nil {
		expiryDate := time.Time(*secretProperties.Attributes.Expires)
		if expiryDate.Before(currentTimeUTC) {
			markForDeletion = true
		}
	}

	// Check if secret is disabled.
	isEnabled := *secretProperties.Attributes.Enabled
	if !isEnabled {
		markForDeletion = true
	}
	return markForDeletion
}

// Get SecretValue from KeyVault if Secret is enabled.
// If secret is disabled, return empty string.
func (k *Keyvault) getSecretValue(secretName string) (value string, err error) {
	secretResp, err := k.basicClient.GetSecret(context.Background(), "https://"+k.VaultName+".vault.azure.net", secretName, "")
	if err != nil {
		return "", err
	}

	return *secretResp.Value, nil
}

func (k *Keyvault) GetSecrets(env *data.Env) (secretList map[string]data.SecretAttribute, err error) {
	err = k.initialize()
	if err != nil {
		return nil, err
	}
	secretList, err = k.listSecrets(env)
	if err != nil {
		return nil, err
	}
	return secretList, nil
}
