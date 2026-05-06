package util

import (
	grpccommand "github.com/bcurnow/go-plugin-command-line/grpc/command"
	grpcservice "github.com/bcurnow/go-plugin-command-line/grpc/service"
	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

type GRPCCommandRegister struct {
}

func (r *GRPCCommandRegister) Register(pluginDir string, cobraCmd *cobra.Command, serviceInfo map[string]plugin.Reattach) error {
	return command.RegisterDir(pluginDir, cobraCmd, &grpccommand.Plugin{}, []goplugin.Protocol{goplugin.ProtocolGRPC}, serviceInfo)
}

type GRPCServiceRegister struct {
}

func (r *GRPCServiceRegister) Register(serviceDir string) (map[string]plugin.Reattach, error) {
	return service.RegisterDir(serviceDir, &grpcservice.Plugin{}, []goplugin.Protocol{goplugin.ProtocolGRPC}, buildReattach)
}

func buildReattach(pluginClient *goplugin.Client, pluginName string, svc service.Service) plugin.Reattach {
	service.Logger.Debug("Building reattach", "Plugin Name", pluginName, "Protocol", pluginClient.ReattachConfig().Protocol)
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

}
