package model

import (
	"fmt"
	"io"
)

type Disassembler struct {
	Config *ConfigDisassembler

	Program Block
	HRAM    Block
	WRAM    Block

	PC Addr

	stack             []Addr
	stackIdx          int
	cachedDisassembly *Disassembly
}

type Block struct {
	CanExplore bool
	Begin      Addr
	Source     []Data8
	Decoded    []DisInstruction
}

type DisInstruction struct {
	Raw     [3]Data8
	Address Addr
	Opcode  Opcode
	Visited bool
}

func (di *DisInstruction) Size() Size16 {
	if !di.Visited {
		return 0
	}
	return instSize[di.Opcode]
}

func (di *DisInstruction) Asm() string {
	str := di.Opcode.String()
	ln := len(str)
	switch di.Opcode {
	default:
	case OpcodeLDAnn:
		return fmt.Sprintf("LD A, [$%s]", join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeLDBCnn, OpcodeLDDEnn, OpcodeLDHLnn, OpcodeLDSPnn:
		return fmt.Sprintf("LD %s, [$%s]", str[ln-4:ln-2], join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeLDnnA:
		return fmt.Sprintf("LD [$%s], A", join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeLDnnSP:
		return fmt.Sprintf("LD [$%s], SP", join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeLDHLAInc:
		return "LD (HL+), A"
	case OpcodeLDHLADec:
		return "LD (HL-), A"
	case OpcodeLDAHLInc:
		return "LD A, (HL+)"
	case OpcodeLDAHLDec:
		return "LD A, (HL-)"
	case OpcodeLDSPHL:
		return "LD SP, HL"
	case OpcodeADDSPe:
		return fmt.Sprintf("ADD SP, PC$%s", fmtSignedOffset(di.Raw[1]))
	case OpcodeDECHLInd:
		return "DEC (HL)"
	case OpcodeINCHLInd:
		return "DEC (HL)"
	case OpcodeRST0x00, OpcodeRST0x08, OpcodeRST0x10, OpcodeRST0x18, OpcodeRST0x20, OpcodeRST0x28, OpcodeRST0x30, OpcodeRST0x38:
		return fmt.Sprintf("RST $00%sh", str[ln-2:])
	case OpcodeUndefD3, OpcodeUndefDB, OpcodeUndefDD, OpcodeUndefE3, OpcodeUndefE4, OpcodeUndefEB, OpcodeUndefEC, OpcodeUndefED, OpcodeUndefF4, OpcodeUndefFC, OpcodeUndefFD:
		return fmt.Sprintf("Undefined instruction %s", di.Raw[0].Hex())
	case OpcodeLDBCA, OpcodeLDDEA:
		return fmt.Sprintf("LD (%s), A", str[ln-3:ln-1])
	case OpcodeLDHLn:
		return fmt.Sprintf("LD (HL), $%s", di.Raw[1].Hex())
	case OpcodeRET, OpcodeNop, OpcodeRLA, OpcodeRLCA, OpcodeRRA, OpcodeRRCA, OpcodeDAA, OpcodeDI, OpcodeEI,
		OpcodeSTOP, OpcodeCCF, OpcodeSCF, OpcodeRETI, OpcodeHALT:
		return str
	case OpcodeRETZ, OpcodeRETC:
		return fmt.Sprintf("RET %s", str[ln-1:])
	case OpcodeRETNZ, OpcodeRETNC:
		return fmt.Sprintf("RET %s", str[ln-2:])
	case OpcodePUSHBC, OpcodePUSHDE, OpcodePUSHHL, OpcodePUSHAF, OpcodePOPBC, OpcodePOPDE, OpcodePOPHL, OpcodePOPAF:
		return fmt.Sprintf("%s %s", str[:ln-2], str[ln-2:])
	case OpcodeCPLaka2f:
		return "CPL"
	case OpcodeXORn, OpcodeADDn, OpcodeANDn, OpcodeORn, OpcodeADCn, OpcodeSBCn, OpcodeCPn, OpcodeSUBn:
		return fmt.Sprintf("%s A, $%s", str[:ln-1], di.Raw[1].Hex())
	case OpcodeLDAA, OpcodeLDAB, OpcodeLDAC, OpcodeLDAD, OpcodeLDAE, OpcodeLDAH, OpcodeLDAL,
		OpcodeLDBA, OpcodeLDBB, OpcodeLDBC, OpcodeLDBD, OpcodeLDBE, OpcodeLDBH, OpcodeLDBL,
		OpcodeLDCA, OpcodeLDCB, OpcodeLDCC, OpcodeLDCD, OpcodeLDCE, OpcodeLDCH, OpcodeLDCL,
		OpcodeLDDA, OpcodeLDDB, OpcodeLDDC, OpcodeLDDD, OpcodeLDDE, OpcodeLDDH, OpcodeLDDL,
		OpcodeLDEA, OpcodeLDEB, OpcodeLDEC, OpcodeLDED, OpcodeLDEE, OpcodeLDEH, OpcodeLDEL,
		OpcodeLDHA, OpcodeLDHB, OpcodeLDHC, OpcodeLDHD, OpcodeLDHE, OpcodeLDHH, OpcodeLDHL,
		OpcodeLDLA, OpcodeLDLB, OpcodeLDLC, OpcodeLDLD, OpcodeLDLE, OpcodeLDLH, OpcodeLDLL:
		return fmt.Sprintf("LD %s, %s", str[ln-2:ln-1], str[ln-1:])
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
		return fmt.Sprintf("LD %s, (HL)", str[ln-3:ln-2])
	case OpcodeLDHLA, OpcodeLDHLB, OpcodeLDHLC, OpcodeLDHLD, OpcodeLDHLE, OpcodeLDHLH, OpcodeLDHLL:
		return fmt.Sprintf("LD (HL), %s", str[ln-1:])
	case OpcodeDECA, OpcodeDECB, OpcodeDECC, OpcodeDECD, OpcodeDECE, OpcodeDECH, OpcodeDECL:
		return fmt.Sprintf("DEC %c", str[ln-1])
	case OpcodeINCA, OpcodeINCB, OpcodeINCC, OpcodeINCD, OpcodeINCE, OpcodeINCH, OpcodeINCL:
		return fmt.Sprintf("INC %c", str[ln-1])
	case OpcodeINCBC, OpcodeINCDE, OpcodeINCHL, OpcodeINCSP:
		return fmt.Sprintf("INC %s", str[ln-2:])
	case OpcodeDECBC, OpcodeDECDE, OpcodeDECHL, OpcodeDECSP:
		return fmt.Sprintf("DEC %s", str[ln-2:])
	case OpcodeJPnn, OpcodeCALLnn:
		return fmt.Sprintf("%s $%s", str[:ln-2], join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeJPCnn, OpcodeJPZnn:
		return fmt.Sprintf("JP %s, $%s", str[ln-3:ln-2], join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeJPNZnn, OpcodeJPNCnn:
		return fmt.Sprintf("JP %s, $%s", str[ln-4:ln-2], join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeCALLZnn, OpcodeCALLCnn:
		return fmt.Sprintf("CALL %s, $%s", str[ln-3:ln-2], join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeCALLNZnn, OpcodeCALLNCnn:
		return fmt.Sprintf("CALL %s, $%s", str[ln-4:ln-2], join16(di.Raw[2], di.Raw[1]).Hex())
	case OpcodeJRNZe, OpcodeJRNCe:
		return fmt.Sprintf("JR %s, PC%s", str[ln-3:ln-1], fmtSignedOffset(di.Raw[1]))
	case OpcodeJRZe, OpcodeJRCe:
		return fmt.Sprintf("JR %s, PC%s", str[ln-2:ln-1], fmtSignedOffset(di.Raw[1]))
	case OpcodeJRe:
		return fmt.Sprintf("JR PC%s", fmtSignedOffset(di.Raw[1]))
	case OpcodeLDAn, OpcodeLDBn, OpcodeLDCn, OpcodeLDDn, OpcodeLDEn, OpcodeLDHn, OpcodeLDLn:
		return fmt.Sprintf("LD %s, $%s", str[ln-2:ln-1], di.Raw[1].Hex())
	case OpcodeADDHLHL, OpcodeADDHLDE, OpcodeADDHLBC, OpcodeADDHLSP:
		return fmt.Sprintf("ADD %s, %s", str[ln-4:ln-2], str[ln-2:])
	case OpcodeLDHnA:
		return fmt.Sprintf("LDH ($%s),A; eff. LD $%s, A", di.Raw[1].Hex(), join16(0xff, di.Raw[1]).Hex())
	case OpcodeLDHAn:
		return fmt.Sprintf("LDH A,($%s); eff. LD A,$%s", di.Raw[1].Hex(), join16(0xff, di.Raw[1]).Hex())
	case OpcodeLDHAC:
		return "LDH A,(C); eff. LD A, [C+$0xff00]"
	case OpcodeLDHCA:
		return "LDH (C), A; eff. LD [C+$0xff00], a"
	case OpcodeLDHLSPe:
		return fmt.Sprintf("LD HL,SP%s", fmtSignedOffset(di.Raw[1]))
	case OpcodeCB:
		cbop := CBOp{Op: cb((di.Raw[1] & 0xf8) >> 3), Target: CBTarget(di.Raw[1] & 0x7)}
		return fmt.Sprintf("%s %s", cbop.Op, cbop.Target)
	case OpcodeJPHL:
		return "JP HL"
	}
	panicf("CAN'T FMT %v\n", di.Opcode)
	return ""
}

func fmtSignedOffset(offs Data8) string {
	if offs.SignBit() {
		return fmt.Sprintf("-$%s", offs.SignedAbs().Hex())
	}
	return fmt.Sprintf("$%s", offs.Hex())
}

type DataSection struct {
	Raw     []Data8
	Address Addr
}

type CodeSection struct {
	Instructions []DisInstruction
}

func (cs CodeSection) Address() Addr {
	return cs.Instructions[0].Address
}

type Disassembly struct {
	PC   Addr
	Code []CodeSection
	Data []DataSection
}

func (d *Disassembly) Print(w io.Writer) {
	for _, section := range d.Code {
		fmt.Fprintf(w, "\nCode section at %s\n", section.Address().Hex())
		for _, inst := range section.Instructions {
			if inst.Address == d.PC {
				fmt.Fprintf(w, "[%s]->%s\n", inst.Address.Hex(), inst.Asm())
			} else {
				fmt.Fprintf(w, "%sh | %s\n", inst.Address.Hex(), inst.Asm())
			}
		}
	}
	data := splitSections(d.Data)
	prevEndAddr := Addr(0xffff)
	for _, section := range data {
		if prevEndAddr != section.Address {
			fmt.Fprintf(w, "\nData section at %s\n", section.Address.Hex())
		}
		prevEndAddr = section.Address + Addr(len(section.Raw))
		allEqual := true
		testByte := section.Raw[0]
		for _, b := range section.Raw {
			if b != testByte {
				allEqual = false
				break
			}
		}
		if allEqual {
			fmt.Fprintf(w, "%0x bytes of %s\n", len(section.Raw), testByte.Hex())
			continue
		}

		i := 0
		for line := range (len(section.Raw) + 15) / 16 {
			fmt.Fprintf(w, "%s | ", (section.Address + Addr(line*16)).Hex())
			for range 16 {
				if i >= len(section.Raw) {
					break
				}
				fmt.Fprintf(w, "%s ", section.Raw[i].Hex())
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
					Address: section.Address + Addr(runStart),
					Raw:     raw[runStart : runStart+runLen],
				})
				// Continue after the run
				raw = raw[runStart+runLen:]
				section.Address += Addr(runStart + runLen)
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

func NewDisassembler(config *ConfigDisassembler) Disassembler {
	dis := Disassembler{
		Config:  config,
		Program: Block{Begin: 0},
		HRAM:    Block{Begin: AddrHRAMBegin, Decoded: make([]DisInstruction, SizeHRAM)},
		WRAM:    Block{Begin: AddrWRAMBegin, Decoded: make([]DisInstruction, SizeWRAM)},
	}
	return dis
}

func (dis *Disassembler) SetProgram(program []byte) {
	dis.Program.CanExplore = true
	dis.Program.Source = Data8Slice(program)
	dis.Program.Decoded = make([]DisInstruction, len(program))
	dis.cachedDisassembly = nil
}

func (dis *Disassembler) printf(format string, args ...any) {
	if !dis.Config.Trace {
		return
	}
	for range dis.stackIdx {
		fmt.Printf("  ")
	}
	if dis.stackIdx > 0 {
		fmt.Printf("%s ", dis.stack[dis.stackIdx-1].Hex())
	}
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}

func (dis *Disassembler) SetPC(address Addr) {
	if address >= AddrHRAMBegin && address <= AddrHRAMEnd {
		dis.HRAM.CanExplore = true
	}
	if address >= AddrWRAMBegin && address <= AddrWRAMEnd {
		dis.WRAM.CanExplore = true
	}
	dis.ExploreFrom(address)
}

func (dis *Disassembler) ExploreFrom(address Addr) {
	//dis.printf("SetPC %s", address.Hex())
	dis.stack = append(dis.stack, address)
	dis.stackIdx++
	defer func() {
		dis.stackIdx--
		dis.stack = dis.stack[:dis.stackIdx]
		dis.PC = address
	}()

	for {
		var block *Block
		if int(address) < len(dis.Program.Source) {
			block = &dis.Program
		} else if address >= AddrHRAMBegin && address <= AddrHRAMEnd {
			if !dis.HRAM.CanExplore {
				dis.printf("reached HRAM address before setting HRAM")
				return
			}
			block = &dis.HRAM
		} else if address >= AddrWRAMBegin && address <= AddrWRAMEnd {
			if !dis.WRAM.CanExplore {
				dis.printf("reached WRAM address before setting WRAM")
				return
			}
			block = &dis.WRAM
		} else {
			dis.printf("reached address outside of program, WRAM and HRAM")
			return
		}

		if block.Decoded[address-block.Begin].Size() > 0 {
			//dis.printf("already seen this address")
			return
		}

		maxConsecutiveNops := Addr(10)
		for i := range maxConsecutiveNops {
			next := address - block.Begin + i
			if next < block.Begin || int(next) >= len(block.Source) {
				break
			}
			if block.Source[next] == 0x00 {
				if i == maxConsecutiveNops-1 {
					dis.printf("too many nops ahead, probably not code")
					return
				}
			} else {
				break
			}
		}

		di, err := dis.readNewInstruction(address, block)
		if err != nil {
			panicf("failed decoding at %v: %v", address.Hex(), err)
			return
		}

		for i := range Addr(di.Size()) {
			if block.Decoded[address+i-block.Begin].Size() > 0 {
				dis.printf("at offset %s there is an existing instruction, returning", i.Dec())
				return
			}
		}

		dis.insert(di, block)
		address += Addr(di.Size())
		dis.stack[dis.stackIdx-1] = address

		if dis.checkBranches(di, block) {
			return
		}
	}
}

func (dis *Disassembler) Disassembly() *Disassembly {
	if dis.cachedDisassembly != nil && len(dis.cachedDisassembly.Code) > 0 {
		return dis.cachedDisassembly
	}

	var out Disassembly
	for _, block := range []*Block{&dis.Program, &dis.HRAM, &dis.WRAM} {
		if block.Source == nil {
			continue
		}
		currCodeSection := CodeSection{}
		currDataSection := DataSection{}
		for offs := 0; offs < len(block.Decoded); {
			di := block.Decoded[offs]
			if di.Size() > 0 {
				if currDataSection.Raw != nil {
					out.Data = append(out.Data, currDataSection)
				}
				currCodeSection.Instructions = append(currCodeSection.Instructions, di)
				offs += int(di.Size())

				currDataSection.Raw = nil
			} else {
				if currCodeSection.Instructions != nil {
					out.Code = append(out.Code, currCodeSection)
				}

				if currDataSection.Raw == nil {
					currDataSection.Address = Addr(offs) + block.Begin
				}
				currDataSection.Raw = append(currDataSection.Raw, block.Source[offs])
				offs++
				currCodeSection.Instructions = nil
			}
		}
		if currCodeSection.Instructions != nil {
			out.Code = append(out.Code, currCodeSection)
		}
		if currDataSection.Raw != nil {
			out.Data = append(out.Data, currDataSection)
		}
	}

	out.PC = dis.PC
	dis.cachedDisassembly = &out
	return &out
}

func (dis *Disassembler) insert(di DisInstruction, block *Block) {
	dis.printf("insert %v at %s:%s", di.Opcode, di.Address.Hex(), (di.Address + Addr(di.Size()) - 1).Hex())
	for addr := di.Address; addr != di.Address+Addr(di.Size()); addr++ {
		block.Decoded[addr-block.Begin] = di
	}
	dis.cachedDisassembly = nil
}

// returns true if unconditional branch, i.e. can never fall through
func (dis *Disassembler) checkBranches(di DisInstruction, block *Block) bool {
	switch di.Opcode {
	case OpcodeJRNZe, OpcodeJRZe, OpcodeJRe:
		e := int8(di.Raw[1])

		// branch taken
		dis.printf("checking branch-taken for relative jump")
		if e > 0 {
			dis.ExploreFrom(di.Address + Addr(di.Size()) + Addr(e))
		} else {
			dis.ExploreFrom(di.Address + Addr(di.Size()) - Addr(-e))
		}

		// branch not taken will be inspected after this
		if di.Opcode == OpcodeJRe {
			dis.printf("unconditional branch, can never fall through")
			return true
		}
	case OpcodeJPnn, OpcodeJPCnn, OpcodeJPNCnn, OpcodeJPZnn, OpcodeJPNZnn, OpcodeCALLnn:
		addr := Addr(join16(di.Raw[2], di.Raw[1]))
		dis.printf("checking branch-taken for absolute jump")
		dis.ExploreFrom(addr)
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
		block.Decoded[di.Address-block.Begin] = DisInstruction{}
		return true
	case OpcodeRET, OpcodeRETI, OpcodeJPHL:
		return true
	}

	return false
}

func (dis *Disassembler) doRST(vec Addr) {
	dis.printf("unconditional function call to %s", vec.Hex())
	dis.ExploreFrom(vec)
	dis.printf("will probably return here")
}

func (dis *Disassembler) readNewInstruction(addr Addr, block *Block) (DisInstruction, error) {
	di := DisInstruction{
		Address: addr,
		Visited: true,
	}
	if int(addr-block.Begin) >= len(block.Source) {
		return di, fmt.Errorf("address out of bounds")
	}
	di.Opcode = Opcode(block.Source[addr-block.Begin])
	ok := di.Opcode.IsValid()
	if !ok {
		return di, fmt.Errorf("invalid opcode %v", di.Opcode)
	}
	if di.Size() == 0 {
		return di, fmt.Errorf("no size set for opcode %v (0x%x)", di.Opcode, int(di.Opcode))
	}
	for i := range Addr(di.Size()) {
		di.Raw[i] = block.Source[di.Address+i-block.Begin]
	}
	return di, nil
}
