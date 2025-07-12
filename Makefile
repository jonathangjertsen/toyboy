all: emulator

.PHONY: run

emulator:
	go generate ./...
	go build -o bin/emulator github.com/jonathangjertsen/gameboy/cmd/emulator

run: emulator
	bin/emulator
