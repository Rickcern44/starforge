package database

import (
	"log/slog"
	"time"

	"github.com/bouncy/bouncy-api/internal/domain/models"
	"github.com/bouncy/bouncy-api/internal/infrastructure/persistence"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (dbs *Service) Seed() error {
	var count int64
	dbs.Database.Model(&persistence.User{}).Count(&count)

	if count > 0 {
		slog.Info("Database already has data, skipping seed.")
		return nil
	}

	slog.Info("Seeding database with initial data...")

	return dbs.Database.Transaction(func(tx *gorm.DB) error {
		// 1. Create a Default Admin User
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123"), 14)
		adminID := uuid.NewString()
		adminUser := &persistence.User{
			Base:         persistence.Base{ID: adminID},
			Email:        "admin@bouncy.com",
			Name:         "League Admin",
			PasswordHash: string(hashedPassword),
			Roles:        []string{"admin", "user"},
		}
		if err := tx.Create(adminUser).Error; err != nil {
			return err
		}

		// 1b. Create a Normal User
		userID := uuid.NewString()
		normalUser := &persistence.User{
			Base:         persistence.Base{ID: userID},
			Email:        "user@bouncy.com",
			Name:         "Regular Player",
			PasswordHash: string(hashedPassword),
			Roles:        []string{"user"},
		}
		if err := tx.Create(normalUser).Error; err != nil {
			return err
		}

		// 2. Create a Sample League
		leagueID := uuid.NewString()
		league := &persistence.League{
			Base:     persistence.Base{ID: leagueID},
			Name:     "Tuesday Night Hoops",
			IsActive: true,
		}
		if err := tx.Create(league).Error; err != nil {
			return err
		}

		// 3. Enroll users in league
		members := []persistence.LeagueMember{
			{LeagueID: leagueID, UserID: adminID, Role: models.Role("owner"), JoinedAt: time.Now()},
			{LeagueID: leagueID, UserID: userID, Role: models.Role("user"), JoinedAt: time.Now()},
		}
		for _, m := range members {
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}

		// 4. Create an Upcoming Game
		gameID := uuid.NewString()
		game := &persistence.Game{
			Base:        persistence.Base{ID: gameID},
			LeagueID:    leagueID,
			StartTime:   time.Now().Add(48 * time.Hour),
			Location:    "Community Center Court 1",
			CostInCents: 1000,
			IsCanceled:  false,
		}
		if err := tx.Create(game).Error; err != nil {
			return err
		}

		slog.Info("Database seeded successfully", "admin_email", "admin@bouncy.com", "password", "Password123")
		return nil
	})
}
