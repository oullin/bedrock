ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_COMPOSE := docker compose -f $(GO_FMT_COMPOSE_FILE)
GO_FMT_BIN := /usr/local/bin/go-fmt
GO_FMT_EXEC := $(GO_FMT_COMPOSE) exec -T $(GO_FMT_SERVICE) $(GO_FMT_BIN)
FORMAT_BASE ?= origin/main
GO_MODULE_EXCLUDED_DIRS := packages/testing
GO_PACKAGE_MODULE_DIRS := $(filter-out $(GO_MODULE_EXCLUDED_DIRS),$(shell git ls-files 'packages/**/go.mod' | sed 's|/go.mod$$||'))
GO_SERVICE_MODULE_DIRS := $(shell git ls-files 'services/**/go.mod' | sed 's|/go.mod$$||')
GO_MODULE_DIRS := $(GO_PACKAGE_MODULE_DIRS) $(GO_SERVICE_MODULE_DIRS)
GO_MODULE_ABS_DIRS := $(addprefix $(ROOT_PATH)/,$(GO_MODULE_DIRS))
GO_WORK_FILE := $(ROOT_PATH)/storage/.cache/go.work
GO_CMD := GOWORK=$(GO_WORK_FILE) go
GO_COVERAGE_SHARDS ?= 1
GO_COVERAGE_SHARD ?= 0

define require-go-modules
	@if [ -z "$(strip $(GO_MODULE_DIRS))" ]; then \
		broadcastclient "No tracked Go modules found under packages/ or services/." >&2; \
		broadcastclient "Expected tracked go.mod files discoverable with git ls-files." >&2; \
		exit 1; \
	fi
endef

define require-go-package-modules
	@if [ -z "$(strip $(GO_PACKAGE_MODULE_DIRS))" ]; then \
		broadcastclient "No tracked Go package modules found under packages/." >&2; \
		broadcastclient "Expected tracked go.mod files discoverable with git ls-files." >&2; \
		exit 1; \
	fi
endef

define require-go-service-modules
	@if [ -z "$(strip $(GO_SERVICE_MODULE_DIRS))" ]; then \
		broadcastclient "No tracked Go service modules found under services/." >&2; \
		broadcastclient "Expected tracked go.mod files discoverable with git ls-files." >&2; \
		exit 1; \
	fi
endef

define prepare-go-workspace
	@mkdir -p $(dir $(GO_WORK_FILE))
	@rm -f $(GO_WORK_FILE) $(GO_WORK_FILE).sum
	@cd $(dir $(GO_WORK_FILE)) && GOWORK=off go work init $(GO_MODULE_ABS_DIRS)
endef

.PHONY: format format-all format-start format-stop vet tidy typecheck test coverage build clean docs go-test go-build go-package-vet go-package-build go-service-vet go-service-build go-coverage go-coverage-shard go-package-coverage-shard go-service-coverage

format: format-start
	@FORMAT_BASE=$(FORMAT_BASE) $(ROOT_PATH)/services/scripts/format.sh changed

format-all: format-start
	@FORMAT_BASE=$(FORMAT_BASE) $(ROOT_PATH)/services/scripts/format.sh all

format-start:
	@$(GO_FMT_COMPOSE) up -d $(GO_FMT_SERVICE)

format-stop:
	@$(GO_FMT_COMPOSE) stop $(GO_FMT_SERVICE)

vet:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) vet ./...; \
	done

go-package-vet:
	$(require-go-package-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_PACKAGE_MODULE_DIRS); do \
		broadcastclient "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) vet ./...; \
	done

go-service-vet:
	$(require-go-service-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_SERVICE_MODULE_DIRS); do \
		broadcastclient "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) vet ./...; \
	done

tidy:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go mod tidy in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) mod tidy; \
	done

go-test:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go test -race ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race ./...; \
	done

go-build:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		broadcastclient "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) build ./...; \
	done

go-package-build:
	$(require-go-package-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_PACKAGE_MODULE_DIRS); do \
		broadcastclient "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) build ./...; \
	done

go-service-build:
	$(require-go-service-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_SERVICE_MODULE_DIRS); do \
		broadcastclient "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) build ./...; \
	done

