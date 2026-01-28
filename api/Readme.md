# Bouncy API

This document provides an overview of the Bouncy API project, its structure, and a guide for developers.

## To-Do List

Here is a list of suggested tasks to improve the codebase, categorized for clarity.

### Code Smells & Refactoring Opportunities
-   **[High] Adhere to Dependency Inversion Principle:** Services in the `application` layer currently depend on concrete repository implementations from the `infrastructure/persistence` layer. They should instead depend on the repository *interfaces* defined in `internal/domain/interfaces`.
-   **[High] Fix Inconsistent Route Parameter Names:** There are inconsistencies between route definitions and their handlers. For example, in `internal/infrastructure/api/handlers/league_handler.go`, the route is defined with `/:leagueId` but the handler reads `c.Param("id")`. This is a bug and should be fixed across the application.
-   **[Medium] Centralize Configuration:** Configuration values, especially secrets like the `JWT_TOKEN`, are read directly from environment variables within the business logic. It would be better to use a dedicated configuration struct, populated at startup, to make configuration management more robust and explicit.
-   **[Low] Remove Dead Code:** The `internal/infrastructure/api/handlers/league_handler.go` file contains commented-out code. This should be either implemented or removed to keep the codebase clean.

### Recommended Features & Updates
-   **[High] Implement a Test Suite:** The project currently lacks automated tests. It is highly recommended to add unit tests for the application services and integration tests for the API endpoints to ensure code quality and prevent regressions.
-   **[Medium] Introduce Structured Logging:** The current logging uses the standard `log` package, which is not ideal for production environments. Integrating a structured logging library like `slog`, `zerolog`, or `zap` would significantly improve log parsing and monitoring.
-   **[Medium] Enhance API Documentation:** The project is set up with `swaggo` for Swagger documentation, but the annotations are minimal. Adding more detailed comments to the handlers will produce a more comprehensive and useful API specification.
-   **[Medium] Implement Pagination:** For endpoints that can return a large number of items (e.g., a list of all leagues), pagination should be implemented to improve performance and usability.
-   **[Low] Establish a CI/CD Pipeline:** A Continuous Integration/Continuous Deployment pipeline (e.g., using GitHub Actions) should be set up to automate testing, building, and deployment processes.

---

## Code Structure

This project follows a [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) pattern, which separates concerns and makes the application more modular, testable, and maintainable. The code is organized into the following main directories:

-   **/cmd/api**: The entry point of the application. The `main.go` file is responsible for initializing the server, database, dependencies, and starting the application.
-   **/internal/domain**: This is the core of the application. It contains the business models and the interfaces that define the behavior of the outer layers (like repositories).
    -   **/models**: Contains the core data structures (e.g., `League`, `Game`).
    -   **/interfaces**: Defines the contracts for repositories that the application layer will use.
-   **/internal/application**: Contains the application-specific business logic. Each service orchestrates the flow of data and calls the domain interfaces to perform operations.
-   **/internal/infrastructure**: This layer contains the implementation details for external concerns, such as the database, the web server, and other third-party services.
    -   **/api**: Contains everything related to the web API, including:
        -   **/handlers**: The HTTP handlers that receive requests, call the appropriate application services, and return responses.
        -   **/routes**: The definition of the API routes.
        -   **/server**: The web server configuration and setup (using Gin).
    -   **/persistence**: The implementation of the repository interfaces defined in the domain layer. This is where the data is actually persisted to the database (using GORM).
        -   **/repositories**: The concrete implementations of the repository interfaces.
        -   **/mappers**: Functions to map data between the domain models and the persistence models.
-   **/bouncy.http**: A file containing HTTP requests for testing the API endpoints, likely for use with a tool like the JetBrains HTTP Client or a similar VS Code extension.

---

## Walkthrough: Adding a New API Endpoint

This guide will walk you through the process of adding a new API endpoint. As an example, we will add an endpoint to **retrieve all games for a specific league**: `GET /api/leagues/{leagueId}/games`.

### Step 1: Define the Repository Method (Domain Layer)
The application services should only depend on interfaces. First, define the method for retrieving games in the `GameRepository` interface.

**File:** `internal/domain/interfaces/repositories.go`
```go
type GameRepository interface {
    GetById(id string) (*models.Game, error)
    Save(game *models.Game) error
    Delete(id string) error
    // Add this new method:
    GetGamesByLeagueID(leagueID string) ([]*models.Game, error)
}
```

### Step 2: Implement the Repository Method (Persistence Layer)
Now, implement the `GetGamesByLeagueID` method in the concrete `GameRepository`. For this example, let's assume you have a `game_repository.go` file.

