package util

import (
	rpccommand "github.com/bcurnow/go-plugin-command-line/rpc/command"
	rpcservice "github.com/bcurnow/go-plugin-command-line/rpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type RPCCommandRegister struct {
}

func (r *RPCCommandRegister) Register(pluginDir string, cobraCmd *cobra.Command, serviceInfo map[string]plugin.Reattach) error {
	return command.RegisterDir(pluginDir, cobraCmd, &rpccommand.Plugin{}, []goplugin.Protocol{goplugin.ProtocolNetRPC}, serviceInfo)
}

type RPCServiceRegister struct {
}

func (r *RPCServiceRegister) Register(serviceDir string) (map[string]plugin.Reattach, error) {
	return service.RegisterDir(serviceDir, &rpcservice.Plugin{}, []goplugin.Protocol{goplugin.ProtocolNetRPC}, func(pluginClient *goplugin.Client, pluginName string, svc service.Service) plugin.Reattach {
		return plugin.Reattach{
			ReattachConfig: plugin.ReattachConfig{
				Protocol:        pluginClient.ReattachConfig().Protocol,
				ProtocolVersion: pluginClient.ReattachConfig().ProtocolVersion,
				Addr: plugin.Addr{
					Net:  pluginClient.ReattachConfig().Addr.Network(),
					Name: pluginClient.ReattachConfig().Addr.String(),
				},
				Pid:  pluginClient.ReattachConfig().Pid,
				Test: pluginClient.ReattachConfig().Test,
			},
			PluginName: pluginName,
		}
	})
}
