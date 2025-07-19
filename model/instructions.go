package model

//go:generate go-enum --marshal --flag --values --nocomments

// ENUM(
// Nop      = 0x00,
// LDBCnn   = 0x01,
// INCBC    = 0x03,
// INCB     = 0x04,
// DECB     = 0x05,
// LDBn     = 0x06,
// ADDHLBC  = 0x09,
// DECBC    = 0x0b,
// INCC     = 0x0c,
// DECC     = 0x0d,
// LDCn     = 0x0e,
// LDDEnn   = 0x11,
// INCDE    = 0x13,
// INCD     = 0x14,
// DECD     = 0x15,
// LDDn     = 0x16,
// RLA      = 0x17,
// JRe      = 0x18,
// ADDHLDE  = 0x19,
// LDADE    = 0x1a,
// DECDE    = 0x1b,
// INCE     = 0x1c,
// DECE     = 0x1d,
// LDEn     = 0x1e,
// JRNZe    = 0x20,
// LDHLnn   = 0x21,
// LDHLAInc = 0x22,
// INCHL    = 0x23,
// INCH     = 0x24,
// DECH     = 0x25,
// LDHn     = 0x26,
// JRZe     = 0x28,
// ADDHLHL  = 0x29,
// DECHL    = 0x2b,
// INCL     = 0x2C,
// DECL     = 0x2D,
// LDLn     = 0x2e,
// JRNCe    = 0x30,
// LDSPnn   = 0x31,
// LDHLADec = 0x32,
// INCSP    = 0x33,
// JRCe     = 0x38,
// ADDHLSP  = 0x39,
// DECSP    = 0x3b,
// INCA     = 0x3c,
// DECA     = 0x3d,
// LDAn     = 0x3e,
// LDBB     = 0x40,
// LDBC     = 0x41,
// LDBD     = 0x42,
// LDBE     = 0x43,
// LDBH     = 0x44,
// LDBL     = 0x45,
// LDBHL    = 0x46,
// LDBA     = 0x47,
// LDCB     = 0x48,
// LDCC     = 0x49,
// LDCD     = 0x4a,
// LDCE     = 0x4b,
// LDCH     = 0x4c,
// LDCL     = 0x4d,
// LDCHL    = 0x4e,
// LDCA     = 0x4f,
// LDDB     = 0x50,
// LDDC     = 0x51,
// LDDD     = 0x52,
// LDDE     = 0x53,
// LDDH     = 0x54,
// LDDL     = 0x55,
// LDDHL    = 0x56,
// LDDA     = 0x57,
// LDEB     = 0x58,
// LDEC     = 0x59,
// LDED     = 0x5a,
// LDEE     = 0x5b,
// LDEH     = 0x5c,
// LDEL     = 0x5d,
// LDEHL    = 0x5e,
// LDEA     = 0x5f,
// LDHB     = 0x60,
// LDHC     = 0x61,
// LDHD     = 0x62,
// LDHE     = 0x63,
// LDHH     = 0x64,
// LDHL     = 0x65,
// LDHHL    = 0x66,
// LDHA     = 0x67,
// LDLB     = 0x68,
// LDLC     = 0x69,
// LDLD     = 0x6a,
// LDLE     = 0x6b,
// LDLH     = 0x6c,
// LDLL     = 0x6d,
// LDLHL    = 0x6e,
// LDLA     = 0x6f,
// LDHLB    = 0x70,
// LDHLC    = 0x71,
// LDHLD    = 0x72,
// LDHLE    = 0x73,
// LDHLH    = 0x74,
// LDHLL    = 0x75,
// HALT     = 0x76,
// LDHLA    = 0x77,
// LDAB     = 0x78,
// LDAC     = 0x79,
// LDAD     = 0x7a,
// LDAE     = 0x7b,
// LDAH     = 0x7c,
// LDAL     = 0x7d,
// LDAHL    = 0x7e,
// LDAA     = 0x7f,
// ADDB     = 0x80,
// ADDC     = 0x81,
// ADDD     = 0x82,
// ADDE     = 0x83,
// ADDH     = 0x84,
// ADDL     = 0x85,
// ADDHL    = 0x86,
// ADDA     = 0x87,
// ADCB     = 0x88,
// ADCC     = 0x89,
// ADCD     = 0x8A,
// ADCE     = 0x8B,
// ADCH     = 0x8c,
// ADCL     = 0x8d,
// ADCHL    = 0x8e,
// ADCA     = 0x8f,
// SUBB     = 0x90,
// SUBC     = 0x91,
// SUBD     = 0x92,
// SUBE     = 0x93,
// SUBH     = 0x94,
// SUBL     = 0x95,
// SUBHL    = 0x96,
// SUBA     = 0x97,
// SBCB     = 0x98,
// SBCC     = 0x99,
// SBCD     = 0x9A,
// SBCE     = 0x9B,
// SBCH     = 0x9c,
// SBCL     = 0x9d,
// SBCHL    = 0x9e,
// SBCA     = 0x9f,
// ANDB     = 0xA0,
// ANDC     = 0xA1,
// ANDD     = 0xA2,
// ANDE     = 0xA3,
// ANDH     = 0xA4,
// ANDL     = 0xA5,
// ANDHL    = 0xA6,
// ANDA     = 0xA7,
// XORB     = 0xA8,
// XORC     = 0xA9,
// XORD     = 0xAA,
// XORE     = 0xAB,
// XORH     = 0xAc,
// XORL     = 0xAd,
// XORHL    = 0xAe,
// XORA     = 0xAf,
// ORB      = 0xB0,
// ORC      = 0xB1,
// ORD      = 0xB2,
// ORE      = 0xB3,
// ORH      = 0xB4,
// ORL      = 0xB5,
// ORHL     = 0xB6,
// ORA      = 0xB7,
// CPB      = 0xB8,
// CPC      = 0xB9,
// CPD      = 0xBA,
// CPE      = 0xBB,
// CPH      = 0xBc,
// CPL      = 0xBd,
// CPHL     = 0xBe,
// CPA      = 0xBf,
// RETNZ    = 0xC0,
// POPBC    = 0xC1,
// JPNZnn   = 0xC2,
// JPnn     = 0xC3,
// PUSHBC   = 0xC5,
// ADDn     = 0xC6,
// RETZ     = 0xC8,
// RET      = 0xC9,
// JPZnn    = 0xCA,
// CB       = 0xCB,
// CALLnn   = 0xCD,
// ADCn     = 0xCE,
// RETNC    = 0xD0,
// RETC     = 0xD8,
// RETI     = 0xD9,
// JPCnn    = 0xDA,
// JPNCnn   = 0xD2,
// SUBn     = 0xD6
// SBCn     = 0xDE,
// LDHnA    = 0xE0,
// LDHCA    = 0xE2,
// ANDn     = 0xE6,
// LDnnA    = 0xEA,
// XORn     = 0xEE,
// LDHAn    = 0xF0,
// LDAnn    = 0xFA,
// DI       = 0xF3,
// ORn      = 0xF6,
// EI       = 0xFB,
// CPn      = 0xFE,
// )
type Opcode uint8

