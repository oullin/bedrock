ROOT_PATH := $(shell pwd)

.PHONY: format typecheck test coverage build clean

format:
	pnpm fmt

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
