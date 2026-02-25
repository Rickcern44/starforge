# Bouncy: League Management System

## Project Overview
Bouncy is a league management application featuring a Go-based REST API and two frontends. Currently, **development is focused on the SvelteKit web application (`bouncy-web`)**, which serves as the primary interface for managing users, leagues, game scheduling, and payments.

### Key Components
- **Web App (`bouncy-web/`) [Primary Focus]**: A SvelteKit application for web-based league management.
- **Backend API (`api/`)**: A Go 1.25+ service providing the core business logic and persistence.
- **Mobile App (`bouncy-ui/`) [Secondary/Future]**: A Flutter application for cross-platform mobile access.

## Tech Stack
- **Web**: SvelteKit (TypeScript, Tailwind CSS v4, Svelte 5)
- **Backend**: Go (Chi, GORM, PostgreSQL, JWT)
- **Infrastructure**: Docker Compose (for PostgreSQL)
- **Mobile**: Flutter (Dart) - *Note: Secondary focus*

## Getting Started

### Prerequisites
- Node.js and npm
- Go 1.25 or higher
- Docker and Docker Compose
- (Optional) Flutter SDK for mobile work

### Building and Running (Development Workflow)

#### 1. Database
Start the PostgreSQL database:
```bash
docker compose up -d
```

#### 2. Backend API
The API requires configuration. Ensure environment variables are set.
```bash
cd api
# export POSTGRES_PASS=Password123
# export JWT_SECRET=your_secret
# export APP_CONFIG_PATH=.config/settings.dev.toml
go run cmd/api/main.go
```

#### 3. Web UI (Main Interface)
```bash
cd bouncy-web
npm install
npm run dev
```

#### 4. Mobile UI (Secondary)
```bash
cd bouncy-ui
flutter pub get
flutter run
```

## Development Conventions

### Web (SvelteKit) - *Current Priority*
- **Framework**: Svelte 5 (using runes like `$state`, `$derived`, etc.).
- **Styling**: Tailwind CSS v4.
- **Typing**: Strict TypeScript usage for all components and utilities.
- **State Management**: Prefer Svelte's built-in reactivity and context for shared state.

### Backend (Go)
- **Architecture**: Layered (CMD -> Infrastructure -> Application -> Domain).
- **Routing**: Centralized in `internal/infrastructure/api/routes`.
- **Config**: Managed via TOML and environment variables.

### Mobile (Flutter)
- **Status**: Currently in maintenance/experimental phase.
- **Architecture**: Service-based API abstraction in `lib/services`.

## Key Files
- `bouncy-web/package.json`: Main web project configuration.
- `api/cmd/api/main.go`: Backend entry point.
- `compose.yaml`: Database infrastructure.
- `bouncy-ui/pubspec.yaml`: Mobile project configuration.
