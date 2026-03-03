# Bouncy API

This document provides an overview of the Bouncy API project, its structure, and a guide for developers.

## To-Do List

Here is a list of suggested tasks to improve the codebase, categorized for clarity.

### Code Smells & Refactoring Opportunities
-   **[High] Adhere to Dependency Inversion Principle:** Services in the `application` layer should depend on the repository *interfaces* defined in `internal/domain/interfaces`. (Status: `AuthService`, `GameService`, `PaymentsService`, `GameAttendanceService` and `UserService` adhere. `LeagueService` still violates this and needs refactoring.)
-   **[Medium] Refactor User Domain Model:** The `User` model in the domain layer currently includes a `PasswordHash` field (though it's ignored in JSON). To strictly separate concerns, the domain model should ideally not have this field at all, and it should only exist in the persistence/DTO layer. (Status: `User` model currently in `internal/domain/models/auth.go` includes `PasswordHash`).

### Recommended Features & Updates

-   **[High] Implement a Test Suite:** The project currently lacks a comprehensive automated test suite. While some unit tests for mappers and config loading exist, it is highly recommended to add tests for application services and integration tests for API endpoints. (Status: 3 test files exist in total).
-   **[Medium] Complete Structured Logging:** Structured logging using `slog` has been integrated into handlers and most services. However, `main.go` and some older parts of the codebase still use the standard `log` package. (Status: Widely used, but incomplete in `main.go`)
-   **[Medium] Enhance API Documentation:** The project uses `Scalar-go` to serve documentation from `.docs/swagger.yaml`. While some handlers (`GameHandler`, `UserHandler`) have Swagger annotations, others (`AuthHandler`, `PaymentsHandler`) are missing them. Adding complete annotations will ensure a comprehensive API specification.
-   **[Medium] Implement Pagination:** For endpoints that return lists (e.g., listing games, payments, or members), pagination should be implemented to improve performance and usability as the database grows.
-   **[Low] Establish a CI/CD Pipeline:** A Continuous Integration/Continuous Deployment pipeline (e.g., using GitHub Actions) should be set up to automate testing, building, and deployment processes.

#### Completed Enhancements
-   **[High] Create "Get Current User" Endpoint:** Added `GET /api/v1/users/me` endpoint. (Status: Completed)
-   **[Medium] Centralize Configuration:** Configuration is now managed via a centralized struct and loaded using `koanf`, supporting defaults, TOML files, and environment variables. (Status: Completed)

---

## Code Structure

This project follows a [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) pattern. The code is organized into:

-   **/cmd/api**: Entry point. `main.go` initializes the server, database, dependencies, and starts the app.
-   **/internal/domain**: Core business models and interfaces.
    -   **/models**: Core data structures (e.g., `League`, `Game`, `User`, `Invitation`).
    -   **/interfaces**: Repository and service contracts.
-   **/internal/application**: Business logic. Services orchestrate data flow and call domain interfaces. Organized by domain (e.g., `leagues`, `payments`, `users`).
-   **/internal/infrastructure**: Implementation details for external concerns.
    -   **/api**: Web API related code.
        -   **/handlers**: HTTP handlers.
        -   **/routes**: API route definitions.
        -   **/middleware**: JWT and Role-based access control.
        -   **/contract**: Shared request/response types.
    -   **/persistence**: Repository implementations using GORM.
    -   **/config**: Centralized configuration management.
    -   **/auth**: JWT claims and token generation.
-   **/.http-requests**: Directory containing `.http` files for testing endpoints (JetBrains/VS Code).

---

## Object Structure

-   **`User`**: A global user.
-   **`Invitation`**: A pending invite to a league, linked to an email.
-   **`League`**: A sports league containing members and games.
-   **`LeagueMember`**: Connects a `User` to a `League` with a specific `Role`.
-   **`Game`**: A scheduled game within a `League`.
-   **`GameAttendance`**: Tracks user check-ins/RSVPs to a `Game`.
-   **`GameCharge`**: Financial record for a user's participation in a `Game`.
-   **`Payment`**: A record of funds received from a user.
-   **`PaymentAllocation`**: Links a `Payment` to a `GameCharge`.

---

## API Endpoints (v1)

### Authentication & Users
-   `POST /api/v1/auth/login`: Authenticate and receive a JWT.
-   `POST /api/v1/auth/register`: Register using an invitation token.
-   `GET /api/v1/users/me`: Get profile of the currently authenticated user.
-   `POST /api/v1/admin/invite`: (Admin) Invite a user by email to a league.
-   `GET /api/v1/admin/league/{leagueId}/invitations`: (Admin) List invitations for a league.
-   `PATCH /api/v1/users/{userId}/roles`: (Admin) Update a user's global roles.

### Leagues
-   `POST /api/v1/league`: Create a new league.
-   `GET /api/v1/league/{leagueId}`: Get league details.
-   `GET /api/v1/me/leagues`: List all leagues the current user is a member of.
-   `DELETE /api/v1/league/{leagueId}`: Delete a league.

### League Members
-   `GET /api/v1/league/{leagueId}/members`: List all members of a league.
-   `POST /api/v1/league/{leagueId}/members`: Add a member to a league.
-   `PATCH /api/v1/league/{leagueId}/members/{memberId}`: Update a member's role.
-   `DELETE /api/v1/league/{leagueId}/members/{memberId}`: Remove a member.

### Games & Attendance
-   `GET /api/v1/league/{leagueId}/games`: List all games for a league.
-   `POST /api/v1/league/{leagueId}/games`: Create a game (supports recurring).
-   `GET /api/v1/game/{gameId}`: Get game details.
-   `PUT /api/v1/game/{gameId}`: Update game details.
-   `DELETE /api/v1/game/{gameId}`: Cancel a game.
-   `POST /api/v1/game/{gameId}/attendance`: RSVP/Check-in to a game.
-   `DELETE /api/v1/game/{gameId}/attendance/{userId}`: Remove attendance.

### Payments & Charges
-   `GET /api/v1/me/payments`: List payments for the current user.
-   `GET /api/v1/me/charges`: List charges for the current user.
-   `GET /api/v1/league/{leagueId}/payments`: (Admin) List all payments in a league.
-   `GET /api/v1/league/{leagueId}/financials`: (Admin) Get a financial summary for the league (total collected, charges, etc.).
-   `POST /api/v1/league/{leagueId}/payments`: (Admin) Record a new payment.
-   `POST /api/v1/payments/{paymentId}/allocations`: (Admin) Allocate payment to a charge.
-   `POST /api/v1/admin/payments/claim`: (Admin) Link unclaimed records to a user.

---

## Walkthrough: Adding a New API Endpoint

Example: Adding `GET /api/v1/admin/users` to list all users.

### 1. Define Repository Interface (Domain)
**File:** `internal/domain/interfaces/repositories.go`
```go
type UserRepository interface {
    // ...
    ListAllUsers() ([]*models.User, error)
}
```

### 2. Implement Repository (Persistence)
**File:** `internal/infrastructure/persistence/repositories/auth_repository.go`
```go
func (r *AuthRepository) ListAllUsers() ([]*models.User, error) {
    var users []persistence.User
    if err := r.db.Find(&users).Error; err != nil {
        return nil, err
    }
    // Map to domain models...
    return mappers.UsersToDomain(users), nil
}
```

### 3. Add Service Method (Application)
**File:** `internal/application/users/user_service.go`
```go
func (s *Service) ListAllUsers() ([]*models.User, error) {
    return s.userRepo.ListAllUsers()
}
```

### 4. Create Handler (Infrastructure/API)
**File:** `internal/infrastructure/api/handlers/user_handler.go`
```go
func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.service.ListAllUsers()
    if err != nil {
        utils.WriteJSON(w, http.StatusInternalServerError, contract.ErrorResponse{Error: err.Error()})
        return
    }
    utils.WriteJSON(w, http.StatusOK, users)
}
```

### 5. Register Route
**File:** `internal/infrastructure/api/handlers/user_handler.go`
```go
func RegisterAdminUserRoutes(r chi.Router, handler *UserHandler) {
    r.Get("/users", handler.ListAllUsers)
    r.Patch("/users/{userId}/roles", handler.UpdateUserRoles)
}
```
*(Ensure `RegisterAdminUserRoutes` is called within an admin-protected group in `routes.go`)*
