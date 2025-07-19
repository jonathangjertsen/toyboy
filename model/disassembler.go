package model

import (
	"fmt"
	"io"
)

type Disassembler struct {
	Program []uint8
	Decoded []DisInstruction
	Trace   bool
	PC      uint16

	stack             []uint16
	stackIdx          int
	cachedDisassembly *Disassembly
}

type DisInstruction struct {
	Raw     [3]uint8
	Size    uint16
	Address uint16
	Opcode  Opcode
}

func (di *DisInstruction) Asm() string {
	str := di.Opcode.String()
	ln := len(str)
	switch di.Opcode {
	default:
	case OpcodeLDAnn:
		return fmt.Sprintf("%s %s, [$%04xh]", str[:ln-3], str[ln-3:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeLDBCnn, OpcodeLDDEnn, OpcodeLDHLnn, OpcodeLDSPnn:
		return fmt.Sprintf("%s %s, [$%04xh]", str[:ln-4], str[ln-4:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeLDnnA:
		return fmt.Sprintf("LD [$%04xh], A", join16(di.Raw[2], di.Raw[1]))
	case OpcodeLDnnSP:
		return fmt.Sprintf("LD [$%04xh], SP", join16(di.Raw[2], di.Raw[1]))
	case OpcodeLDHLAInc:
		return "LD (HL+), A"
	case OpcodeLDHLADec:
		return "LD (HL-), A"
	case OpcodeLDAHLInc:
		return "LD A, (HL+)"
	case OpcodeLDAHLDec:
		return "LD A, (HL-)"
	case OpcodeDECHLInd:
		return "DEC (HL)"
	case OpcodeINCHLInd:
		return "DEC (HL)"
	case OpcodeRST0x00, OpcodeRST0x08, OpcodeRST0x10, OpcodeRST0x18, OpcodeRST0x20, OpcodeRST0x28, OpcodeRST0x30, OpcodeRST0x38:
		return fmt.Sprintf("RST $00%sh", str[ln-2:])
	case OpcodeUndefD3, OpcodeUndefDB, OpcodeUndefDD, OpcodeUndefE3, OpcodeUndefE4, OpcodeUndefEB, OpcodeUndefEC, OpcodeUndefED, OpcodeUndefF4, OpcodeUndefFC, OpcodeUndefFD:
		return fmt.Sprintf("Undefined instruction 0x%02x", di.Raw[0])
	case OpcodeLDBCA, OpcodeLDDEA:
		return fmt.Sprintf("LD (%s), A", str[ln-3:ln-1])
	case OpcodeLDHLn:
		return fmt.Sprintf("LD (HL), $%02xh", di.Raw[1])
	case OpcodeRET, OpcodeNop, OpcodeRLA, OpcodeRLCA, OpcodeRRA, OpcodeRRCA, OpcodeDAA, OpcodeDI, OpcodeEI, OpcodeSTOP, OpcodeCCF, OpcodeRETI:
		return str
	case OpcodeRETZ, OpcodeRETC:
		return fmt.Sprintf("%s %s", str[:ln-1], str[ln-1:])
	case OpcodeRETNZ, OpcodeRETNC:
		return fmt.Sprintf("%s %s", str[:ln-2], str[ln-2:])
	case OpcodePUSHBC, OpcodePUSHDE, OpcodePUSHHL, OpcodePUSHAF, OpcodePOPBC, OpcodePOPDE, OpcodePOPHL, OpcodePOPAF:
		return fmt.Sprintf("%s %s", str[:ln-2], str[ln-2:])
	case OpcodeCPLaka2f:
		return "CPL"
	case OpcodeXORn, OpcodeADDn, OpcodeANDn, OpcodeORn, OpcodeADCn, OpcodeSBCn, OpcodeCPn, OpcodeSUBn:
		return fmt.Sprintf("%s A, $%02xh", str[:ln-1], di.Raw[1])
	case OpcodeLDAA, OpcodeLDAB, OpcodeLDAC, OpcodeLDAD, OpcodeLDAE, OpcodeLDAH, OpcodeLDAL,
		OpcodeLDBA, OpcodeLDBB, OpcodeLDBC, OpcodeLDBD, OpcodeLDBE, OpcodeLDBH, OpcodeLDBL,
		OpcodeLDCA, OpcodeLDCB, OpcodeLDCC, OpcodeLDCD, OpcodeLDCE, OpcodeLDCH, OpcodeLDCL,
		OpcodeLDDA, OpcodeLDDB, OpcodeLDDC, OpcodeLDDD, OpcodeLDDE, OpcodeLDDH, OpcodeLDDL,
		OpcodeLDEA, OpcodeLDEB, OpcodeLDEC, OpcodeLDED, OpcodeLDEE, OpcodeLDEH, OpcodeLDEL,
		OpcodeLDHA, OpcodeLDHB, OpcodeLDHC, OpcodeLDHD, OpcodeLDHE, OpcodeLDHH, OpcodeLDHL,
		OpcodeLDLA, OpcodeLDLB, OpcodeLDLC, OpcodeLDLD, OpcodeLDLE, OpcodeLDLH, OpcodeLDLL:
		return fmt.Sprintf("%s %s, %s", str[:ln-2], str[ln-2:ln-1], str[ln-1:])
	case OpcodeXORA, OpcodeXORB, OpcodeXORC, OpcodeXORD, OpcodeXORE, OpcodeXORH, OpcodeXORL,
		OpcodeORA, OpcodeORB, OpcodeORC, OpcodeORD, OpcodeORE, OpcodeORH, OpcodeORL,
		OpcodeANDA, OpcodeANDB, OpcodeANDC, OpcodeANDD, OpcodeANDE, OpcodeANDH, OpcodeANDL,
		OpcodeADDA, OpcodeADDB, OpcodeADDC, OpcodeADDD, OpcodeADDE, OpcodeADDH, OpcodeADDL,
		OpcodeADCA, OpcodeADCB, OpcodeADCC, OpcodeADCD, OpcodeADCE, OpcodeADCH, OpcodeADCL,
		OpcodeSUBA, OpcodeSUBB, OpcodeSUBC, OpcodeSUBD, OpcodeSUBE, OpcodeSUBH, OpcodeSUBL,
		OpcodeSBCA, OpcodeSBCB, OpcodeSBCC, OpcodeSBCD, OpcodeSBCE, OpcodeSBCH, OpcodeSBCL,
		OpcodeCPA, OpcodeCPB, OpcodeCPC, OpcodeCPD, OpcodeCPE, OpcodeCPH, OpcodeCPL:
		return fmt.Sprintf("%s A, %s", str[:ln-1], str[ln-1:])
	case OpcodeADDHL, OpcodeSUBHL, OpcodeANDHL, OpcodeORHL, OpcodeADCHL, OpcodeSBCHL, OpcodeXORHL, OpcodeCPHL:
		return fmt.Sprintf("%s A, (HL)", str[:ln-2])
	case OpcodeLDABC, OpcodeLDADE,
		OpcodeLDAHL, OpcodeLDBHL, OpcodeLDCHL, OpcodeLDDHL, OpcodeLDEHL, OpcodeLDHHL, OpcodeLDLHL:
		return fmt.Sprintf("%s %s, (HL)", str[:ln-3], str[ln-3:ln-2])
	case OpcodeLDHLA, OpcodeLDHLB, OpcodeLDHLC, OpcodeLDHLD, OpcodeLDHLE, OpcodeLDHLH, OpcodeLDHLL:
		return fmt.Sprintf("%s (HL), %s", str[:ln-3], str[ln-1:])
	case OpcodeDECA, OpcodeDECB, OpcodeDECC, OpcodeDECD, OpcodeDECE, OpcodeDECH, OpcodeDECL,
		OpcodeINCA, OpcodeINCB, OpcodeINCC, OpcodeINCD, OpcodeINCE, OpcodeINCH, OpcodeINCL:
		return fmt.Sprintf("%s %s", str[:ln-1], str[ln-1:])
	case OpcodeINCBC, OpcodeINCDE, OpcodeINCHL, OpcodeINCSP, OpcodeDECBC, OpcodeDECDE, OpcodeDECHL, OpcodeDECSP:
		return fmt.Sprintf("%s %s", str[:ln-2], str[ln-2:])
	case OpcodeJPnn, OpcodeCALLnn:
		return fmt.Sprintf("%s $%04xh", str[:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeJPCnn, OpcodeJPZnn:
		return fmt.Sprintf("%s %s, $%04xh", str[:ln-3], str[ln-3:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeJPNZnn, OpcodeJPNCnn:
		return fmt.Sprintf("%s %s, $%04xh", str[:ln-4], str[ln-4:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeCALLZnn, OpcodeCALLCnn:
		return fmt.Sprintf("CALL %s, $%04xh", str[ln-3:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeCALLNZnn, OpcodeCALLNCnn:
		return fmt.Sprintf("CALL %s, $%04xh", str[ln-4:ln-2], join16(di.Raw[2], di.Raw[1]))
	case OpcodeJRNZe, OpcodeJRNCe:
		return fmt.Sprintf("%s %s, PC+$%02xh", str[:ln-3], str[ln-3:ln-1], int8(di.Raw[1]))
	case OpcodeJRZe, OpcodeJRCe:
		return fmt.Sprintf("%s %s, PC+$%02xh", str[:ln-2], str[ln-2:ln-1], int8(di.Raw[1]))
	case OpcodeJRe:
		return fmt.Sprintf("%s PC+$%02xh", str[:ln-1], int8(di.Raw[1]))
	case OpcodeLDAn, OpcodeLDBn, OpcodeLDCn, OpcodeLDDn, OpcodeLDEn, OpcodeLDHn, OpcodeLDLn:
		return fmt.Sprintf("%s %s, $%02xh", str[:ln-2], str[ln-2:ln-1], join16(di.Raw[2], di.Raw[1]))
	case OpcodeADDHLHL, OpcodeADDHLDE, OpcodeADDHLBC, OpcodeADDHLSP:
		return fmt.Sprintf("%s %s, %s", str[:ln-4], str[ln-4:ln-2], str[ln-2:])
	case OpcodeLDHnA:
		return fmt.Sprintf("LDH ($%02xh),A; eff. LD $%04x, A", di.Raw[1], 0xff00+int(di.Raw[1]))
	case OpcodeLDHAn:
		return fmt.Sprintf("LDH A,($%02xh); eff. LD A,$%04x", di.Raw[1], 0xff00+int(di.Raw[1]))
	case OpcodeLDHAC:
		return "LDH A,(C); eff. LD A, [C+$0xff00]"
	case OpcodeLDHCA:
		return "LDH (C), A; eff. LD [C+$0xff00], a"
	case OpcodeLDHLSPe:
		if di.Raw[1]&0x80 == 0 {
			return fmt.Sprintf("LD HL,SP+$%02xh", di.Raw[1])
		} else {
			return fmt.Sprintf("LD HL,SP-$%02xh", -int8(di.Raw[1]))
		}
	case OpcodeCB:
		cbop := CBOp{Op: cb((di.Raw[0] & 0xf8) >> 3), Target: CBTarget(di.Raw[0] & 0x7)}
		return fmt.Sprintf("%s %s", cbop.Op, cbop.Target)
	case OpcodeJPHL:
		return "JP HL"
	}
	panicf("CAN'T FMT %v\n", di.Opcode)
	return ""
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
	PC   uint16
	Code []CodeSection
	Data []DataSection
}

func (d *Disassembly) Print(w io.Writer) {
	for _, section := range d.Code {
		fmt.Fprintf(w, "\nCode section at 0x%04x\n", section.Address())
		for _, inst := range section.Instructions {
			if inst.Address == d.PC {
				fmt.Fprintf(w, "[%04x]->%s\n", inst.Address, inst.Asm())
			} else {
				fmt.Fprintf(w, "%04xh | %s\n", inst.Address, inst.Asm())
			}
		}
	}
	data := splitSections(d.Data)
	prevEndAddr := uint16(0xffff)
	for _, section := range data {
		if prevEndAddr != section.Address {
			fmt.Fprintf(w, "\nData section at 0x%04x\n", section.Address)
		}
		prevEndAddr = section.Address + uint16(len(section.Raw))
		allEqual := true
		testByte := section.Raw[0]
		for _, b := range section.Raw {
			if b != testByte {
				allEqual = false
				break
			}
		}
		if allEqual {
			fmt.Fprintf(w, "0x%0x bytes of 0x%02x\n", len(section.Raw), testByte)
			continue
		}

		i := 0
		for line := range (len(section.Raw) + 15) / 16 {
			fmt.Fprintf(w, "%04xh | ", int(section.Address)+line*16)
			for range 16 {
				if i >= len(section.Raw) {
					break
				}
				fmt.Fprintf(w, "%02x ", section.Raw[i])
				i++
			}
			fmt.Fprintf(w, "\n")
		}
	}
}
func splitSections(sections []DataSection) []DataSection {
	var result []DataSection
	for _, section := range sections {
		raw := section.Raw
		start := 0

		for start < len(raw) {
			runStart := start
			runByte := raw[start]
			runLen := 1

			// Find run of same byte
			for i := start + 1; i < len(raw) && raw[i] == runByte; i++ {
				runLen++
			}

			if runLen >= 32 {
				// Add preceding non-uniform part
				if runStart > 0 {
					result = append(result, DataSection{
						Address: section.Address,
						Raw:     raw[0:runStart],
					})
				}
				// Add the uniform run
				result = append(result, DataSection{
					Address: section.Address + uint16(runStart),
					Raw:     raw[runStart : runStart+runLen],
				})
				// Continue after the run
				raw = raw[runStart+runLen:]
				section.Address += uint16(runStart + runLen)
				start = 0
			} else {
				start++
			}
		}

		// Add any remaining non-uniform tail
		if len(raw) > 0 {
			result = append(result, DataSection{
				Address: section.Address,
				Raw:     raw,
			})
		}
	}
	return result
}

func NewDisassembler() *Disassembler {
	return &Disassembler{}
}

func (dis *Disassembler) SetProgram(program []uint8) {
	dis.Program = program
	dis.Decoded = make([]DisInstruction, len(program))
	dis.cachedDisassembly = nil
	dis.printf("setProgram with len=%d", len(program))
}

func (dis *Disassembler) printf(format string, args ...any) {
	if !dis.Trace {
		return
	}
	for range dis.stackIdx {
		fmt.Printf("  ")
	}
	if dis.stackIdx > 0 {
		fmt.Printf("%04xh ", dis.stack[dis.stackIdx-1])
	}
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}

func (dis *Disassembler) SetPC(address uint16) {
	dis.printf("SetPC %04xh", address)
	dis.stack = append(dis.stack, address)
	dis.stackIdx++
	defer func() {
		dis.stackIdx--
		dis.stack = dis.stack[:dis.stackIdx]
		dis.PC = address
	}()

	for {
		if int(address) >= len(dis.Program) {
			dis.printf("reached address outside of program (len=%0x)", len(dis.Program))
			return
		}

		if dis.Decoded[address].Size > 0 {
			dis.printf("already seen this address")
			return
		}

		di, err := dis.readNewInstruction(address)
		if err != nil {
			panicf("failed decoding at %x: %v", address, err)
			return
		}

		for i := range di.Size {
			if dis.Decoded[address+i].Size > 0 {
				dis.printf("at offset %d there is an existing instruction, returning", i)
				return
			}
		}

		dis.insert(di)
		address += di.Size
		dis.stack[dis.stackIdx-1] = address

		if dis.checkBranches(di) {
			return
		}
	}
}

func (dis *Disassembler) Disassembly() *Disassembly {
	if dis.cachedDisassembly != nil {
		dis.cachedDisassembly.PC = dis.PC
		return dis.cachedDisassembly
	}

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
	dis.cachedDisassembly = &out
	return &out
}

func (dis *Disassembler) insert(di DisInstruction) {
	dis.printf("insert %v at %04xh:%04xh", di.Opcode, di.Address, di.Address+di.Size-1)
	for addr := di.Address; addr != di.Address+di.Size; addr++ {
		dis.Decoded[addr] = di
	}
	dis.cachedDisassembly = nil
}

// returns true if unconditional branch, i.e. can never fall through
func (dis *Disassembler) checkBranches(di DisInstruction) bool {
	switch di.Opcode {
	case OpcodeJRNZe, OpcodeJRZe, OpcodeJRe:
		e := int8(di.Raw[1])

		// branch taken
		dis.printf("checking branch-taken for relative jump")
		if e > 0 {
			dis.SetPC(di.Address + di.Size + uint16(e))
		} else {
			dis.SetPC(di.Address + di.Size - uint16(-e))
		}

		// branch not taken will be inspected after this
		if di.Opcode == OpcodeJRe {
			dis.printf("unconditional branch, can never fall through")
			return true
		}
	case OpcodeJPnn, OpcodeJPCnn, OpcodeJPNCnn, OpcodeJPZnn, OpcodeJPNZnn, OpcodeCALLnn:
		addr := join16(di.Raw[2], di.Raw[1])
		dis.printf("checking branch-taken for absolute jump")
		dis.SetPC(addr)
		if di.Opcode == OpcodeJPnn || di.Opcode == OpcodeJPHL {
			dis.printf("unconditional branch, can never fall through")
			return true
		} else if di.Opcode == OpcodeCALLnn {
			dis.printf("will probably return here")
			// We don't fall thru here, but structured programs would return here from RET.
			// Unless the code adjusts the return address in stack memory.
			// If that happens AND there is data just after this section, it will be interpreted as code
		}
	case OpcodeRST0x00:
		dis.doRST(0x00)
	case OpcodeRST0x08:
		dis.doRST(0x08)
	case OpcodeRST0x10:
		dis.doRST(0x10)
	case OpcodeRST0x18:
		dis.doRST(0x18)
	case OpcodeRST0x20:
		dis.doRST(0x20)
	case OpcodeRST0x28:
		dis.doRST(0x28)
	case OpcodeRST0x30:
		dis.doRST(0x30)
	case OpcodeRST0x38:
		dis.doRST(0x38)
	case OpcodeUndefD3, OpcodeUndefDB, OpcodeUndefDD, OpcodeUndefE3, OpcodeUndefE4, OpcodeUndefEB, OpcodeUndefEC, OpcodeUndefED, OpcodeUndefF4, OpcodeUndefFC, OpcodeUndefFD:
		dis.printf("dropping undefined instruction %02x", di.Opcode)
		dis.Decoded[di.Address] = DisInstruction{}
		return true
	case OpcodeRET, OpcodeRETI, OpcodeJPHL:
		return true
	}

	return false
}

func (dis *Disassembler) doRST(vec uint16) {
	dis.printf("unconditional function call to %04x", vec)
	dis.SetPC(vec)
	dis.printf("will probably return here")
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
