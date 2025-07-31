package model

import (
	"fmt"
	"slices"
)

type InstructionHandling func(mem []Data8, e int) bool

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

func handlers(cpu *CPU) [256]InstructionHandling {
	return [256]InstructionHandling{
		OpcodeNop:      cpu.singleCycle("NOP", func() {}),
		OpcodeLDAA:     cpu.ld(&cpu.Regs.A, &cpu.Regs.A),
		OpcodeLDAB:     cpu.ld(&cpu.Regs.A, &cpu.Regs.B),
		OpcodeLDAC:     cpu.ld(&cpu.Regs.A, &cpu.Regs.C),
		OpcodeLDAD:     cpu.ld(&cpu.Regs.A, &cpu.Regs.D),
		OpcodeLDAE:     cpu.ld(&cpu.Regs.A, &cpu.Regs.E),
		OpcodeLDAH:     cpu.ld(&cpu.Regs.A, &cpu.Regs.H),
		OpcodeLDAL:     cpu.ld(&cpu.Regs.A, &cpu.Regs.L),
		OpcodeLDAHL:    cpu.ldrhl(&cpu.Regs.A),
		OpcodeLDAHLInc: cpu.ldahlinc,
		OpcodeLDAHLDec: cpu.ldahldec,
		OpcodeLDBA:     cpu.ld(&cpu.Regs.B, &cpu.Regs.A),
		OpcodeLDBB:     cpu.ld(&cpu.Regs.B, &cpu.Regs.B),
		OpcodeLDBC:     cpu.ld(&cpu.Regs.B, &cpu.Regs.C),
		OpcodeLDBD:     cpu.ld(&cpu.Regs.B, &cpu.Regs.D),
		OpcodeLDBE:     cpu.ld(&cpu.Regs.B, &cpu.Regs.E),
		OpcodeLDBH:     cpu.ld(&cpu.Regs.B, &cpu.Regs.H),
		OpcodeLDBL:     cpu.ld(&cpu.Regs.B, &cpu.Regs.L),
		OpcodeLDBHL:    cpu.ldrhl(&cpu.Regs.B),
		OpcodeLDCA:     cpu.ld(&cpu.Regs.C, &cpu.Regs.A),
		OpcodeLDCB:     cpu.ld(&cpu.Regs.C, &cpu.Regs.B),
		OpcodeLDCC:     cpu.ld(&cpu.Regs.C, &cpu.Regs.C),
		OpcodeLDCD:     cpu.ld(&cpu.Regs.C, &cpu.Regs.D),
		OpcodeLDCE:     cpu.ld(&cpu.Regs.C, &cpu.Regs.E),
		OpcodeLDCH:     cpu.ld(&cpu.Regs.C, &cpu.Regs.H),
		OpcodeLDCL:     cpu.ld(&cpu.Regs.C, &cpu.Regs.L),
		OpcodeLDCHL:    cpu.ldrhl(&cpu.Regs.C),
		OpcodeLDDA:     cpu.ld(&cpu.Regs.D, &cpu.Regs.A),
		OpcodeLDDB:     cpu.ld(&cpu.Regs.D, &cpu.Regs.B),
		OpcodeLDDC:     cpu.ld(&cpu.Regs.D, &cpu.Regs.C),
		OpcodeLDDD:     cpu.ld(&cpu.Regs.D, &cpu.Regs.D),
		OpcodeLDDE:     cpu.ld(&cpu.Regs.D, &cpu.Regs.E),
		OpcodeLDDH:     cpu.ld(&cpu.Regs.D, &cpu.Regs.H),
		OpcodeLDDL:     cpu.ld(&cpu.Regs.D, &cpu.Regs.L),
		OpcodeLDDHL:    cpu.ldrhl(&cpu.Regs.D),
		OpcodeLDEA:     cpu.ld(&cpu.Regs.E, &cpu.Regs.A),
		OpcodeLDEB:     cpu.ld(&cpu.Regs.E, &cpu.Regs.B),
		OpcodeLDEC:     cpu.ld(&cpu.Regs.E, &cpu.Regs.C),
		OpcodeLDED:     cpu.ld(&cpu.Regs.E, &cpu.Regs.D),
		OpcodeLDEE:     cpu.ld(&cpu.Regs.E, &cpu.Regs.E),
		OpcodeLDEH:     cpu.ld(&cpu.Regs.E, &cpu.Regs.H),
		OpcodeLDEL:     cpu.ld(&cpu.Regs.E, &cpu.Regs.L),
		OpcodeLDEHL:    cpu.ldrhl(&cpu.Regs.E),
		OpcodeLDHA:     cpu.ld(&cpu.Regs.H, &cpu.Regs.A),
		OpcodeLDHB:     cpu.ld(&cpu.Regs.H, &cpu.Regs.B),
		OpcodeLDHC:     cpu.ld(&cpu.Regs.H, &cpu.Regs.C),
		OpcodeLDHD:     cpu.ld(&cpu.Regs.H, &cpu.Regs.D),
		OpcodeLDHE:     cpu.ld(&cpu.Regs.H, &cpu.Regs.E),
		OpcodeLDHH:     cpu.ld(&cpu.Regs.H, &cpu.Regs.H),
		OpcodeLDHL:     cpu.ld(&cpu.Regs.H, &cpu.Regs.L),
		OpcodeLDHHL:    cpu.ldrhl(&cpu.Regs.H),
		OpcodeLDLA:     cpu.ld(&cpu.Regs.L, &cpu.Regs.A),
		OpcodeLDLB:     cpu.ld(&cpu.Regs.L, &cpu.Regs.B),
		OpcodeLDLC:     cpu.ld(&cpu.Regs.L, &cpu.Regs.C),
		OpcodeLDLD:     cpu.ld(&cpu.Regs.L, &cpu.Regs.D),
		OpcodeLDLE:     cpu.ld(&cpu.Regs.L, &cpu.Regs.E),
		OpcodeLDLH:     cpu.ld(&cpu.Regs.L, &cpu.Regs.H),
		OpcodeLDLL:     cpu.ld(&cpu.Regs.L, &cpu.Regs.L),
		OpcodeLDLHL:    cpu.ldrhl(&cpu.Regs.L),
		OpcodeRLA:      cpu.singleCycle("RLA", func() { cpu.Regs.SetFlagsAndA(RLA(cpu.Regs.A, cpu.Regs.GetFlagC())) }),
		OpcodeRRA:      cpu.singleCycle("RRA", func() { cpu.Regs.SetFlagsAndA(RRA(cpu.Regs.A, cpu.Regs.GetFlagC())) }),
		OpcodeRLCA:     cpu.singleCycle("RLCA", func() { cpu.Regs.SetFlagsAndA(RLCA(cpu.Regs.A)) }),
		OpcodeRRCA:     cpu.singleCycle("RRCA", func() { cpu.Regs.SetFlagsAndA(RRCA(cpu.Regs.A)) }),
		OpcodeORA:      cpu.orreg(&cpu.Regs.A),
		OpcodeORB:      cpu.orreg(&cpu.Regs.B),
		OpcodeORC:      cpu.orreg(&cpu.Regs.C),
		OpcodeORD:      cpu.orreg(&cpu.Regs.D),
		OpcodeORE:      cpu.orreg(&cpu.Regs.E),
		OpcodeORH:      cpu.orreg(&cpu.Regs.H),
		OpcodeORL:      cpu.orreg(&cpu.Regs.L),
		OpcodeORHL:     cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, v)) }),
		OpcodeANDA:     cpu.andreg(&cpu.Regs.A),
		OpcodeANDB:     cpu.andreg(&cpu.Regs.B),
		OpcodeANDC:     cpu.andreg(&cpu.Regs.C),
		OpcodeANDD:     cpu.andreg(&cpu.Regs.D),
		OpcodeANDE:     cpu.andreg(&cpu.Regs.E),
		OpcodeANDH:     cpu.andreg(&cpu.Regs.H),
		OpcodeANDL:     cpu.andreg(&cpu.Regs.L),
		OpcodeANDHL:    cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, v)) }),
		OpcodeXORA:     cpu.xorreg(&cpu.Regs.A),
		OpcodeXORB:     cpu.xorreg(&cpu.Regs.B),
		OpcodeXORC:     cpu.xorreg(&cpu.Regs.C),
		OpcodeXORD:     cpu.xorreg(&cpu.Regs.D),
		OpcodeXORE:     cpu.xorreg(&cpu.Regs.E),
		OpcodeXORH:     cpu.xorreg(&cpu.Regs.H),
		OpcodeXORL:     cpu.xorreg(&cpu.Regs.L),
		OpcodeXORHL:    cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, v)) }),
		OpcodeSUBA:     cpu.subreg(&cpu.Regs.A),
		OpcodeSUBB:     cpu.subreg(&cpu.Regs.B),
		OpcodeSUBC:     cpu.subreg(&cpu.Regs.C),
		OpcodeSUBD:     cpu.subreg(&cpu.Regs.D),
		OpcodeSUBE:     cpu.subreg(&cpu.Regs.E),
		OpcodeSUBH:     cpu.subreg(&cpu.Regs.H),
		OpcodeSUBL:     cpu.subreg(&cpu.Regs.L),
		OpcodeSUBHL:    cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, v, false)) }),
		OpcodeSBCA:     cpu.sbcreg(&cpu.Regs.A),
		OpcodeSBCB:     cpu.sbcreg(&cpu.Regs.B),
		OpcodeSBCC:     cpu.sbcreg(&cpu.Regs.C),
		OpcodeSBCD:     cpu.sbcreg(&cpu.Regs.D),
		OpcodeSBCE:     cpu.sbcreg(&cpu.Regs.E),
		OpcodeSBCH:     cpu.sbcreg(&cpu.Regs.H),
		OpcodeSBCL:     cpu.sbcreg(&cpu.Regs.L),
		OpcodeSBCHL:    cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, v, cpu.Regs.GetFlagC())) }),
		OpcodeCPA:      cpu.cpreg(&cpu.Regs.A),
		OpcodeCPB:      cpu.cpreg(&cpu.Regs.B),
		OpcodeCPC:      cpu.cpreg(&cpu.Regs.C),
		OpcodeCPD:      cpu.cpreg(&cpu.Regs.D),
		OpcodeCPE:      cpu.cpreg(&cpu.Regs.E),
		OpcodeCPH:      cpu.cpreg(&cpu.Regs.H),
		OpcodeCPL:      cpu.cpreg(&cpu.Regs.L),
		OpcodeCPHL:     cpu.aluhl(func(v Data8) { cpu.Regs.SetFlags(SUB(cpu.Regs.A, v, false)) }),
		OpcodeADDA:     cpu.addreg(&cpu.Regs.A),
		OpcodeADDB:     cpu.addreg(&cpu.Regs.B),
		OpcodeADDC:     cpu.addreg(&cpu.Regs.C),
		OpcodeADDD:     cpu.addreg(&cpu.Regs.D),
		OpcodeADDE:     cpu.addreg(&cpu.Regs.E),
		OpcodeADDH:     cpu.addreg(&cpu.Regs.H),
		OpcodeADDL:     cpu.addreg(&cpu.Regs.L),
		OpcodeADDHL:    cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, v, false)) }),
		OpcodeADDSPe:   cpu.addspe,
		OpcodeADCA:     cpu.adcreg(&cpu.Regs.A),
		OpcodeADCB:     cpu.adcreg(&cpu.Regs.B),
		OpcodeADCC:     cpu.adcreg(&cpu.Regs.C),
		OpcodeADCD:     cpu.adcreg(&cpu.Regs.D),
		OpcodeADCE:     cpu.adcreg(&cpu.Regs.E),
		OpcodeADCH:     cpu.adcreg(&cpu.Regs.H),
		OpcodeADCL:     cpu.adcreg(&cpu.Regs.L),
		OpcodeADCHL:    cpu.aluhl(func(v Data8) { cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, v, cpu.Regs.GetFlagC())) }),
		OpcodeDAA: cpu.singleCycle("DAA", func() {
			cpu.Regs.SetFlagsAndA(DAA(cpu.Regs.A, cpu.Regs.GetFlagC(), cpu.Regs.GetFlagN(), cpu.Regs.GetFlagH()))
		}),
		OpcodeCPLaka2f: cpu.singleCycle("CPL", func() {
			cpu.Regs.A ^= 0xff
			cpu.Regs.SetFlagN(true)
			cpu.Regs.SetFlagH(true)
		}),
		OpcodeCCF: cpu.singleCycle("CCF", func() {
			cpu.Regs.SetFlagC(!cpu.Regs.GetFlagC())
			cpu.Regs.SetFlagN(false)
			cpu.Regs.SetFlagH(false)
		}),
		OpcodeSCF: cpu.singleCycle("SCF", func() {
			cpu.Regs.SetFlagC(true)
			cpu.Regs.SetFlagN(false)
			cpu.Regs.SetFlagH(false)
		}),
		OpcodeDECA:     cpu.decreg(&cpu.Regs.A),
		OpcodeDECB:     cpu.decreg(&cpu.Regs.B),
		OpcodeDECC:     cpu.decreg(&cpu.Regs.C),
		OpcodeDECD:     cpu.decreg(&cpu.Regs.D),
		OpcodeDECE:     cpu.decreg(&cpu.Regs.E),
		OpcodeDECH:     cpu.decreg(&cpu.Regs.H),
		OpcodeDECL:     cpu.decreg(&cpu.Regs.L),
		OpcodeINCA:     cpu.increg(&cpu.Regs.A),
		OpcodeINCB:     cpu.increg(&cpu.Regs.B),
		OpcodeINCC:     cpu.increg(&cpu.Regs.C),
		OpcodeINCD:     cpu.increg(&cpu.Regs.D),
		OpcodeINCE:     cpu.increg(&cpu.Regs.E),
		OpcodeINCH:     cpu.increg(&cpu.Regs.H),
		OpcodeINCL:     cpu.increg(&cpu.Regs.L),
		OpcodeINCHLInd: cpu.inchlind,
		OpcodeDECHLInd: cpu.dechlind,
		OpcodeDI: cpu.singleCycle("DI", func() {
			if cpu.Interrupts == nil {
				return
			}
			cpu.Interrupts.setIMENextCycle = false
			cpu.Interrupts.SetIME(false)
		}),
		OpcodeEI: cpu.singleCycle("EI", func() {
			if cpu.Interrupts == nil {
				return
			}
			cpu.Interrupts.setIMENextCycle = true
		}),
		OpcodeHALT:     cpu.halt,
		OpcodeJRe:      cpu.jre,
		OpcodeJPnn:     cpu.jpnn,
		OpcodeJPHL:     cpu.jphl,
		OpcodeJRZe:     cpu.jrcce(func() bool { return cpu.Regs.GetFlagZ() }),
		OpcodeJRCe:     cpu.jrcce(func() bool { return cpu.Regs.GetFlagC() }),
		OpcodeJRNZe:    cpu.jrcce(func() bool { return !cpu.Regs.GetFlagZ() }),
		OpcodeJRNCe:    cpu.jrcce(func() bool { return !cpu.Regs.GetFlagC() }),
		OpcodeJPCnn:    cpu.jpccnn(func() bool { return cpu.Regs.GetFlagC() }),
		OpcodeJPNCnn:   cpu.jpccnn(func() bool { return !cpu.Regs.GetFlagC() }),
		OpcodeJPZnn:    cpu.jpccnn(func() bool { return cpu.Regs.GetFlagZ() }),
		OpcodeJPNZnn:   cpu.jpccnn(func() bool { return !cpu.Regs.GetFlagZ() }),
		OpcodeINCBC:    cpu.iduOp(func() { cpu.SetBC(cpu.GetBC() + 1) }),
		OpcodeINCDE:    cpu.iduOp(func() { cpu.SetDE(cpu.GetDE() + 1) }),
		OpcodeINCHL:    cpu.iduOp(func() { cpu.SetHL(cpu.GetHL() + 1) }),
		OpcodeINCSP:    cpu.iduOp(func() { cpu.Regs.SP++ }),
		OpcodeDECBC:    cpu.iduOp(func() { cpu.SetBC(cpu.GetBC() - 1) }),
		OpcodeDECDE:    cpu.iduOp(func() { cpu.SetDE(cpu.GetDE() - 1) }),
		OpcodeDECHL:    cpu.iduOp(func() { cpu.SetHL(cpu.GetHL() - 1) }),
		OpcodeDECSP:    cpu.iduOp(func() { cpu.Regs.SP-- }),
		OpcodeCALLnn:   cpu.callnn,
		OpcodeCALLNZnn: cpu.callccnn(func() bool { return !cpu.Regs.GetFlagZ() }),
		OpcodeCALLZnn:  cpu.callccnn(func() bool { return cpu.Regs.GetFlagZ() }),
		OpcodeCALLNCnn: cpu.callccnn(func() bool { return !cpu.Regs.GetFlagC() }),
		OpcodeCALLCnn:  cpu.callccnn(func() bool { return cpu.Regs.GetFlagC() }),
		OpcodeRET:      cpu.ret,
		OpcodeRETI:     cpu.reti,
		OpcodeRETZ:     cpu.retcc(func() bool { return cpu.Regs.GetFlagZ() }),
		OpcodeRETNZ:    cpu.retcc(func() bool { return !cpu.Regs.GetFlagZ() }),
		OpcodeRETC:     cpu.retcc(func() bool { return cpu.Regs.GetFlagC() }),
		OpcodeRETNC:    cpu.retcc(func() bool { return !cpu.Regs.GetFlagC() }),
		OpcodePUSHBC:   cpu.push(&cpu.Regs.B, &cpu.Regs.C),
		OpcodePUSHDE:   cpu.push(&cpu.Regs.D, &cpu.Regs.E),
		OpcodePUSHHL:   cpu.push(&cpu.Regs.H, &cpu.Regs.L),
		OpcodePUSHAF:   cpu.push(&cpu.Regs.A, &cpu.Regs.F),
		OpcodePOPBC:    cpu.pop(&cpu.Regs.B, &cpu.Regs.C),
		OpcodePOPDE:    cpu.pop(&cpu.Regs.D, &cpu.Regs.E),
		OpcodePOPHL:    cpu.pop(&cpu.Regs.H, &cpu.Regs.L),
		OpcodePOPAF:    cpu.pop(&cpu.Regs.A, &cpu.Regs.F),
		OpcodeADDHLHL:  cpu.addhlrr(&cpu.Regs.H, &cpu.Regs.L),
		OpcodeADDHLBC:  cpu.addhlrr(&cpu.Regs.B, &cpu.Regs.C),
		OpcodeADDHLDE:  cpu.addhlrr(&cpu.Regs.D, &cpu.Regs.E),
		OpcodeADDHLSP:  cpu.addhlsp,
		OpcodeLDBCnn:   cpu.ldxxnn(func(wz Data16) { cpu.SetBC(wz) }),
		OpcodeLDDEnn:   cpu.ldxxnn(func(wz Data16) { cpu.SetDE(wz) }),
		OpcodeLDHLnn:   cpu.ldxxnn(func(wz Data16) { cpu.SetHL(wz) }),
		OpcodeLDSPnn:   cpu.ldxxnn(func(wz Data16) { cpu.SetSP(Addr(wz)) }),
		OpcodeLDHLn:    cpu.ldhln,
		OpcodeLDHLSPe:  cpu.ldhlspe,
		OpcodeLDSPHL:   cpu.ldsphl,
		OpcodeLDHLAInc: cpu.ldhlr(&cpu.Regs.A, +1),
		OpcodeLDHLADec: cpu.ldhlr(&cpu.Regs.A, -1),
		OpcodeLDHLA:    cpu.ldhlr(&cpu.Regs.A, 0),
		OpcodeLDHLB:    cpu.ldhlr(&cpu.Regs.B, 0),
		OpcodeLDHLC:    cpu.ldhlr(&cpu.Regs.C, 0),
		OpcodeLDHLD:    cpu.ldhlr(&cpu.Regs.D, 0),
		OpcodeLDHLE:    cpu.ldhlr(&cpu.Regs.E, 0),
		OpcodeLDHLH:    cpu.ldhlr(&cpu.Regs.H, 0),
		OpcodeLDHLL:    cpu.ldhlr(&cpu.Regs.L, 0),
		OpcodeLDBCA:    cpu.ldrra(&cpu.Regs.B, &cpu.Regs.C),
		OpcodeLDDEA:    cpu.ldrra(&cpu.Regs.D, &cpu.Regs.E),
		OpcodeLDHCA:    cpu.ldhca,
		OpcodeLDHAC:    cpu.ldhac,
		OpcodeLDnnSP:   cpu.ldnnsp,
		OpcodeLDnnA:    cpu.ldnna,
		OpcodeLDAnn:    cpu.ldann,
		OpcodeCPn:      cpu.alun(func(imm Data8) { cpu.Regs.SetFlags(SUB(cpu.Regs.A, imm, false)) }),
		OpcodeSUBn:     cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, imm, false)) }),
		OpcodeORn:      cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, imm)) }),
		OpcodeANDn:     cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, imm)) }),
		OpcodeADDn:     cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, imm, false)) }),
		OpcodeADCn:     cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, imm, cpu.Regs.GetFlagC())) }),
		OpcodeSBCn:     cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, imm, cpu.Regs.GetFlagC())) }),
		OpcodeXORn:     cpu.alun(func(imm Data8) { cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, imm)) }),
		OpcodeLDHnA:    cpu.ldhna,
		OpcodeLDHAn:    cpu.ldhan,
		OpcodeLDADE:    cpu.ldarr(&cpu.Regs.D, &cpu.Regs.E),
		OpcodeLDABC:    cpu.ldarr(&cpu.Regs.B, &cpu.Regs.C),
		OpcodeLDAn:     cpu.ldrn(&cpu.Regs.A),
		OpcodeLDBn:     cpu.ldrn(&cpu.Regs.B),
		OpcodeLDCn:     cpu.ldrn(&cpu.Regs.C),
		OpcodeLDDn:     cpu.ldrn(&cpu.Regs.D),
		OpcodeLDEn:     cpu.ldrn(&cpu.Regs.E),
		OpcodeLDHn:     cpu.ldrn(&cpu.Regs.H),
		OpcodeLDLn:     cpu.ldrn(&cpu.Regs.L),
		OpcodeCB:       cpu.cb,
		OpcodeRST0x00:  cpu.rst(0x00),
		OpcodeRST0x08:  cpu.rst(0x08),
		OpcodeRST0x10:  cpu.rst(0x10),
		OpcodeRST0x18:  cpu.rst(0x18),
		OpcodeRST0x20:  cpu.rst(0x20),
		OpcodeRST0x28:  cpu.rst(0x28),
		OpcodeRST0x30:  cpu.rst(0x30),
		OpcodeRST0x38:  cpu.rst(0x38),
		OpcodeUndefD3:  cpu.notImplemented,
		OpcodeUndefDB:  cpu.notImplemented,
		OpcodeUndefDD:  cpu.notImplemented,
		OpcodeUndefE3:  cpu.notImplemented,
		OpcodeUndefE4:  cpu.notImplemented,
		OpcodeUndefEB:  cpu.notImplemented,
		OpcodeUndefEC:  cpu.notImplemented,
		OpcodeUndefED:  cpu.notImplemented,
		OpcodeUndefF4:  cpu.notImplemented,
		OpcodeUndefFC:  cpu.notImplemented,
		OpcodeUndefFD:  cpu.notImplemented,
	}
}

