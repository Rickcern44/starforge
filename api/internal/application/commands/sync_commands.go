package commands

import (
	"log/slog"
	"os"

	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/container"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
	"github.com/bouncy/bouncy-api/internal/infrastructure/sheets"
	"github.com/spf13/cobra"
)

var (
	filePath        string
	SheetName       string
	LeagueId        string
	GameCostInCents int
)

func getContainer() (*container.AppContainer, error) {
	settings, err := config.LoadConfig(os.Getenv("APP_CONFIG_PATH"))
	if err != nil {
		return nil, err
	}

	dbServer := database.NewDatabaseService(settings)
	if err := dbServer.Connect(); err != nil {
		return nil, err
	}

	return container.NewAppContainer(dbServer.Database, settings), nil
}

var syncRootCommand = &cobra.Command{
	Use:   "sync",
	Short: "Commands for synchronizing sheet data to the application database",
}

var syncAllCommand = &cobra.Command{
	Use:   "all",
	Short: "Commands for synchronizing all sheet data to the application database",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := getContainer()
		if err != nil {
			slog.Error("Failed to initialize application", "error", err)
			return
		}

		service := sheets.NewSheetsService(filePath, SheetName, LeagueId, GameCostInCents)
		attendanceData, err := service.LoadAttendanceData()

		if err != nil {
			slog.Error(err.Error())
			return
		}

		slog.Info("AttendanceData loaded", "count", len(attendanceData))

		dates, err := service.GetTotalDates()
		if err != nil {
			slog.Error("Failed to get dates", "error", err)
			return
		}

		slog.Info("Detected game dates from sheet", "dates", dates)
		slog.Info("Creating Games", "Number of Games", len(dates))

		if err := service.CreateMissingGames(dates, app.GameService); err != nil {
			slog.Error(err.Error())
		}
	},
}

func init() {
	RootCommand.AddCommand(syncRootCommand)

	syncRootCommand.AddCommand(syncAllCommand)

	syncAllCommand.Flags().StringVarP(&filePath, "path", "f", "", "path to file to sync")
	syncAllCommand.MarkFlagRequired("path")

	syncAllCommand.Flags().StringVarP(&LeagueId, "league-id", "l", "", "league id to sync")
	syncAllCommand.MarkFlagRequired("league-id")

	syncAllCommand.Flags().IntVarP(&GameCostInCents, "cost", "c", 0, "game cost in cents")
	syncAllCommand.MarkFlagRequired("cost")

	syncRootCommand.Flags().StringVarP(&SheetName, "sheet", "s", "Sheet1", "name of sheet to sync")

}
