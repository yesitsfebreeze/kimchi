.PHONY: build

MAKEFLAGS += --no-print-directory

all: release

release:
	@cd daemon && $(MAKE) -B release	
	@cd editor && $(MAKE) -B release

