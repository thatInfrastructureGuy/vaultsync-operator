package vault

import (
	"errors"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault/providers/aws/secretsmanager"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault/providers/azure/keyvault"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault/providers/local"
)

type Vaults interface {
	GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error)
}

type Vault struct {
	Provider Vaults
}

func (v *Vault) GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error) {
	return v.Provider.GetSecrets(env)
}

func SelectProvider(env *data.Env, lastUpdated time.Time) (v *Vault, err error) {
	switch env.Provider {
	case "azure":
		v = &Vault{&keyvault.Keyvault{DestinationLastUpdated: lastUpdated, VaultName: env.VaultName}}
	case "aws":
		v = &Vault{&secretsmanager.SecretsManager{DestinationLastUpdated: lastUpdated, VaultName: env.VaultName}}
	case "gcp":
		return nil, errors.New("Google Secrets Manager: Not implemented yet!")
	case "hashicorp":
		return nil, errors.New("Hashicorp Vault: Not implemented yet!")
	case "local":
		v = &Vault{&local.Local{DestinationLastUpdated: lastUpdated}}
	default:
		return nil, errors.New("Please specify valid vault provider: azure, aws. (Coming soon: gcp, hashicorp)")
	}
	return v, nil
}
