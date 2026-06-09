.DEFAULT_GOAL := gogen

# Compatibility wrapper. The canonical task runner is justfile.

gogen:
	just gogen

tidy:
	just tidy

test:
	just test

check:
	just check

lint:
	just lint

clean:
	just clean

proto:
	just proto

help:
	just --list

.PHONY: gogen tidy test check lint clean proto help
