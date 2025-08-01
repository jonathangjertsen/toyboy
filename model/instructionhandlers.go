package model

import (
	"fmt"
	"slices"
)

type InstructionHandling func(gb *Gameboy, e int) bool

type HandlerArray [256]InstructionHandling

type CBOp struct {
	Op     cb
	Target CBTarget
}

func (cbv cb) Is3Cycles() bool {
	return slices.Contains([]cb{
		CbBit0,
		CbBit1,
		CbBit2,
		CbBit3,
		CbBit4,
		CbBit5,
		CbBit6,
		CbBit7,
	}, cbv)
}

func (cb CBOp) String() string {
	return fmt.Sprintf("%s %s", cb.Op, cb.Target)
}

func Handlers(gb *Gameboy) HandlerArray {
	return HandlerArray{
		OpcodeNop:      singleCycle("NOP", func(gb *Gameboy) {}),
		OpcodeLDAA:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.A),
		OpcodeLDAB:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.B),
		OpcodeLDAC:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.C),
		OpcodeLDAD:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.D),
		OpcodeLDAE:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.E),
		OpcodeLDAH:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.H),
		OpcodeLDAL:     ld(&gb.CPU.Regs.A, &gb.CPU.Regs.L),
		OpcodeLDAHL:    ldrhl(&gb.CPU.Regs.A),
		OpcodeLDAHLInc: ldahlinc,
		OpcodeLDAHLDec: ldahldec,
		OpcodeLDBA:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.A),
		OpcodeLDBB:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.B),
		OpcodeLDBC:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.C),
		OpcodeLDBD:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.D),
		OpcodeLDBE:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.E),
		OpcodeLDBH:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.H),
		OpcodeLDBL:     ld(&gb.CPU.Regs.B, &gb.CPU.Regs.L),
		OpcodeLDBHL:    ldrhl(&gb.CPU.Regs.B),
		OpcodeLDCA:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.A),
		OpcodeLDCB:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.B),
		OpcodeLDCC:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.C),
		OpcodeLDCD:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.D),
		OpcodeLDCE:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.E),
		OpcodeLDCH:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.H),
		OpcodeLDCL:     ld(&gb.CPU.Regs.C, &gb.CPU.Regs.L),
		OpcodeLDCHL:    ldrhl(&gb.CPU.Regs.C),
		OpcodeLDDA:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.A),
		OpcodeLDDB:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.B),
		OpcodeLDDC:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.C),
		OpcodeLDDD:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.D),
		OpcodeLDDE:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.E),
		OpcodeLDDH:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.H),
		OpcodeLDDL:     ld(&gb.CPU.Regs.D, &gb.CPU.Regs.L),
		OpcodeLDDHL:    ldrhl(&gb.CPU.Regs.D),
		OpcodeLDEA:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.A),
		OpcodeLDEB:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.B),
		OpcodeLDEC:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.C),
		OpcodeLDED:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.D),
		OpcodeLDEE:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.E),
		OpcodeLDEH:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.H),
		OpcodeLDEL:     ld(&gb.CPU.Regs.E, &gb.CPU.Regs.L),
		OpcodeLDEHL:    ldrhl(&gb.CPU.Regs.E),
		OpcodeLDHA:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.A),
		OpcodeLDHB:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.B),
		OpcodeLDHC:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.C),
		OpcodeLDHD:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.D),
		OpcodeLDHE:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.E),
		OpcodeLDHH:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.H),
		OpcodeLDHL:     ld(&gb.CPU.Regs.H, &gb.CPU.Regs.L),
		OpcodeLDHHL:    ldrhl(&gb.CPU.Regs.H),
		OpcodeLDLA:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.A),
		OpcodeLDLB:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.B),
		OpcodeLDLC:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.C),
		OpcodeLDLD:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.D),
		OpcodeLDLE:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.E),
		OpcodeLDLH:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.H),
		OpcodeLDLL:     ld(&gb.CPU.Regs.L, &gb.CPU.Regs.L),
		OpcodeLDLHL:    ldrhl(&gb.CPU.Regs.L),
		OpcodeRLA:      singleCycle("RLA", func(gb *Gameboy) { gb.CPU.Regs.SetFlagsAndA(RLA(gb.CPU.Regs.A, gb.CPU.Regs.GetFlagC())) }),
		OpcodeRRA:      singleCycle("RRA", func(gb *Gameboy) { gb.CPU.Regs.SetFlagsAndA(RRA(gb.CPU.Regs.A, gb.CPU.Regs.GetFlagC())) }),
		OpcodeRLCA:     singleCycle("RLCA", func(gb *Gameboy) { gb.CPU.Regs.SetFlagsAndA(RLCA(gb.CPU.Regs.A)) }),
		OpcodeRRCA:     singleCycle("RRCA", func(gb *Gameboy) { gb.CPU.Regs.SetFlagsAndA(RRCA(gb.CPU.Regs.A)) }),
		OpcodeORA:      orreg(&gb.CPU.Regs.A),
		OpcodeORB:      orreg(&gb.CPU.Regs.B),
		OpcodeORC:      orreg(&gb.CPU.Regs.C),
		OpcodeORD:      orreg(&gb.CPU.Regs.D),
		OpcodeORE:      orreg(&gb.CPU.Regs.E),
		OpcodeORH:      orreg(&gb.CPU.Regs.H),
		OpcodeORL:      orreg(&gb.CPU.Regs.L),
		OpcodeORHL:     aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(OR(gb.CPU.Regs.A, v)) }),
		OpcodeANDA:     andreg(&gb.CPU.Regs.A),
		OpcodeANDB:     andreg(&gb.CPU.Regs.B),
		OpcodeANDC:     andreg(&gb.CPU.Regs.C),
		OpcodeANDD:     andreg(&gb.CPU.Regs.D),
		OpcodeANDE:     andreg(&gb.CPU.Regs.E),
		OpcodeANDH:     andreg(&gb.CPU.Regs.H),
		OpcodeANDL:     andreg(&gb.CPU.Regs.L),
		OpcodeANDHL:    aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(AND(gb.CPU.Regs.A, v)) }),
		OpcodeXORA:     xorreg(&gb.CPU.Regs.A),
		OpcodeXORB:     xorreg(&gb.CPU.Regs.B),
		OpcodeXORC:     xorreg(&gb.CPU.Regs.C),
		OpcodeXORD:     xorreg(&gb.CPU.Regs.D),
		OpcodeXORE:     xorreg(&gb.CPU.Regs.E),
		OpcodeXORH:     xorreg(&gb.CPU.Regs.H),
		OpcodeXORL:     xorreg(&gb.CPU.Regs.L),
		OpcodeXORHL:    aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(XOR(gb.CPU.Regs.A, v)) }),
		OpcodeSUBA:     subreg(&gb.CPU.Regs.A),
		OpcodeSUBB:     subreg(&gb.CPU.Regs.B),
		OpcodeSUBC:     subreg(&gb.CPU.Regs.C),
		OpcodeSUBD:     subreg(&gb.CPU.Regs.D),
		OpcodeSUBE:     subreg(&gb.CPU.Regs.E),
		OpcodeSUBH:     subreg(&gb.CPU.Regs.H),
		OpcodeSUBL:     subreg(&gb.CPU.Regs.L),
		OpcodeSUBHL:    aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, v, false)) }),
		OpcodeSBCA:     sbcreg(&gb.CPU.Regs.A),
		OpcodeSBCB:     sbcreg(&gb.CPU.Regs.B),
		OpcodeSBCC:     sbcreg(&gb.CPU.Regs.C),
		OpcodeSBCD:     sbcreg(&gb.CPU.Regs.D),
		OpcodeSBCE:     sbcreg(&gb.CPU.Regs.E),
		OpcodeSBCH:     sbcreg(&gb.CPU.Regs.H),
		OpcodeSBCL:     sbcreg(&gb.CPU.Regs.L),
		OpcodeSBCHL:    aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, v, gb.CPU.Regs.GetFlagC())) }),
		OpcodeCPA:      cpreg(&gb.CPU.Regs.A),
		OpcodeCPB:      cpreg(&gb.CPU.Regs.B),
		OpcodeCPC:      cpreg(&gb.CPU.Regs.C),
		OpcodeCPD:      cpreg(&gb.CPU.Regs.D),
		OpcodeCPE:      cpreg(&gb.CPU.Regs.E),
		OpcodeCPH:      cpreg(&gb.CPU.Regs.H),
		OpcodeCPL:      cpreg(&gb.CPU.Regs.L),
		OpcodeCPHL:     aluhl(func(v Data8) { gb.CPU.Regs.SetFlags(SUB(gb.CPU.Regs.A, v, false)) }),
		OpcodeADDA:     addreg(&gb.CPU.Regs.A),
		OpcodeADDB:     addreg(&gb.CPU.Regs.B),
		OpcodeADDC:     addreg(&gb.CPU.Regs.C),
		OpcodeADDD:     addreg(&gb.CPU.Regs.D),
		OpcodeADDE:     addreg(&gb.CPU.Regs.E),
		OpcodeADDH:     addreg(&gb.CPU.Regs.H),
		OpcodeADDL:     addreg(&gb.CPU.Regs.L),
		OpcodeADDHL:    aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, v, false)) }),
		OpcodeADDSPe:   addspe,
		OpcodeADCA:     adcreg(&gb.CPU.Regs.A),
		OpcodeADCB:     adcreg(&gb.CPU.Regs.B),
		OpcodeADCC:     adcreg(&gb.CPU.Regs.C),
		OpcodeADCD:     adcreg(&gb.CPU.Regs.D),
		OpcodeADCE:     adcreg(&gb.CPU.Regs.E),
		OpcodeADCH:     adcreg(&gb.CPU.Regs.H),
		OpcodeADCL:     adcreg(&gb.CPU.Regs.L),
		OpcodeADCHL:    aluhl(func(v Data8) { gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, v, gb.CPU.Regs.GetFlagC())) }),
		OpcodeDAA: singleCycle("DAA", func(gb *Gameboy) {
			gb.CPU.Regs.SetFlagsAndA(DAA(gb.CPU.Regs.A, gb.CPU.Regs.GetFlagC(), gb.CPU.Regs.GetFlagN(), gb.CPU.Regs.GetFlagH()))
		}),
		OpcodeCPLaka2f: singleCycle("CPL", func(gb *Gameboy) {
			gb.CPU.Regs.A ^= 0xff
			gb.CPU.Regs.SetFlagN(true)
			gb.CPU.Regs.SetFlagH(true)
		}),
		OpcodeCCF: singleCycle("CCF", func(gb *Gameboy) {
			gb.CPU.Regs.SetFlagC(!gb.CPU.Regs.GetFlagC())
			gb.CPU.Regs.SetFlagN(false)
			gb.CPU.Regs.SetFlagH(false)
		}),
		OpcodeSCF: singleCycle("SCF", func(gb *Gameboy) {
			gb.CPU.Regs.SetFlagC(true)
			gb.CPU.Regs.SetFlagN(false)
			gb.CPU.Regs.SetFlagH(false)
		}),
		OpcodeDECA:     decreg(&gb.CPU.Regs.A),
		OpcodeDECB:     decreg(&gb.CPU.Regs.B),
		OpcodeDECC:     decreg(&gb.CPU.Regs.C),
		OpcodeDECD:     decreg(&gb.CPU.Regs.D),
		OpcodeDECE:     decreg(&gb.CPU.Regs.E),
		OpcodeDECH:     decreg(&gb.CPU.Regs.H),
		OpcodeDECL:     decreg(&gb.CPU.Regs.L),
		OpcodeINCA:     increg(&gb.CPU.Regs.A),
		OpcodeINCB:     increg(&gb.CPU.Regs.B),
		OpcodeINCC:     increg(&gb.CPU.Regs.C),
		OpcodeINCD:     increg(&gb.CPU.Regs.D),
		OpcodeINCE:     increg(&gb.CPU.Regs.E),
		OpcodeINCH:     increg(&gb.CPU.Regs.H),
		OpcodeINCL:     increg(&gb.CPU.Regs.L),
		OpcodeINCHLInd: inchlind,
		OpcodeDECHLInd: dechlind,
		OpcodeDI: func(gb *Gameboy, e int) bool {
			gb.Interrupts.SetIMENextCycle = false
			gb.Interrupts.SetIME(gb.Mem, false)
			return true
		},
		OpcodeEI: func(gb *Gameboy, e int) bool {
			gb.Interrupts.SetIMENextCycle = true
			return true
		},
		OpcodeHALT:     halt,
		OpcodeJRe:      jre,
		OpcodeJPnn:     jpnn,
		OpcodeJPHL:     jphl,
		OpcodeJRZe:     jrcce(func() bool { return gb.CPU.Regs.GetFlagZ() }),
		OpcodeJRCe:     jrcce(func() bool { return gb.CPU.Regs.GetFlagC() }),
		OpcodeJRNZe:    jrcce(func() bool { return !gb.CPU.Regs.GetFlagZ() }),
		OpcodeJRNCe:    jrcce(func() bool { return !gb.CPU.Regs.GetFlagC() }),
		OpcodeJPCnn:    jpccnn(func() bool { return gb.CPU.Regs.GetFlagC() }),
		OpcodeJPNCnn:   jpccnn(func() bool { return !gb.CPU.Regs.GetFlagC() }),
		OpcodeJPZnn:    jpccnn(func() bool { return gb.CPU.Regs.GetFlagZ() }),
		OpcodeJPNZnn:   jpccnn(func() bool { return !gb.CPU.Regs.GetFlagZ() }),
		OpcodeINCBC:    iduOp(func() { gb.CPU.SetBC(gb.CPU.GetBC() + 1) }),
		OpcodeINCDE:    iduOp(func() { gb.CPU.SetDE(gb.CPU.GetDE() + 1) }),
		OpcodeINCHL:    iduOp(func() { gb.CPU.SetHL(gb.CPU.GetHL() + 1) }),
		OpcodeINCSP:    iduOp(func() { gb.CPU.Regs.SP++ }),
		OpcodeDECBC:    iduOp(func() { gb.CPU.SetBC(gb.CPU.GetBC() - 1) }),
		OpcodeDECDE:    iduOp(func() { gb.CPU.SetDE(gb.CPU.GetDE() - 1) }),
		OpcodeDECHL:    iduOp(func() { gb.CPU.SetHL(gb.CPU.GetHL() - 1) }),
		OpcodeDECSP:    iduOp(func() { gb.CPU.Regs.SP-- }),
		OpcodeCALLnn:   callnn,
		OpcodeCALLNZnn: callccnn(func() bool { return !gb.CPU.Regs.GetFlagZ() }),
		OpcodeCALLZnn:  callccnn(func() bool { return gb.CPU.Regs.GetFlagZ() }),
		OpcodeCALLNCnn: callccnn(func() bool { return !gb.CPU.Regs.GetFlagC() }),
		OpcodeCALLCnn:  callccnn(func() bool { return gb.CPU.Regs.GetFlagC() }),
		OpcodeRET:      ret,
		OpcodeRETI:     reti,
		OpcodeRETZ:     retcc(func() bool { return gb.CPU.Regs.GetFlagZ() }),
		OpcodeRETNZ:    retcc(func() bool { return !gb.CPU.Regs.GetFlagZ() }),
		OpcodeRETC:     retcc(func() bool { return gb.CPU.Regs.GetFlagC() }),
		OpcodeRETNC:    retcc(func() bool { return !gb.CPU.Regs.GetFlagC() }),
		OpcodePUSHBC:   push(&gb.CPU.Regs.B, &gb.CPU.Regs.C),
		OpcodePUSHDE:   push(&gb.CPU.Regs.D, &gb.CPU.Regs.E),
		OpcodePUSHHL:   push(&gb.CPU.Regs.H, &gb.CPU.Regs.L),
		OpcodePUSHAF:   push(&gb.CPU.Regs.A, &gb.CPU.Regs.F),
		OpcodePOPBC:    pop(&gb.CPU.Regs.B, &gb.CPU.Regs.C),
		OpcodePOPDE:    pop(&gb.CPU.Regs.D, &gb.CPU.Regs.E),
		OpcodePOPHL:    pop(&gb.CPU.Regs.H, &gb.CPU.Regs.L),
		OpcodePOPAF:    pop(&gb.CPU.Regs.A, &gb.CPU.Regs.F),
		OpcodeADDHLHL:  addhlrr(&gb.CPU.Regs.H, &gb.CPU.Regs.L),
		OpcodeADDHLBC:  addhlrr(&gb.CPU.Regs.B, &gb.CPU.Regs.C),
		OpcodeADDHLDE:  addhlrr(&gb.CPU.Regs.D, &gb.CPU.Regs.E),
		OpcodeADDHLSP:  addhlsp,
		OpcodeLDBCnn:   ldxxnn(func(wz Data16) { gb.CPU.SetBC(wz) }),
		OpcodeLDDEnn:   ldxxnn(func(wz Data16) { gb.CPU.SetDE(wz) }),
		OpcodeLDHLnn:   ldxxnn(func(wz Data16) { gb.CPU.SetHL(wz) }),
		OpcodeLDSPnn:   ldxxnn(func(wz Data16) { gb.CPU.SetSP(Addr(wz)) }),
		OpcodeLDHLn:    ldhln,
		OpcodeLDHLSPe:  ldhlspe,
		OpcodeLDSPHL:   ldsphl,
		OpcodeLDHLAInc: ldhlr(&gb.CPU.Regs.A, +1),
		OpcodeLDHLADec: ldhlr(&gb.CPU.Regs.A, -1),
		OpcodeLDHLA:    ldhlr(&gb.CPU.Regs.A, 0),
		OpcodeLDHLB:    ldhlr(&gb.CPU.Regs.B, 0),
		OpcodeLDHLC:    ldhlr(&gb.CPU.Regs.C, 0),
		OpcodeLDHLD:    ldhlr(&gb.CPU.Regs.D, 0),
		OpcodeLDHLE:    ldhlr(&gb.CPU.Regs.E, 0),
		OpcodeLDHLH:    ldhlr(&gb.CPU.Regs.H, 0),
		OpcodeLDHLL:    ldhlr(&gb.CPU.Regs.L, 0),
		OpcodeLDBCA:    ldrra(&gb.CPU.Regs.B, &gb.CPU.Regs.C),
		OpcodeLDDEA:    ldrra(&gb.CPU.Regs.D, &gb.CPU.Regs.E),
		OpcodeLDHCA:    ldhca,
		OpcodeLDHAC:    ldhac,
		OpcodeLDnnSP:   ldnnsp,
		OpcodeLDnnA:    ldnna,
		OpcodeLDAnn:    ldann,
		OpcodeCPn:      alun(func(imm Data8) { gb.CPU.Regs.SetFlags(SUB(gb.CPU.Regs.A, imm, false)) }),
		OpcodeSUBn:     alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, imm, false)) }),
		OpcodeORn:      alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(OR(gb.CPU.Regs.A, imm)) }),
		OpcodeANDn:     alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(AND(gb.CPU.Regs.A, imm)) }),
		OpcodeADDn:     alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, imm, false)) }),
		OpcodeADCn:     alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, imm, gb.CPU.Regs.GetFlagC())) }),
		OpcodeSBCn:     alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, imm, gb.CPU.Regs.GetFlagC())) }),
		OpcodeXORn:     alun(func(imm Data8) { gb.CPU.Regs.SetFlagsAndA(XOR(gb.CPU.Regs.A, imm)) }),
		OpcodeLDHnA:    ldhna,
		OpcodeLDHAn:    ldhan,
		OpcodeLDADE:    ldarr(&gb.CPU.Regs.D, &gb.CPU.Regs.E),
		OpcodeLDABC:    ldarr(&gb.CPU.Regs.B, &gb.CPU.Regs.C),
		OpcodeLDAn:     ldrn(&gb.CPU.Regs.A),
		OpcodeLDBn:     ldrn(&gb.CPU.Regs.B),
		OpcodeLDCn:     ldrn(&gb.CPU.Regs.C),
		OpcodeLDDn:     ldrn(&gb.CPU.Regs.D),
		OpcodeLDEn:     ldrn(&gb.CPU.Regs.E),
		OpcodeLDHn:     ldrn(&gb.CPU.Regs.H),
		OpcodeLDLn:     ldrn(&gb.CPU.Regs.L),
		OpcodeCB:       runCB,
		OpcodeRST0x00:  rst(0x00),
		OpcodeRST0x08:  rst(0x08),
		OpcodeRST0x10:  rst(0x10),
		OpcodeRST0x18:  rst(0x18),
		OpcodeRST0x20:  rst(0x20),
		OpcodeRST0x28:  rst(0x28),
		OpcodeRST0x30:  rst(0x30),
		OpcodeRST0x38:  rst(0x38),
		OpcodeUndefD3:  notImplemented,
		OpcodeUndefDB:  notImplemented,
		OpcodeUndefDD:  notImplemented,
		OpcodeUndefE3:  notImplemented,
		OpcodeUndefE4:  notImplemented,
		OpcodeUndefEB:  notImplemented,
		OpcodeUndefEC:  notImplemented,
		OpcodeUndefED:  notImplemented,
		OpcodeUndefF4:  notImplemented,
		OpcodeUndefFC:  notImplemented,
		OpcodeUndefFD:  notImplemented,
	}
}