// ENUM(
// RLC,
// RRC,
// RL,
// RR,
// SLA,
// SRA,
// SWAP,
// SRL,
// Bit0,
// Bit1,
// Bit2,
// Bit3,
// Bit4,
// Bit5,
// Bit6,
// Bit7,
// Res0,
// Res1,
// Res2,
// Res3,
// Res4,
// Res5,
// Res6,
// Res7,
// Set0,
// Set1,
// Set2,
// Set3,
// Set4,
// Set5,
// Set6,
// Set7,
// )
type cb uint8

// ENUM(B, C, D, E, H, L, IndirectHL, A)
type CBTarget uint8

type CBOp struct {
	Op     cb
	Target CBTarget
}

type edge struct {
	Cycle   int
	Falling bool
}

type InstructionHandling func(e edge) bool

var instSize = [256]uint16{
	OpcodeLDSPnn: 3,
	OpcodeLDHLnn: 3,
	OpcodeLDBCnn: 3,
	OpcodeLDDEnn: 3,
	OpcodeLDAnn:  3,

	OpcodeCALLnn: 3,
	OpcodeJPnn:   3,
	OpcodeJPCnn:  3,
	OpcodeJPNCnn: 3,
	OpcodeJPZnn:  3,
	OpcodeJPNZnn: 3,
	OpcodeJRe:    2,
	OpcodeJRNZe:  2,
	OpcodeJRZe:   2,

	OpcodeLDADE: 1,
	OpcodeLDAn:  2,
	OpcodeLDBn:  2,
	OpcodeLDCn:  2,
	OpcodeLDDn:  2,
	OpcodeLDEn:  2,
	OpcodeLDHn:  2,
	OpcodeLDLn:  2,
	OpcodeLDHnA: 2,
	OpcodeLDHAn: 2,

	OpcodeXORA:  1,
	OpcodeXORB:  1,
	OpcodeXORC:  1,
	OpcodeXORD:  1,
	OpcodeXORE:  1,
	OpcodeXORH:  1,
	OpcodeXORL:  1,
	OpcodeXORHL: 1,

	OpcodeORA:  1,
	OpcodeORB:  1,
	OpcodeORC:  1,
	OpcodeORD:  1,
	OpcodeORE:  1,
	OpcodeORH:  1,
	OpcodeORL:  1,
	OpcodeORHL: 1,

	OpcodeANDA:  1,
	OpcodeANDB:  1,
	OpcodeANDC:  1,
	OpcodeANDD:  1,
	OpcodeANDE:  1,
	OpcodeANDH:  1,
	OpcodeANDL:  1,
	OpcodeANDHL: 1,

	OpcodeLDHLA: 1,
	OpcodeLDHLB: 1,
	OpcodeLDHLC: 1,
	OpcodeLDHLD: 1,
	OpcodeLDHLE: 1,
	OpcodeLDHLH: 1,
	OpcodeLDHLL: 1,

	OpcodeLDHLADec: 1,
	OpcodeLDHLAInc: 1,

	OpcodeCB: 2,

	OpcodeLDHCA: 1,

	OpcodeCPA:  1,
	OpcodeCPB:  1,
	OpcodeCPC:  1,
	OpcodeCPD:  1,
	OpcodeCPE:  1,
	OpcodeCPH:  1,
	OpcodeCPL:  1,
	OpcodeCPHL: 1,

	OpcodeSUBA:  1,
	OpcodeSUBB:  1,
	OpcodeSUBC:  1,
	OpcodeSUBD:  1,
	OpcodeSUBE:  1,
	OpcodeSUBH:  1,
	OpcodeSUBL:  1,
	OpcodeSUBHL: 1,

	OpcodeSBCA:  1,
	OpcodeSBCB:  1,
	OpcodeSBCC:  1,
	OpcodeSBCD:  1,
	OpcodeSBCE:  1,
	OpcodeSBCH:  1,
	OpcodeSBCL:  1,
	OpcodeSBCHL: 1,

	OpcodeADDA:  1,
	OpcodeADDB:  1,
	OpcodeADDC:  1,
	OpcodeADDD:  1,
	OpcodeADDE:  1,
	OpcodeADDH:  1,
	OpcodeADDL:  1,
	OpcodeADDHL: 1,

	OpcodeADCA:  1,
	OpcodeADCB:  1,
	OpcodeADCC:  1,
	OpcodeADCD:  1,
	OpcodeADCE:  1,
	OpcodeADCH:  1,
	OpcodeADCL:  1,
	OpcodeADCHL: 1,

	OpcodeINCA:  1,
	OpcodeINCB:  1,
	OpcodeINCC:  1,
	OpcodeINCD:  1,
	OpcodeINCE:  1,
	OpcodeINCH:  1,
	OpcodeINCL:  1,
	OpcodeINCHL: 1,
	OpcodeINCBC: 1,
	OpcodeINCDE: 1,
	OpcodeINCSP: 1,

	OpcodeDECA:  1,
	OpcodeDECB:  1,
	OpcodeDECC:  1,
	OpcodeDECD:  1,
	OpcodeDECE:  1,
	OpcodeDECH:  1,
	OpcodeDECL:  1,
	OpcodeDECHL: 1,
	OpcodeDECBC: 1,
	OpcodeDECDE: 1,
	OpcodeDECSP: 1,

	OpcodeADDHLBC: 1,
	OpcodeADDHLDE: 1,
	OpcodeADDHLHL: 1,
	OpcodeADDHLSP: 1,

	OpcodeLDAA:  1,
	OpcodeLDAB:  1,
	OpcodeLDAC:  1,
	OpcodeLDAD:  1,
	OpcodeLDAE:  1,
	OpcodeLDAH:  1,
	OpcodeLDAL:  1,
	OpcodeLDAHL: 1,
	OpcodeLDBA:  1,
	OpcodeLDBB:  1,
	OpcodeLDBC:  1,
	OpcodeLDBD:  1,
	OpcodeLDBE:  1,
	OpcodeLDBH:  1,
	OpcodeLDBL:  1,
	OpcodeLDBHL: 1,
	OpcodeLDCA:  1,
	OpcodeLDCB:  1,
	OpcodeLDCC:  1,
	OpcodeLDCD:  1,
	OpcodeLDCE:  1,
	OpcodeLDCH:  1,
	OpcodeLDCL:  1,
	OpcodeLDCHL: 1,
	OpcodeLDDA:  1,
	OpcodeLDDB:  1,
	OpcodeLDDC:  1,
	OpcodeLDDD:  1,
	OpcodeLDDE:  1,
	OpcodeLDDH:  1,
	OpcodeLDDL:  1,
	OpcodeLDDHL: 1,
	OpcodeLDEA:  1,
	OpcodeLDEB:  1,
	OpcodeLDEC:  1,
	OpcodeLDED:  1,
	OpcodeLDEE:  1,
	OpcodeLDEH:  1,
	OpcodeLDEL:  1,
	OpcodeLDEHL: 1,
	OpcodeLDHA:  1,
	OpcodeLDHB:  1,
	OpcodeLDHC:  1,
	OpcodeLDHD:  1,
	OpcodeLDHE:  1,
	OpcodeLDHH:  1,
	OpcodeLDHL:  1,
	OpcodeLDHHL: 1,
	OpcodeLDLA:  1,
	OpcodeLDLB:  1,
	OpcodeLDLC:  1,
	OpcodeLDLD:  1,
	OpcodeLDLE:  1,
	OpcodeLDLH:  1,
	OpcodeLDLL:  1,
	OpcodeLDLHL: 1,

	OpcodePUSHBC: 1,
	OpcodePOPBC:  1,
	OpcodeRLA:    1,
	OpcodeRET:    1,
	OpcodeRETZ:   1,
	OpcodeRETNZ:  1,
	OpcodeRETC:   1,
	OpcodeRETNC:  1,

	OpcodeADDn: 2,
	OpcodeSUBn: 2,
	OpcodeANDn: 2,
	OpcodeORn:  2,
	OpcodeADCn: 2,
	OpcodeSBCn: 2,
	OpcodeXORn: 2,
	OpcodeCPn:  2,

	OpcodeLDnnA: 3,
}

