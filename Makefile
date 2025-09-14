include .env

export $(shell sed 's/=.*//' .env)

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATE_CMD = migrate -path=/app/migrations -database $(DB_URL)
MOCKS_DIR=internal/mocks

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

restart:
	down up

logs:
	docker compose logs -f



migrate-up:
	docker-compose run --rm app $(MIGRATE_CMD) up

migrate-down:
	docker-compose run --rm app $(MIGRATE_CMD) down 1


mocks:
	mkdir -p $(MOCKS_DIR)
	mockgen -source=internal/service/service.go -destination=internal/mocks/mock_repository.go -package=mocks
	mockgen -source=internal/repository/cached_db_repository.go -destination=internal/mocks/mock_storage.go -package=mocks -mock_names=Repository=MockDB



