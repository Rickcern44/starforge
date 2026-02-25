package commands

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "bouncy",
	Short: "Commands for the bouncy system",
}

func Execute() {
	err := RootCommand.Execute()

	if err != nil {
		slog.Error(err.Error())
	}
}
