package util

import (
	"encoding/gob"

	"github.com/bcurnow/go-plugin-command-line/rpc/service"
)

// The entire function of this package is to ensure that the right types are registered with gob
func init() {
	gob.Register(&service.RPCClientInfo{})
}
