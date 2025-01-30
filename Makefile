include scripts/*.mk

DEV_COMPOSE_ARGS=--env-file .env.dev -f docker-compose.dev.yaml
DEV_COMPOSE_ENV=docker compose $(DEV_COMPOSE_ARGS)
DEV_COMPOSE=docker compose $(DEV_COMPOSE_ARGS)

dev-build:
	$(DEV_COMPOSE) build
dev-up: api_docker_build dev-build
	$(DEV_COMPOSE) up -d
dev-api-run:
	go run cmd/api/api.go