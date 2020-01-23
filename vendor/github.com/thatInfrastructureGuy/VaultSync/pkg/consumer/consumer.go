package consumer

import (
	"errors"
	"log"
	"time"

	"github.com/thatInfrastructureGuy/VaultSync/pkg/common/data"
	"github.com/thatInfrastructureGuy/VaultSync/pkg/consumer/kubernetes"
)

type Consumers interface {
	GetLastUpdatedDate() (date time.Time, err error)
	PostSecrets(secretList map[string]data.SecretAttribute) error
}

type Consumer struct {
	Destination Consumers
}

func (c *Consumer) PostSecrets(secretList map[string]data.SecretAttribute) (err error, updatedDestination bool) {
	if len(secretList) == 0 {
		log.Println("Nothing to update!")
		return nil, false
	}
	return c.Destination.PostSecrets(secretList), true
}

func (c *Consumer) GetLastUpdatedDate() (date time.Time, err error) {
	return c.Destination.GetLastUpdatedDate()
}

func SelectConsumer(env *data.Env) (c *Consumer, err error) {
	switch env.ConsumerType {
	case "kubernetes":
		if env.SecretName == "" {
			return nil, errors.New("Invalid secret name!")
		}
		c = &Consumer{&kubernetes.Config{
			SecretName: env.SecretName,
			Namespace:  env.Namespace,
		}}
	default:
		return nil, errors.New("No consumer provided.")
	}
	return c, nil
}
