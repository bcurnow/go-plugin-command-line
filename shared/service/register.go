package service

import (
	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	goplugin "github.com/hashicorp/go-plugin"
)

var (
	Logger = logging.Logger().Named("service")

	// This is the go-plugin handshake information that needs to be used for all plugins
	HandshakeConfig = goplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "SERVICE_PLUGIN",
		MagicCookieValue: "a44a5c99-fac2-4ff4-9a47-3eec6283725a",
	}
)

// Traverses the dir and registers any executable found as a Service
func RegisterDir(dir string, service goplugin.Plugin, allowedProtocols []goplugin.Protocol, serviceInfo func(client *goplugin.Client, name string, service Service) plugin.Reattach) (map[string]plugin.Reattach, error) {
	serviceInfos := make(map[string]plugin.Reattach)
	err := plugin.Register("services", dir, func(pluginName string, pluginCmd string) error {
		client := plugin.Client("service", &plugin.CommandClientConfigBuilder{
			BaseClientConfigBuilder: plugin.BaseClientConfigBuilder{
				HandshakeConfig: HandshakeConfig,
				Plugins: map[string]goplugin.Plugin{
					pluginName: service,
				},
				Name:             pluginName,
				Logger:           Logger,
				AllowedProtocols: allowedProtocols,
			},
			PluginCmd: pluginCmd,
		})

		Logger.Debug("Retrieved the client")
		// Cast the raw plugin to the Service interface so we have access to the methods
		plugin, err := ToService(client, pluginName)
		if err != nil {
			return err
		}
		Logger.Debug("Converted to Service")

		// Instead of storing the actual client in the map, we're going to store information that will allow the command to reattach to the services Server
		Logger.Debug("About to create serviceInfo")
		serviceInfos[pluginName] = serviceInfo(client, pluginName, plugin)
		Logger.Debug("Converted to serviceInfo")

		return nil
	})
	if err != nil {
		return nil, err
	}

	return serviceInfos, nil
}

// Get a raw interface from the client and converts to a Service
func ToService(client *goplugin.Client, pluginName string) (Service, error) {
	raw, err := plugin.Interface(client, pluginName)
	if err != nil {
		return nil, err
	}

	// Cast the raw plugin to the Service interface so we have access to the methods
	return raw.(Service), nil
}
