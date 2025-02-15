package logging

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

// Create a root logger for the application, using hclog for consistency with go-plugin
var logger = hclog.New(&hclog.LoggerOptions{
	Name:   "gpcl",
	Output: os.Stdout,
	Level:  hclog.Trace,
})

// Gets the main logger for the application
func Logger() hclog.Logger {
	return logger
}
