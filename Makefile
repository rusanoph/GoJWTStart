.PHONY: build, run
build:
	go build main.go

run:
	go run main.go


.DEFAULT_GOAL := build