func handlers(cpu *CPU) [256]InstructionHandling {
	return [256]InstructionHandling{
		OpcodeNop: cpu.singleCycle(func() {
		}),
		OpcodeLDAA:  cpu.ld(&cpu.Regs.A, &cpu.Regs.A),
		OpcodeLDAB:  cpu.ld(&cpu.Regs.A, &cpu.Regs.B),
		OpcodeLDAC:  cpu.ld(&cpu.Regs.A, &cpu.Regs.C),
		OpcodeLDAD:  cpu.ld(&cpu.Regs.A, &cpu.Regs.D),
		OpcodeLDAE:  cpu.ld(&cpu.Regs.A, &cpu.Regs.E),
		OpcodeLDAH:  cpu.ld(&cpu.Regs.A, &cpu.Regs.H),
		OpcodeLDAL:  cpu.ld(&cpu.Regs.A, &cpu.Regs.L),
		OpcodeLDAHL: cpu.ldrhl(&cpu.Regs.A),
		OpcodeLDBA:  cpu.ld(&cpu.Regs.B, &cpu.Regs.A),
		OpcodeLDBB:  cpu.ld(&cpu.Regs.B, &cpu.Regs.B),
		OpcodeLDBC:  cpu.ld(&cpu.Regs.B, &cpu.Regs.C),
		OpcodeLDBD:  cpu.ld(&cpu.Regs.B, &cpu.Regs.D),
		OpcodeLDBE:  cpu.ld(&cpu.Regs.B, &cpu.Regs.E),
		OpcodeLDBH:  cpu.ld(&cpu.Regs.B, &cpu.Regs.H),
		OpcodeLDBL:  cpu.ld(&cpu.Regs.B, &cpu.Regs.L),
		OpcodeLDBHL: cpu.ldrhl(&cpu.Regs.B),
		OpcodeLDCA:  cpu.ld(&cpu.Regs.C, &cpu.Regs.A),
		OpcodeLDCB:  cpu.ld(&cpu.Regs.C, &cpu.Regs.B),
		OpcodeLDCC:  cpu.ld(&cpu.Regs.C, &cpu.Regs.C),
		OpcodeLDCD:  cpu.ld(&cpu.Regs.C, &cpu.Regs.D),
		OpcodeLDCE:  cpu.ld(&cpu.Regs.C, &cpu.Regs.E),
		OpcodeLDCH:  cpu.ld(&cpu.Regs.C, &cpu.Regs.H),
		OpcodeLDCL:  cpu.ld(&cpu.Regs.C, &cpu.Regs.L),
		OpcodeLDCHL: cpu.ldrhl(&cpu.Regs.C),
		OpcodeLDDA:  cpu.ld(&cpu.Regs.D, &cpu.Regs.A),
		OpcodeLDDB:  cpu.ld(&cpu.Regs.D, &cpu.Regs.B),
		OpcodeLDDC:  cpu.ld(&cpu.Regs.D, &cpu.Regs.C),
		OpcodeLDDD:  cpu.ld(&cpu.Regs.D, &cpu.Regs.D),
		OpcodeLDDE:  cpu.ld(&cpu.Regs.D, &cpu.Regs.E),
		OpcodeLDDH:  cpu.ld(&cpu.Regs.D, &cpu.Regs.H),
		OpcodeLDDL:  cpu.ld(&cpu.Regs.D, &cpu.Regs.L),
		OpcodeLDDHL: cpu.ldrhl(&cpu.Regs.D),
		OpcodeLDEA:  cpu.ld(&cpu.Regs.E, &cpu.Regs.A),
		OpcodeLDEB:  cpu.ld(&cpu.Regs.E, &cpu.Regs.B),
		OpcodeLDEC:  cpu.ld(&cpu.Regs.E, &cpu.Regs.C),
		OpcodeLDED:  cpu.ld(&cpu.Regs.E, &cpu.Regs.D),
		OpcodeLDEE:  cpu.ld(&cpu.Regs.E, &cpu.Regs.E),
		OpcodeLDEH:  cpu.ld(&cpu.Regs.E, &cpu.Regs.H),
		OpcodeLDEL:  cpu.ld(&cpu.Regs.E, &cpu.Regs.L),
		OpcodeLDEHL: cpu.ldrhl(&cpu.Regs.E),
		OpcodeLDHA:  cpu.ld(&cpu.Regs.H, &cpu.Regs.A),
		OpcodeLDHB:  cpu.ld(&cpu.Regs.H, &cpu.Regs.B),
		OpcodeLDHC:  cpu.ld(&cpu.Regs.H, &cpu.Regs.C),
		OpcodeLDHD:  cpu.ld(&cpu.Regs.H, &cpu.Regs.D),
		OpcodeLDHE:  cpu.ld(&cpu.Regs.H, &cpu.Regs.E),
		OpcodeLDHH:  cpu.ld(&cpu.Regs.H, &cpu.Regs.H),
		OpcodeLDHL:  cpu.ld(&cpu.Regs.H, &cpu.Regs.L),
		OpcodeLDHHL: cpu.ldrhl(&cpu.Regs.H),
		OpcodeLDLA:  cpu.ld(&cpu.Regs.L, &cpu.Regs.A),
		OpcodeLDLB:  cpu.ld(&cpu.Regs.L, &cpu.Regs.B),
		OpcodeLDLC:  cpu.ld(&cpu.Regs.L, &cpu.Regs.C),
		OpcodeLDLD:  cpu.ld(&cpu.Regs.L, &cpu.Regs.D),
		OpcodeLDLE:  cpu.ld(&cpu.Regs.L, &cpu.Regs.E),
		OpcodeLDLH:  cpu.ld(&cpu.Regs.L, &cpu.Regs.H),
		OpcodeLDLL:  cpu.ld(&cpu.Regs.L, &cpu.Regs.L),
		OpcodeLDLHL: cpu.ldrhl(&cpu.Regs.L),
		OpcodeRLA: cpu.singleCycle(func() {
			a := cpu.Regs.A
			bit7 := a & 0x80
			a <<= 1
			if cpu.Regs.GetFlagC() {
				a |= 0x01
			}
			cpu.Regs.SetFlagZ(false)
			cpu.Regs.SetFlagN(false)
			cpu.Regs.SetFlagH(false)
			cpu.Regs.SetFlagC(bit7 != 0)
			cpu.Regs.A = a
		}),
		OpcodeORA:  cpu.orreg(&cpu.Regs.A),
		OpcodeORB:  cpu.orreg(&cpu.Regs.B),
		OpcodeORC:  cpu.orreg(&cpu.Regs.C),
		OpcodeORD:  cpu.orreg(&cpu.Regs.D),
		OpcodeORE:  cpu.orreg(&cpu.Regs.E),
		OpcodeORH:  cpu.orreg(&cpu.Regs.H),
		OpcodeORL:  cpu.orreg(&cpu.Regs.L),
		OpcodeANDA: cpu.andreg(&cpu.Regs.A),
		OpcodeANDB: cpu.andreg(&cpu.Regs.B),
		OpcodeANDC: cpu.andreg(&cpu.Regs.C),
		OpcodeANDD: cpu.andreg(&cpu.Regs.D),
		OpcodeANDE: cpu.andreg(&cpu.Regs.E),
		OpcodeANDH: cpu.andreg(&cpu.Regs.H),
		OpcodeANDL: cpu.andreg(&cpu.Regs.L),
		OpcodeXORA: cpu.xorreg(&cpu.Regs.A),
		OpcodeXORB: cpu.xorreg(&cpu.Regs.B),
		OpcodeXORC: cpu.xorreg(&cpu.Regs.C),
		OpcodeXORD: cpu.xorreg(&cpu.Regs.D),
		OpcodeXORE: cpu.xorreg(&cpu.Regs.E),
		OpcodeXORH: cpu.xorreg(&cpu.Regs.H),
		OpcodeXORL: cpu.xorreg(&cpu.Regs.L),
		OpcodeSUBA: cpu.subreg(&cpu.Regs.A),
		OpcodeSUBB: cpu.subreg(&cpu.Regs.B),
		OpcodeSUBC: cpu.subreg(&cpu.Regs.C),
		OpcodeSUBD: cpu.subreg(&cpu.Regs.D),
		OpcodeSUBE: cpu.subreg(&cpu.Regs.E),
		OpcodeSUBH: cpu.subreg(&cpu.Regs.H),
		OpcodeSUBL: cpu.subreg(&cpu.Regs.L),
		OpcodeCPA:  cpu.cpreg(&cpu.Regs.A),
		OpcodeCPB:  cpu.cpreg(&cpu.Regs.B),
		OpcodeCPC:  cpu.cpreg(&cpu.Regs.C),
		OpcodeCPD:  cpu.cpreg(&cpu.Regs.D),
		OpcodeCPE:  cpu.cpreg(&cpu.Regs.E),
		OpcodeCPH:  cpu.cpreg(&cpu.Regs.H),
		OpcodeCPL:  cpu.cpreg(&cpu.Regs.L),
		OpcodeADDA: cpu.addreg(&cpu.Regs.A),
		OpcodeADDB: cpu.addreg(&cpu.Regs.B),
		OpcodeADDC: cpu.addreg(&cpu.Regs.C),
		OpcodeADDD: cpu.addreg(&cpu.Regs.D),
		OpcodeADDE: cpu.addreg(&cpu.Regs.E),
		OpcodeADDH: cpu.addreg(&cpu.Regs.H),
		OpcodeADDL: cpu.addreg(&cpu.Regs.L),
		OpcodeADCA: cpu.adcreg(&cpu.Regs.A),
		OpcodeADCB: cpu.adcreg(&cpu.Regs.B),
		OpcodeADCC: cpu.adcreg(&cpu.Regs.C),
		OpcodeADCD: cpu.adcreg(&cpu.Regs.D),
		OpcodeADCE: cpu.adcreg(&cpu.Regs.E),
		OpcodeADCH: cpu.adcreg(&cpu.Regs.H),
		OpcodeADCL: cpu.adcreg(&cpu.Regs.L),
		OpcodeDECA: cpu.decreg(&cpu.Regs.A),
		OpcodeDECB: cpu.decreg(&cpu.Regs.B),
		OpcodeDECC: cpu.decreg(&cpu.Regs.C),
		OpcodeDECD: cpu.decreg(&cpu.Regs.D),
		OpcodeDECE: cpu.decreg(&cpu.Regs.E),
		OpcodeDECH: cpu.decreg(&cpu.Regs.H),
		OpcodeDECL: cpu.decreg(&cpu.Regs.L),
		OpcodeINCA: cpu.increg(&cpu.Regs.A),
		OpcodeINCB: cpu.increg(&cpu.Regs.B),
		OpcodeINCC: cpu.increg(&cpu.Regs.C),
		OpcodeINCD: cpu.increg(&cpu.Regs.D),
		OpcodeINCE: cpu.increg(&cpu.Regs.E),
		OpcodeINCH: cpu.increg(&cpu.Regs.H),
		OpcodeINCL: cpu.increg(&cpu.Regs.L),
		OpcodeDI: cpu.singleCycle(func() {
			cpu.Interrupts.setIMENextCycle = false
			cpu.Interrupts.IME = false
		}),
		OpcodeEI: cpu.singleCycle(func() {
			cpu.Interrupts.setIMENextCycle = true
		}),
		OpcodeJRe: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			// TODO: this impl is not exactly correct
			case edge{2, false}:
			case edge{2, true}:
			case edge{3, false}:
				cpu.SetPC(uint16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.TempZ))))
				return true
			case edge{3, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeJPnn: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{2, true}:
				cpu.Regs.TempW = cpu.Bus.Data
			case edge{3, false}:
			case edge{3, true}:
			case edge{4, false}:
				cpu.SetPC(cpu.Regs.GetWZ())
				return true
			case edge{4, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeJRZe:   cpu.jrcce(func() bool { return cpu.Regs.GetFlagZ() }),
		OpcodeJRCe:   cpu.jrcce(func() bool { return cpu.Regs.GetFlagC() }),
		OpcodeJRNZe:  cpu.jrcce(func() bool { return !cpu.Regs.GetFlagZ() }),
		OpcodeJRNCe:  cpu.jrcce(func() bool { return !cpu.Regs.GetFlagC() }),
		OpcodeJPCnn:  cpu.jpccnn(func() bool { return cpu.Regs.GetFlagC() }),
		OpcodeJPNCnn: cpu.jpccnn(func() bool { return !cpu.Regs.GetFlagC() }),
		OpcodeJPZnn:  cpu.jpccnn(func() bool { return cpu.Regs.GetFlagZ() }),
		OpcodeJPNZnn: cpu.jpccnn(func() bool { return !cpu.Regs.GetFlagZ() }),
		OpcodeINCBC:  cpu.iduOp(func() { cpu.SetBC(cpu.GetBC() + 1) }),
		OpcodeINCDE:  cpu.iduOp(func() { cpu.SetDE(cpu.GetDE() + 1) }),
		OpcodeINCHL:  cpu.iduOp(func() { cpu.SetHL(cpu.GetHL() + 1) }),
		OpcodeINCSP:  cpu.iduOp(func() { cpu.Regs.SP++ }),
		OpcodeDECBC:  cpu.iduOp(func() { cpu.SetBC(cpu.GetBC() - 1) }),
		OpcodeDECDE:  cpu.iduOp(func() { cpu.SetDE(cpu.GetDE() - 1) }),
		OpcodeDECHL:  cpu.iduOp(func() { cpu.SetHL(cpu.GetHL() - 1) }),
		OpcodeDECSP:  cpu.iduOp(func() { cpu.Regs.SP-- }),
		OpcodeCALLnn: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{2, true}:
				cpu.Regs.TempW = cpu.Bus.Data
			case edge{3, false}:
				cpu.SetSP(cpu.Regs.SP - 1)
			case edge{3, true}:
			case edge{4, false}:
				cpu.writeAddressBus(cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP - 1)
			case edge{4, true}:
				cpu.Bus.WriteData(msb(cpu.Regs.PC))
			case edge{5, false}:
				cpu.writeAddressBus(cpu.Regs.SP)
			case edge{5, true}:
				cpu.Bus.WriteData(lsb(cpu.Regs.PC))
			case edge{6, false}:
				cpu.SetPC(cpu.Regs.GetWZ())
				return true
			case edge{6, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeRET: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP + 1)
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP + 1)
			case edge{2, true}:
				cpu.Regs.TempW = cpu.Bus.Data
			case edge{3, false}:
				cpu.SetPC(cpu.Regs.GetWZ())
			case edge{3, true}:
			case edge{4, false}, edge{4, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeRETZ: func(e edge) bool {
			f := func() bool { return cpu.Regs.GetFlagZ() }
			switch e {
			case edge{1, false}:
			case edge{1, true}:
			case edge{2, false}:
				if f() {
					cpu.lastBranchResult = +1
					cpu.SetSP(cpu.Regs.SP + 1)
					cpu.writeAddressBus(cpu.Regs.SP)
				} else {
					cpu.lastBranchResult = -1
					return true
				}
			case edge{2, true}:
				if cpu.lastBranchResult == +1 {
					cpu.Regs.TempW = cpu.Bus.Data
				} else {
					return true
				}
			case edge{3, false}:
				if cpu.lastBranchResult == +1 {
					cpu.SetSP(cpu.Regs.SP + 1)
					cpu.writeAddressBus(cpu.Regs.SP)
				} else {
					panicv(e)
				}
			case edge{3, true}:
				if cpu.lastBranchResult == +1 {
					cpu.Regs.TempZ = cpu.Bus.Data
				} else {
					panicv(e)
				}
			case edge{4, false}:
				if cpu.lastBranchResult == +1 {
					cpu.SetPC(cpu.Regs.GetWZ())
				} else {
					panicv(e)
				}
			case edge{4, true}:
				if cpu.lastBranchResult == +1 {
				} else {
					panicv(e)
				}
			case edge{5, false}:
				if cpu.lastBranchResult == +1 {
					return true
				} else {
					panicv(e)
				}
			case edge{5, true}:
				if cpu.lastBranchResult == +1 {
					return true
				} else {
					panicv(e)
				}
			default:
				panicv(e)
			}
			return false
		},
		OpcodePUSHBC: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.SetSP(cpu.Regs.SP - 1)
				cpu.writeAddressBus(cpu.Regs.SP)
			case edge{1, true}:
				cpu.Bus.WriteData(cpu.Regs.B)
			case edge{2, false}:
				cpu.SetSP(cpu.Regs.SP - 1)
				cpu.writeAddressBus(cpu.Regs.SP)
			case edge{2, true}:
				cpu.Bus.WriteData(cpu.Regs.C)
			case edge{3, false}:
			case edge{3, true}:
			case edge{4, false}:
				return true
			case edge{4, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodePOPBC: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP + 1)
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(cpu.Regs.SP)
				cpu.SetSP(cpu.Regs.SP + 1)
			case edge{2, true}:
				cpu.Regs.TempW = cpu.Bus.Data
			case edge{3, false}:
				cpu.SetBC(cpu.Regs.GetWZ())
				return true
			case edge{3, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeADDHLHL:  cpu.addhlrr(&cpu.Regs.H, &cpu.Regs.L),
		OpcodeADDHLBC:  cpu.addhlrr(&cpu.Regs.B, &cpu.Regs.C),
		OpcodeADDHLDE:  cpu.addhlrr(&cpu.Regs.D, &cpu.Regs.E),
		OpcodeLDBCnn:   cpu.ldxxnn(func(wz uint16) { cpu.SetBC(wz) }),
		OpcodeLDDEnn:   cpu.ldxxnn(func(wz uint16) { cpu.SetDE(wz) }),
		OpcodeLDHLnn:   cpu.ldxxnn(func(wz uint16) { cpu.SetHL(wz) }),
		OpcodeLDSPnn:   cpu.ldxxnn(func(wz uint16) { cpu.SetSP(wz) }),
		OpcodeLDHLAInc: cpu.ldhla(func() { cpu.SetHL(cpu.GetHL() + 1) }),
		OpcodeLDHLADec: cpu.ldhla(func() { cpu.SetHL(cpu.GetHL() - 1) }),
		OpcodeLDHLA:    cpu.ldhla(func() {}),
		OpcodeLDHCA: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(join16(0xff, cpu.Regs.C))
			case edge{1, true}:
				cpu.Bus.WriteData(cpu.Regs.A)
			case edge{2, false}, edge{2, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeADDHL: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.GetHL())
			case edge{1, true}:
				data := cpu.Bus.Data
				carry := uint16(cpu.Regs.A)+uint16(data) > 256
				result := cpu.Regs.A + data
				cpu.Regs.A = result
				cpu.Regs.SetFlagZ(result == 0)
				cpu.Regs.SetFlagN(false)
				cpu.Regs.TODOFlagH()
				cpu.Regs.SetFlagC(carry)
			case edge{2, false}, edge{2, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeCPHL: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.GetHL())
			case edge{1, true}:
				data := cpu.Bus.Data
				carry := data > cpu.Regs.A
				result := cpu.Regs.A - data
				cpu.Regs.SetFlagZ(result == 0)
				cpu.Regs.SetFlagN(true)
				cpu.Regs.TODOFlagH()
				cpu.Regs.SetFlagC(carry)
			case edge{2, false}, edge{2, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeLDnnA: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{2, true}:
				cpu.Regs.TempW = cpu.Bus.Data
			case edge{3, false}:
				cpu.writeAddressBus(cpu.Regs.GetWZ())
			case edge{3, true}:
				cpu.Bus.WriteData(cpu.Regs.A)
			case edge{4, false}, edge{4, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeLDAnn: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{2, true}:
				cpu.Regs.TempW = cpu.Bus.Data
			case edge{3, false}:
				cpu.writeAddressBus(cpu.Regs.GetWZ())
			case edge{3, true}:
				cpu.Regs.A = cpu.Bus.Data
			case edge{4, false}, edge{4, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeCPn: cpu.alun(func(imm uint8) {
			carry := cpu.Regs.A < imm
			result := cpu.Regs.A - imm
			cpu.Regs.SetFlagZ(result == 0)
			cpu.Regs.SetFlagN(true)
			cpu.Regs.TODOFlagH()
			cpu.Regs.SetFlagC(carry)
		}),
		OpcodeSUBn: cpu.alun(func(imm uint8) {
			carry := cpu.Regs.A < imm
			cpu.Regs.A -= imm
			cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
			cpu.Regs.SetFlagN(true)
			cpu.Regs.TODOFlagH()
			cpu.Regs.SetFlagC(carry)
		}),
		OpcodeORn: cpu.alun(func(imm uint8) {
			cpu.Regs.A |= imm
			cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
			cpu.Regs.SetFlagN(false)
			cpu.Regs.SetFlagH(false)
			cpu.Regs.SetFlagC(false)
		}),
		OpcodeANDn: cpu.alun(func(imm uint8) {
			cpu.Regs.A &= imm
			cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
			cpu.Regs.SetFlagN(false)
			cpu.Regs.SetFlagH(false)
			cpu.Regs.SetFlagC(false)
		}),
		OpcodeADDn: cpu.alun(func(imm uint8) {
			carry := uint16(cpu.Regs.A)+uint16(imm) > 0xff
			cpu.Regs.A += imm
			cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
			cpu.Regs.SetFlagN(false)
			cpu.Regs.TODOFlagH()
			cpu.Regs.SetFlagC(carry)
		}),
		OpcodeLDHnA: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(join16(0xff, cpu.Regs.TempZ))
				cpu.IncPC()
			case edge{2, true}:
				cpu.Bus.WriteData(cpu.Regs.A)
			case edge{3, false}, edge{3, true}:
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeLDHAn: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				cpu.writeAddressBus(join16(0xff, cpu.Regs.TempZ))
			case edge{2, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{3, false}:
				return true
			case edge{3, true}:
				cpu.Regs.A = cpu.Regs.TempZ
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeLDADE: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(join16(cpu.Regs.D, cpu.Regs.E))
			case edge{1, true}:
				cpu.Regs.TempZ = cpu.Bus.Data
			case edge{2, false}:
				return true
			case edge{2, true}:
				cpu.Regs.A = cpu.Regs.TempZ
				return true
			default:
				panicv(e)
			}
			return false
		},
		OpcodeLDAn: cpu.ldrn(&cpu.Regs.A),
		OpcodeLDBn: cpu.ldrn(&cpu.Regs.B),
		OpcodeLDCn: cpu.ldrn(&cpu.Regs.C),
		OpcodeLDDn: cpu.ldrn(&cpu.Regs.D),
		OpcodeLDEn: cpu.ldrn(&cpu.Regs.E),
		OpcodeLDHn: cpu.ldrn(&cpu.Regs.H),
		OpcodeLDLn: cpu.ldrn(&cpu.Regs.L),
		OpcodeCB: func(e edge) bool {
			switch e {
			case edge{1, false}:
				cpu.writeAddressBus(cpu.Regs.PC)
				cpu.IncPC()
			case edge{1, true}:
				opcode := cpu.Bus.Data
				cpu.CBOp = CBOp{Op: cb((opcode & 0xf8) >> 3), Target: CBTarget(opcode & 0x7)}
			case edge{2, false}:
				return true
			case edge{2, true}:
				var val uint8
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
					panic("indirect thru HL not implemented")
				case CBTargetA:
					val = cpu.Regs.A
				default:
					panic("unknown CBOp target")
				}
				switch cpu.CBOp.Op {
				case CbRL:
					bit7 := val & 0x80
					val <<= 1
					if cpu.Regs.GetFlagC() {
						val |= 0x01
					}
					cpu.Regs.SetFlagZ(val == 0)
					cpu.Regs.SetFlagN(false)
					cpu.Regs.SetFlagH(false)
					cpu.Regs.SetFlagC(bit7 != 0)
				case CbSRL:
					bit0 := val & 1
					val >>= 1
					cpu.Regs.SetFlagZ(val == 0)
					cpu.Regs.SetFlagN(false)
					cpu.Regs.SetFlagH(false)
					cpu.Regs.SetFlagC(bit0 != 0)
				case CbSWAP:
					val = ((val & 0x0f) << 4) | ((val & 0xf0) >> 4)
					cpu.Regs.SetFlagZ(val == 0)
					cpu.Regs.SetFlagN(false)
					cpu.Regs.SetFlagH(false)
					cpu.Regs.SetFlagC(false)
				case CbBit0:
					cbbit(cpu, val, 0x01)
				case CbBit1:
					cbbit(cpu, val, 0x02)
				case CbBit2:
					cbbit(cpu, val, 0x04)
				case CbBit3:
					cbbit(cpu, val, 0x08)
				case CbBit4:
					cbbit(cpu, val, 0x10)
				case CbBit5:
					cbbit(cpu, val, 0x20)
				case CbBit6:
					cbbit(cpu, val, 0x40)
				case CbBit7:
					cbbit(cpu, val, 0x80)
				default:
					panicf("unknown op = %+v", cpu.CBOp)
				}
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
					panic("indirect thru HL not implemented")
				case CBTargetA:
					cpu.Regs.A = val
				default:
					panic("unknown CBOp target")
				}
				return true
			default:
				panicv(e)
			}
			return false
		},
	}
}

func (cpu *CPU) singleCycle(f func()) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			return true
		case edge{1, true}:
			f()
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) jrcce(f func() bool) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Regs.TempZ = cpu.Bus.Data
		case edge{2, false}:
			if f() {
				cpu.lastBranchResult = +1
			} else {
				cpu.lastBranchResult = -1
				return true
			}
		case edge{2, true}:
			if cpu.lastBranchResult == +1 {
				newPC := uint16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.TempZ)))
				cpu.Regs.SetWZ(newPC)
			} else {
				return true
			}
		case edge{3, false}:
			if cpu.lastBranchResult == +1 {
				cpu.SetPC(cpu.Regs.GetWZ())
				return true
			} else {
				panicv(e)
			}
		case edge{3, true}:
			if cpu.lastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) jpccnn(f func() bool) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Regs.TempZ = cpu.Bus.Data
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.Regs.TempW = cpu.Bus.Data
		case edge{3, false}:
			if f() {
				cpu.lastBranchResult = +1
				cpu.SetPC(cpu.Regs.GetWZ())
			} else {
				cpu.lastBranchResult = -1
				return true
			}
		case edge{3, true}:
			if cpu.lastBranchResult == +1 {
			} else {
				return true
			}
		case edge{4, false}:
			if cpu.lastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		case edge{4, true}:
			if cpu.lastBranchResult == +1 {
				return true
			} else {
				panicv(e)
			}
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) ld(dst *uint8, src *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		*dst = *src
	})
}

