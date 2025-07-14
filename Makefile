all: emulator testrom

.PHONY: run

emulator:
	go generate ./...
	go build -o bin/emulator github.com/jonathangjertsen/toyboy/cmd/emulator

testrom:
	cd assets/cartridges/asm;\
	rgbasm -o hello-world.o hello-world.asm;\
	rgblink -o ../hello-world.gb hello-world.o;\
	rgbfix -v -p 0xFF ../hello-world.gb

run: emulator
	bin/emulator
