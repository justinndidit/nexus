dev-compose-build:
	@docker compose -f infra/docker-compose.dev.yml up --build

dev-compose-up:
	@docker compose -f infra/docker-compose.dev.yml up