func notImplemented(gb *Gameboy, e int) bool {
	panicf("not implemented opcode %v", gb.CPU.Regs.IR)
	return false
}

func checkCycleNamed(name string, e, max int) {
	if e == 0 || e > max {
		panicf("%s: %v", name, e)
	}
}

func checkCycle(e, max int) {
	if e == 0 || e > max {
		panicf("%v", e)
	}
}

func singleCycle(name string, f func(gb *Gameboy)) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycleNamed(name, e, 1)
		f(gb)
		return true
	}
}

func jphl(gb *Gameboy, e int) bool {
	checkCycle(e, 1)
	gb.CPU.SetPC(Addr(gb.CPU.GetHL()))
	return true
}

func halt(gb *Gameboy, e int) bool {
	checkCycle(e, 1)
	gb.CPU.Halted = true
	return true
}

func jre(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	// TODO: this impl is not exactly correct
	case 2:
	case 3:
		if gb.CPU.Regs.TempZ&SignBit8 != 0 {
			gb.CPU.SetPC(gb.CPU.Regs.PC - Addr(gb.CPU.Regs.TempZ.SignedAbs()))
		} else {
			gb.CPU.SetPC(gb.CPU.Regs.PC + Addr(gb.CPU.Regs.TempZ))
		}
		return true
	}
	return false
}

