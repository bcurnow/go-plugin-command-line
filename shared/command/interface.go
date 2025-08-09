package command

import (
	"github.com/bcurnow/go-plugin-command-line/shared/plugin"
	"github.com/spf13/cobra"
)

// This is the interface for all command plugins
type Command interface {
	Help() string
	Execute(args []string) error
	SetServices(serviceInfo map[string]plugin.Reattach) error
}

type Register interface {
	Register(pluginDir string, cmd *cobra.Command, serviceInfo map[string]plugin.Reattach) error
}
