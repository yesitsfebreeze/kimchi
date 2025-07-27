.PHONY: help grammars build run clean install

MAKEFLAGS += --no-print-directory
BIN = kitsuned
OUT := ../bin/$(BIN)

RELEASE_FLAGS = -Zlocation-detail=none

help: # Show this help
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

release:
	@killall kitsuned || true
	@RUSTFLAGS="$(RELEASE_FLAGS)" cargo +nightly build --release
	@mkdir -p ../bin
	@upx --best --lzma ./bin/release/$(BIN)
	@cp ./bin/release/$(BIN) $(OUT)

build: # builds the executeable
	@echo "bridge build"

grammars:
	@./scripts/grammars.sh

launch: # [args] -> builds and runs the BIN
	@killall kitsuned || true
	@RUSTFLAGS="$(RELEASE_FLAGS)" cargo +nightly build --release
	@mkdir -p ../bin
	@upx --best --lzma ./bin/release/$(BIN)
	@cp ./bin/release/$(BIN) $(OUT)
	@$(OUT)

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
	@echo '{"Highlight":{"lang":"cpp","path":"/mnt/d/dev/kitsune/editor/src/index.js"}}' | socat - UNIX-CONNECT:/tmp/kitsuned.sock
	@echo '{"Connect":{"user":"feb"}}' | socat - UNIX-CONNECT:/tmp/kitsuned.sock
