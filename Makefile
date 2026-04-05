ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_RUN := docker compose -f $(GO_FMT_COMPOSE_FILE) run --rm $(GO_FMT_SERVICE)
PACKAGE_FMT := pnpm fmt
MARKDOWN_FILES := $(shell git ls-files '*.md')

.PHONY: format typecheck test coverage build clean demo

format:
	$(PACKAGE_FMT)
	$(GO_FMT_RUN) format .
	@if [ -n "$(MARKDOWN_FILES)" ]; then \
		pnpm exec oxfmt --ignore-path .gitignore $(MARKDOWN_FILES); \
	fi

typecheck:
	pnpm typecheck

test:
	pnpm test

coverage:
	pnpm test:coverage

build:
	pnpm build

demo:
	cd $(ROOT_PATH)/demo && node scripts/build-assets.mjs
	cd $(ROOT_PATH)/demo && go run ./public

clean:
	rm -rf $(ROOT_PATH)/packages/storage/.cache
	rm -rf $(ROOT_PATH)/packages/storage/.pnpm-store
	rm -rf $(ROOT_PATH)/packages/storage/.turbo
	rm -rf $(ROOT_PATH)/packages/storage/dist
	rm -rf $(ROOT_PATH)/packages/storage/coverage
	rm -rf $(ROOT_PATH)/demo/public/build
	mkdir -p $(ROOT_PATH)/packages/storage/.cache
	mkdir -p $(ROOT_PATH)/packages/storage/.pnpm-store
	mkdir -p $(ROOT_PATH)/packages/storage/.turbo
	mkdir -p $(ROOT_PATH)/packages/storage/dist
	mkdir -p $(ROOT_PATH)/packages/storage/coverage/go
	mkdir -p $(ROOT_PATH)/packages/storage/coverage/playwright
	mkdir -p $(ROOT_PATH)/demo/public/build
	touch $(ROOT_PATH)/packages/storage/.cache/.gitkeep
	touch $(ROOT_PATH)/packages/storage/.turbo/.gitkeep
	touch $(ROOT_PATH)/packages/storage/dist/.gitkeep
	touch $(ROOT_PATH)/packages/storage/coverage/.gitkeep
	touch $(ROOT_PATH)/packages/storage/coverage/go/.gitkeep
	touch $(ROOT_PATH)/packages/storage/coverage/playwright/.gitkeep
