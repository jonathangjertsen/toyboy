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

	stack    []Addr
	stackIdx int
}

type Block struct {
	CanExplore bool
	Name       string
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
		return fmt.Sprintf("UNDEF %s", di.Raw[0].Hex())
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
	case OpcodeLDABC, OpcodeLDADE:
		return fmt.Sprintf("LD A, %s", str[ln-2:])
	case OpcodeLDAHL, OpcodeLDBHL, OpcodeLDCHL, OpcodeLDDHL, OpcodeLDEHL, OpcodeLDHHL, OpcodeLDLHL:
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
		return fmt.Sprintf("LDH ($ff00+%s)", di.Raw[1].Hex())
	case OpcodeLDHAn:
		return fmt.Sprintf("LDH A,($ff00+%s)", di.Raw[1].Hex())
	case OpcodeLDHAC:
		return "LDH A,($ff00+C)"
	case OpcodeLDHCA:
		return "LDH (C+$ff00), A"
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
	if offs&SignBit8 != 0 {
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
	data := splitSections(d.Data)

	nCodeSections := len(d.Code)
	nDataSections := len(data)
	iData := 0
	iCode := 0

	prevDataEndAddr := Addr(0xffff)
	for iData < nDataSections && iCode < nCodeSections {
		var codeSection *CodeSection
		var dataSection *DataSection
		var selectCode = false
		if iCode < nCodeSections {
			codeSection = &d.Code[iCode]
		}
		if iData < nDataSections {
			dataSection = &data[iData]
		}
		if codeSection != nil && dataSection != nil {
			selectCode = codeSection.Address() < dataSection.Address
		} else if codeSection != nil {
			selectCode = true
		} else {
			selectCode = false
		}
		if selectCode {
			printCodeSection(w, codeSection, d.PC)
			iCode++
		} else {
			if prevDataEndAddr != dataSection.Address {
				fmt.Fprintf(w, "\nData section at %s\n", dataSection.Address.Hex())
			}
			prevDataEndAddr = dataSection.Address + Addr(len(dataSection.Raw))
			printDataSection(w, dataSection)
			iData++
		}
	}
}

func printCodeSection(w io.Writer, section *CodeSection, pc Addr) {
	fmt.Fprintf(w, "\nCode section at %s\n", section.Address().Hex())
	for _, inst := range section.Instructions {
		if inst.Address == pc {
			fmt.Fprintf(w, "[%s]->%s\n", inst.Address.Hex(), inst.Asm())
		} else {
			fmt.Fprintf(w, "%sh | %s\n", inst.Address.Hex(), inst.Asm())
		}
	}
}

func printDataSection(w io.Writer, section *DataSection) {
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
		return
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
		Program: Block{Name: "Program", Begin: 0},
		HRAM:    Block{Name: "HRAM", Begin: AddrHRAMBegin, Decoded: make([]DisInstruction, SizeHRAM)},
		WRAM:    Block{Name: "WRAM", Begin: AddrWRAMBegin, Decoded: make([]DisInstruction, SizeWRAM)},
	}
	return dis
}

func (dis *Disassembler) SetProgram(program []byte) {
	dis.Program.CanExplore = true
	dis.Program.Source = Data8Slice(program)
	dis.Program.Decoded = make([]DisInstruction, len(program))
}

func (dis *Disassembler) print(msg string) {
	if !dis.Config.Trace {
		return
	}
	for range dis.stackIdx {
		fmt.Printf("  ")
	}
	if dis.stackIdx > 0 {
		fmt.Printf("%s ", dis.stack[dis.stackIdx-1].Hex())
	}
	fmt.Println(msg)
}

func (dis *Disassembler) SetPC(address Addr) {
	if address >= AddrHRAMBegin && address <= AddrHRAMEnd {
		dis.HRAM.CanExplore = true
	}
	if address >= AddrWRAMBegin && address <= AddrWRAMEnd {
		dis.WRAM.CanExplore = true
	}
	dis.PC = address
	dis.ExploreFrom(address)
}

func (dis *Disassembler) ExploreFrom(address Addr) {
	dis.stack = append(dis.stack, address)
	dis.stackIdx++

	dis.exploreFromInner(address)

	dis.stackIdx--
	dis.stack = dis.stack[:dis.stackIdx]
}

