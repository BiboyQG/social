include .envrc
MIGRATIONS_DIR=./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@echo "Creating migration file..."
	@migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@echo "Migrating database..."
	@migrate -database $(DB_ADDR) -path $(MIGRATIONS_DIR) up

.PHONY: migrate-down
migrate-down:
	@echo "Migrating database down..."
	@migrate -database $(DB_ADDR) -path $(MIGRATIONS_DIR) down

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