func jpnn(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempW = gb.Data
	case 3:
	case 4:
		gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
		return true
	}
	return false
}

func jrcce(f func() bool) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 3)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			if f() {
				gb.CPU.LastBranchResult = +1
				newPC := Data16(int16(gb.CPU.Regs.PC) + int16(int8(gb.CPU.Regs.TempZ)))
				gb.CPU.Regs.SetWZ(newPC)
			} else {
				gb.CPU.LastBranchResult = -1
				return true
			}
		case 3:
			if gb.CPU.LastBranchResult == +1 {
				gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func jpccnn(f func() bool) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 4)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempW = gb.Data
		case 3:
			if f() {
				gb.CPU.LastBranchResult = +1
				gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
			} else {
				gb.CPU.LastBranchResult = -1
				return true
			}
		case 4:
			if gb.CPU.LastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func push(msb, lsb *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 4)
		switch e {
		case 1:
			gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
			gb.WriteAddress(gb.CPU.Regs.SP)
			gb.WriteData(*msb)
		case 2:
			gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
			gb.WriteAddress(gb.CPU.Regs.SP)
			gb.WriteData(*lsb)
		case 3:
		case 4:
			return true
		}
		return false
	}
}

func pop(msb, lsb *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 3)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.SP)
			gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			gb.WriteAddress(gb.CPU.Regs.SP)
			gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
			gb.CPU.Regs.TempW = gb.Data
		case 3:
			*msb = gb.CPU.Regs.TempW
			if lsb == &gb.CPU.Regs.F {
				*lsb = gb.CPU.Regs.TempZ & 0xf0
			} else {
				*lsb = gb.CPU.Regs.TempZ
			}
			return true
		}
		return false
	}
}

