.PHONY: run
run: postgres-up redis-up
	@go tool air

.PHONY: run-in-memory
run-in-memory:
	@go tool air --build.cmd "go build -tags in_memory -buildvcs=false -o ./tmp/main ./cmd/api"

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

.PHONY: redis-up
redis-up:
	@docker compose up -d redis

.PHONY: redis-down
redis-down:
	@docker compose down redis

.PHONY: redis-volumes-down
redis-volumes-down:
	@docker compose down -v redis

.PHONY: redis-logs
redis-logs:
	@docker compose logs -f redis

.PHONY: redis-exec
redis-exec:
	@docker compose exec redis redis-cli	

.PHONY: openapi-generator
openapi-generator:
	@docker compose run --rm openapi-generator

.PHONY: tbls
tbls:
	@docker compose run --rm tbls
