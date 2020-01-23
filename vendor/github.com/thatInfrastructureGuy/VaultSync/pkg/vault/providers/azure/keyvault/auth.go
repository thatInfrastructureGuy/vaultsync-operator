package keyvault

import (
	"time"

	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

type Keyvault struct {
	basicClient            keyvault.BaseClient
	DestinationLastUpdated time.Time
	VaultName              string
}

// initialize creates KeyVault instance
func (k *Keyvault) initialize() (err error) {
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	k.logger()
	k.basicClient = keyvault.New()
	k.basicClient.Authorizer = authorizer
	return nil
}