func callnn(gb *Gameboy, e int) bool {
	checkCycle(e, 6)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempW = gb.Data
	case 3:
		gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	case 4:
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
		gb.WriteData(gb.CPU.Regs.PC.MSB())
	case 5:
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.WriteData(gb.CPU.Regs.PC.LSB())
	case 6:
		gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
		return true
	}
	return false
}

func callccnn(f func() bool) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 6)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempW = gb.Data
		case 3:
			if f() {
				gb.CPU.LastBranchResult = +1
				gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
			} else {
				gb.CPU.LastBranchResult = -1
				return true
			}
		case 4:
			if gb.CPU.LastBranchResult == +1 {
				gb.WriteAddress(gb.CPU.Regs.SP)
				gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
			} else {
				panicv(e)
			}
			if gb.CPU.LastBranchResult == +1 {
				gb.WriteData(gb.CPU.Regs.PC.MSB())
			} else {
				panicv(e)
			}
		case 5:
			if gb.CPU.LastBranchResult == +1 {
				gb.WriteAddress(gb.CPU.Regs.SP)
			} else {
				panicv(e)
			}
			if gb.CPU.LastBranchResult == +1 {
				gb.WriteData(gb.CPU.Regs.PC.LSB())
			} else {
				panicv(e)
			}
		case 6:
			if gb.CPU.LastBranchResult == +1 {
				gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func ret(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
		gb.CPU.Regs.TempW = gb.Data
	case 3:
		gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	case 4:
		return true
	}
	return false
}

