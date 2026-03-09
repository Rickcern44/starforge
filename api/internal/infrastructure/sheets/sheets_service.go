package sheets

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/application/users"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type Service struct {
	FilePath        string
	SheetName       string
	LeagueId        string
	GameCostInCents int
}

type AttendanceData struct {
	Name              string
	Dates             map[string]string
	TotalAttendance   int
	TotalDueInCents   int
	AmountPaidInCents int
	BalanceInCents    int
}

type SheetIndices struct {
	NameIndex            int
	TotalAttendanceIndex int
	TotalDueIndex        int
	AmountPaidIndex      int
	BalanceIndex         int
	DateIndices          map[int]string // Maps column index -> Date header string
}

func NewSheetsService(filePath, sheetName, leagueId string, gameCostInCents int) *Service {
	return &Service{
		FilePath:        filePath,
		SheetName:       sheetName,
		LeagueId:        leagueId,
		GameCostInCents: gameCostInCents,
	}
}

func (service *Service) LoadAttendanceData() ([]AttendanceData, error) {
	f, err := excelize.OpenFile(service.FilePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			slog.Error(err.Error())
		}
	}()

	rows, err := f.GetRows(service.SheetName)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet %s is empty", service.SheetName)
	}

	headerRow := rows[0]
	indices := getSheetIndexes(headerRow)

	var allAttendance []AttendanceData
	for _, row := range rows[1:] {
		if len(row) == 0 || len(row) <= indices.NameIndex || row[indices.NameIndex] == "" {
			continue
		}

		data := AttendanceData{
			Name:  row[indices.NameIndex],
			Dates: make(map[string]string),
		}

		// Parse Total Attendance (Int)
		if indices.TotalAttendanceIndex < len(row) {
			val, _ := strconv.Atoi(row[indices.TotalAttendanceIndex])
			data.TotalAttendance = val
		}

		// Parse Financials (Cents)
		if indices.TotalDueIndex < len(row) {
			data.TotalDueInCents = parseCurrency(row[indices.TotalDueIndex])
		}
		if indices.AmountPaidIndex < len(row) {
			data.AmountPaidInCents = parseCurrency(row[indices.AmountPaidIndex])
		}
		if indices.BalanceIndex < len(row) {
			data.BalanceInCents = parseCurrency(row[indices.BalanceIndex])
		}

		// Apply dynamic date columns to the map
		for colIdx, dateLabel := range indices.DateIndices {
			if colIdx < len(row) {
				data.Dates[dateLabel] = row[colIdx]
			}
		}

		allAttendance = append(allAttendance, data)
	}

	return allAttendance, nil
}

