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
