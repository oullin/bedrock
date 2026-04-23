ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_COMPOSE := docker compose -f $(GO_FMT_COMPOSE_FILE)
GO_FMT_BIN := /usr/local/bin/go-fmt
GO_FMT_EXEC := $(GO_FMT_COMPOSE) exec -T $(GO_FMT_SERVICE) $(GO_FMT_BIN)
PACKAGE_FMT := pnpm fmt
MARKDOWN_FILES := $(shell git ls-files '*.md')
GO_MODULE_EXCLUDED_DIRS := packages/testing
GO_MODULE_DIRS := $(filter-out $(GO_MODULE_EXCLUDED_DIRS),$(shell git ls-files 'packages/**/go.mod' 'services/**/go.mod' | sed 's|/go.mod$$||'))
GO_MODULE_ABS_DIRS := $(addprefix $(ROOT_PATH)/,$(GO_MODULE_DIRS))
GO_WORK_FILE := $(ROOT_PATH)/storage/.cache/go.work
GO_CMD := GOWORK=$(GO_WORK_FILE) go

define require-go-modules
	@if [ -z "$(strip $(GO_MODULE_DIRS))" ]; then \
		echo "No tracked Go modules found under packages/ or services/." >&2; \
		echo "Expected tracked go.mod files discoverable with git ls-files." >&2; \
		exit 1; \
	fi
endef

define prepare-go-workspace
	@mkdir -p $(dir $(GO_WORK_FILE))
	@rm -f $(GO_WORK_FILE) $(GO_WORK_FILE).sum
	@cd $(dir $(GO_WORK_FILE)) && GOWORK=off go work init $(GO_MODULE_ABS_DIRS)
endef

.PHONY: format format-start format-stop vet tidy typecheck test coverage build clean docs compliance go-test go-build go-coverage

format: format-start
	$(PACKAGE_FMT)
	@echo "go-fmt format in $(ROOT_PATH)"; \
	$(GO_FMT_EXEC) format --cwd $(ROOT_PATH) --host-path $(ROOT_PATH)
	@if [ -n "$(MARKDOWN_FILES)" ]; then \
		pnpm exec oxfmt --ignore-path .gitignore $(MARKDOWN_FILES); \
	fi

format-start:
	@$(GO_FMT_COMPOSE) up -d $(GO_FMT_SERVICE)

format-stop:
	@$(GO_FMT_COMPOSE) stop $(GO_FMT_SERVICE)

vet:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		echo "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) vet ./...; \
	done

tidy:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		echo "go mod tidy in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) mod tidy; \
	done

go-test:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		echo "go test ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test ./...; \
	done

go-build:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		echo "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) build ./...; \
	done

go-coverage:
	$(require-go-modules)
	$(prepare-go-workspace)
	@mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		safe=$$(echo "$$pkg" | tr '/.' '__'); \
		report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/$$safe"; \
		mkdir -p "$$report_dir"; \
		echo "go test -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -coverprofile=$$report_dir/coverage.out ./...; \
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

compliance:
	services/scripts/laravel-compliance.sh check

clean:
	rm -rf $(ROOT_PATH)/services/storage/.cache
	rm -rf $(ROOT_PATH)/services/storage/.turbo
	rm -rf $(ROOT_PATH)/services/storage/.bin
	rm -rf $(ROOT_PATH)/storage/.cache/coverage
	mkdir -p $(ROOT_PATH)/services/storage/.cache/.pnpm-store
	mkdir -p $(ROOT_PATH)/services/storage/.turbo
	mkdir -p $(ROOT_PATH)/services/storage/.bin
	mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go
	touch $(ROOT_PATH)/services/storage/.cache/.gitkeep
	touch $(ROOT_PATH)/services/storage/.turbo/.gitkeep
