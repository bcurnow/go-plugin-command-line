package service

import (
	"encoding/gob"

	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
)

func init() {
	// This is needed for the RPC implementations
	gob.Register(map[string]ReconnectInfo{})
}

// Marker interface for all Services
type Service interface {
	Name() string
	Type() plugin.Type
}

// Logger type of service
type LoggerService interface {
	Service
	Log(val string)
}

// Marker interface for an object that provides information necessary to reconnect to a service
type ReconnectInfo interface {
	Name() string
	Type() plugin.Type
}

type Register interface {
	Register(serviceDir string) (map[string]ReconnectInfo, error)
}
