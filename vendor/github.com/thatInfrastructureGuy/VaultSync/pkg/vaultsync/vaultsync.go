package vaultsync

import (
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/consumer"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/vault"
)

func Synchronize(env *data.Env) (error, bool) {
	// Select the destination
	destination, err := consumer.SelectConsumer(env)
	if err != nil {
		return err, false
	}

	// Get lastUpdated date timestamp from consumer
	destinationlastUpdated, err := destination.GetLastUpdatedDate()
	if err != nil {
		return err, false
	}

	// Select the source
	source, err := vault.SelectProvider(env, destinationlastUpdated)
	if err != nil {
		return err, false
	}
	// Poll secrets from vault which were updated since lastUpdated value
	secretList, err := source.GetSecrets(env)
	if err != nil {
		return err, false
	}

	// Update kuberenetes secrets
	err, updatedDestination := destination.PostSecrets(secretList)
	if err != nil {
		return err, updatedDestination
	}
	return nil, updatedDestination
}
