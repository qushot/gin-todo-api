.PHONY: run
run: postgres-up
	@go tool air

.PHONY: build-mcp
build-mcp:
	@go build -o mcp ./mcp

.PHONY: postgres-up
postgres-up:
	@docker compose up -d postgres

.PHONY: postgres-down
postgres-down:
	@docker compose down postgres

.PHONY: postgres-volumes-down
postgres-volumes-down:
	@docker compose down -v postgres

.PHONY: postgres-logs
postgres-logs:
	@docker compose logs -f postgres

.PHONY: postgres-exec
postgres-exec:
	@docker compose exec postgres psql -U postgres -d postgres

.PHONY: openapi-generator
openapi-generator:
	@docker compose run --rm openapi-generator

.PHONY: tbls
tbls:
	@docker compose run --rm tbls
