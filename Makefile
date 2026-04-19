ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_COMPOSE := docker compose -f $(GO_FMT_COMPOSE_FILE)
GO_FMT_BIN := /usr/local/bin/go-fmt
GO_FMT_EXEC := $(GO_FMT_COMPOSE) exec -T $(GO_FMT_SERVICE) $(GO_FMT_BIN)
PACKAGE_FMT := pnpm fmt
MARKDOWN_FILES := $(shell git ls-files '*.md')
GO_MODULE_DIRS := $(shell awk 'BEGIN { in_use = 0 } /^use \(/ { in_use = 1; next } in_use && /^\)/ { in_use = 0; next } in_use { gsub(/^\.\//, "", $$1); print $$1 }' go.work)

.PHONY: format format-start format-stop vet tidy typecheck test coverage build clean docs go-test go-build go-coverage

format: format-start
	$(PACKAGE_FMT)
	@broadcastclient "go-fmt format in $(ROOT_PATH)"; \
	$(GO_FMT_EXEC) format --cwd $(ROOT_PATH) --host-path $(ROOT_PATH)
	@if [ -n "$(MARKDOWN_FILES)" ]; then \
		pnpm exec oxfmt --ignore-path .gitignore $(MARKDOWN_FILES); \
	fi

format-start:
	@$(GO_FMT_COMPOSE) up -d $(GO_FMT_SERVICE)

format-stop:
	@$(GO_FMT_COMPOSE) stop $(GO_FMT_SERVICE)

vet:
	@for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go vet ./...; \
	done

tidy:
	@for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go mod tidy in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go mod tidy; \
	done

go-test:
	@for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go test ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go test ./...; \
	done

go-build:
	@for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go build ./...; \
	done

go-coverage:
	@for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go test -cover ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go test -cover ./...; \
	done

typecheck:
	pnpm typecheck

test:
	pnpm test

coverage:
	pnpm test:coverage

build:
	pnpm build

docs:
	pnpm install
	pnpm --filter=@bedrock/docs run dev

clean:
	rm -rf $(ROOT_PATH)/service/storage/.cache
	rm -rf $(ROOT_PATH)/service/storage/.turbo
	rm -rf $(ROOT_PATH)/service/storage/.bin
	mkdir -p $(ROOT_PATH)/service/storage/.cache/.pnpm-store
	mkdir -p $(ROOT_PATH)/service/storage/.cache/coverage/go
	mkdir -p $(ROOT_PATH)/service/storage/.cache/coverage/playwright
	mkdir -p $(ROOT_PATH)/service/storage/.turbo
	mkdir -p $(ROOT_PATH)/service/storage/.bin
	touch $(ROOT_PATH)/service/storage/.cache/.gitkeep
	touch $(ROOT_PATH)/service/storage/.turbo/.gitkeep
