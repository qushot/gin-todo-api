TOOLS=\
	github.com/cosmtrek/air@v1.45.0 \
	github.com/pressly/goose/v3/cmd/goose@v3.15.0

install-tools:
	@for tool in $(TOOLS); do \
		echo "Installing $$tool"; \
		go install $$tool; \
	done

run:
	@air
