package commands

import (
	"github.com/bouncy/bouncy-api/internal/infrastructure/config"
	"github.com/bouncy/bouncy-api/internal/infrastructure/container"
	"github.com/bouncy/bouncy-api/internal/infrastructure/database"
)

func getContainer() (*container.AppContainer, error) {
	settings, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	dbServer := database.NewDatabaseService(settings)
	if err := dbServer.Connect(); err != nil {
		return nil, err
	}

	return container.NewAppContainer(dbServer.Database, settings), nil
}
