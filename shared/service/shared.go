package service

import (
	"encoding/gob"

	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/bcurnow/go-plugin-command-line/shared/util"
	"github.com/hashicorp/go-plugin"
)

func init() {
	gob.Register(map[string]ReconnectInfo{})
}

var (
	Logger = logging.Logger().Named("service")

	// This is the go-plugin handshake information that needs to be used for all plugins
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "SERVICE_PLUGIN",
		MagicCookieValue: "a44a5c99-fac2-4ff4-9a47-3eec6283725a",
	}
)

// This is the interface for all service plugins
type Service interface {
	Name() string
	Type() string
	Execute(arg string) error
}

// Marker interface for an object that provides information necessary to reconnect to a service
type ReconnectInfo interface {
	Name() string
	Type() string
}

type ServiceRegister interface {
	RegisterServices(serviceDir string) (map[string]ReconnectInfo, error)
}

// The plugin.Plugin implementation for a Service which returns the RCP Client or Server
type ServicePlugin struct{ Impl Service }

// Traverses the dir and registers any executable found as a Service
func RegisterServices(dir string, service plugin.Plugin, createServiceInfo func(pluginClient *plugin.Client, pluginName string, service Service) ReconnectInfo) (map[string]ReconnectInfo, error) {
	serviceInfos := make(map[string]ReconnectInfo)
	err := util.RegisterPlugins("services", dir, func(pluginName string, pluginCmd string) error {
		// Build a new plugin definition
		pluginClient := util.GetPluginClient("service", pluginName, pluginCmd, HandshakeConfig, service, Logger)

		// Cast the raw plugin to the Service interface so we have access to the methods
		plugin, err := ToService(pluginClient, pluginName)
		if err != nil {
			return err
		}

		// Instead of storing the actual client in the map, we're going to store information that will allow the command to reattach to the services Server
		serviceInfo := createServiceInfo(pluginClient, pluginName, plugin)

		// Add the service to the list of registered services
		serviceInfos[pluginName] = serviceInfo

		return nil
	})
	if err != nil {
		return nil, err
	}

	return serviceInfos, nil
}

func ToService(pluginDef *plugin.Client, pluginName string) (Service, error) {
	raw, err := util.ToInterface(pluginDef, pluginName)
	if err != nil {
		return nil, err
	}

	// Cast the raw plugin to the Service interface so we have access to the methods
	return raw.(Service), nil
}
