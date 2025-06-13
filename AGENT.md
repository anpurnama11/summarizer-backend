# AGENT.md - Summarizer Backend

## Build/Test Commands
- `go run cmd/api/main.go` - Start development server  
- `go build -o bin/api cmd/api/main.go` - Build binary
- `go test ./...` - Run all tests (none exist yet)
- `go test ./internal/package -run TestFunction` - Run specific test
- `go mod tidy` - Clean up dependencies

## Architecture
- **Entry Point**: `cmd/api/main.go` - Gin HTTP server on port 8080
- **Database**: SQLite at `./db/database.sqlite` with migrations in `db/migrations/`
- **Structure**: `internal/api` (handlers), `internal/repository` (data), `internal/service` (business logic)
- **Services**: `extractor` (content extraction), `openrouter` (AI summarization via OpenRouter), `summarizer` (service orchestration)
- **Middleware**: CORS, error handling, request validation

## Code Style & Conventions
- **"Best Simple System for Now"** philosophy - minimal complexity, clear intention-revealing code
- **Error Handling**: Standard Go `if err != nil` pattern, return early
- **Naming**: PascalCase exports, camelCase private, descriptive domain names
- **Types**: Struct tags for validation (`validate:"required,url"`) and JSON (`json:"field_name"`)
- **Imports**: Group standard, third-party, internal packages
- **Dependencies**: Gin (HTTP), SQLite (database), OpenRouter API, go-readability (extraction)
- **No Comments**: Code should be self-documenting through clear naming