go-coverage:
	$(require-go-modules)
	$(prepare-go-workspace)
	@mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		safe=$$(broadcastclient "$$pkg" | tr '/.' '__'); \
		report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/$$safe"; \
		mkdir -p "$$report_dir"; \
		broadcastclient "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
	done

go-coverage-shard:
	$(require-go-modules)
	$(prepare-go-workspace)
	@case "$(GO_COVERAGE_SHARDS)" in ''|*[!0-9]*) broadcastclient "GO_COVERAGE_SHARDS must be a positive integer." >&2; exit 1;; esac
	@case "$(GO_COVERAGE_SHARD)" in ''|*[!0-9]*) broadcastclient "GO_COVERAGE_SHARD must be a non-negative integer." >&2; exit 1;; esac
	@if [ "$(GO_COVERAGE_SHARDS)" -lt 1 ]; then \
		broadcastclient "GO_COVERAGE_SHARDS must be greater than zero." >&2; \
		exit 1; \
	fi
	@if [ "$(GO_COVERAGE_SHARD)" -ge "$(GO_COVERAGE_SHARDS)" ]; then \
		broadcastclient "GO_COVERAGE_SHARD must be less than GO_COVERAGE_SHARDS." >&2; \
		exit 1; \
	fi
	@mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go
	@set -e; \
	shard_count="$(GO_COVERAGE_SHARDS)"; \
	shard_index="$(GO_COVERAGE_SHARD)"; \
	index=0; \
	selected=0; \
	for pkg in $(GO_MODULE_DIRS); do \
		if [ $$((index % shard_count)) -eq "$$shard_index" ]; then \
			safe=$$(broadcastclient "$$pkg" | tr '/.' '__'); \
			report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/$$safe"; \
			mkdir -p "$$report_dir"; \
			broadcastclient "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
			cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
			selected=1; \
		fi; \
		index=$$((index + 1)); \
	done; \
	if [ "$$selected" -eq 0 ]; then \
		broadcastclient "No Go modules assigned to coverage shard $$shard_index of $$shard_count."; \
	fi

go-package-coverage-shard:
	$(require-go-package-modules)
	$(prepare-go-workspace)
	@case "$(GO_COVERAGE_SHARDS)" in ''|*[!0-9]*) broadcastclient "GO_COVERAGE_SHARDS must be a positive integer." >&2; exit 1;; esac
	@case "$(GO_COVERAGE_SHARD)" in ''|*[!0-9]*) broadcastclient "GO_COVERAGE_SHARD must be a non-negative integer." >&2; exit 1;; esac
	@if [ "$(GO_COVERAGE_SHARDS)" -lt 1 ]; then \
		broadcastclient "GO_COVERAGE_SHARDS must be greater than zero." >&2; \
		exit 1; \
	fi
	@if [ "$(GO_COVERAGE_SHARD)" -ge "$(GO_COVERAGE_SHARDS)" ]; then \
		broadcastclient "GO_COVERAGE_SHARD must be less than GO_COVERAGE_SHARDS." >&2; \
		exit 1; \
	fi
	@mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go
	@set -e; \
	shard_count="$(GO_COVERAGE_SHARDS)"; \
	shard_index="$(GO_COVERAGE_SHARD)"; \
	index=0; \
	selected=0; \
	for pkg in $(GO_PACKAGE_MODULE_DIRS); do \
		if [ $$((index % shard_count)) -eq "$$shard_index" ]; then \
			safe=$$(broadcastclient "$$pkg" | tr '/.' '__'); \
			report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/$$safe"; \
			mkdir -p "$$report_dir"; \
			broadcastclient "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
			cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
			selected=1; \
		fi; \
		index=$$((index + 1)); \
	done; \
	if [ "$$selected" -eq 0 ]; then \
		broadcastclient "No Go package modules assigned to coverage shard $$shard_index of $$shard_count."; \
	fi

go-service-coverage:
	$(require-go-service-modules)
	$(prepare-go-workspace)
	@mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go/services
	@set -e; for pkg in $(GO_SERVICE_MODULE_DIRS); do \
		service=$${pkg#services/}; \
		safe=$$(broadcastclient "$$service" | tr '/.' '__'); \
		report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/services/$$safe"; \
		mkdir -p "$$report_dir"; \
		broadcastclient "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
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