func (cpu *CPU) notImplemented(mem []Data8, e int) bool {
	panicf("not implemented opcode %v", cpu.Regs.IR)
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

func (cpu *CPU) singleCycle(name string, f func()) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycleNamed(name, e, 1)
		f()
		return true
	}
}

func (cpu *CPU) jphl(mem []Data8, e int) bool {
	checkCycle(e, 1)
	cpu.SetPC(Addr(cpu.GetHL()))
	return true
}

func (cpu *CPU) halt(mem []Data8, e int) bool {
	checkCycle(e, 1)
	cpu.halted = true
	return true
}

func (cpu *CPU) jre(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	// TODO: this impl is not exactly correct
	case 2:
	case 3:
		if cpu.Regs.TempZ&SignBit8 != 0 {
			cpu.SetPC(cpu.Regs.PC - Addr(cpu.Regs.TempZ.SignedAbs()))
		} else {
			cpu.SetPC(cpu.Regs.PC + Addr(cpu.Regs.TempZ))
		}
		return true
	}
	return false
}

func (cpu *CPU) jpnn(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
	case 4:
		cpu.SetPC(Addr(cpu.Regs.GetWZ()))
		return true
	}
	return false
}

func (cpu *CPU) jrcce(f func() bool) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 3)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			if f() {
				cpu.lastBranchResult = +1
				newPC := Data16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.TempZ)))
				cpu.Regs.SetWZ(newPC)
			} else {
				cpu.lastBranchResult = -1
				return true
			}
		case 3:
			if cpu.lastBranchResult == +1 {
				cpu.SetPC(Addr(cpu.Regs.GetWZ()))
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func (cpu *CPU) jpccnn(f func() bool) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 4)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempW = cpu.Bus.GetData()
		case 3:
			if f() {
				cpu.lastBranchResult = +1
				cpu.SetPC(Addr(cpu.Regs.GetWZ()))
			} else {
				cpu.lastBranchResult = -1
				return true
			}
		case 4:
			if cpu.lastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func (cpu *CPU) push(msb, lsb *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 4)
		switch e {
		case 1:
			cpu.SetSP(cpu.Regs.SP - 1)
			cpu.writeAddressBus(mem, cpu.Regs.SP)
			cpu.Bus.WriteData(mem, *msb)
		case 2:
			cpu.SetSP(cpu.Regs.SP - 1)
			cpu.writeAddressBus(mem, cpu.Regs.SP)
			cpu.Bus.WriteData(mem, *lsb)
		case 3:
		case 4:
			return true
		}
		return false
	}
}

