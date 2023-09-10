.PHONY: build, run
build:
	go build ./cmd/gojwtstart

run:
	go run ./cmd/gojwtstart/main.go


.DEFAULT_GOAL := build