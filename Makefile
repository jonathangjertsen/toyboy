all: emulator roms disassembler

roms: rom-hello-world rom-empty rom-unbricked

.PHONY: run

gen:
	go generate ./...

disassembler: gen
	go build -o bin/disassembler github.com/jonathangjertsen/toyboy/cmd/disassembler

emulator: gen
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

rom-unbricked:
	cd assets/cartridges/asm;\
	rgbasm -o unbricked.o unbricked.asm;\
	rgblink -o ../unbricked.gb -m unbricked.map -n unbricked.sym unbricked.o;\
	rgbfix -v -p 0xFF ../unbricked.gb

install-gio:
	go install gioui.org/cmd/gogio@latest
	sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev

run: emulator rom-empty rom-hello-world rom-unbricked
	APP_ENV=development bin/emulator

run-dis: disassembler
	bin/disassembler
