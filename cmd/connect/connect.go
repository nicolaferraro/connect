package main

import (
	"github.com/nicolaferraro/connect/pkg/commands"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	var cmd = cobra.Command{
		Use:          "connect",
		Short:        "Connect allows to connect cloud applications to Kubernetes",
		SilenceUsage: true,
	}

	cmd.AddCommand(commands.NewCmdRegister())
	cmd.AddCommand(commands.NewCmdRefresh())

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
