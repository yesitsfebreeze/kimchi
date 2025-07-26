.PHONY: help build run clean install

MAKEFLAGS += --no-print-directory
NAME = kitsune
BIN = kitsune
OUT := ./bin/$(BIN)
GO_VERSION = 1.24.5
GO_TARBALL = go$(GO_VERSION).linux-amd64.tar.gz
GO_URL = https://go.dev/dl/$(GO_TARBALL)
GO_INSTALL_DIR = /usr/local/go

RELEASE_FLAGS = -s -w
DEBUG_FLAGS = -N -l

help: # Show this help
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

release:
	@make install -B
	@go mod tidy
	@go build -C src -o ../$(OUT) -ldflags="$(RELEASE_FLAGS)"
	@upx --best --lzma ./$(OUT)
	@cp ./$(OUT) ../$(OUT)


build: # builds the executeable
	@make install -B
	@go mod tidy
	@go build -C src -o ../$(OUT) -gcflags=all="$(DEBUG_FLAGS)"

launch: # [args] -> builds and runs the BIN
	@go mod tidy
	@go build -C src -o ../$(OUT) -ldflags="$(RELEASE_FLAGS)"
	@./$(OUT) $(filter-out $@,$(MAKECMDGOALS))

run: # [args] runs the BIN via go run
	@go run -C . ./src $(filter-out $@,$(MAKECMDGOALS))

clean: # removes the BIN and any temporary files
	@rm -f $(OUT)
	@rm -rf *~

install: # installs go and dependencies
	@bash ./scripts/install_go.sh $(GO_VERSION)
