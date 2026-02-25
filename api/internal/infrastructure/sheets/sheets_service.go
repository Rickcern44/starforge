package sheets

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bouncy/bouncy-api/internal/application"
	"github.com/bouncy/bouncy-api/internal/domain/models"
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

func (service *Service) CreateMissingGames(dates []string, gameService *application.GameService) error {
	for _, date := range dates {

		parsedTime := parseSpreadsheetDate(date)

		game := &models.Game{
			LeagueID:    service.LeagueId,
			StartTime:   parsedTime,
			Location:    "Fairview Wellness Center",
			CostInCents: service.GameCostInCents,
			IsCanceled:  false,
			Attendance:  nil,
			Charges:     nil,
		}

		if _, err := gameService.Create(game); err != nil {
			slog.Error(err.Error())
			return err
		}
	}
	return nil
}

func parseSpreadsheetDate(date string) time.Time {
	layout := "02-Jan"
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
