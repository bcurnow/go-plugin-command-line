package command

import (
	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/bcurnow/go-plugin-command-line/shared/util"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

// Create a specific logger to be used by the plugin system as a whole
var (
	Logger = logging.Logger().Named("command")

	// This is the go-plugin handshake information that needs to be used for all plugins
	HandshakeConfig = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "COMMAND_PLUGIN",
		MagicCookieValue: "ec8c34d0-0185-4c22-83d4-67ad032a3dae",
	}
)

// This is the interface for all command plugins
type Command interface {
	Help() string
	Execute(args []string) error
	SetServices(serviceInfo map[string]service.ReconnectInfo) error
}

type CommandRegister interface {
	RegisterCommands(pluginDir string, cmd *cobra.Command, serviceInfo map[string]service.ReconnectInfo) error
}

func ToCommand(pluginDef *plugin.Client, pluginName string) (Command, error) {
	raw, err := util.ToInterface(pluginDef, pluginName)
	if err != nil {
		return nil, err
	}

	// Cast the raw plugin to the Command interface so we have access to the methods
	return raw.(Command), nil
}
