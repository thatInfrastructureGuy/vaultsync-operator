package checks

import (
	"strings"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
)

func CommonProviderChecks(env *data.Env, originalSecretName string, sourceDate time.Time, destinationDate time.Time) (updatedSecretName string, skipUpdate bool) {
	// Set updatedName as original name
	updatedSecretName = originalSecretName
	// Check if destination keys are outdated.
	if !sourceDate.After(destinationDate) {
		skipUpdate = true
	}
	// Check if ALL hyphers should be converted to underscores
	if env.ConvertHyphenToUnderscores {
		updatedSecretName = strings.ReplaceAll(originalSecretName, "-", "_")
	}
	return updatedSecretName, skipUpdate
}
