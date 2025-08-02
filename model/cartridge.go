package model

import (
	"fmt"
	"io"
)

type Cartridge struct {
	ROM              [][ROMBankSize]Data8
	RAM              [][RAMBankSize]Data8
	MBCFeatures      MBCFeatures
	ExtRAMEnabled    bool
	RegLow           Data8
	RegHigh          Data8
	RegSelect        Data8
	SelectedRAMBank  Data8
	SelectedROMBank0 Data8
	SelectedROMBank1 Data8
}

func (gb *Gameboy) SetROMBank0(which Data8) {
	if which >= Data8(gb.Cartridge.MBCFeatures.NROMBanks) {
		return
	}

	fmt.Printf("Switched bank 0\n")
	copy(gb.Mem[AddrZero:AddrCartridgeBank0End], gb.Cartridge.ROM[which][:])
	gb.Cartridge.SelectedROMBank0 = which
}

func (gb *Gameboy) SetROMBank1(which Data8) {
	if which >= Data8(gb.Cartridge.MBCFeatures.NROMBanks) {
		return
	}
	fmt.Printf("Switched bank 1\n")
	copy(gb.Mem[AddrCartridgeBankNBegin:AddrCartridgeBankNEnd], gb.Cartridge.ROM[which][:])
	gb.Cartridge.SelectedROMBank1 = which
}

func (gb *Gameboy) SetRAMBank(which Data8) {
	if which >= Data8(gb.Cartridge.MBCFeatures.NRAMBanks) {
		return
	}
	fmt.Printf("Switched RAM bank\n")

	// Store current RAM contents to bank
	copy(gb.Cartridge.RAM[gb.Cartridge.SelectedRAMBank][:], gb.Mem[AddrCartridgeRAMBegin:AddrCartridgeRAMEnd])

	// Load from bank to RAM
	copy(gb.Mem[AddrCartridgeRAMBegin:AddrCartridgeRAMEnd], gb.Cartridge.RAM[which][:])
	gb.Cartridge.SelectedRAMBank = which
}

func (gb *Gameboy) WriteCartridge(addr Addr, v Data8) {
	switch gb.Cartridge.MBCFeatures.ID {
	case MBCIDNone:
		return
	case MBCID1:
		gb.writeCartridgeMBC1(addr, v)
	default:
		panic("not implemented MBC")
	}
}

func (gb *Gameboy) writeCartridgeMBC1(addr Addr, v Data8) {
	if addr <= 0x1fff {
		gb.Cartridge.ExtRAMEnabled = v&0x0f == 0x0a
	} else if addr <= 0x3fff {
		if v == 0x00 {
			v = 0x01
		}
		gb.Cartridge.RegLow = v
		gb.updateBankMBC1()
	} else if addr <= 0x5fff {
		gb.Cartridge.RegHigh = v
		gb.updateBankMBC1()
	} else if addr <= 0x7fff {
		gb.Cartridge.RegSelect = v
		gb.updateBankMBC1()
	}
}

func (gb *Gameboy) updateBankMBC1() {
	// Peeked from https://git.sr.ht/~dajolly/dmgl/tree/master/item/src/bus/memory/cartridge/mapper_1.c. Thanks!
	count := Data16(gb.Cartridge.MBCFeatures.NROMBanks)

	var ramBank, romBank0, romBank1 Data16
	high, low := Data16(gb.Cartridge.RegHigh), Data16(gb.Cartridge.RegLow)

	if count >= 64 {
		ramBank = 0
		if (gb.Cartridge.RegSelect & 1) != 0 {
			romBank0 = (high & 0x3) << 5
		} else {
			romBank0 = 0
		}
		romBank1 = ((high & 0x3) << 5) | (low & 0x1f)
	} else {
		if (gb.Cartridge.RegSelect & 1) != 0 {
			ramBank = high & 0x3
		} else {
			ramBank = 0
		}
		romBank0 = 0
		romBank1 = low & 0x1f
	}
	switch romBank1 {
	case 0, 32, 64, 96:
		romBank1++
	}
	romBank0 &= (count - 1)
	romBank1 &= (count - 1)

	count = Data16(gb.Cartridge.MBCFeatures.NRAMBanks)
	ramBank &= (count - 1)

	gb.SetRAMBank(Data8(ramBank))
	gb.SetROMBank0(Data8(romBank0))
	gb.SetROMBank1(Data8(romBank1))
}

func (gb *Gameboy) PrintCartridgeInfo(f io.Writer) {
	cart := &gb.Cartridge
	mbc := gb.Cartridge.MBCFeatures
	fmt.Fprintf(f, "Title: %s\n", gb.CartridgeTitle())
	fmt.Fprintf(f, "MBC: %s\n", mbc.ID.String())
	fmt.Fprintf(f, "ROM size: %d kB (%d banks)\n", mbc.TotalROMSize()/1024, mbc.NROMBanks)
	fmt.Fprintf(f, "RAM size: %d kB (%d banks)\n", mbc.TotalRAMSize()/1024, mbc.NRAMBanks)
	fmt.Fprintf(f, "Features: %s\n", mbc.Features())
	fmt.Fprintf(f, "HIGH=%s LOW=%s SEL=%s\n", cart.RegHigh.Hex(), cart.RegLow.Hex(), cart.RegSelect.Hex())

}

func (gb *Gameboy) CartridgeTitle() string {
	titleStart := 0x134
	titleEnd := titleStart + 16
	for i := range 16 {
		titleEnd = titleStart + i
		if gb.Mem[titleStart+i] == 0 {
			break
		}
	}
	return string(gb.Mem[titleStart:titleEnd])
}
