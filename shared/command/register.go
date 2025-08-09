package command

import (
	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

// Create a specific logger to be used by the plugin system as a whole
var (
	Logger = logging.Logger().Named("command")

	// This is the go-plugin handshake information that needs to be used for all plugins
	HandshakeConfig = goplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "COMMAND_PLUGIN",
		MagicCookieValue: "ec8c34d0-0185-4c22-83d4-67ad032a3dae",
	}
)

// Traverses the dir and registers any executable found as a subcommand
func RegisterDir(dir string, cmd *cobra.Command, impl goplugin.Plugin, allowedProtocols []goplugin.Protocol, serviceInfo map[string]plugin.Reattach) error {
	err := plugin.Register("commands", dir, func(pluginName string, pluginCmd string) error {
		// Build a new plugin definition
		pluginDef := plugin.Client("command", &plugin.CommandClientConfigBuilder{
			BaseClientConfigBuilder: plugin.BaseClientConfigBuilder{
				HandshakeConfig: HandshakeConfig,
				Plugins: map[string]goplugin.Plugin{
					pluginName: impl,
				},
				Name:             pluginName,
				Logger:           Logger,
				AllowedProtocols: allowedProtocols,
			},
			PluginCmd: pluginCmd,
		})

		// Cast the raw plugin to the Command interface so we have access to the methods
		plugin, err := ToCommand(pluginDef, pluginName)
		if err != nil {
			return err
		}

		// Set the services on the command
		Logger.Debug("Setting services", "Plugin Name", pluginName, "ServiceInfo", serviceInfo)
		err = plugin.SetServices(serviceInfo)
		if err != nil {
			return err
		}

		err = addCommand(cmd, pluginName, plugin)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Adds the supplied pluginName as a subcommand to the supplied cmd
func addCommand(cmd *cobra.Command, pluginName string, plugin Command) error {

	// Create a new command which executes this plugin
	cmd.AddCommand(&cobra.Command{
		Use: pluginName,
		Run: func(cmd *cobra.Command, args []string) {
			plugin.Execute(args)
		},
	})

	return nil
}

func ToCommand(pluginDef *goplugin.Client, pluginName string) (Command, error) {
	raw, err := plugin.Interface(pluginDef, pluginName)
	if err != nil {
		return nil, err
	}

	// Cast the raw plugin to the Command interface so we have access to the methods
	return raw.(Command), nil
}
