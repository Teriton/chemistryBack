MIGRATE_PATH ?="./sql/migrations"
MIGRATE_BIN ?= $(shell which migrate)
DATABASE_URL ?= "Error"

exportEnv:
	@echo set -a
	@echo source ./env
	@echo set +a

migrateUp:
	$(MIGRATE_BIN) -database $(POSTGRESQL_URL) -path sql/migrations up $(N)

migrateDown:
	$(MIGRATE_BIN) -database $(POSTGRESQL_URL) -path sql/migrations down $(N)

createMigration: 
	migrate create -ext sql -dir ./sql/migrations -seq $(name)
