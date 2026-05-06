package main

import (
	"github.com/bcurnow/go-plugin-command-line/grpc/util"
	"github.com/bcurnow/go-plugin-command-line/shared/cmd"
)

// Simply call the Execute function in cmd (the Cobra package)
func main() {
	cmd.Execute(&util.GRPCServiceRegister{}, &util.GRPCCommandRegister{})
}
