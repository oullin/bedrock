ROOT_PATH := $(shell pwd)
GO_FMT_COMPOSE_FILE := go-fmt.compose.yaml
GO_FMT_SERVICE := go-fmt
GO_FMT_RUN := docker compose -f $(GO_FMT_COMPOSE_FILE) run --rm $(GO_FMT_SERVICE)
TURBO_FMT := node ./scripts/run-turbo.mjs fmt
JS_FMT_FILTERS := \
	--filter=@gollin/billing \
	--filter=@gollin/files \
	--filter=@gollin/media \
	--filter=@gollin/notification \
	--filter=@gollin/queue
MARKDOWN_FILES := $(shell git ls-files '*.md')

.PHONY: format typecheck test coverage build clean

format:
	$(TURBO_FMT) $(JS_FMT_FILTERS)
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

clean:
	rm -rf $(ROOT_PATH)/storage/.cache
	rm -rf $(ROOT_PATH)/storage/.turbo
	rm -rf $(ROOT_PATH)/storage/dist
	rm -rf $(ROOT_PATH)/storage/coverage
	find $(ROOT_PATH)/packages -maxdepth 2 -type d \( -name dist -o -name coverage -o -name .turbo \) -prune -exec rm -rf {} +
	mkdir -p $(ROOT_PATH)/storage/.cache
	mkdir -p $(ROOT_PATH)/storage/.turbo
	mkdir -p $(ROOT_PATH)/storage/dist
	mkdir -p $(ROOT_PATH)/storage/coverage/packages
	mkdir -p $(ROOT_PATH)/storage/coverage/go
	mkdir -p $(ROOT_PATH)/storage/coverage/playwright
	touch $(ROOT_PATH)/storage/.cache/.gitkeep
	touch $(ROOT_PATH)/storage/.turbo/.gitkeep
	touch $(ROOT_PATH)/storage/dist/.gitkeep
	touch $(ROOT_PATH)/storage/coverage/.gitkeep
	touch $(ROOT_PATH)/storage/coverage/packages/.gitkeep
	touch $(ROOT_PATH)/storage/coverage/go/.gitkeep
	touch $(ROOT_PATH)/storage/coverage/playwright/.gitkeep
