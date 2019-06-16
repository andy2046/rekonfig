.PHONY: all help test build

all: help

help:				## Show this help
	@scripts/help.sh

test:				## Test potential bugs and race conditions
	@scripts/test.sh

build:				## Build image
	@scripts/build.sh
