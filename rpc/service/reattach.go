package service

import (
	"net"
	"os"

	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/runner"
)

type RPCClientInfo struct {
	ReattachConfig ReattachConfig
	PluginName     string
}

func (i *RPCClientInfo) Name() string {
	return i.PluginName
}

// This struct mirrors plugin.ReattachConfig but only includes types that are registered with gob to avoid issues
type ReattachConfig struct {
	Protocol        goplugin.Protocol
	ProtocolVersion int
	Addr            Addr
	Pid             int
	ReattachFunc    runner.ReattachFunc
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
		ReattachFunc:    rc.ReattachFunc,
		Test:            rc.Test,
	}
}

// Reattach to the existing RPC service and return a Service
func Services(serviceInfos map[string]service.ReconnectInfo, svc goplugin.Plugin) (map[string]service.Service, error) {
	command.Logger.Debug("Reconstituting services", "ReconnectInfo", serviceInfos)
	services := make(map[string]service.Service)

	for name, serviceInfo := range serviceInfos {
		rpcClientInfo := serviceInfo.(*RPCClientInfo)
		// Create a new PluginDef using the ReattachConfig instead of a Cmd
		client := goplugin.NewClient(&goplugin.ClientConfig{
			HandshakeConfig: service.HandshakeConfig,
			Plugins: map[string]goplugin.Plugin{
				name: svc,
			},
			Reattach:   rpcClientInfo.ReattachConfig.ReattachConfig(),
			Logger:     service.Logger,
			Managed:    true,      // Allow the plugin runtime to manage this plugin
			SyncStdout: os.Stdout, // Print any extra output to Stdout from the plugin to the host processes Stdout
			// I'd love to use this but I haven't yet figured out out to get it to work with Reattach
			// AutoMTLS:    true,      // Ensure that we're using MTLS for communication between the plugin and the host
			SkipHostEnv: true, // Don't pass the host environment to the plugin to avoid security issues
		})

		service, err := service.ToService(client, name)
		if err != nil {
			return nil, err
		}

		services[service.Name()] = service
	}
	return services, nil
}