func (cpu *CPU) pop(msb, lsb *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 3)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP + 1)
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			cpu.writeAddressBus(mem, cpu.Regs.SP)
			cpu.SetSP(cpu.Regs.SP + 1)
			cpu.Regs.TempW = cpu.Bus.GetData()
		case 3:
			*msb = cpu.Regs.TempW
			if lsb == &cpu.Regs.F {
				*lsb = cpu.Regs.TempZ & 0xf0
			} else {
				*lsb = cpu.Regs.TempZ
			}
			return true
		}
		return false
	}
}

func (cpu *CPU) callnn(mem []Data8, e int) bool {
	checkCycle(e, 6)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
		cpu.SetSP(cpu.Regs.SP - 1)
	case 4:
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.SetSP(cpu.Regs.SP - 1)
		cpu.Bus.WriteData(mem, cpu.Regs.PC.MSB())
	case 5:
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.Bus.WriteData(mem, cpu.Regs.PC.LSB())
	case 6:
		cpu.SetPC(Addr(cpu.Regs.GetWZ()))
		return true
	}
	return false
}

func (cpu *CPU) callccnn(f func() bool) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 6)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempW = cpu.Bus.GetData()
		case 3:
			if f() {
				cpu.lastBranchResult = +1
				cpu.SetSP(cpu.Regs.SP - 1)
			} else {
				cpu.lastBranchResult = -1
				return true
			}
		case 4:
			if cpu.lastBranchResult == +1 {
				cpu.writeAddressBus(mem, cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP - 1)
			} else {
				panicv(e)
			}
			if cpu.lastBranchResult == +1 {
				cpu.Bus.WriteData(mem, cpu.Regs.PC.MSB())
			} else {
				panicv(e)
			}
		case 5:
			if cpu.lastBranchResult == +1 {
				cpu.writeAddressBus(mem, cpu.Regs.SP)
			} else {
				panicv(e)
			}
			if cpu.lastBranchResult == +1 {
				cpu.Bus.WriteData(mem, cpu.Regs.PC.LSB())
			} else {
				panicv(e)
			}
		case 6:
			if cpu.lastBranchResult == +1 {
				cpu.SetPC(Addr(cpu.Regs.GetWZ()))
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func (cpu *CPU) ret(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.SetSP(cpu.Regs.SP + 1)
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.SetSP(cpu.Regs.SP + 1)
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
		cpu.SetPC(Addr(cpu.Regs.GetWZ()))
	case 4:
		return true
	}
	return false
}