func (service *Service) GetTotalDates() ([]string, error) {
	f, err := excelize.OpenFile(service.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(service.SheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("sheet %s is empty", service.SheetName)
	}

	indices := getSheetIndexes(rows[0])

	dates := make([]string, 0, len(indices.DateIndices))
	for _, label := range indices.DateIndices {
		dates = append(dates, label)
	}
	return dates, nil
}

func (service *Service) CreateMissingGames(dates []string, gameService *application.GameService) ([]*models.Game, error) {
	existingGames, err := gameService.GetGamesForLeague(service.LeagueId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing games: %w", err)
	}

	var gamesToCreate []*models.Game
	for _, date := range dates {
		parsedTime := parseSpreadsheetDate(date)

		exists := false
		for _, existing := range existingGames {
			// Check if a game exists within the same hour (simple duplicate check)
			if existing.StartTime.Truncate(time.Hour).Equal(parsedTime.Truncate(time.Hour)) {
				exists = true
				break
			}
		}

		if !exists {
			game := models.CreateGameFromData(service.LeagueId, "Fairview Wellness Center", service.GameCostInCents, parsedTime)
			gamesToCreate = append(gamesToCreate, game)
		}
	}

	return gamesToCreate, nil
}

func (service *Service) ExecuteGameCreation(games []*models.Game, gameService *application.GameService) error {
	for _, game := range games {
		if _, err := gameService.Create(game); err != nil {
			return fmt.Errorf("failed to create game for %v: %w", game.StartTime, err)
		}
	}
	return nil
}

func (service *Service) IdentifyNewCharges(
	attendanceData []AttendanceData,
	games []*models.Game,
	gameService *application.GameService,
	paymentsService *payments.Service,
	userService *users.Service,
) ([]*models.GameCharge, error) {
	var chargesToCreate []*models.GameCharge

	// Create a map of date -> Game ID for quick lookup
	// The spreadsheet uses labels like "5-Mar"
	gameMap := make(map[string]string)
	for _, g := range games {
		key := g.StartTime.Format("2-Jan")
		gameMap[key] = g.ID
	}

	// Cache user lookups by name
	userCache := make(map[string]*string)

	for _, data := range attendanceData {
		// 1. Try to find User ID from name
		var userID *string
		if cachedID, ok := userCache[data.Name]; ok {
			userID = cachedID
		} else {
			user, err := userService.FindByName(data.Name)
			if err == nil && user != nil {
				id := user.ID
				userID = &id
			}
			userCache[data.Name] = userID
		}

		// 2. Get existing charges for this specific user/name to avoid duplicates
		var existingCharges []models.GameCharge
		if userID != nil {
			ec, err := paymentsService.ListChargesByUser(*userID)
			if err == nil {
				existingCharges = append(existingCharges, ec...)
			}
		}

		// ALWAYS check for unclaimed charges by name as well
		unclaimed, err := paymentsService.ListChargesByExternalName(data.Name)
		if err == nil {
			existingCharges = append(existingCharges, unclaimed...)
		}

		for dateLabel, value := range data.Dates {
			trimmedVal := strings.ToLower(strings.TrimSpace(value))
			if trimmedVal == "x" {
				gameID, ok := gameMap[dateLabel]
				if !ok {
					// This might happen if the game exists in the spreadsheet but not in our DB list
					// though Phase 1 should have caught it.
					continue
				}

				// Check if charge already exists
				duplicate := false
				for _, ec := range existingCharges {
					if ec.GameID == gameID {
						duplicate = true
						break
					}
				}

				if !duplicate {
					charge := models.CreateGameCharge(gameID, userID, data.Name, service.GameCostInCents)
					chargesToCreate = append(chargesToCreate, charge)
				}
			}
		}
	}

	return chargesToCreate, nil
}

func (service *Service) IdentifyNewPayments(
	attendanceData []AttendanceData,
	paymentsService *payments.Service,
	userService *users.Service,
) ([]*models.Payment, error) {
	var paymentsToCreate []*models.Payment

	// Cache user lookups by name
	userCache := make(map[string]*string)

	for _, data := range attendanceData {
		if data.AmountPaidInCents <= 0 {
			continue
		}

		// 1. Try to find User ID from name
		var userID *string
		if cachedID, ok := userCache[data.Name]; ok {
			userID = cachedID
		} else {
			user, err := userService.FindByName(data.Name)
			if err == nil && user != nil {
				id := user.ID
				userID = &id
			}
			userCache[data.Name] = userID
		}

		// 2. Calculate current total payments in DB
		currentPaidInCents := 0
		if userID != nil {
			p, err := paymentsService.ListPaymentsByUser(*userID)
			if err == nil {
				for _, pay := range p {
					currentPaidInCents += pay.AmountCents
				}
			}
		}

		unclaimed, err := paymentsService.ListPaymentsByExternalName(data.Name)
		if err == nil {
			for _, pay := range unclaimed {
				currentPaidInCents += pay.AmountCents
			}
		}

		// 3. Compare with spreadsheet
		diff := data.AmountPaidInCents - currentPaidInCents
		if diff > 0 {
			ref := fmt.Sprintf("Spreadsheet Sync: %s", data.Name)
			payment := models.CreatePayment(userID, data.Name, service.LeagueId, diff, models.PaymentMethodCash, &ref)
			// Ensure it has an ID
			id, _ := uuid.NewV7()
			payment.ID = id.String()

			paymentsToCreate = append(paymentsToCreate, payment)
		}
	}

	return paymentsToCreate, nil
}

func (service *Service) ExecuteChargeCreation(charges []*models.GameCharge, paymentsService *payments.Service) error {
	for _, c := range charges {
		if err := paymentsService.CreateCharge(c); err != nil {
			return err
		}
	}
	return nil
}

func (service *Service) ExecutePaymentCreation(payments []*models.Payment, paymentsService *payments.Service) error {
	for _, p := range payments {
		if err := paymentsService.Add(p); err != nil {
			return err
		}
	}
	return nil
}

func parseSpreadsheetDate(date string) time.Time {
	layout := "2-Jan"
	now := time.Now()
	parsed, err := time.Parse(layout, date)

	if err != nil {
		slog.Error(err.Error())
		return now
	}

	return time.Date(now.Year(), parsed.Month(), parsed.Day(), 19, 30, 0, 0, now.Location())
}

// parseCurrency converts strings like "$10.00" or "10.50" to cents (int)
func parseCurrency(val string) int {
	if val == "" {
		return 0
	}
	// Remove currency symbols, commas, etc.
	reg := regexp.MustCompile(`[^0-9.]`)
	clean := reg.ReplaceAllString(val, "")

	f, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return int(f * 100)
}

func getSheetIndexes(headerRow []string) *SheetIndices {
	indices := &SheetIndices{
		DateIndices: make(map[int]string),
	}

	for index, col := range headerRow {
		trimmed := strings.TrimSpace(col)
		switch {
		case trimmed == "Name":
			indices.NameIndex = index
		case trimmed == "Attendance":
			indices.TotalAttendanceIndex = index
		case strings.Contains(trimmed, "Total Due"):
			indices.TotalDueIndex = index
		case trimmed == "Amount Paid":
			indices.AmountPaidIndex = index
		case trimmed == "Balance":
			indices.BalanceIndex = index
		default:
			// If it's not a known column, check if it's a date header (e.g. 5-Mar)
			if isDateHeader(trimmed) {
				indices.DateIndices[index] = trimmed
			}
		}
	}
	return indices
}

func isDateHeader(col string) bool {
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	lower := strings.ToLower(col)
	for _, m := range months {
		if strings.Contains(lower, strings.ToLower(m)) {
			return true
		}
	}
	return false
}