**File:** `internal/infrastructure/persistence/repositories/game_repository.go`
```go
// (Assuming you have a GameRepository struct similar to LeagueRepository)

func (gr *GameRepository) GetGamesByLeagueID(leagueID string) ([]*models.Game, error) {
    var gameModels []persistence.Game
    // Assuming 'LeagueID' is a field on the persistence.Game model
    if err := gr.db.Where("league_id = ?", leagueID).Find(&gameModels).Error; err != nil {
        return nil, err
    }

    // Map persistence models to domain models
    var domainGames []*models.Game
    for _, gameModel := range gameModels {
        domainGame := mappers.GameToDomain(gameModel) // You'll need a Game mapper
        domainGames = append(domainGames, &domainGame)
    }

    return domainGames, nil
}
```
*Note: You may need to create a `GameRepository` and `game_mapper.go` similar to what exists for `League`.*

### Step 3: Create the Application Service
Create a new service that will contain the business logic for this operation.

**File:** `internal/application/games/game_service.go` (new file)
```go
package games

import (
    "github.com/bouncy/bouncy-api/internal/domain/interfaces"
    "github.com/bouncy/bouncy-api/internal/domain/models"
)

type GameService struct {
    gameRepo interfaces.GameRepository // Depend on the interface
}

func NewGameService(gameRepo interfaces.GameRepository) *GameService {
    return &GameService{gameRepo: gameRepo}
}

func (s *GameService) GetGamesForLeague(leagueID string) ([]*models.Game, error) {
    return s.gameRepo.GetGamesByLeagueID(leagueID)
}
```

### Step 4: Create the Handler (API Layer)
Add a new handler function in the `GameHandler` to process the incoming HTTP request.

**File:** `internal/infrastructure/api/handlers/game_handler.go`
```go
// (Assuming you have a GameHandler struct similar to LeagueHandler)

func (h *GameHandler) GetGamesForLeague(c *gin.Context) {
    leagueID := c.Param("leagueId") // Make sure this matches the route param

    games, err := h.service.GetGamesForLeague(leagueID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch games for the league"})
        return
    }

    c.JSON(http.StatusOK, games)
}
```

### Step 5: Register the New Route
Now, register the new route and connect it to the handler. A good place for this might be in the `league_handler.go` or a new `game_handler.go` registration function. For this example, let's add it to a hypothetical `RegisterGameRoutes`.

**File:** `internal/infrastructure/api/handlers/game_handler.go`
```go
func RegisterGameRoutes(rg *gin.RouterGroup, handler *GameHandler) {
	games := rg.Group("/games")
	// ... other game routes
}
```
**File:** `internal/infrastructure/api/routes/routes.go`
```go
func RegisterRoutes(
	engine *gin.Engine,
	deps *dependencies.Dependencies,
) {
	api := engine.Group("/api")

	// ... other route registrations

	// Register a sub-route on leagues
	leagues := api.Group("/leagues")
	leagues.GET("/:leagueId/games", deps.GameHandler.GetGamesForLeague) // Add the new route here
}
```
*Note: This approach assumes a `GameHandler` is added to your dependencies. You might choose to add this route under the `/leagues` group or create a new top-level `/games` group depending on your API design.*

### Step 6: Wire Up Dependencies
Finally, update the dependency injection setup in `main.go` and `dependencies.go` to include the new service and repository.

**File:** `internal/infrastructure/api/dependencies/dependencies.go`
```go
// Add GameHandler to Dependencies struct
type Dependencies struct {
    // ... other handlers
    GameHandler         *handlers.GameHandler
}

// Update BuildDependencies to create and include GameHandler
func BuildDependencies(
    // ... other services
    gameService *games.GameService,
) *Dependencies {
    return &Dependencies{
        // ...
        GameHandler: handlers.NewGameHandler(gameService),
    }
}
```

**File:** `cmd/api/main.go`
```go
func BuildApplication(db *gorm.DB) *dependencies.Dependencies {
    // ... other repository initializations
    gameRepo := repositories.NewGameRepository(db) // Assuming NewGameRepository exists

    // ... other service initializations
    gameService := games.NewGameService(gameRepo)

    // Pass the new service to the dependency builder
    return dependencies.BuildDependencies(
        // ... other services
        gameService,

    )
}

// And in main(), register the new handler with the routes
func main() {
    // ...
    deps := BuildApplication(dbServer.Database)
	routes.RegisterRoutes(ginServer.Engine(), deps) // Make sure your RegisterRoutes wires everything up
    // ...
}
```
And with that, your new endpoint should be ready to be tested!
