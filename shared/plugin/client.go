package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
)

type ClientConfigBuilder interface {
	Build() *goplugin.ClientConfig
	LoggingFields() []interface{}
}

type BaseClientConfigBuilder struct {
	HandshakeConfig  goplugin.HandshakeConfig
	Plugins          map[string]goplugin.Plugin
	Name             string
	Logger           hclog.Logger
	AllowedProtocols []goplugin.Protocol
}

func (b *BaseClientConfigBuilder) BuildBase() *goplugin.ClientConfig {
	syncStderr := io.Discard
	if logging.PluginDebug() {
		syncStderr = os.Stderr
	}

	return &goplugin.ClientConfig{
		HandshakeConfig:  b.HandshakeConfig,
		Plugins:          b.Plugins,
		Managed:          true,       // Allow the plugin runtime to manage this plugin
		SyncStdout:       os.Stdout,  // Print any extra output to Stdout from the plugin to the host processes Stdout
		SyncStderr:       syncStderr, // Print any extra output to Stderr from the plugin to the host processes Stderr
		AllowedProtocols: b.AllowedProtocols,
		Logger:           b.Logger,
		SkipHostEnv:      true, // Don't pass the host environment to the plugin to avoid security issues
		// I'd love to use this but I haven't yet figured out out to get it to work with Reattach
		// AutoMTLS:    true,      // Ensure that we're using MTLS for communication between the plugin and the host
	}
}

func (b *BaseClientConfigBuilder) LoggingFieldsBase() []interface{} {
	return []interface{}{
		"Plugin Name", b.Name,
		"AllowedProtocols", b.AllowedProtocols,
	}
}

type CommandClientConfigBuilder struct {
	ClientConfigBuilder
	BaseClientConfigBuilder
	PluginCmd string
}

func (b *CommandClientConfigBuilder) Build() *goplugin.ClientConfig {
	clientConfig := b.BuildBase()
	clientConfig.Cmd = exec.Command(b.PluginCmd)
	return clientConfig
}

func (b *CommandClientConfigBuilder) LoggingFields() []interface{} {
	return append(b.LoggingFieldsBase(), "Cmd", b.PluginCmd)
}

type ReattachClientConfigBuilder struct {
	ClientConfigBuilder
	BaseClientConfigBuilder
	ReattachConfig ReattachConfig
}

func (b *ReattachClientConfigBuilder) Build() *goplugin.ClientConfig {
	clientConfig := b.BuildBase()
	clientConfig.Reattach = b.ReattachConfig.ReattachConfig()
	return clientConfig
}

func (b *ReattachClientConfigBuilder) LoggingFields() []interface{} {
	return append(b.LoggingFieldsBase(),
		"Protocol", b.ReattachConfig.Protocol,
		"Addr Name", b.ReattachConfig.Addr.Name,
		"Addr Net", b.ReattachConfig.Addr.Net,
	)
}

func Client(pluginType string, clientConfigBuilder ClientConfigBuilder) *goplugin.Client {
	logger.Debug(fmt.Sprintf("Configuring %s", pluginType), clientConfigBuilder.LoggingFields()...)
	client := goplugin.NewClient(clientConfigBuilder.Build())
	logger.Debug("Plugin client created", append(clientConfigBuilder.LoggingFields(), "Protocol", client.Protocol())...)
	return client
}

// Dispenses the plugin from the ClientProtocal and returns the raw interface
func Interface(client *goplugin.Client, pluginName string) (interface{}, error) {
	logger.Trace("Getting the ClientProtocol from the client", "Plugin Name", pluginName)
	// Get the RPC Client from the plugin definition
	clientProtocol, err := client.Client()
	if err != nil {
		return nil, err
	}

	logger.Trace("Dispensing", "Plugin Name", pluginName)
	// Get the actual client so we can use it
	raw, err := clientProtocol.Dispense(pluginName)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
