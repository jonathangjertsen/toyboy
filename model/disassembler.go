package model

import "fmt"

type Disassembler struct {
	Program []uint8
	Decoded []DisInstruction
	Trace   bool

	addr      uint16
	recursion int
}

type DisInstruction struct {
	Raw     [3]uint8
	Size    uint16
	Address uint16
	Opcode  Opcode
}

type DataSection struct {
	Raw     []uint8
	Address uint16
}

type CodeSection struct {
	Instructions []DisInstruction
}

func (cs CodeSection) Address() uint16 {
	return cs.Instructions[0].Address
}

type Disassembly struct {
	Code []CodeSection
	Data []DataSection
}

func NewDisassembler(program []uint8) *Disassembler {
	return &Disassembler{
		Program: program,
		Decoded: make([]DisInstruction, len(program)),
	}
}

func (dis *Disassembler) printf(format string, args ...any) {
	if !dis.Trace {
		return
	}
	for range dis.recursion {
		fmt.Printf("  ")
	}
	fmt.Printf(format, args...)
}

func (dis *Disassembler) SetPC(address uint16) {
	dis.recursion++
	defer func() {
		dis.recursion--
	}()

	for {
		if int(address) >= len(dis.Program) {
			dis.printf("reached address outside of program (addr=%04x, len=%0x)\n", address, len(dis.Program))
			return
		}

		if dis.Decoded[address].Size > 0 {
			dis.printf("already decoded %x\n", address)
			return // already disassembled
		}

		dis.addr = address
		di, err := dis.readNewInstruction(address)
		if err != nil {
			panicf("failed decoding at %x: %v\n", address, err)
			return // failed to decode at address
		}
		dis.add(di)
		address += di.Size

		if dis.checkBranches(di) {
			return
		}
	}
}

func (dis *Disassembler) Disassembly() Disassembly {
	var out Disassembly
	currCodeSection := CodeSection{}
	currDataSection := DataSection{}
	for addr := 0; addr < len(dis.Decoded); {
		di := dis.Decoded[addr]
		if di.Size > 0 {
			if currDataSection.Raw != nil {
				out.Data = append(out.Data, currDataSection)
			}

			currCodeSection.Instructions = append(currCodeSection.Instructions, di)
			addr += int(di.Size)

			currDataSection.Raw = nil
		} else {
			if currCodeSection.Instructions != nil {
				out.Code = append(out.Code, currCodeSection)
			}

			if currDataSection.Raw == nil {
				currDataSection.Address = uint16(addr)
			}
			currDataSection.Raw = append(currDataSection.Raw, dis.Program[addr])
			addr++
			currCodeSection.Instructions = nil
		}
	}
	if currCodeSection.Instructions != nil {
		out.Code = append(out.Code, currCodeSection)
	}
	if currDataSection.Raw != nil {
		out.Data = append(out.Data, currDataSection)
	}
	return out
}

func (dis *Disassembler) add(di DisInstruction) {
	dis.printf("added %v at %v\n", di.Opcode, di.Address)
	for addr := di.Address; addr != di.Address+di.Size; addr++ {
		dis.Decoded[addr] = di
	}
}

// returns true if unconditional branch, i.e. can never fall through
func (dis *Disassembler) checkBranches(di DisInstruction) bool {
	switch di.Opcode {
	case OpcodeJRNZe, OpcodeJRZe, OpcodeJRe:
		e := int8(di.Raw[1])

		// branch taken
		if e > 0 {
			dis.SetPC(di.Address + di.Size + uint16(e))
		} else {
			dis.SetPC(di.Address + di.Size - uint16(-e))
		}

		// branch not taken will be inspected after this
		if di.Opcode == OpcodeJRe {
			// unconditional branch, can never fall through
			return true
		}
	case OpcodeJPnn, OpcodeJPCnn, OpcodeJPNCnn, OpcodeJPZnn, OpcodeJPNZnn, OpcodeCALLnn:
		addr := join16(di.Raw[2], di.Raw[1])
		dis.SetPC(addr)
		if di.Opcode == OpcodeJPnn {
			// unconditional branch, can never fall through
			return true
		} else if di.Opcode == OpcodeCALLnn {
			// We don't fall thru here, but structured programs would return here from RET.
			// Unless the code adjusts the return address in stack memory.
			// If that happens AND there is data just after this section, it will be interpreted as code
		}
	case OpcodeRET:
		return true
	}

	return false
}

func (dis *Disassembler) readNewInstruction(addr uint16) (DisInstruction, error) {
	di := DisInstruction{
		Address: addr,
	}
	if int(addr) >= len(dis.Program) {
		return di, fmt.Errorf("address out of bounds")
	}
	di.Opcode = Opcode(dis.Program[addr])
	ok := di.Opcode.IsValid()
	if !ok {
		return di, fmt.Errorf("invalid opcode %v", di.Opcode)
	}
	di.Size = instSize[di.Opcode]
	if di.Size == 0 {
		return di, fmt.Errorf("no size set for opcode %v (0x%x)", di.Opcode, int(di.Opcode))
	}
	for i := range di.Size {
		di.Raw[i] = dis.Program[di.Address+i]
	}
	return di, nil
}
