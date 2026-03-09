package database

import (
	"log/slog"
	"time"

	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	connectionString string
	Database         *gorm.DB
}

func NewDatabaseService(settings *config.Config) *Service {
	return &Service{
		connectionString: settings.Database.ConnectionString,
	}
}

func (dbs *Service) Connect() error {
	dsn := dbs.connectionString

	var db *gorm.DB
	var err error
	maxRetries := 10
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		slog.Warn("Failed to connect to database, retrying...",
			"attempt", i+1,
			"max_retries", maxRetries,
			"error", err.Error())
		time.Sleep(retryDelay)
	}

	if err != nil {
		slog.Error("Database connection failed after retries", "error", err.Error())
		return err
	}

	dbs.Database = db
	slog.Info("Connected to database")
	return nil
}

func (dbs *Service) UpdateDatabase() error {
	slog.Info("Updating database...")
	err := dbs.Database.AutoMigrate(
		&persistence.User{},
		&persistence.League{},
		&persistence.LeagueMember{},
		&persistence.Game{},
		&persistence.GameAttendance{},
		&persistence.Payment{},
		&persistence.PaymentAllocation{},
		&persistence.GameCharge{},
		&persistence.Invitation{},
		&persistence.FeatureFlag{},
	)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("Database Schema has been migrated successfully")
	return nil
}