func (cpu *CPU) andreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		cpu.Regs.A &= *reg
		cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
		cpu.Regs.SetFlagN(false)
		cpu.Regs.SetFlagH(true)
		cpu.Regs.SetFlagC(false)
	})
}

func (cpu *CPU) xorreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		cpu.Regs.A ^= *reg
		cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
		cpu.Regs.SetFlagN(false)
		cpu.Regs.SetFlagH(false)
		cpu.Regs.SetFlagC(false)
	})
}

func (cpu *CPU) orreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		cpu.Regs.A |= *reg
		cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
		cpu.Regs.SetFlagN(false)
		cpu.Regs.SetFlagH(false)
		cpu.Regs.SetFlagC(false)
	})
}

func (cpu *CPU) addreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		carry := uint16(*reg)+uint16(cpu.Regs.A) > 256
		cpu.Regs.A += *reg
		cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
		cpu.Regs.SetFlagN(false)
		cpu.Regs.TODOFlagH()
		cpu.Regs.SetFlagC(carry)
	})
}

func (cpu *CPU) adcreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		vCarry := 0
		if cpu.Regs.GetFlagC() {
			vCarry = 1
		}
		carry := uint16(*reg)+uint16(cpu.Regs.A)+uint16(vCarry) > 256
		cpu.Regs.A += *reg + uint8(vCarry)
		cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
		cpu.Regs.SetFlagN(false)
		cpu.Regs.TODOFlagH()
		cpu.Regs.SetFlagC(carry)
	})
}

