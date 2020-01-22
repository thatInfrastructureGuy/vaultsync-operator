package controller

import (
	"github.com/thatinfrastructureguy/vaultsync-operator/pkg/controller/vaultsyncer"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, vaultsyncer.Add)
}
