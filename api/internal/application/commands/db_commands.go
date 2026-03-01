package commands

import (
	"log/slog"
	"os"

	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
	"github.com/spf13/cobra"
)

var databaseRootCommand = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations (AutoMigrate)",
	Run: func(cmd *cobra.Command, args []string) {
		settings, err := config.LoadConfig()
		if err != nil {
			slog.Error("Failed to load config", "error", err)
			os.Exit(1)
		}

		dbServer := database.NewDatabaseService(settings)
		if err := dbServer.Connect(); err != nil {
			slog.Error("Failed to connect to database", "error", err)
			os.Exit(1)
		}

		if err := dbServer.UpdateDatabase(); err != nil {
			slog.Error("Failed to migrate database", "error", err)
			os.Exit(1)
		}

		slog.Info("Database migration completed successfully")
	},
}

func init() {
	RootCommand.AddCommand(databaseRootCommand)
	databaseRootCommand.AddCommand(migrateCommand)
}