func (cpu *CPU) subreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		carry := *reg > cpu.Regs.A
		cpu.Regs.A -= *reg
		cpu.Regs.SetFlagZ(cpu.Regs.A == 0)
		cpu.Regs.SetFlagN(true)
		cpu.Regs.TODOFlagH()
		cpu.Regs.SetFlagC(carry)
	})
}

func (cpu *CPU) cpreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		carry := *reg > cpu.Regs.A
		result := cpu.Regs.A - *reg
		cpu.Regs.SetFlagZ(result == 0)
		cpu.Regs.SetFlagN(true)
		cpu.Regs.TODOFlagH()
		cpu.Regs.SetFlagC(carry)
	})
}

func (cpu *CPU) decreg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		*reg -= 1
		cpu.Regs.SetFlagZ(*reg == 0)
		cpu.Regs.SetFlagN(true)
		cpu.Regs.TODOFlagH()
	})
}

func (cpu *CPU) increg(reg *uint8) func(e edge) bool {
	return cpu.singleCycle(func() {
		*reg += 1
		cpu.Regs.SetFlagZ(*reg == 0)
		cpu.Regs.SetFlagN(false)
		cpu.Regs.TODOFlagH()
	})
}

func (cpu *CPU) iduOp(f func()) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			f()
		case edge{1, true}:
		case edge{2, false}, edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) alun(f func(imm uint8)) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Regs.TempZ = cpu.Bus.Data
		case edge{2, false}:
			return true
		case edge{2, true}:
			f(cpu.Regs.TempZ)
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) ldrn(reg *uint8) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
		case edge{1, true}:
			cpu.Regs.TempZ = cpu.Bus.Data
		case edge{2, false}:
			cpu.IncPC()
			return true
		case edge{2, true}:
			*reg = cpu.Regs.TempZ
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) ldrhl(reg *uint8) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.GetHL())
		case edge{1, true}:
		case edge{2, false}:
			*reg = cpu.Bus.Data
			return true
		case edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) ldhla(f func()) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.GetHL())
			f()
		case edge{1, true}:
			cpu.Bus.WriteData(cpu.Regs.A)
		case edge{2, false}:
			return true
		case edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) addhlrr(hi, lo *uint8) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			carry := uint16(cpu.Regs.L)+uint16(*lo) > 256
			cpu.Regs.L += *lo
			cpu.Regs.SetFlagN(false)
			cpu.Regs.TODOFlagH()
			cpu.Regs.SetFlagC(carry)
		case edge{1, true}:
		case edge{2, false}:
			vCarry := 0
			if cpu.Regs.GetFlagC() {
				vCarry = 1
			}
			carry := uint16(cpu.Regs.H)+uint16(*hi)+uint16(vCarry) > 256
			cpu.Regs.H += *hi + uint8(vCarry)
			cpu.Regs.SetFlagN(false)
			cpu.Regs.TODOFlagH()
			cpu.Regs.SetFlagC(carry)
			return true
		case edge{2, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func (cpu *CPU) ldxxnn(f func(wz uint16)) func(e edge) bool {
	return func(e edge) bool {
		switch e {
		case edge{1, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{1, true}:
			cpu.Regs.TempZ = cpu.Bus.Data
		case edge{2, false}:
			cpu.writeAddressBus(cpu.Regs.PC)
			cpu.IncPC()
		case edge{2, true}:
			cpu.Regs.TempW = cpu.Bus.Data
		case edge{3, false}:
			f(join16(cpu.Regs.TempW, cpu.Regs.TempZ))
			return true
		case edge{3, true}:
			return true
		default:
			panicv(e)
		}
		return false
	}
}

func cbbit(cpu *CPU, val, mask uint8) {
	cpu.Regs.SetFlagZ(val&mask == 0)
	cpu.Regs.SetFlagN(false)
	cpu.Regs.SetFlagH(true)
}
