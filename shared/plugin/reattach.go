package plugin

import (
	"encoding/gob"
	"net"

	goplugin "github.com/hashicorp/go-plugin"
)

func init() {
	// This is needed for the RPC implementations
	gob.Register(map[string]Reattach{})
}

type Reattach struct {
	ReattachConfig ReattachConfig
	PluginName     string
}

// This struct mirrors plugin.ReattachConfig but only includes types that are registered with gob to avoid issues
type ReattachConfig struct {
	Protocol        goplugin.Protocol
	ProtocolVersion int
	Addr            Addr
	Pid             int
	Test            bool
}

// A version of net.UnixAddr without using interfaces so it works properly with gob
type Addr struct {
	Net  string
	Name string
}

func (rc *ReattachConfig) ReattachConfig() *goplugin.ReattachConfig {
	return &goplugin.ReattachConfig{
		Protocol:        rc.Protocol,
		ProtocolVersion: rc.ProtocolVersion,
		Addr:            &net.UnixAddr{Name: rc.Addr.Name, Net: rc.Addr.Net},
		Pid:             rc.Pid,
		Test:            rc.Test,
	}
}
