all: emulator rom-hello-world rom-empty

.PHONY: run

emulator:
	go generate ./...
	go build -o bin/emulator github.com/jonathangjertsen/toyboy/cmd/emulator

rom-empty:
	cd assets/cartridges/asm;\
	rgbasm -o empty.o empty.asm;\
	rgblink -o ../empty.gb empty.o;\
	rgbfix -v -p 0xFF ../empty.gb

rom-hello-world:
	cd assets/cartridges/asm;\
	rgbasm -o hello-world.o hello-world.asm;\
	rgblink -o ../hello-world.gb hello-world.o;\
	rgbfix -v -p 0xFF ../hello-world.gb

install-gio:
	go install gioui.org/cmd/gogio@latest
	sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

run: emulator rom-empty rom-hello-world
	APP_ENV=development bin/emulator
