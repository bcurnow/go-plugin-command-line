package plugin

import (
	"os"

	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// Utility function to start a plugin, this should only be called from a command/services main method
// as it will not return
// If a grpcServer is passed in (not nil) then it will active gRPC mode for the plugin
func Start(pluginImpl goplugin.Plugin, name string, handshakeConfig goplugin.HandshakeConfig, grpc func([]grpc.ServerOption) *grpc.Server) {
	// This is the logger that will be used inside the plugin, it needs to be configured to use
	// Stderr because Stdout is used to communicate back to the host program, if this is not configured correctly
	// the plugin will fail to start with an "Unrecognized remote plugin message" error
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "plugin",
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	}).Named(name)

	logger.Debug("Starting plugin", "Plugin Name", name, "GRPC", grpc != nil)

	// Start the server
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]goplugin.Plugin{
			name: pluginImpl,
		},
		Logger:     logger,
		GRPCServer: grpc,
	})
}