func (cpu *CPU) reti(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.SetSP(cpu.Regs.SP + 1)
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.SP)
		cpu.SetSP(cpu.Regs.SP + 1)
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
		cpu.SetPC(Addr(cpu.Regs.GetWZ()))
		// TODO verify if this is the right cycle
		if cpu.Interrupts != nil {
			cpu.Interrupts.IME = true
		}
	case 4:
		return true
	}
	return false
}

func (cpu *CPU) retcc(cond func() bool) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 5)
		switch e {
		case 1:
		case 2:
			if cond() {
				cpu.lastBranchResult = +1
				cpu.writeAddressBus(mem, cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP + 1)
				cpu.Regs.TempZ = cpu.Bus.GetData()
			} else {
				cpu.lastBranchResult = -1
				return true
			}
		case 3:
			if cpu.lastBranchResult == +1 {
				cpu.writeAddressBus(mem, cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP + 1)
				cpu.Regs.TempW = cpu.Bus.GetData()
			} else {
				panicv(e)
			}
		case 4:
			if cpu.lastBranchResult == +1 {
				cpu.SetPC(Addr(cpu.Regs.GetWZ()))
			} else {
				panicv(e)
			}
		case 5:
			if cpu.lastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		}
		return false
	}
}

func (cpu *CPU) ld(dst *Data8, src *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("LD r, r", func() {
		*dst = *src
	})
}

