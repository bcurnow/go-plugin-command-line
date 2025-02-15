package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

var logger = logging.Logger().Named("util")

// Traverses the dir and calls the register function for any executable
func RegisterPlugins(pluginType string, dir string, register func(pluginName string, pluginCmd string) error) error {
	logger.Debug(fmt.Sprintf("Loading %s", pluginType), "Dir", dir)
	executables, err := discoverPlugins(dir)
	if err != nil {
		return err
	}

	for pluginName, pluginCmd := range executables {
		err = register(pluginName, pluginCmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetPluginClient(pluginType string, pluginName string, pluginCmd string, handshakeConfig plugin.HandshakeConfig, impl plugin.Plugin, logger hclog.Logger) *plugin.Client {
	logger.Debug(fmt.Sprintf("Registering %s", pluginType), "Name", pluginName, "Cmd", pluginCmd)
	pluginDef := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			pluginName: impl,
		},
		Cmd:        exec.Command(pluginCmd),
		Logger:     logger,
		Managed:    true,      // Allow the plugin runtime to manage this plugin
		SyncStdout: os.Stdout, // Print any extra output to Stdout from the plugin to the host processes Stdout
		// AutoMTLS:    true,      // Ensure that we're using MTLS for communication between the plugin and the host
		SkipHostEnv: true, // Don't pass the host environment to the plugin to avoid security issues
	})

	return pluginDef
}

func ToInterface(pluginDef *plugin.Client, pluginName string) (interface{}, error) {
	// Get the RPC Client from the plugin definition
	clientProtocol, err := pluginDef.Client()
	if err != nil {
		return nil, err
	}

	// Get the actual client so we can use it
	raw, err := clientProtocol.Dispense(pluginName)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

// Utility function to start a plugin, this should only be called from a command/services main method
// as it will not return
func StartPlugin(pluginImpl plugin.Plugin, name string, handshakeConfig plugin.HandshakeConfig) {
	// This is the logger that will be used inside the plugin, it needs to be configured to use
	// Stderr because Stdout is used to communicate back to the host program, if this is not configured correctly
	// the plugin will fail to start with an "Unrecognized remote plugin message" error
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       name,
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	logger.Debug("Starting plugin", "Name", name)

	// Start the server
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			name: pluginImpl,
		}})
}

// Until plugin.Discover is updated to check for the executable bit, this is our own implementation
func discoverPlugins(dir string) (map[string]string, error) {
	var executables map[string]string = make(map[string]string)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Don't traverse sub-directories, this is arbitrary but we are keeping it simple
		if d.IsDir() && path != dir {
			logger.Warn("Subdirectories are not supported", "Subdirectory", path)
			return filepath.SkipDir
		}

		// Because we're using WalkDir, we need to get the FileInfo from the DirEntry
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Check if this is a file and if the file is executable
		if info.Mode().IsRegular() {
			// 0111 checks for the execute bit to be set
			if info.Mode()&0111 == 0 {
				logger.Warn("Skipping non-executable file", "File", path)
				return nil
			}

			// Get the absolute path of the file so we can provide the best debugging information
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			executables[filepath.Base(path)] = absPath
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return executables, nil
}
