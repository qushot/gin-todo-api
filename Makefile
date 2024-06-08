TOOLS=\
	github.com/air-verse/air@v1.52.2 \
	github.com/pressly/goose/v3/cmd/goose@v3.15.0

install-tools:
	@for tool in $(TOOLS); do \
		echo "Installing $$tool"; \
		go install $$tool; \
	done

run:
	@air
