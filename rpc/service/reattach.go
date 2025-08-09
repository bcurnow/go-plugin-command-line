package service

import (
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
)

// Reattach to the existing RPC service and return a Service
func Services(reattaches map[string]plugin.Reattach, impl goplugin.Plugin) (map[string]service.Service, error) {
	command.Logger.Debug("Reconstituting services", "Reattaches", reattaches)
	services := make(map[string]service.Service)

	for name, reattach := range reattaches {
		client := plugin.Client("service", &plugin.ReattachClientConfigBuilder{
			BaseClientConfigBuilder: plugin.BaseClientConfigBuilder{
				HandshakeConfig: service.HandshakeConfig,
				Plugins: map[string]goplugin.Plugin{
					name: impl,
				},
				Name:             name,
				Logger:           service.Logger,
				AllowedProtocols: []goplugin.Protocol{goplugin.ProtocolNetRPC},
			},
			ReattachConfig: reattach.ReattachConfig,
		})

		service, err := service.ToService(client, name)
		if err != nil {
			return nil, err
		}

		services[service.Name()] = service
	}
	return services, nil
}
