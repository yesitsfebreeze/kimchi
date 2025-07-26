.PHONY: help build run clean install

MAKEFLAGS += --no-print-directory
BIN = kitsuned
OUT := ../bin/$(BIN)

RELEASE_FLAGS = -s -w
DEBUG_FLAGS = -N -l

help: # Show this help
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

release:
	@cargo build --release
	@mkdir -p ../bin
	@cp ./bin/release/$(BIN) $(OUT)

build: # builds the executeable
	@echo "bridge build"

launch: # [args] -> builds and runs the BIN
	@echo "bridge launch"

run: # [args] runs the BIN via go run
	@cargo run --release

clean: # removes the BIN and any temporary files
	@echo "bridge clean"

install: # installs go and dependencies
	@sudo apt update
	@sudo apt install -y build-essential pkg-config libssl-dev cmake
	@curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
	@. "$(HOME)/.cargo/env"

test:
	@echo '{"lang":"go","job":"ast","path":"/mnt/d/dev/kitsune/editor/src/buffer.go"}' | socat - UNIX-CONNECT:/tmp/kitd.sock
