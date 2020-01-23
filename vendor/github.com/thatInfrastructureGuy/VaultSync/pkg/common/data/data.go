package data

import "time"

// SecretAttribute is constructed after querying Vault for each secret.
// It contains various attributes of secret other than values.
type SecretAttribute struct {
	DateUpdated       time.Time
	Value             string
	MarkedForDeletion bool
}
