package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bouncy/bouncy-api/internal/application/features"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence/repositories"
	"github.com/spf13/cobra"
)

var featureRootCommand = &cobra.Command{
	Use:   "feature",
	Short: "Feature flag management commands",
}

var featureListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all feature flags",
	Run: func(cmd *cobra.Command, args []string) {
		service := getFeatureService()
		flags, err := service.GetAll()
		if err != nil {
			slog.Error("Failed to list feature flags", "error", err)
			os.Exit(1)
		}

		fmt.Printf("%-20s %-20s %-10s %s\n", "KEY", "NAME", "ENABLED", "DESCRIPTION")
		fmt.Println("--------------------------------------------------------------------------------")
		for _, f := range flags {
			enabledStr := "off"
			if f.Enabled {
				enabledStr = "on"
			}
			fmt.Printf("%-20s %-20s %-10s %s\n", f.Key, f.Name, enabledStr, f.Description)
		}
	},
}

var featureAddCommand = &cobra.Command{
	Use:   "add [key] [name]",
	Short: "Add a new feature flag",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		name := args[1]
		description, _ := cmd.Flags().GetString("description")
		enabled, _ := cmd.Flags().GetBool("enabled")

		service := getFeatureService()
		err := service.Create(&models.FeatureFlag{
			Key:         key,
			Name:        name,
			Description: description,
			Enabled:     enabled,
		})

		if err != nil {
			slog.Error("Failed to add feature flag", "error", err)
			os.Exit(1)
		}

		slog.Info("Feature flag added successfully", "key", key)
	},
}

var featureRemoveCommand = &cobra.Command{
	Use:   "remove [key]",
	Short: "Remove a feature flag",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		service := getFeatureService()
		err := service.Delete(key)

		if err != nil {
			slog.Error("Failed to remove feature flag", "error", err)
			os.Exit(1)
		}

		slog.Info("Feature flag removed successfully", "key", key)
	},
}

var featureToggleCommand = &cobra.Command{
	Use:   "toggle [key] [on|off]",
	Short: "Toggle a feature flag on or off",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		state := args[1]

		var enabled bool
		if state == "on" {
			enabled = true
		} else if state == "off" {
			enabled = false
		} else {
			slog.Error("Invalid state. Use 'on' or 'off'")
			os.Exit(1)
		}

		service := getFeatureService()
		err := service.Update(key, enabled)

		if err != nil {
			slog.Error("Failed to toggle feature flag", "error", err)
			os.Exit(1)
		}

		slog.Info("Feature flag toggled successfully", "key", key, "enabled", enabled)
	},
}

func getFeatureService() *features.FeatureFlagService {
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

	repo := repositories.NewFeatureFlagRepository(dbServer.Database)
	return features.NewFeatureFlagService(repo)
}

func init() {
	RootCommand.AddCommand(featureRootCommand)
	featureRootCommand.AddCommand(featureListCommand)
	featureRootCommand.AddCommand(featureAddCommand)
	featureRootCommand.AddCommand(featureRemoveCommand)
	featureRootCommand.AddCommand(featureToggleCommand)

	featureAddCommand.Flags().StringP("description", "d", "", "Description of the feature flag")
	featureAddCommand.Flags().BoolP("enabled", "e", false, "Whether the feature flag is enabled by default")
}
