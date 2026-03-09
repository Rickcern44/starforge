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

func (dbs *Service) seedFeatureFlags() {
	flags := []persistence.FeatureFlag{
		{
			Base:        persistence.Base{ID: uuid.NewString()},
			Key:         "league_creation",
			Name:        "Public League Creation",
			Description: "Allow users to create their own leagues",
			Enabled:     true,
		},
		{
			Base:        persistence.Base{ID: uuid.NewString()},
			Key:         "payments",
			Name:        "Payment Processing",
			Description: "Enable tracking of user payments and financial summaries",
			Enabled:     true,
		},
		{
			Base:        persistence.Base{ID: uuid.NewString()},
			Key:         "admin_invites",
			Name:        "Admin Invitations",
			Description: "Enable the admin invitation system",
			Enabled:     true,
		},
		{
			Base:        persistence.Base{ID: uuid.NewString()},
			Key:         "notifications",
			Name:        "Notifications",
			Description: "Enable the in-app notification system",
			Enabled:     false,
		},
	}

	for _, flag := range flags {
		var existing persistence.FeatureFlag
		if err := dbs.Database.Where("key = ?", flag.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				slog.Info("Seeding feature flag", "key", flag.Key)
				_ = dbs.Database.Create(&flag).Error
			}
		}
	}
}

func (dbs *Service) Seed() error {
	// Seed Feature Flags regardless of user count (idempotent)
	dbs.seedFeatureFlags()

	var count int64
	dbs.Database.Model(&persistence.User{}).Count(&count)

	if count > 0 {
		slog.Info("Database already has data, skipping seed.")
		return nil
	}

	slog.Info("Seeding database with initial data...")

	return dbs.Database.Transaction(func(tx *gorm.DB) error {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123"), 14)
		
		// 1. Platform Admin (Global access)
		platformAdminID := uuid.NewString()
		platformAdmin := &persistence.User{
			Base:         persistence.Base{ID: platformAdminID},
			Email:        "admin@bouncy.com",
			Name:         "Platform Admin",
			PasswordHash: string(hashedPassword),
			Roles:        []string{"admin", "user"},
		}
		if err := tx.Create(platformAdmin).Error; err != nil {
			return err
		}

		// 2. League Admin (Can create leagues, but not global platform settings)
		leagueAdminID := uuid.NewString()
		leagueAdmin := &persistence.User{
			Base:         persistence.Base{ID: leagueAdminID},
			Email:        "league-admin@bouncy.com",
			Name:         "League Organizer",
			PasswordHash: string(hashedPassword),
			Roles:        []string{"league_creator", "user"},
		}
		if err := tx.Create(leagueAdmin).Error; err != nil {
			return err
		}

		// 3. Normal User (Just a player)
		regularUserID := uuid.NewString()
		regularUser := &persistence.User{
			Base:         persistence.Base{ID: regularUserID},
			Email:        "user@bouncy.com",
			Name:         "Regular Player",
			PasswordHash: string(hashedPassword),
			Roles:        []string{"user"},
		}
		if err := tx.Create(regularUser).Error; err != nil {
			return err
		}

		// 4. Create Sample Leagues
		tuesdayLeagueID := uuid.NewString()
		tuesdayLeague := &persistence.League{
			Base:     persistence.Base{ID: tuesdayLeagueID},
			Name:     "Tuesday Night Hoops",
			IsActive: true,
		}
		if err := tx.Create(tuesdayLeague).Error; err != nil {
			return err
		}

		thursdayLeagueID := uuid.NewString()
		thursdayLeague := &persistence.League{
			Base:     persistence.Base{ID: thursdayLeagueID},
			Name:     "Elite Volleyball",
			IsActive: true,
		}
		if err := tx.Create(thursdayLeague).Error; err != nil {
			return err
		}

		// 5. Enroll users in leagues
		members := []persistence.LeagueMember{
			// Tuesday Hoops: Platform Admin is Owner, Regular User is Player
			{LeagueID: tuesdayLeagueID, UserID: platformAdminID, Role: models.Role("owner"), JoinedAt: time.Now()},
			{LeagueID: tuesdayLeagueID, UserID: regularUserID, Role: models.Role("user"), JoinedAt: time.Now()},
			
			// Elite Volleyball: League Admin is Owner, Platform Admin is also Admin, Regular User is Player
			{LeagueID: thursdayLeagueID, UserID: leagueAdminID, Role: models.Role("owner"), JoinedAt: time.Now()},
			{LeagueID: thursdayLeagueID, UserID: platformAdminID, Role: models.Role("admin"), JoinedAt: time.Now()},
			{LeagueID: thursdayLeagueID, UserID: regularUserID, Role: models.Role("user"), JoinedAt: time.Now()},
		}
		for _, m := range members {
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}

		// 6. Create Upcoming Games
		gameIDs := []string{uuid.NewString(), uuid.NewString()}
		games := []persistence.Game{
			{
				Base:        persistence.Base{ID: gameIDs[0]},
				LeagueID:    tuesdayLeagueID,
				StartTime:   time.Now().Add(48 * time.Hour),
				Location:    "Main Stadium Court",
				CostInCents: 1000,
			},
			{
				Base:        persistence.Base{ID: gameIDs[1]},
				LeagueID:    thursdayLeagueID,
				StartTime:   time.Now().Add(72 * time.Hour),
				Location:    "Beach Park Court 4",
				CostInCents: 1500,
			},
		}
		for _, g := range games {
			if err := tx.Create(&g).Error; err != nil {
				return err
			}
		}

		slog.Info("Database seeded successfully with testing suite")
		slog.Info("Users: admin@bouncy.com, league-admin@bouncy.com, user@bouncy.com")
		slog.Info("Password for all: Password123")
		
		return nil
	})
}
