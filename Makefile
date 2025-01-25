TOOLS=\
	github.com/air-verse/air@v1.52.2 \
	github.com/pressly/goose/v3/cmd/goose@v3.15.0

install-tools:
	@for tool in $(TOOLS); do \
		echo "Installing $$tool"; \
		go install $$tool; \
	done

postgres-up:
	@docker compose up -d postgres
postgres-down:
	@docker compose down postgres
postgres-volumes-down:
	@docker compose down -v postgres
postgres-logs:
	@docker compose logs -f postgres
postgres-exec:
	@docker compose exec postgres psql -U postgres -d postgres
.PHONY: postgres-up postgres-down postgres-volumes-down postgres-logs postgres-exec

run:
	@air
