package database

import (
	"fmt"
	"log/slog"

	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	host     string
	port     string
	user     string
	pass     string
	dbName   string
	Database *gorm.DB
}

func NewDatabaseService(settings *config.Config) *Service {
	return &Service{
		host:   settings.Database.Host,
		port:   settings.Database.Port,
		user:   settings.Database.Username,
		pass:   settings.Database.Password,
		dbName: settings.Database.Database,
	}
}

func (dbs *Service) Connect() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbs.host, dbs.user, dbs.pass, dbs.dbName, dbs.port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	dbs.Database = db
	slog.With("database", dbs.dbName).Info("Connected to database")
	return nil
}

func (dbs *Service) UpdateDatabase() error {
	fmt.Println("Updating database...")
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
	)

	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Info("Database Schema has been migrated successfully")
	return nil
}
