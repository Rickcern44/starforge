package sheets

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/interfaces"
	"github.com/bouncy/bouncy-api/internal/application/payments"
	"github.com/bouncy/bouncy-api/internal/application/users"
	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/xuri/excelize/v2"
)

// --- Mocks ---

type mockUserRepo struct {
	interfaces.UserRepository
	users map[string]*models.User
}

func (m *mockUserRepo) FindByName(name string) (*models.User, error) {
	if u, ok := m.users[name]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user not found")
}

type mockPaymentsRepo struct {
	interfaces.PaymentsRepository
	payments  map[string][]models.Payment // map[userID]
	unclaimed map[string][]models.Payment // map[externalName]
}

func (m *mockPaymentsRepo) ListPaymentsByUser(userID string) ([]models.Payment, error) {
	return m.payments[userID], nil
}

func (m *mockPaymentsRepo) ListPaymentsByExternalName(name string) ([]models.Payment, error) {
	return m.unclaimed[name], nil
}

// --- Tests ---

func TestParseCurrency(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"$10.00", 1000},
		{"10.50", 1050},
		{"$1,234.56", 123456},
		{"0", 0},
		{"", 0},
		{"invalid", 0},
		{"$5", 500},
	}

	for _, tt := range tests {
		got := parseCurrency(tt.input)
		if got != tt.expected {
			t.Errorf("parseCurrency(%q) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

func TestIsDateHeader(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5-Mar", true},
		{"12-Jan", true},
		{"Name", false},
		{"Attendance", false},
		{"25-Dec", true},
		{"random", false},
	}

	for _, tt := range tests {
		got := isDateHeader(tt.input)
		if got != tt.expected {
			t.Errorf("isDateHeader(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}

func TestParseSpreadsheetDate(t *testing.T) {
	dateStr := "5-Mar"
	parsed := parseSpreadsheetDate(dateStr)

	if parsed.Month() != time.March {
		t.Errorf("Expected March, got %v", parsed.Month())
	}
	if parsed.Day() != 5 {
		t.Errorf("Expected day 5, got %d", parsed.Day())
	}
	if parsed.Hour() != 19 || parsed.Minute() != 30 {
		t.Errorf("Expected 19:30, got %02d:%02d", parsed.Hour(), parsed.Minute())
	}
}

func TestGetSheetIndexes(t *testing.T) {
	header := []string{"Name", "5-Mar", "12-Mar", "Attendance", "Total Due", "Balance"}
	indices := getSheetIndexes(header)

	if indices.NameIndex != 0 {
		t.Errorf("Expected NameIndex 0, got %d", indices.NameIndex)
	}
	if indices.TotalAttendanceIndex != 3 {
		t.Errorf("Expected TotalAttendanceIndex 3, got %d", indices.TotalAttendanceIndex)
	}
	if len(indices.DateIndices) != 2 {
		t.Errorf("Expected 2 date columns, got %d", len(indices.DateIndices))
	}
	if indices.DateIndices[1] != "5-Mar" || indices.DateIndices[2] != "12-Mar" {
		t.Errorf("Date mapping incorrect: %v", indices.DateIndices)
	}
}

func TestLoadAttendanceData(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.xlsx")
	sheetName := "Sheet1"

	f := excelize.NewFile()
	defer f.Close()

	f.NewSheet(sheetName)
	f.SetCellValue(sheetName, "A1", "Name")
	f.SetCellValue(sheetName, "B1", "5-Mar")
	f.SetCellValue(sheetName, "C1", "Attendance")
	f.SetCellValue(sheetName, "D1", "Total Due")
	f.SetCellValue(sheetName, "E1", "Amount Paid")

	f.SetCellValue(sheetName, "A2", "John Doe")
	f.SetCellValue(sheetName, "B2", "x")
	f.SetCellValue(sheetName, "C2", "1")
	f.SetCellValue(sheetName, "D2", "$10.00")
	f.SetCellValue(sheetName, "E2", "$5.00")

	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to save test excel file: %v", err)
	}

	service := NewSheetsService(filePath, sheetName, "league-1", 1000)
	data, err := service.LoadAttendanceData()
	if err != nil {
		t.Fatalf("LoadAttendanceData failed: %v", err)
	}

	if len(data) != 1 {
		t.Fatalf("Expected 1 row of data, got %d", len(data))
	}
	if data[0].Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', got '%s'", data[0].Name)
	}
	if data[0].AmountPaidInCents != 500 {
		t.Errorf("Expected AmountPaid 500, got %d", data[0].AmountPaidInCents)
	}
}

func TestIdentifyNewPayments(t *testing.T) {
	// Setup Services with Mocks
	uRepo := &mockUserRepo{users: make(map[string]*models.User)}
	pRepo := &mockPaymentsRepo{
		payments:  make(map[string][]models.Payment),
		unclaimed: make(map[string][]models.Payment),
	}

	uRepo.users["John Doe"] = &models.User{ID: "u1", Name: "John Doe"}
	pRepo.payments["u1"] = []models.Payment{
		{AmountInCents: 1000}, // Already paid $10
	}

	userService := users.NewUserService(uRepo)
	paymentsService := payments.NewPaymentsService(pRepo)

	service := NewSheetsService("fake.xlsx", "Sheet1", "league-1", 1000)

	tests := []struct {
		name           string
		attendance     []AttendanceData
		expectedNew    int
		expectedAmount int
	}{
		{
			name: "Partial Payment - Create Diff",
			attendance: []AttendanceData{
				{Name: "John Doe", AmountPaidInCents: 1500}, // Sheet says $15
			},
			expectedNew:    1,
			expectedAmount: 500, // Should create $5 payment
		},
		{
			name: "Full Payment - No Diff",
			attendance: []AttendanceData{
				{Name: "John Doe", AmountPaidInCents: 1000}, // Sheet says $10
			},
			expectedNew: 0,
		},
		{
			name: "New Unclaimed User",
			attendance: []AttendanceData{
				{Name: "Jane Doe", AmountPaidInCents: 2000}, // Sheet says $20, user not in DB
			},
			expectedNew:    1,
			expectedAmount: 2000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.IdentifyNewPayments(tt.attendance, paymentsService, userService)
			if err != nil {
				t.Fatalf("IdentifyNewPayments failed: %v", err)
			}

			if len(got) != tt.expectedNew {
				t.Errorf("Expected %d new payments, got %d", tt.expectedNew, len(got))
			}

			if tt.expectedNew > 0 && got[0].AmountInCents != tt.expectedAmount {
				t.Errorf("Expected amount %d, got %d", tt.expectedAmount, got[0].AmountInCents)
			}
		})
	}
}
