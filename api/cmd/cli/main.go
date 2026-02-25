package main

import (
	"log/slog"

	"github.com/bouncy/bouncy-api/internal/application/commands"
)

func main() {
	slog.Info("Starting the CLI")

	commands.Execute()
}