func reti(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
		gb.CPU.Regs.TempW = gb.Data
	case 3:
		gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
		// TODO verify if this is the right cycle
		gb.Interrupts.IME = true
	case 4:
		return true
	}
	return false
}

func retcc(cond func() bool) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 5)
		switch e {
		case 1:
		case 2:
			if cond() {
				gb.CPU.LastBranchResult = +1
				gb.WriteAddress(gb.CPU.Regs.SP)
				gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
				gb.CPU.Regs.TempZ = gb.Data
			} else {
				gb.CPU.LastBranchResult = -1
				return true
			}
		case 3:
			if gb.CPU.LastBranchResult == +1 {
				gb.WriteAddress(gb.CPU.Regs.SP)
				gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
				gb.CPU.Regs.TempW = gb.Data
			} else {
				panicv(e)
			}
		case 4:
			if gb.CPU.LastBranchResult == +1 {
				gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
			} else {
				panicv(e)
			}
		case 5:
			if gb.CPU.LastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func ld(dst *Data8, src *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("LD r, r", func(gb *Gameboy) {
		*dst = *src
	})
}

func andreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("AND r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(AND(gb.CPU.Regs.A, *reg))
	})
}

func xorreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("XOR r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(XOR(gb.CPU.Regs.A, *reg))
	})
}

func orreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("OR r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(OR(gb.CPU.Regs.A, *reg))
	})
}

func addreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("ADD r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, *reg, false))
	})
}

func adcreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("ADC r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, *reg, gb.CPU.Regs.GetFlagC()))
	})
}

func inchlind(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		res := ADD(gb.CPU.Regs.TempZ, 1, false)
		// apparently doesn't set C.
		gb.CPU.Regs.SetFlagH(res.H)
		gb.CPU.Regs.SetFlagZ(res.Z())
		gb.CPU.Regs.SetFlagN(res.N)
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.CPU.Regs.TempZ = res.Value
		gb.WriteData(gb.CPU.Regs.TempZ)
	case 3:
		return true
	}
	return false
}

func dechlind(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		res := SUB(gb.CPU.Regs.TempZ, 1, false)
		// apparently doesn't set C.
		gb.CPU.Regs.SetFlagH(res.H)
		gb.CPU.Regs.SetFlagZ(res.Z())
		gb.CPU.Regs.SetFlagN(res.N)
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.CPU.Regs.TempZ = res.Value
		gb.WriteData(gb.CPU.Regs.TempZ)
	case 3:
		return true
	}
	return false
}

func aluhl(f func(v Data8)) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(Addr(gb.CPU.GetHL()))
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			f(gb.CPU.Regs.TempZ)
			return true
		}
		return false
	}
}

func subreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("SUB r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, *reg, false))
	})
}

func sbcreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("SBC r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, *reg, gb.CPU.Regs.GetFlagC()))
	})
}

func cpreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("CP r", func(gb *Gameboy) {
		gb.CPU.Regs.SetFlags(SUB(gb.CPU.Regs.A, *reg, false))
	})
}

func decreg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("DEC r", func(gb *Gameboy) {
		result := SUB(*reg, 1, false)
		*reg = result.Value
		gb.CPU.Regs.SetFlagZ(result.Z())
		gb.CPU.Regs.SetFlagH(result.H)
		gb.CPU.Regs.SetFlagN(result.N)
	})
}

func increg(reg *Data8) func(gb *Gameboy, e int) bool {
	return singleCycle("INC r", func(gb *Gameboy) {
		result := ADD(*reg, 1, false)
		*reg = result.Value
		gb.CPU.Regs.SetFlagZ(result.Z())
		gb.CPU.Regs.SetFlagH(result.H)
		gb.CPU.Regs.SetFlagN(result.N)
	})
}

func iduOp(f func()) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			f()
		case 2:
			return true
		}
		return false
	}
}

func alun(f func(imm Data8)) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			f(gb.CPU.Regs.TempZ)
			return true
		}
		return false
	}
}

func ldrn(reg *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			gb.CPU.IncPC()
			*reg = gb.CPU.Regs.TempZ
			return true
		}
		return false
	}
}

func ldrhl(reg *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(Addr(gb.CPU.GetHL()))
		case 2:
			*reg = gb.Data
			return true
		}
		return false
	}
}

func ldahlinc(gb *Gameboy, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.CPU.SetHL(gb.CPU.GetHL() + 1)
	case 2:
		gb.CPU.Regs.A = gb.Data
		return true
	}
	return false
}

func ldahldec(gb *Gameboy, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.CPU.SetHL(gb.CPU.GetHL() - 1)
	case 2:
		gb.CPU.Regs.A = gb.Data
		return true
	}
	return false
}

func ldhlr(reg *Data8, inc int) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(Addr(gb.CPU.GetHL()))
			if inc == +1 {
				gb.CPU.SetHL(gb.CPU.GetHL() + 1)
			} else if inc == -1 {
				gb.CPU.SetHL(gb.CPU.GetHL() - 1)
			}
			gb.WriteData(*reg)
		case 2:
			return true
		}
		return false
	}
}

func ldrra(msb, lsb *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(Addr(join16(*msb, *lsb)))
			gb.WriteData(gb.CPU.Regs.A)
		case 2:
			return true
		}
		return false
	}
}

func ldhca(gb *Gameboy, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.C)))
		gb.WriteData(gb.CPU.Regs.A)
	case 2:
		return true
	}
	return false
}

func ldhac(gb *Gameboy, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.C)))
		gb.CPU.Regs.A = gb.Data
	case 2:
		return true
	}
	return false
}