func (cpu *CPU) andreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("AND r", func() {
		cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, *reg))
	})
}

func (cpu *CPU) xorreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("XOR r", func() {
		cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, *reg))
	})
}

func (cpu *CPU) orreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("OR r", func() {
		cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, *reg))
	})
}

func (cpu *CPU) addreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("ADD r", func() {
		cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, *reg, false))
	})
}

func (cpu *CPU) adcreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("ADC r", func() {
		cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, *reg, cpu.Regs.GetFlagC()))
	})
}

func (cpu *CPU) inchlind(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		res := ADD(cpu.Regs.TempZ, 1, false)
		// apparently doesn't set C.
		cpu.Regs.SetFlagH(res.H)
		cpu.Regs.SetFlagZ(res.Z())
		cpu.Regs.SetFlagN(res.N)
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.Regs.TempZ = res.Value
		cpu.Bus.WriteData(mem, cpu.Regs.TempZ)
	case 3:
		return true
	}
	return false
}

func (cpu *CPU) dechlind(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		res := SUB(cpu.Regs.TempZ, 1, false)
		// apparently doesn't set C.
		cpu.Regs.SetFlagH(res.H)
		cpu.Regs.SetFlagZ(res.Z())
		cpu.Regs.SetFlagN(res.N)
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.Regs.TempZ = res.Value
		cpu.Bus.WriteData(mem, cpu.Regs.TempZ)
	case 3:
		return true
	}
	return false
}

