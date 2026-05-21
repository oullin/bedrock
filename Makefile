ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_COMPOSE := docker compose -f $(GO_FMT_COMPOSE_FILE)
GO_FMT_BIN := /usr/local/bin/go-fmt
GO_FMT_EXEC := $(GO_FMT_COMPOSE) exec -T $(GO_FMT_SERVICE) $(GO_FMT_BIN)
PACKAGE_FMT := pnpm fmt
MARKDOWN_FILES := $(shell git ls-files '*.md')
FORMAT_BASE ?= origin/main
GIT_CHANGED_FN = $(shell { \
		git diff --name-only --diff-filter=ACMRT $(FORMAT_BASE)...HEAD -- $(1) 2>/dev/null; \
		git diff --name-only --diff-filter=ACMRT -- $(1) 2>/dev/null; \
		git ls-files --others --exclude-standard -- $(1) 2>/dev/null; \
	} | sort -u)
CHANGED_MD := $(call GIT_CHANGED_FN,'*.md')
CHANGED_GO := $(call GIT_CHANGED_FN,'*.go')
CHANGED_GO_DIRS := $(sort $(patsubst %/,%,$(dir $(CHANGED_GO))))
CHANGED_GO_MODULES := $(shell \
		{ git diff --name-only --diff-filter=ACMRT $(FORMAT_BASE)...HEAD -- '*.go' 2>/dev/null; \
		  git diff --name-only --diff-filter=ACMRT -- '*.go' 2>/dev/null; \
		  git ls-files --others --exclude-standard -- '*.go' 2>/dev/null; } \
		| sort -u \
		| while read f; do \
		    d=$$(dirname $$f); \
		    while [ "$$d" != "." ] && [ ! -f "$$d/go.mod" ]; do d=$$(dirname $$d); done; \
		    [ "$$d" != "." ] && echo $$d; \
		  done \
		| sort -u)
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
		echo "No tracked Go modules found under packages/ or services/." >&2; \
		echo "Expected tracked go.mod files discoverable with git ls-files." >&2; \
		exit 1; \
	fi
endef

define require-go-package-modules
	@if [ -z "$(strip $(GO_PACKAGE_MODULE_DIRS))" ]; then \
		echo "No tracked Go package modules found under packages/." >&2; \
		echo "Expected tracked go.mod files discoverable with git ls-files." >&2; \
		exit 1; \
	fi
endef

define require-go-service-modules
	@if [ -z "$(strip $(GO_SERVICE_MODULE_DIRS))" ]; then \
		echo "No tracked Go service modules found under services/." >&2; \
		echo "Expected tracked go.mod files discoverable with git ls-files." >&2; \
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
	@$(PACKAGE_FMT) & pnpm_pid=$$!; \
	go_fmt_status=0; \
	if [ -n "$(strip $(CHANGED_GO_MODULES))" ]; then \
		paths=""; \
		for dir in $(CHANGED_GO_MODULES); do paths="$$paths /work/$$dir"; done; \
		echo "go-fmt format ($(words $(CHANGED_GO_MODULES)) module(s))"; \
		$(GO_FMT_EXEC) format --cwd /work $$paths || go_fmt_status=$$?; \
	else \
		echo "go-fmt: no changed Go modules"; \
	fi; \
	wait $$pnpm_pid; pnpm_status=$$?; \
	if [ $$go_fmt_status -ne 0 ]; then exit $$go_fmt_status; fi; \
	if [ $$pnpm_status -ne 0 ]; then exit $$pnpm_status; fi
	@if [ -n "$(strip $(CHANGED_MD))" ]; then \
		echo "oxfmt ($(words $(CHANGED_MD)) markdown files)"; \
		pnpm exec oxfmt --ignore-path .gitignore $(CHANGED_MD); \
	else \
		echo "oxfmt: no changed markdown files"; \
	fi

format-all: format-start
	@$(PACKAGE_FMT) & pnpm_pid=$$!; \
	go_fmt_status=0; \
	echo "go-fmt format in $(ROOT_PATH)"; \
	$(GO_FMT_EXEC) format --cwd /work --host-path $(ROOT_PATH) || go_fmt_status=$$?; \
	wait $$pnpm_pid; pnpm_status=$$?; \
	if [ $$go_fmt_status -ne 0 ]; then exit $$go_fmt_status; fi; \
	if [ $$pnpm_status -ne 0 ]; then exit $$pnpm_status; fi
	@if [ -n "$(MARKDOWN_FILES)" ]; then \
		echo "oxfmt ($(words $(MARKDOWN_FILES)) markdown files)"; \
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

