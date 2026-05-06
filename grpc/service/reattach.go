package service

import (
	"github.com/bcurnow/go-plugin-command-line/grpc/proto"
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
				Name:   name,
				Logger: service.Logger.Named(name),
				AllowedProtocols: []goplugin.Protocol{
					goplugin.ProtocolGRPC,
				},
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

func ToProtoReattaches(reattaches map[string]plugin.Reattach) map[string]*proto.Reattach {
	protoReattaches := make(map[string]*proto.Reattach)
	for name, reattach := range reattaches {
		protoReattaches[name] = ToProtoReattach(reattach)
	}
	return protoReattaches
}

func ToProtoReattach(reattach plugin.Reattach) *proto.Reattach {
	protoReattach := &proto.Reattach{
		ReattachConfig: &proto.ReattachConfig{
			Protocol:        string(reattach.ReattachConfig.Protocol),
			ProtocolVersion: int64(reattach.ReattachConfig.ProtocolVersion),
			Addr: &proto.Addr{
				Net:  reattach.ReattachConfig.Addr.Net,
				Name: reattach.ReattachConfig.Addr.Name,
			},
			Pid:  int64(reattach.ReattachConfig.Pid),
			Test: reattach.ReattachConfig.Test,
		},
		PluginName: reattach.PluginName,
	}

	return protoReattach
}

func ToReattaches(protoReattaches map[string]*proto.Reattach) map[string]plugin.Reattach {
	reattaches := make(map[string]plugin.Reattach)
	for name, protoReattach := range protoReattaches {
		reattaches[name] = ToReattach(protoReattach)
	}
	return reattaches
}

func ToReattach(protoReattach *proto.Reattach) plugin.Reattach {
	reattach := plugin.Reattach{
		ReattachConfig: plugin.ReattachConfig{
			Protocol:        goplugin.ProtocolGRPC,
			ProtocolVersion: int(protoReattach.ReattachConfig.ProtocolVersion),
			Addr: plugin.Addr{
				Net:  protoReattach.ReattachConfig.Addr.Net,
				Name: protoReattach.ReattachConfig.Addr.Name,
			},
			Pid:  int(protoReattach.ReattachConfig.Pid),
			Test: protoReattach.ReattachConfig.Test,
		},
		PluginName: protoReattach.PluginName,
	}

	return reattach
}