func (cpu *CPU) aluhl(f func(v Data8)) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			f(cpu.Regs.TempZ)
			return true
		}
		return false
	}
}

func (cpu *CPU) subreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("SUB r", func() {
		cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, *reg, false))
	})
}

func (cpu *CPU) sbcreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("SBC r", func() {
		cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, *reg, cpu.Regs.GetFlagC()))
	})
}

func (cpu *CPU) cpreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("CP r", func() {
		cpu.Regs.SetFlags(SUB(cpu.Regs.A, *reg, false))
	})
}

func (cpu *CPU) decreg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("DEC r", func() {
		result := SUB(*reg, 1, false)
		*reg = result.Value
		cpu.Regs.SetFlagZ(result.Z())
		cpu.Regs.SetFlagH(result.H)
		cpu.Regs.SetFlagN(result.N)
	})
}

func (cpu *CPU) increg(reg *Data8) func(mem []Data8, e int) bool {
	return cpu.singleCycle("INC r", func() {
		result := ADD(*reg, 1, false)
		*reg = result.Value
		cpu.Regs.SetFlagZ(result.Z())
		cpu.Regs.SetFlagH(result.H)
		cpu.Regs.SetFlagN(result.N)
	})
}

func (cpu *CPU) iduOp(f func()) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
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

func (cpu *CPU) alun(f func(imm Data8)) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			f(cpu.Regs.TempZ)
			return true
		}
		return false
	}
}

func (cpu *CPU) ldrn(reg *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			cpu.IncPC()
			*reg = cpu.Regs.TempZ
			return true
		}
		return false
	}
}

func (cpu *CPU) ldrhl(reg *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		case 2:
			*reg = cpu.Bus.GetData()
			return true
		}
		return false
	}
}

func (cpu *CPU) ldahlinc(mem []Data8, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.SetHL(cpu.GetHL() + 1)
	case 2:
		cpu.Regs.A = cpu.Bus.GetData()
		return true
	}
	return false
}

func (cpu *CPU) ldahldec(mem []Data8, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.SetHL(cpu.GetHL() - 1)
	case 2:
		cpu.Regs.A = cpu.Bus.GetData()
		return true
	}
	return false
}