go-package-vet:
	$(require-go-package-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_PACKAGE_MODULE_DIRS); do \
		echo "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) vet ./...; \
	done

go-service-vet:
	$(require-go-service-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_SERVICE_MODULE_DIRS); do \
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
		echo "go test -race ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race ./...; \
	done

go-build:
	$(require-go-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_MODULE_DIRS); do \
		echo "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) build ./...; \
	done

go-package-build:
	$(require-go-package-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_PACKAGE_MODULE_DIRS); do \
		echo "go build ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) build ./...; \
	done

go-service-build:
	$(require-go-service-modules)
	$(prepare-go-workspace)
	@set -e; for pkg in $(GO_SERVICE_MODULE_DIRS); do \
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
		echo "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
	done

go-coverage-shard:
	$(require-go-modules)
	$(prepare-go-workspace)
	@case "$(GO_COVERAGE_SHARDS)" in ''|*[!0-9]*) echo "GO_COVERAGE_SHARDS must be a positive integer." >&2; exit 1;; esac
	@case "$(GO_COVERAGE_SHARD)" in ''|*[!0-9]*) echo "GO_COVERAGE_SHARD must be a non-negative integer." >&2; exit 1;; esac
	@if [ "$(GO_COVERAGE_SHARDS)" -lt 1 ]; then \
		echo "GO_COVERAGE_SHARDS must be greater than zero." >&2; \
		exit 1; \
	fi
	@if [ "$(GO_COVERAGE_SHARD)" -ge "$(GO_COVERAGE_SHARDS)" ]; then \
		echo "GO_COVERAGE_SHARD must be less than GO_COVERAGE_SHARDS." >&2; \
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
			safe=$$(echo "$$pkg" | tr '/.' '__'); \
			report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/$$safe"; \
			mkdir -p "$$report_dir"; \
			echo "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
			cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
			selected=1; \
		fi; \
		index=$$((index + 1)); \
	done; \
	if [ "$$selected" -eq 0 ]; then \
		echo "No Go modules assigned to coverage shard $$shard_index of $$shard_count."; \
	fi

go-package-coverage-shard:
	$(require-go-package-modules)
	$(prepare-go-workspace)
	@case "$(GO_COVERAGE_SHARDS)" in ''|*[!0-9]*) echo "GO_COVERAGE_SHARDS must be a positive integer." >&2; exit 1;; esac
	@case "$(GO_COVERAGE_SHARD)" in ''|*[!0-9]*) echo "GO_COVERAGE_SHARD must be a non-negative integer." >&2; exit 1;; esac
	@if [ "$(GO_COVERAGE_SHARDS)" -lt 1 ]; then \
		echo "GO_COVERAGE_SHARDS must be greater than zero." >&2; \
		exit 1; \
	fi
	@if [ "$(GO_COVERAGE_SHARD)" -ge "$(GO_COVERAGE_SHARDS)" ]; then \
		echo "GO_COVERAGE_SHARD must be less than GO_COVERAGE_SHARDS." >&2; \
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
			safe=$$(echo "$$pkg" | tr '/.' '__'); \
			report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/$$safe"; \
			mkdir -p "$$report_dir"; \
			echo "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
			cd $(ROOT_PATH)/$$pkg && $(GO_CMD) test -race -coverprofile=$$report_dir/coverage.out ./...; \
			selected=1; \
		fi; \
		index=$$((index + 1)); \
	done; \
	if [ "$$selected" -eq 0 ]; then \
		echo "No Go package modules assigned to coverage shard $$shard_index of $$shard_count."; \
	fi

go-service-coverage:
	$(require-go-service-modules)
	$(prepare-go-workspace)
	@mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go/services
	@set -e; for pkg in $(GO_SERVICE_MODULE_DIRS); do \
		service=$${pkg#services/}; \
		safe=$$(echo "$$service" | tr '/.' '__'); \
		report_dir="$(ROOT_PATH)/storage/.cache/coverage/go/services/$$safe"; \
		mkdir -p "$$report_dir"; \
		echo "go test -race -coverprofile=$$report_dir/coverage.out ./... in $$pkg"; \
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
