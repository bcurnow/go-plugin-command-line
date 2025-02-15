package util

import (
	"encoding/gob"

	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

func init() {
	gob.Register(&rpcservice.RPCClientInfo{})
}

type RPCCommandRegister struct {
}

func (r *RPCCommandRegister) Register(pluginDir string, cobraCmd *cobra.Command, serviceInfo map[string]service.ReconnectInfo) error {
	return command.RegisterDir(pluginDir, cobraCmd, &rpccommand.Plugin{}, serviceInfo)
}

type RPCServiceRegister struct {
}

func (r *RPCServiceRegister) Register(serviceDir string) (map[string]service.ReconnectInfo, error) {
	return service.RegisterDir(serviceDir, &rpcservice.Plugin{}, func(pluginClient *goplugin.Client, pluginName string, service service.Service) service.ReconnectInfo {
		return &rpcservice.RPCClientInfo{
			ReattachConfig: rpcservice.ReattachConfig{
				Protocol:        pluginClient.ReattachConfig().Protocol,
				ProtocolVersion: pluginClient.ReattachConfig().ProtocolVersion,
				Addr: rpcservice.Addr{
					Net:  pluginClient.ReattachConfig().Addr.Network(),
					Name: pluginClient.ReattachConfig().Addr.String(),
				},
				Pid:          pluginClient.ReattachConfig().Pid,
				ReattachFunc: pluginClient.ReattachConfig().ReattachFunc,
				Test:         pluginClient.ReattachConfig().Test,
			},
			PluginName: pluginName,
			PluginType: service.Type(),
		}
	})
}
