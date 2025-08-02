package model

import (
	"fmt"
	"os"
)

func LoadROM(
	filename string,
	gb *Gameboy,
) error {
	// Read file
	rom, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Printf("LEN=%d\n", len(gb.Cartridge.ROM))
	// Check ROM size
	if len(rom)%ROMBankSize != 0 {
		return fmt.Errorf("ROM size %d is not a multiple of %d", len(rom), ROMBankSize)
	}

	// Load into ROM banks
	for i := range len(rom) / ROMBankSize {
		fmt.Printf("Creating ROM Bank %d\n", i)
		copy(gb.Cartridge.ROM[i][:], Data8Slice(rom[i*ROMBankSize:(i+1)*ROMBankSize]))
	}

	// If BootROM is done, map in Bank 0
	// If BootROM isn't done yet, don't overwrite that
	if gb.BootROMLock.BootOff {
		copy(gb.Mem[:AddrCartridgeBank0End], gb.Cartridge.ROM[0][:])
	} else {
		copy(gb.Mem[SizeBootROM:AddrCartridgeBank0End], gb.Cartridge.ROM[0][SizeBootROM:])
	}

	// Configure cartridge MCB features
	gb.Cartridge.MBCFeatures = GetMBCFeatures(
		rom[AddrCartridgeType],
		rom[AddrROMSize],
		rom[AddrRAMSize],
	)
	if r := gb.Cartridge.MBCFeatures.TotalROMSize(); r != len(rom) {
		gb.PrintCartridgeInfo(os.Stdout)
		panicf("expected ROM size=%d got %d (header=%v)", len(rom), r, rom[AddrCartridgeType:AddrRAMSize+1])
	}

	// Map in initial Bank 1
	gb.SetROMBank1(1)

	return nil
}
