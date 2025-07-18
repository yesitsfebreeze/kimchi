.PHONY: help build run clean install

help: # Show this help
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

OUT := bin/kim

GO_VERSION = 1.24.5
GO_TARBALL = go$(GO_VERSION).linux-amd64.tar.gz
GO_URL = https://go.dev/dl/$(GO_TARBALL)
GO_INSTALL_DIR = /usr/local/go


all: build

build: # builds the executeable
	@go mod tidy
	@go build -C kimchi -o ../$(OUT)

run: # runs the executable (make run -- ...args)
	@make build -B
	@./$(OUT) $(filter-out $@,$(MAKECMDGOALS))

clean: # removes the executable and any temporary files
	@rm -f $(OUT)
	@rm -rf *~

install: # installs go and dependencies
	@bash ./dev/install_go.sh $(GO_VERSION)