func (dis *Disassembler) exploreFromInner(address Addr) {
	for {
		var block *Block
		if int(address) < len(dis.Program.Source) {
			block = &dis.Program
		} else if address >= AddrHRAMBegin && address <= AddrHRAMEnd {
			if !dis.HRAM.CanExplore {
				dis.print("reached HRAM address before setting HRAM")
				return
			}
			block = &dis.HRAM
		} else if address >= AddrWRAMBegin && address <= AddrWRAMEnd {
			if !dis.WRAM.CanExplore {
				dis.print("reached WRAM address before setting WRAM")
				return
			}
			block = &dis.WRAM
		} else {
			dis.print("reached address outside of program, WRAM and HRAM")
			return
		}

		if block.Decoded[address-block.Begin].Size() > 0 {
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
					dis.print("too many nops ahead, probably not code")
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
				dis.print("slices an existing instruction, returning")
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

func (dis *Disassembler) Disassembly(start, end Addr) *Disassembly {
	var out Disassembly
	for _, block := range []*Block{&dis.Program, &dis.HRAM, &dis.WRAM} {
		if block.Source == nil {
			fmt.Printf("%s source is nil\n", block.Name)
			continue
		}
		if block.Begin > end {
			fmt.Printf("%s block.Begin=%s > end=%s\n", block.Name, block.Begin.Hex(), end.Hex())
			continue
		}
		if block.Begin > start {
			start = block.Begin
		}
		blockEnd := block.Begin + Addr(len(block.Source))
		if blockEnd < start {
			fmt.Printf("%s blockEnd=%s < start=%s\n", block.Name, blockEnd.Hex(), start.Hex())
			continue
		}
		if blockEnd < end {
			end = blockEnd
		}

		if start >= end {
			fmt.Printf("%s start=%s >= end=%s\n", block.Name, start.Hex(), end.Hex())
			continue
		}

		beginOffs := start - block.Begin
		// align to instruction
		beginOffs -= start - block.Decoded[beginOffs].Address

		endOffs := end - block.Begin

		currCodeSection := CodeSection{}
		currDataSection := DataSection{}
		fmt.Printf("%s slurp start=%s end=%s\n", block.Name, start.Hex(), end.Hex())
		for offs := beginOffs; offs < endOffs; {
			di := block.Decoded[offs]
			if di.Size() > 0 {
				if currDataSection.Raw != nil {
					out.Data = append(out.Data, currDataSection)
				}
				currCodeSection.Instructions = append(currCodeSection.Instructions, di)
				offs += Addr(di.Size())

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
	return &out
}

func (dis *Disassembler) insert(di DisInstruction, block *Block) {
	if dis.Config.Trace {
		dis.print(fmt.Sprintf("insert %s at %s:%s", di.Asm(), di.Address.Hex(), (di.Address + Addr(di.Size()) - 1).Hex()))
	}
	if di.Address == 0xa8 {
		panic("here")
	}
	for addr := di.Address; addr != di.Address+Addr(di.Size()); addr++ {
		block.Decoded[addr-block.Begin] = di
	}
}

// returns true if unconditional branch, i.e. can never fall through
func (dis *Disassembler) checkBranches(di DisInstruction, block *Block) bool {
	switch di.Opcode {
	case OpcodeJRNZe, OpcodeJRZe, OpcodeJRNCe, OpcodeJRCe, OpcodeJRe:
		e := int8(di.Raw[1])

		// branch taken
		if dis.Config.Trace {
			dis.print("checking branch-taken for relative jump @ " + di.Asm())
		}
		if e > 0 {
			dis.ExploreFrom(di.Address + Addr(di.Size()) + Addr(e))
		} else {
			dis.ExploreFrom(di.Address + Addr(di.Size()) - Addr(-e))
		}

		// branch not taken will be inspected after this
		if di.Opcode == OpcodeJRe {
			dis.print("unconditional branch, can never fall through")
			return true
		}
	case OpcodeJPnn, OpcodeJPCnn, OpcodeJPNCnn, OpcodeJPZnn, OpcodeJPNZnn,
		OpcodeCALLnn, OpcodeCALLCnn, OpcodeCALLNCnn, OpcodeCALLZnn, OpcodeCALLNZnn:
		addr := Addr(join16(di.Raw[2], di.Raw[1]))
		if dis.Config.Trace {
			dis.print("checking branch-taken for absolute jump @ " + di.Asm())
		}
		dis.ExploreFrom(addr)
		if di.Opcode == OpcodeJPnn || di.Opcode == OpcodeJPHL {
			dis.print("unconditional branch, can never fall through")
			return true
		} else if di.Opcode == OpcodeCALLnn {
			dis.print("will probably return here")
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
		dis.print("dropping undefined instruction " + di.Opcode.String())
		block.Decoded[di.Address-block.Begin] = DisInstruction{}
		return true
	case OpcodeRET, OpcodeRETI, OpcodeJPHL:
		// will likely go back to somewhere already explored
		return true
	case OpcodeRETC, OpcodeRETNC, OpcodeRETZ, OpcodeRETNZ:
		// will likely go back to somewhere already explored
		return false
	case OpcodeLDSPnn, OpcodeLDHLnn, OpcodeLDBCnn, OpcodeLDDEnn,
		OpcodeLDHLADec, OpcodeLDHLAInc,
		OpcodeLDAHLInc, OpcodeLDAHLDec,
		OpcodeCB,
		OpcodeLDAn, OpcodeLDBn, OpcodeLDCn, OpcodeLDDn, OpcodeLDEn, OpcodeLDHn, OpcodeLDLn, OpcodeLDHLn,
		OpcodeLDHCA, OpcodeLDnnA, OpcodeLDnnSP, OpcodeLDAnn,
		OpcodeLDHnA, OpcodeLDHAn,
		OpcodeLDABC, OpcodeLDADE,
		OpcodeXORA, OpcodeXORB, OpcodeXORC, OpcodeXORD, OpcodeXORE, OpcodeXORL, OpcodeXORH, OpcodeXORHL, OpcodeXORn,
		OpcodeINCA, OpcodeINCB, OpcodeINCC, OpcodeINCD, OpcodeINCE, OpcodeINCL, OpcodeINCH, OpcodeINCHL,
		OpcodeDECA, OpcodeDECB, OpcodeDECC, OpcodeDECD, OpcodeDECE, OpcodeDECL, OpcodeDECH, OpcodeDECHL,
		OpcodeANDA, OpcodeANDB, OpcodeANDC, OpcodeANDD, OpcodeANDE, OpcodeANDL, OpcodeANDH, OpcodeANDHL, OpcodeANDn,
		OpcodeORA, OpcodeORB, OpcodeORC, OpcodeORD, OpcodeORE, OpcodeORL, OpcodeORH, OpcodeORHL, OpcodeORn,
		OpcodeADDA, OpcodeADDB, OpcodeADDC, OpcodeADDD, OpcodeADDE, OpcodeADDL, OpcodeADDH, OpcodeADDHL, OpcodeADDn,
		OpcodeADCA, OpcodeADCB, OpcodeADCC, OpcodeADCD, OpcodeADCE, OpcodeADCL, OpcodeADCH, OpcodeADCHL, OpcodeADCn,
		OpcodeSUBA, OpcodeSUBB, OpcodeSUBC, OpcodeSUBD, OpcodeSUBE, OpcodeSUBL, OpcodeSUBH, OpcodeSUBHL, OpcodeSUBn,
		OpcodeSBCA, OpcodeSBCB, OpcodeSBCC, OpcodeSBCD, OpcodeSBCE, OpcodeSBCL, OpcodeSBCH, OpcodeSBCHL, OpcodeSBCn,
		OpcodeCPA, OpcodeCPB, OpcodeCPC, OpcodeCPD, OpcodeCPE, OpcodeCPL, OpcodeCPH, OpcodeCPHL, OpcodeCPn,
		OpcodeLDAA, OpcodeLDAB, OpcodeLDAC, OpcodeLDAD, OpcodeLDAE, OpcodeLDAL, OpcodeLDAH, OpcodeLDAHL,
		OpcodeLDBA, OpcodeLDBB, OpcodeLDBC, OpcodeLDBD, OpcodeLDBE, OpcodeLDBL, OpcodeLDBH, OpcodeLDBHL,
		OpcodeLDCA, OpcodeLDCB, OpcodeLDCC, OpcodeLDCD, OpcodeLDCE, OpcodeLDCL, OpcodeLDCH, OpcodeLDCHL,
		OpcodeLDDA, OpcodeLDDB, OpcodeLDDC, OpcodeLDDD, OpcodeLDDE, OpcodeLDDL, OpcodeLDDH, OpcodeLDDHL,
		OpcodeLDEA, OpcodeLDEB, OpcodeLDEC, OpcodeLDED, OpcodeLDEE, OpcodeLDEL, OpcodeLDEH, OpcodeLDEHL,
		OpcodeLDHA, OpcodeLDHB, OpcodeLDHC, OpcodeLDHD, OpcodeLDHE, OpcodeLDHL, OpcodeLDHH, OpcodeLDHHL,
		OpcodeLDLA, OpcodeLDLB, OpcodeLDLC, OpcodeLDLD, OpcodeLDLE, OpcodeLDLL, OpcodeLDLH, OpcodeLDLHL,
		OpcodeLDSPHL,
		OpcodeADDHLBC, OpcodeADDHLDE, OpcodeADDHLHL, OpcodeADDHLSP,
		OpcodeLDBCA, OpcodeLDDEA,
		OpcodeLDHLA, OpcodeLDHLB, OpcodeLDHLC, OpcodeLDHLD, OpcodeLDHLE, OpcodeLDHLL, OpcodeLDHLH, OpcodeLDHLSPe,
		OpcodePUSHAF, OpcodePUSHBC, OpcodePUSHDE, OpcodePUSHHL,
		OpcodePOPAF, OpcodePOPBC, OpcodePOPDE, OpcodePOPHL,
		OpcodeRLA, OpcodeRLCA, OpcodeRRCA, OpcodeRRA,
		OpcodeINCDE, OpcodeINCBC, OpcodeINCSP, OpcodeINCHLInd,
		OpcodeDECDE, OpcodeDECBC, OpcodeDECSP, OpcodeDECHLInd,
		OpcodeDAA, OpcodeCCF, OpcodeSCF,
		OpcodeCPLaka2f,
		OpcodeNop, OpcodeDI:
		// non-branching
		return false
	case OpcodeEI, OpcodeHALT, OpcodeSTOP:
		// likely to continue in an interrupt, then return here
		return false
	default:
		panicf("check %s", di.Opcode)
	}

	return false
}

func (dis *Disassembler) doRST(vec Addr) {
	dis.print("unconditional function call to %s" + vec.Hex())
	dis.ExploreFrom(vec)
	dis.print("will probably return here")
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