func ldnnsp(gb *Gameboy, e int) bool {
	checkCycle(e, 5)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempW = gb.Data
	case 3:
		gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
		gb.WriteData(gb.CPU.Regs.SP.LSB())
		gb.CPU.Regs.SetWZ(gb.CPU.Regs.GetWZ() + 1)
	case 4:
		gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
		gb.WriteData(gb.CPU.Regs.SP.MSB())
	case 5:
		return true
	}
	return false
}

func ldnna(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempW = gb.Data
	case 3:
		gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
		gb.WriteData(gb.CPU.Regs.A)
	case 4:
		return true
	}
	return false
}

func ldann(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempW = gb.Data
	case 3:
		gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
		gb.CPU.Regs.A = gb.Data
	case 4:
		return true
	}
	return false
}

func ldhna(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.TempZ)))
		gb.CPU.IncPC()
		gb.WriteData(gb.CPU.Regs.A)
	case 3:
		return true
	}
	return false
}

func ldhan(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.TempZ)))
		gb.CPU.Regs.TempZ = gb.Data
	case 3:
		gb.CPU.Regs.A = gb.CPU.Regs.TempZ
		return true
	}
	return false
}

func ldarr(msb, lsb *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			gb.WriteAddress(Addr(join16(*msb, *lsb)))
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			gb.CPU.Regs.A = gb.CPU.Regs.TempZ
			return true
		}
		return false
	}
}

func addhlrr(hi, lo *Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			result := ADD(gb.CPU.Regs.L, *lo, false)
			gb.CPU.Regs.L = result.Value
			gb.CPU.Regs.SetFlagC(result.C)
			gb.CPU.Regs.SetFlagH(result.H)
			gb.CPU.Regs.SetFlagN(result.N)
		case 2:
			result := ADD(gb.CPU.Regs.H, *hi, gb.CPU.Regs.GetFlagC())
			gb.CPU.Regs.H = result.Value
			gb.CPU.Regs.SetFlagC(result.C)
			gb.CPU.Regs.SetFlagH(result.H)
			gb.CPU.Regs.SetFlagN(result.N)
			return true
		}
		return false
	}
}

func addhlsp(gb *Gameboy, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		result := ADD(gb.CPU.Regs.L, gb.CPU.Regs.SP.LSB(), false)
		gb.CPU.Regs.L = result.Value
		gb.CPU.Regs.SetFlagC(result.C)
		gb.CPU.Regs.SetFlagH(result.H)
		gb.CPU.Regs.SetFlagN(result.N)
	case 2:
		result := ADD(gb.CPU.Regs.H, gb.CPU.Regs.SP.MSB(), gb.CPU.Regs.GetFlagC())
		gb.CPU.Regs.H = result.Value
		gb.CPU.Regs.SetFlagC(result.C)
		gb.CPU.Regs.SetFlagH(result.H)
		gb.CPU.Regs.SetFlagN(result.N)
		return true
	}
	return false
}

func addspe(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		zSign := gb.CPU.Regs.TempZ&Bit7 != 0
		result := ADD(gb.CPU.Regs.SP.LSB(), gb.CPU.Regs.TempZ, false)
		gb.CPU.Regs.TempZ = result.Value
		gb.CPU.Regs.TempW = 0
		gb.CPU.Regs.SetFlags(result)
		gb.CPU.Regs.SetFlagZ(false)
		if c := gb.CPU.Regs.GetFlagC(); c && !zSign {
			gb.CPU.Regs.TempW = 1
		} else if !c && zSign {
			gb.CPU.Regs.TempW = 0xff
		}
	case 3:
		res := gb.CPU.Regs.SP.MSB()
		if gb.CPU.Regs.TempW == 1 {
			res++
		} else if gb.CPU.Regs.TempW == 0xff {
			res--
		}
		gb.CPU.Regs.TempW = res
	case 4:
		gb.CPU.SetSP(Addr(gb.CPU.Regs.GetWZ()))
		return true
	}
	return false
}

func ldxxnn(f func(wz Data16)) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 3)
		switch e {
		case 1:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempZ = gb.Data
		case 2:
			gb.WriteAddress(gb.CPU.Regs.PC)
			gb.CPU.IncPC()
			gb.CPU.Regs.TempW = gb.Data
		case 3:
			f(join16(gb.CPU.Regs.TempW, gb.CPU.Regs.TempZ))
			return true
		}
		return false
	}
}

func ldhln(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.WriteData(gb.CPU.Regs.TempZ)
	case 3:
		return true
	}
	return false
}

func ldsphl(gb *Gameboy, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
	case 2:
		gb.CPU.Regs.SP = Addr(gb.CPU.GetHL())
		return true
	}
	return false
}

func ldhlspe(gb *Gameboy, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.Regs.TempZ = gb.Data
	case 2:
		res := ADD(gb.CPU.Regs.SP.LSB(), gb.CPU.Regs.TempZ, false)
		gb.CPU.Regs.L = res.Value
		res.Z0 = true
		gb.CPU.Regs.SetFlags(res)
	case 3:
		adj := Data8(0x00)
		if gb.CPU.Regs.TempZ&Bit7 != 0 {
			adj = 0xff
		}
		res := ADD(gb.CPU.Regs.SP.MSB(), adj, gb.CPU.Regs.GetFlagC())
		gb.CPU.Regs.H = res.Value
		return true
	}
	return false
}

