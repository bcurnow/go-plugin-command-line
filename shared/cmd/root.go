package cmd

import (
	"fmt"
	"os"

	"github.com/bcurnow/go-plugin-command-line/shared/command"
	"github.com/bcurnow/go-plugin-command-line/shared/logging"
	"github.com/bcurnow/go-plugin-command-line/shared/service"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
)

var (
	// Define the root command
	// This will use PreRunE to register the plugins and RunE to execute the command
	rootCmd = &cobra.Command{
		Use:   "gpcl",
		Short: "gpcl a pluggable command line program",
		Long:  "gpcl is an example program showing how to build a pluggable command line program using github.com/hashicorp/go-plugin",
		RunE:  run,
	}

	logger = logging.Logger()

	// Flags for the root command
	pluginDir  string
	serviceDir string
	logLevel   string
)

func init() {
	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&pluginDir, "plugin-dir", "d", "./commands", "The directory where the plugins are located")
	rootCmd.PersistentFlags().StringVarP(&serviceDir, "service-dir", "s", "./services", "The directory where the services are located")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "warn", "The log level (trace, debug, info, warn, error, fatal), not case sensitive")
}

// This is the main entry point for Cobra
func Execute(serviceRegister service.Register, commandRegister command.Register) {
	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return preRun(cmd, serviceRegister, commandRegister)

	}
	if err := rootCmd.Execute(); err != nil {
		// Don't use the logger, we want this to be standard formatting
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, serviceRegister service.Register, commandRegister command.Register) error {
	// Very first thing to do, set the log level
	logger.SetLevel(hclog.LevelFromString(logLevel))
	if logger.GetLevel() == hclog.NoLevel {
		// Default to Warn
		logger.SetLevel(hclog.Warn)
		logger.Error("Invalid log level specified, defaulting to Warn", "LogLevel", logLevel)
	}

	serviceInfo, err := serviceRegister.Register(serviceDir)
	if err != nil {
		return err
	}

	err = commandRegister.Register(pluginDir, cmd, serviceInfo)
	if err != nil {
		return err
	}
	return nil
}

func run(cmd *cobra.Command, args []string) error {
	defer cleanup()

	// No args meaning no command was specified, show the help and exit
	if len(args) == 0 {
		cmd.Help()
		os.Exit(0)
	}

	// See if the args match any of the subcommands
	foundCmd, remainingArgs, err := cmd.Traverse(args)
	if err != nil {
		return err
	}

	// If we were unable to find a matching command, show the help and exit
	// It is possible for the Traverse method to return the current command, treat
	// this as a command not found
	if foundCmd == nil || foundCmd == cmd {
		cmd.Help()
		os.Exit(0)
	}

	//If we found the command, use the remaining args as arguments to that command and execute it
	foundCmd.SetArgs(remainingArgs)
	logger.Debug("Executing command", "Name", foundCmd.Name(), "Args", remainingArgs)
	err = foundCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func cleanup() {
	logger.Debug("Cleaning up the clients...")
	goplugin.CleanupClients()
}
