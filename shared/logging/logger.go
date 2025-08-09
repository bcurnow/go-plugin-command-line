package logging

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

// Create a root logger for the application, using hclog for consistency with go-plugin
var logger = hclog.New(&hclog.LoggerOptions{
	Name:   "gpcl",
	Output: os.Stderr,
	Level:  hclog.Trace,
})

var pluginDebug = false

// Gets the main logger for the application
func Logger() hclog.Logger {
	return logger

}

func EnablePluginDebug() {
	pluginDebug = true
}

func PluginDebug() bool {
	return pluginDebug
}
