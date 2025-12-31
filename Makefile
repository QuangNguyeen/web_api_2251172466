.PHONY: run build swagger test clean

# Run the application
run:
	go run cmd/api/main.go

# Build the binary
build:
	go build -o bin/library-api cmd/api/main.go

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build files
clean:
	rm -rf bin/ docs/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run with hot reload (requires air)
dev:
	air

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Create database
db-create:
	psql -U postgres -c "CREATE DATABASE library_db;"

# Run database migration
db-migrate:
	psql -U postgres -d library_db -f database/schema.sql

# Seed database
db-seed:
	psql -U postgres -d library_db -f database/seed.sql

# Drop database
db-drop:
	psql -U postgres -c "DROP DATABASE IF library_db EXISTS;"

# Reset database (drop, create, migrate, seed)
db-reset: db-drop db-create db-migrate db-seed
	@echo "Database reset complete!"