func (cpu *CPU) ldhlr(reg *Data8, inc int) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
			if inc == +1 {
				cpu.SetHL(cpu.GetHL() + 1)
			} else if inc == -1 {
				cpu.SetHL(cpu.GetHL() - 1)
			}
			cpu.Bus.WriteData(mem, *reg)
		case 2:
			return true
		}
		return false
	}
}

func (cpu *CPU) ldrra(msb, lsb *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, Addr(join16(*msb, *lsb)))
			cpu.Bus.WriteData(mem, cpu.Regs.A)
		case 2:
			return true
		}
		return false
	}
}

func (cpu *CPU) ldhca(mem []Data8, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, Addr(join16(0xff, cpu.Regs.C)))
		cpu.Bus.WriteData(mem, cpu.Regs.A)
	case 2:
		return true
	}
	return false
}

func (cpu *CPU) ldhac(mem []Data8, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, Addr(join16(0xff, cpu.Regs.C)))
		cpu.Regs.A = cpu.Bus.GetData()
	case 2:
		return true
	}
	return false
}

func (cpu *CPU) ldnnsp(mem []Data8, e int) bool {
	checkCycle(e, 5)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
		cpu.writeAddressBus(mem, Addr(cpu.Regs.GetWZ()))
		cpu.Bus.WriteData(mem, cpu.Regs.SP.LSB())
		cpu.Regs.SetWZ(cpu.Regs.GetWZ() + 1)
	case 4:
		cpu.writeAddressBus(mem, Addr(cpu.Regs.GetWZ()))
		cpu.Bus.WriteData(mem, cpu.Regs.SP.MSB())
	case 5:
		return true
	}
	return false
}

func (cpu *CPU) ldnna(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
		cpu.writeAddressBus(mem, Addr(cpu.Regs.GetWZ()))
		cpu.Bus.WriteData(mem, cpu.Regs.A)
	case 4:
		return true
	}
	return false
}

func (cpu *CPU) ldann(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempW = cpu.Bus.GetData()
	case 3:
		cpu.writeAddressBus(mem, Addr(cpu.Regs.GetWZ()))
		cpu.Regs.A = cpu.Bus.GetData()
	case 4:
		return true
	}
	return false
}

func (cpu *CPU) ldhna(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, Addr(join16(0xff, cpu.Regs.TempZ)))
		cpu.IncPC()
		cpu.Bus.WriteData(mem, cpu.Regs.A)
	case 3:
		return true
	}
	return false
}

func (cpu *CPU) ldhan(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, Addr(join16(0xff, cpu.Regs.TempZ)))
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 3:
		cpu.Regs.A = cpu.Regs.TempZ
		return true
	}
	return false
}

func (cpu *CPU) ldarr(msb, lsb *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, Addr(join16(*msb, *lsb)))
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			cpu.Regs.A = cpu.Regs.TempZ
			return true
		}
		return false
	}
}

func (cpu *CPU) addhlrr(hi, lo *Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 2)
		switch e {
		case 1:
			result := ADD(cpu.Regs.L, *lo, false)
			cpu.Regs.L = result.Value
			cpu.Regs.SetFlagC(result.C)
			cpu.Regs.SetFlagH(result.H)
			cpu.Regs.SetFlagN(result.N)
		case 2:
			result := ADD(cpu.Regs.H, *hi, cpu.Regs.GetFlagC())
			cpu.Regs.H = result.Value
			cpu.Regs.SetFlagC(result.C)
			cpu.Regs.SetFlagH(result.H)
			cpu.Regs.SetFlagN(result.N)
			return true
		}
		return false
	}
}

func (cpu *CPU) addhlsp(mem []Data8, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
		result := ADD(cpu.Regs.L, cpu.Regs.SP.LSB(), false)
		cpu.Regs.L = result.Value
		cpu.Regs.SetFlagC(result.C)
		cpu.Regs.SetFlagH(result.H)
		cpu.Regs.SetFlagN(result.N)
	case 2:
		result := ADD(cpu.Regs.H, cpu.Regs.SP.MSB(), cpu.Regs.GetFlagC())
		cpu.Regs.H = result.Value
		cpu.Regs.SetFlagC(result.C)
		cpu.Regs.SetFlagH(result.H)
		cpu.Regs.SetFlagN(result.N)
		return true
	}
	return false
}

func (cpu *CPU) addspe(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		zSign := cpu.Regs.TempZ&Bit7 != 0
		result := ADD(cpu.Regs.SP.LSB(), cpu.Regs.TempZ, false)
		cpu.Regs.TempZ = result.Value
		cpu.Regs.TempW = 0
		cpu.Regs.SetFlags(result)
		cpu.Regs.SetFlagZ(false)
		if c := cpu.Regs.GetFlagC(); c && !zSign {
			cpu.Regs.TempW = 1
		} else if !c && zSign {
			cpu.Regs.TempW = 0xff
		}
	case 3:
		res := cpu.Regs.SP.MSB()
		if cpu.Regs.TempW == 1 {
			res++
		} else if cpu.Regs.TempW == 0xff {
			res--
		}
		cpu.Regs.TempW = res
	case 4:
		cpu.SetSP(Addr(cpu.Regs.GetWZ()))
		return true
	}
	return false
}

func (cpu *CPU) ldxxnn(f func(wz Data16)) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 3)
		switch e {
		case 1:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempZ = cpu.Bus.GetData()
		case 2:
			cpu.writeAddressBus(mem, cpu.Regs.PC)
			cpu.IncPC()
			cpu.Regs.TempW = cpu.Bus.GetData()
		case 3:
			f(join16(cpu.Regs.TempW, cpu.Regs.TempZ))
			return true
		}
		return false
	}
}

