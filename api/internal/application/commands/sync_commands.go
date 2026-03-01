package commands

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/bouncy/bouncy-api/internal/infrastructure/sheets"
	"github.com/spf13/cobra"
)

var (
	filePath        string
	SheetName       string
	LeagueId        string
	GameCostInCents int
)

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

		gamesToCreate, err := service.CreateMissingGames(dates, app.GameService)
		if err != nil {
			slog.Error("Failed to identify missing games", "error", err)
			return
		}

		reader := bufio.NewReader(os.Stdin)

		if len(gamesToCreate) > 0 {
			fmt.Println("\n--- Sync Preview ---")
			fmt.Printf("The following %d games will be created:\n", len(gamesToCreate))
			for _, g := range gamesToCreate {
				fmt.Printf(" - %s at %s ($%.2f)\n", g.StartTime.Format("Jan 02, 2006"), g.Location, float64(g.CostInCents)/100.0)
			}
			fmt.Println("--------------------")

			fmt.Print("\nDo you want to proceed with creating these games? (y/N): ")
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "y" || response == "yes" {
				slog.Info("Proceeding with game creation...")
				if err := service.ExecuteGameCreation(gamesToCreate, app.GameService); err != nil {
					slog.Error("Failed to create games", "error", err)
					return
				}
				slog.Info("Successfully created games", "count", len(gamesToCreate))
			} else {
				slog.Info("Game creation skipped.")
			}
		} else {
			slog.Info("No new games to create.")
		}

		// --- PHASE 2: Charges ---
		// Refresh games from DB to get actual IDs
		dbGames, _ := app.GameService.GetGamesForLeague(LeagueId)

		chargesToCreate, err := service.IdentifyNewCharges(attendanceData, dbGames, app.GameService, app.PaymentsService, app.UserService)
		if err != nil {
			slog.Error("Failed to identify charges", "error", err)
			return
		}

		if len(chargesToCreate) == 0 {
			slog.Info("No new charges to apply.")
			return
		}

		fmt.Println("\n--- Charge Preview ---")
		fmt.Printf("The following %d charges will be applied based on 'x' marks:\n", len(chargesToCreate))
		for _, c := range chargesToCreate {
			displayName := c.ExternalName
			status := "UNCLAIMED"

			if c.UserID != nil {
				user, err := app.UserService.GetUserByID(*c.UserID)
				if err == nil {
					displayName = user.Name
					status = "CLAIMED"
				}
			}

			game, _ := app.GameService.GetGameById(c.GameID)
			fmt.Printf(" - [%s] %s: Game on %s ($%.2f)\n", status, displayName, game.StartTime.Format("Jan 02"), float64(c.AmountCents)/100.0)
		}
		fmt.Println("--------------------")

		fmt.Print("\nDo you want to apply these charges? (y/N): ")
		chargeResponse, _ := reader.ReadString('\n')
		chargeResponse = strings.ToLower(strings.TrimSpace(chargeResponse))

		if chargeResponse == "y" || chargeResponse == "yes" {
			slog.Info("Applying charges...")
			if err := service.ExecuteChargeCreation(chargesToCreate, app.PaymentsService); err != nil {
				slog.Error("Failed to apply charges", "error", err)
				return
			}
			slog.Info("Successfully applied charges", "count", len(chargesToCreate))
		} else {
			slog.Info("Charge application aborted.")
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
