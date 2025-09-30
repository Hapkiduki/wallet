.PHONY: up down migrateup migratedown docs

# Starts all docker containers in the background
up:
	docker-compose up -d

# Stops all docker containers
down:
	docker-compose down

# Applies all database migrations
migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/wallet?sslmode=disable" -verbose up

# Rolls back all database migrations
migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/wallet?sslmode=disable" -verbose down

# Generate swagger documentation
docs:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal