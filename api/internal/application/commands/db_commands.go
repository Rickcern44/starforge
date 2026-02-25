package commands

import "github.com/spf13/cobra"

var databaseRootCommand = &cobra.Command{
	Use:   "database",
	Short: "Manage database migrations",
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
}

func init() {
	RootCommand.AddCommand(databaseRootCommand)

	databaseRootCommand.AddCommand(migrateCommand)
}