func (cpu *CPU) ldhln(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.Bus.WriteData(mem, cpu.Regs.TempZ)
	case 3:
		return true
	}
	return false
}

func (cpu *CPU) ldsphl(mem []Data8, e int) bool {
	checkCycle(e, 2)
	switch e {
	case 1:
	case 2:
		cpu.Regs.SP = Addr(cpu.GetHL())
		return true
	}
	return false
}

func (cpu *CPU) ldhlspe(mem []Data8, e int) bool {
	checkCycle(e, 3)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.Regs.TempZ = cpu.Bus.GetData()
	case 2:
		res := ADD(cpu.Regs.SP.LSB(), cpu.Regs.TempZ, false)
		cpu.Regs.L = res.Value
		res.Z0 = true
		cpu.Regs.SetFlags(res)
	case 3:
		adj := Data8(0x00)
		if cpu.Regs.TempZ&Bit7 != 0 {
			adj = 0xff
		}
		res := ADD(cpu.Regs.SP.MSB(), adj, cpu.Regs.GetFlagC())
		cpu.Regs.H = res.Value
		return true
	}
	return false
}

func (cpu *CPU) rst(vec Data8) func(mem []Data8, e int) bool {
	return func(mem []Data8, e int) bool {
		checkCycle(e, 4)
		switch e {
		case 1:
			cpu.SetSP(cpu.Regs.SP - 1)
			cpu.writeAddressBus(mem, cpu.Regs.SP)
			cpu.Bus.WriteData(mem, cpu.Regs.PC.MSB())
		case 2:
			cpu.SetSP(cpu.Regs.SP - 1)
			cpu.writeAddressBus(mem, cpu.Regs.SP)
			cpu.Bus.WriteData(mem, cpu.Regs.PC.LSB())
		case 3:
			cpu.SetPC(Addr(join16(0x00, vec)))
		case 4:
			return true
		}
		return false
	}
}

func NewCBOp(v Data8) CBOp {
	return CBOp{Op: cb((v & 0xf8) >> 3), Target: CBTarget(v & 0x7)}
}

func (cpu *CPU) cb(mem []Data8, e int) bool {
	checkCycle(e, 4)
	switch e {
	case 1:
		cpu.writeAddressBus(mem, cpu.Regs.PC)
		cpu.IncPC()
		cpu.CBOp = NewCBOp(cpu.Bus.GetData())
	case 2:
		if cpu.CBOp.Target == CBTargetIndirectHL {
			cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		}
		var val Data8
		switch cpu.CBOp.Target {
		case CBTargetB:
			val = cpu.Regs.B
		case CBTargetC:
			val = cpu.Regs.C
		case CBTargetD:
			val = cpu.Regs.D
		case CBTargetE:
			val = cpu.Regs.E
		case CBTargetH:
			val = cpu.Regs.H
		case CBTargetL:
			val = cpu.Regs.L
		case CBTargetIndirectHL:
			val = cpu.Bus.GetData()
		case CBTargetA:
			val = cpu.Regs.A
		default:
			panic("unknown CBOp target")
		}
		val = cpu.doCBOp(val)
		switch cpu.CBOp.Target {
		case CBTargetB:
			cpu.Regs.B = val
		case CBTargetC:
			cpu.Regs.C = val
		case CBTargetD:
			cpu.Regs.D = val
		case CBTargetE:
			cpu.Regs.E = val
		case CBTargetH:
			cpu.Regs.H = val
		case CBTargetL:
			cpu.Regs.L = val
		case CBTargetIndirectHL:
			cpu.Regs.TempZ = val
			return false
		case CBTargetA:
			cpu.Regs.A = val
		default:
			panic("unknown CBOp target")
		}
		return true
	case 3:
		if cpu.CBOp.Target != CBTargetIndirectHL {
			panicv(e)
		}
		if cpu.CBOp.Op.Is3Cycles() {
			return true
		}
		cpu.writeAddressBus(mem, Addr(cpu.GetHL()))
		cpu.Bus.WriteData(mem, cpu.Regs.TempZ)
	case 4:
		if cpu.CBOp.Target != CBTargetIndirectHL {
			panicv(e)
		}
		return true
	}
	return false
}

func (cpu *CPU) doCBOp(val Data8) Data8 {
	switch cpu.CBOp.Op {
	case CbRL:
		res := RL(val, cpu.Regs.GetFlagC())
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbRLC:
		res := RLC(val)
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbRR:
		res := RR(val, cpu.Regs.GetFlagC())
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbRRC:
		res := RRC(val)
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbSRL:
		res := SRL(val)
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbSLA:
		res := SLA(val)
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbSRA:
		res := SRA(val)
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbSWAP:
		res := SWAP(val)
		val = res.Value
		cpu.Regs.SetFlags(res)
	case CbBit0:
		cpu.cbbit(val, 0x01)
	case CbBit1:
		cpu.cbbit(val, 0x02)
	case CbBit2:
		cpu.cbbit(val, 0x04)
	case CbBit3:
		cpu.cbbit(val, 0x08)
	case CbBit4:
		cpu.cbbit(val, 0x10)
	case CbBit5:
		cpu.cbbit(val, 0x20)
	case CbBit6:
		cpu.cbbit(val, 0x40)
	case CbBit7:
		cpu.cbbit(val, 0x80)
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
		panicf("unknown op = %+v", cpu.CBOp)
	}
	return val
}

func (cpu *CPU) cbbit(val, mask Data8) {
	cpu.Regs.SetFlagZ(val&mask == 0)
	cpu.Regs.SetFlagN(false)
	cpu.Regs.SetFlagH(true)
}
