ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_RUN := docker compose -f $(GO_FMT_COMPOSE_FILE) run --rm $(GO_FMT_SERVICE)
PACKAGE_FMT := pnpm fmt
MARKDOWN_FILES := $(shell git ls-files '*.md')

GO_PACKAGES := \
	packages/auth \
	packages/bus \
	packages/cache \
	packages/contracts/auth \
	packages/cookie \
	packages/fortify \
	packages/jetstream \
	packages/queue \
	packages/routing \
	packages/session \
	packages/spark \
	dist/app

.PHONY: format vet tidy typecheck test coverage build clean demo

format:
	$(PACKAGE_FMT)
	@for pkg in $(GO_PACKAGES); do \
		echo "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go vet ./...; \
	done
	@for pkg in $(GO_PACKAGES); do \
		echo "go-fmt format in $$pkg"; \
		$(GO_FMT_RUN) format --host-path $(ROOT_PATH)/$$pkg; \
	done
	@if [ -n "$(MARKDOWN_FILES)" ]; then \
		pnpm exec oxfmt --ignore-path .gitignore $(MARKDOWN_FILES); \
	fi

vet:
	@for pkg in $(GO_PACKAGES); do \
		echo "go vet ./... in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go vet ./...; \
	done

tidy:
	@for pkg in $(GO_PACKAGES); do \
		echo "go mod tidy in $$pkg"; \
		cd $(ROOT_PATH)/$$pkg && go mod tidy; \
	done

typecheck:
	pnpm typecheck

test:
	pnpm test

coverage:
	pnpm test:coverage

build:
	pnpm build

demo:
	cd $(ROOT_PATH)/dist/app && node scripts/build-assets.mjs
	cd $(ROOT_PATH)/dist/app && go run ./public

clean:
	rm -rf $(ROOT_PATH)/storage/.cache
	rm -rf $(ROOT_PATH)/storage/.turbo
	rm -rf $(ROOT_PATH)/dist/app/public/build
	rm -rf $(ROOT_PATH)/dist/bin
	mkdir -p $(ROOT_PATH)/storage/.cache/.pnpm-store
	mkdir -p $(ROOT_PATH)/storage/.cache/coverage/go
	mkdir -p $(ROOT_PATH)/storage/.cache/coverage/playwright
	mkdir -p $(ROOT_PATH)/storage/.turbo
	mkdir -p $(ROOT_PATH)/dist/bin
	mkdir -p $(ROOT_PATH)/dist/app/public/build
	touch $(ROOT_PATH)/storage/.cache/.gitkeep
	touch $(ROOT_PATH)/storage/.turbo/.gitkeep
