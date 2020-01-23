package local

import (
	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"time"
)

type Local struct {
	DestinationLastUpdated time.Time
}

func (l *Local) GetSecrets(env *data.Env) (map[string]data.SecretAttribute, error) {
	sampleData := map[string]data.SecretAttribute{
		"key1": {
			DateUpdated: time.Now(),
			Value:       "value1",
		},
	}

	return sampleData, nil
}