func rst(vec Data8) func(gb *Gameboy, e int) bool {
	return func(gb *Gameboy, e int) bool {
		checkCycle(e, 4)
		switch e {
		case 1:
			gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
			gb.WriteAddress(gb.CPU.Regs.SP)
			gb.WriteData(gb.CPU.Regs.PC.MSB())
		case 2:
			gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
			gb.WriteAddress(gb.CPU.Regs.SP)
			gb.WriteData(gb.CPU.Regs.PC.LSB())
		case 3:
			gb.CPU.SetPC(Addr(join16(0x00, vec)))
		case 4:
			return true
		}
		return false
	}
}

func NewCBOp(v Data8) CBOp {
	return CBOp{Op: cb((v & 0xf8) >> 3), Target: CBTarget(v & 0x7)}
}

func runCB(gb *Gameboy, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		gb.WriteAddress(gb.CPU.Regs.PC)
		gb.CPU.IncPC()
		gb.CPU.CBOp = NewCBOp(gb.Data)
	case 2:
		if gb.CPU.CBOp.Target == CBTargetIndirectHL {
			gb.WriteAddress(Addr(gb.CPU.GetHL()))
		}
		var val Data8
		switch gb.CPU.CBOp.Target {
		case CBTargetB:
			val = gb.CPU.Regs.B
		case CBTargetC:
			val = gb.CPU.Regs.C
		case CBTargetD:
			val = gb.CPU.Regs.D
		case CBTargetE:
			val = gb.CPU.Regs.E
		case CBTargetH:
			val = gb.CPU.Regs.H
		case CBTargetL:
			val = gb.CPU.Regs.L
		case CBTargetIndirectHL:
			val = gb.Data
		case CBTargetA:
			val = gb.CPU.Regs.A
		default:
			panic("unknown CBOp target")
		}
		val = doCBOp(gb, val)
		switch gb.CPU.CBOp.Target {
		case CBTargetB:
			gb.CPU.Regs.B = val
		case CBTargetC:
			gb.CPU.Regs.C = val
		case CBTargetD:
			gb.CPU.Regs.D = val
		case CBTargetE:
			gb.CPU.Regs.E = val
		case CBTargetH:
			gb.CPU.Regs.H = val
		case CBTargetL:
			gb.CPU.Regs.L = val
		case CBTargetIndirectHL:
			gb.CPU.Regs.TempZ = val
			return false
		case CBTargetA:
			gb.CPU.Regs.A = val
		default:
			panic("unknown CBOp target")
		}
		return true
	case 3:
		if gb.CPU.CBOp.Target != CBTargetIndirectHL {
			panicv(e)
		}
		if gb.CPU.CBOp.Op.Is3Cycles() {
			return true
		}
		gb.WriteAddress(Addr(gb.CPU.GetHL()))
		gb.WriteData(gb.CPU.Regs.TempZ)
	case 4:
		if gb.CPU.CBOp.Target != CBTargetIndirectHL {
			panicv(e)
		}
		return true
	}
	return false
}

func doCBOp(gb *Gameboy, val Data8) Data8 {
	switch gb.CPU.CBOp.Op {
	case CbRL:
		res := RL(val, gb.CPU.Regs.GetFlagC())
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbRLC:
		res := RLC(val)
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbRR:
		res := RR(val, gb.CPU.Regs.GetFlagC())
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbRRC:
		res := RRC(val)
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbSRL:
		res := SRL(val)
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbSLA:
		res := SLA(val)
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbSRA:
		res := SRA(val)
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbSWAP:
		res := SWAP(val)
		val = res.Value
		gb.CPU.Regs.SetFlags(res)
	case CbBit0:
		cbbit(gb, val, 0x01)
	case CbBit1:
		cbbit(gb, val, 0x02)
	case CbBit2:
		cbbit(gb, val, 0x04)
	case CbBit3:
		cbbit(gb, val, 0x08)
	case CbBit4:
		cbbit(gb, val, 0x10)
	case CbBit5:
		cbbit(gb, val, 0x20)
	case CbBit6:
		cbbit(gb, val, 0x40)
	case CbBit7:
		cbbit(gb, val, 0x80)
	// RES, SET doesn't set flags apparently
	case CbRes0:
		val &= ^Data8(0x01)
	case CbRes1:
		val &= ^Data8(0x02)
	case CbRes2:
		val &= ^Data8(0x04)
	case CbRes3:
		val &= ^Data8(0x08)
	case CbRes4:
		val &= ^Data8(0x10)
	case CbRes5:
		val &= ^Data8(0x20)
	case CbRes6:
		val &= ^Data8(0x40)
	case CbRes7:
		val &= ^Data8(0x80)
	case CbSet0:
		val |= Data8(0x01)
	case CbSet1:
		val |= Data8(0x02)
	case CbSet2:
		val |= Data8(0x04)
	case CbSet3:
		val |= Data8(0x08)
	case CbSet4:
		val |= Data8(0x10)
	case CbSet5:
		val |= Data8(0x20)
	case CbSet6:
		val |= Data8(0x40)
	case CbSet7:
		val |= Data8(0x80)
	default:
		panicf("unknown op = %+v", gb.CPU.CBOp)
	}
	return val
}

func cbbit(gb *Gameboy, val, mask Data8) {
	gb.CPU.Regs.SetFlagZ(val&mask == 0)
	gb.CPU.Regs.SetFlagN(false)
	gb.CPU.Regs.SetFlagH(true)
}
