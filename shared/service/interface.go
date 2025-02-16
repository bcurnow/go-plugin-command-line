package service

import (
	"encoding/gob"
)

func init() {
	// This is needed for the RPC implementations
	gob.Register(map[string]ReconnectInfo{})
}

// Marker interface for all Services
type Service interface {
	Name() string
	// Ideally, this would default a more useful set of methods, this is just an example
	Log(val string)
}

// Marker interface for an object that provides information necessary to reconnect to a service
type ReconnectInfo interface {
	Name() string
}

type Register interface {
	Register(serviceDir string) (map[string]ReconnectInfo, error)
}
