package database

import (
	"fmt"
	"os"

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

func NewDatabaseService() *Service {
	return &Service{
		host:   os.Getenv("POSTGRES_HOST"),
		port:   os.Getenv("POSTGRES_PORT"),
		user:   os.Getenv("POSTGRES_USER"),
		pass:   os.Getenv("POSTGRES_PASS"),
		dbName: os.Getenv("POSTGRES_DB"),
	}
}

func (dbs *Service) Connect() error {
	fmt.Println("Connecting to database...")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbs.host, dbs.user, dbs.pass, dbs.dbName, dbs.port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	dbs.Database = db
	fmt.Printf("Connected to %v:%v\n", dbs.host, dbs.dbName)
	return nil
}

func (dbs *Service) UpdateDatabase() error {
	fmt.Println("Updating database...")
	err := dbs.Database.AutoMigrate(
		&persistence.League{},
		&persistence.LeagueMember{},
		&persistence.Game{},
		&persistence.GameAttendance{},
		&persistence.Payment{},
		&persistence.PaymentAllocation{},
		&persistence.GameCharge{},
	)

	if err != nil {
		return err
	}
	fmt.Println("database schema has been updated")
	return nil
}
