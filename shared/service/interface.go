package service

import "github.com/bcurnow/go-plugin-command-line/shared/plugin"

// Marker interface for all Services
type Service interface {
	Name() string
	// Ideally, this would default a more useful set of methods, this is just an example
	Log(val string)
}
type Register interface {
	Register(serviceDir string) (map[string]plugin.Reattach, error)
}
