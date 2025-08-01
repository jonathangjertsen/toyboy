package model

import (
	"fmt"
	"os"
)

func LoadROM(
	filename string,
	mem []Data8,
	cart *Cartridge,
	bootROMLock *BootROMLock,
) error {
	// Read file
	rom, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Check ROM size
	if len(rom)%ROMBankSize != 0 {
		return fmt.Errorf("ROM size %d is not a multiple of %d", len(rom), ROMBankSize)
	}

	// Load into ROM banks
	for i := range len(rom) / ROMBankSize {
		copy(cart.ROM[i][:], Data8Slice(rom[i*ROMBankSize:(i+1)*ROMBankSize]))
	}

	// If BootROM is done, map in Bank 0
	// If BootROM isn't done yet, don't overwrite that
	if bootROMLock.BootOff {
		copy(mem[:AddrCartridgeBank0End], cart.ROM[0][:])
	} else {
		copy(mem[SizeBootROM:AddrCartridgeBank0End], cart.ROM[0][SizeBootROM:])
	}

	// Map in initial Bank 1
	cart.SetROMBank(mem, 1)

	// Configure cartridge MCB features
	cart.MBCFeatures = GetMBCFeatures(
		rom[AddrCartridgeType],
		rom[AddrROMSize],
		rom[AddrRAMSize],
	)

	return nil
}
