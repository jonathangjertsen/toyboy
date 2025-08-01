package model

import (
	"fmt"
	"slices"
)

type InstructionHandling any

const MaxOpcodesPerInsr = 6

type HandlerArray [256][6]InstructionHandling

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
		OpcodeNop:      {endNoop},
		OpcodeLDAA:     {ldAA},
		OpcodeLDAB:     {ldAB},
		OpcodeLDAC:     {ldAC},
		OpcodeLDAD:     {ldAD},
		OpcodeLDAE:     {ldAE},
		OpcodeLDAH:     {ldAH},
		OpcodeLDAL:     {ldAL},
		OpcodeLDAHL:    {ldrhl_1, ldAhl_2},
		OpcodeLDAHLInc: {ldahlinc_1, ldahlinc_2},
		OpcodeLDAHLDec: {ldahldec_1, ldahldec_2},
		OpcodeLDBA:     {ldBA},
		OpcodeLDBB:     {ldBB},
		OpcodeLDBC:     {ldBC},
		OpcodeLDBD:     {ldBD},
		OpcodeLDBE:     {ldBE},
		OpcodeLDBH:     {ldBH},
		OpcodeLDBL:     {ldBL},
		OpcodeLDBHL:    {ldrhl_1, ldBhl_2},
		OpcodeLDCA:     {ldCA},
		OpcodeLDCB:     {ldCB},
		OpcodeLDCC:     {ldCC},
		OpcodeLDCD:     {ldCD},
		OpcodeLDCE:     {ldCE},
		OpcodeLDCH:     {ldCH},
		OpcodeLDCL:     {ldCL},
		OpcodeLDCHL:    {ldrhl_1, ldChl_2},
		OpcodeLDDA:     {ldDA},
		OpcodeLDDB:     {ldDB},
		OpcodeLDDC:     {ldDC},
		OpcodeLDDD:     {ldDD},
		OpcodeLDDE:     {ldDE},
		OpcodeLDDH:     {ldDH},
		OpcodeLDDL:     {ldDL},
		OpcodeLDDHL:    {ldrhl_1, ldDhl_2},
		OpcodeLDEA:     {ldEA},
		OpcodeLDEB:     {ldEB},
		OpcodeLDEC:     {ldEC},
		OpcodeLDED:     {ldED},
		OpcodeLDEE:     {ldEE},
		OpcodeLDEH:     {ldEH},
		OpcodeLDEL:     {ldEL},
		OpcodeLDEHL:    {ldrhl_1, ldEhl_2},
		OpcodeLDHA:     {ldHA},
		OpcodeLDHB:     {ldHB},
		OpcodeLDHC:     {ldHC},
		OpcodeLDHD:     {ldHD},
		OpcodeLDHE:     {ldHE},
		OpcodeLDHH:     {ldHH},
		OpcodeLDHL:     {ldHL},
		OpcodeLDHHL:    {ldrhl_1, ldHhl_2},
		OpcodeLDLA:     {ldLA},
		OpcodeLDLB:     {ldLB},
		OpcodeLDLC:     {ldLC},
		OpcodeLDLD:     {ldLD},
		OpcodeLDLE:     {ldLE},
		OpcodeLDLH:     {ldLH},
		OpcodeLDLL:     {ldLL},
		OpcodeLDLHL:    {ldrhl_1, ldLhl_2},
		OpcodeRLA:      {rla},
		OpcodeRRA:      {rra},
		OpcodeRLCA:     {rlca},
		OpcodeRRCA:     {rrca},
		OpcodeORA:      {orreg_A},
		OpcodeORB:      {orreg_B},
		OpcodeORC:      {orreg_C},
		OpcodeORD:      {orreg_D},
		OpcodeORE:      {orreg_E},
		OpcodeORH:      {orreg_H},
		OpcodeORL:      {orreg_L},
		OpcodeORHL:     {aluHL_1, ORHL_2},
		OpcodeANDA:     {andreg_A},
		OpcodeANDB:     {andreg_B},
		OpcodeANDC:     {andreg_C},
		OpcodeANDD:     {andreg_D},
		OpcodeANDE:     {andreg_E},
		OpcodeANDH:     {andreg_H},
		OpcodeANDL:     {andreg_L},
		OpcodeANDHL:    {aluHL_1, ANDHL_2},
		OpcodeXORA:     {xorreg_A},
		OpcodeXORB:     {xorreg_B},
		OpcodeXORC:     {xorreg_C},
		OpcodeXORD:     {xorreg_D},
		OpcodeXORE:     {xorreg_E},
		OpcodeXORH:     {xorreg_H},
		OpcodeXORL:     {xorreg_L},
		OpcodeXORHL:    {aluHL_1, XORHL_2},
		OpcodeSUBA:     {subreg_A},
		OpcodeSUBB:     {subreg_B},
		OpcodeSUBC:     {subreg_C},
		OpcodeSUBD:     {subreg_D},
		OpcodeSUBE:     {subreg_E},
		OpcodeSUBH:     {subreg_H},
		OpcodeSUBL:     {subreg_L},
		OpcodeSUBHL:    {aluHL_1, SUBHL_2},
		OpcodeSBCA:     {sbcreg_A},
		OpcodeSBCB:     {sbcreg_B},
		OpcodeSBCC:     {sbcreg_C},
		OpcodeSBCD:     {sbcreg_D},
		OpcodeSBCE:     {sbcreg_E},
		OpcodeSBCH:     {sbcreg_H},
		OpcodeSBCL:     {sbcreg_L},
		OpcodeSBCHL:    {aluHL_1, SBCHL_2},
		OpcodeCPA:      {cpreg_A},
		OpcodeCPB:      {cpreg_B},
		OpcodeCPC:      {cpreg_C},
		OpcodeCPD:      {cpreg_D},
		OpcodeCPE:      {cpreg_E},
		OpcodeCPH:      {cpreg_H},
		OpcodeCPL:      {cpreg_L},
		OpcodeCPHL:     {aluHL_1, CPHL_2},
		OpcodeADDA:     {addreg_A},
		OpcodeADDB:     {addreg_B},
		OpcodeADDC:     {addreg_C},
		OpcodeADDD:     {addreg_D},
		OpcodeADDE:     {addreg_E},
		OpcodeADDH:     {addreg_H},
		OpcodeADDL:     {addreg_L},
		OpcodeADDHL:    {aluHL_1, ADDHL_2},
		OpcodeADDSPe:   {addspe_1, addspe_2, addspe_3, addspe_4},
		OpcodeADCA:     {adcreg_A},
		OpcodeADCB:     {adcreg_B},
		OpcodeADCC:     {adcreg_C},
		OpcodeADCD:     {adcreg_D},
		OpcodeADCE:     {adcreg_E},
		OpcodeADCH:     {adcreg_H},
		OpcodeADCL:     {adcreg_L},
		OpcodeADCHL:    {aluHL_1, ADCHL_2},
		OpcodeDAA:      {daa},
		OpcodeCPLaka2f: {cpl},
		OpcodeCCF:      {ccf},
		OpcodeSCF:      {scf},
		OpcodeDECA:     {decreg_A},
		OpcodeDECB:     {decreg_B},
		OpcodeDECC:     {decreg_C},
		OpcodeDECD:     {decreg_D},
		OpcodeDECE:     {decreg_E},
		OpcodeDECH:     {decreg_H},
		OpcodeDECL:     {decreg_L},
		OpcodeINCA:     {increg_A},
		OpcodeINCB:     {increg_B},
		OpcodeINCC:     {increg_C},
		OpcodeINCD:     {increg_D},
		OpcodeINCE:     {increg_E},
		OpcodeINCH:     {increg_H},
		OpcodeINCL:     {increg_L},
		OpcodeINCHLInd: {inchlind_1, inchlind_2, inchlind_3},
		OpcodeDECHLInd: {dechlind_1, dechlind_2, dechlind_3},
		OpcodeDI:       {di},
		OpcodeEI:       {ei},
		OpcodeHALT:     {halt},
		OpcodeJRe:      {jre_1, jre_2, jre_3},
		OpcodeJPnn:     {jpnn_1, jpnn_2, jpnn_3, jpnn_4},
		OpcodeJPHL:     {jphl},
		OpcodeJRZe:     {jrcce_1, jrZe_2, jrcce_3},
		OpcodeJRCe:     {jrcce_1, jrCe_2, jrcce_3},
		OpcodeJRNZe:    {jrcce_1, jrNZe_2, jrcce_3},
		OpcodeJRNCe:    {jrcce_1, jrNCe_2, jrcce_3},
		OpcodeJPCnn:    {jpccnn_1, jpccnn_2, jpCnn_3, jpccnn_4},
		OpcodeJPNCnn:   {jpccnn_1, jpccnn_2, jpNCnn_3, jpccnn_4},
		OpcodeJPZnn:    {jpccnn_1, jpccnn_2, jpZnn_3, jpccnn_4},
		OpcodeJPNZnn:   {jpccnn_1, jpccnn_2, jpNZnn_3, jpccnn_4},
		OpcodeINCBC:    {INCBC_1, iduOp_2},
		OpcodeINCDE:    {INCDE_1, iduOp_2},
		OpcodeINCHL:    {INCHL_1, iduOp_2},
		OpcodeINCSP:    {INCSP_1, iduOp_2},
		OpcodeDECBC:    {DECBC_1, iduOp_2},
		OpcodeDECDE:    {DECDE_1, iduOp_2},
		OpcodeDECHL:    {DECHL_1, iduOp_2},
		OpcodeDECSP:    {DECSP_1, iduOp_2},
		OpcodeCALLnn:   {callnn_1, callnn_2, callnn_3, callnn_4, callnn_5, callnn_6},
		OpcodeCALLNZnn: {callccnn_1, callccnn_2, callNZnn_3, callccnn_4, callccnn_5, callccnn_6},
		OpcodeCALLZnn:  {callccnn_1, callccnn_2, callZnn_3, callccnn_4, callccnn_5, callccnn_6},
		OpcodeCALLNCnn: {callccnn_1, callccnn_2, callNCnn_3, callccnn_4, callccnn_5, callccnn_6},
		OpcodeCALLCnn:  {callccnn_1, callccnn_2, callCnn_3, callccnn_4, callccnn_5, callccnn_6},
		OpcodeRET:      {ret_1, ret_2, ret_3, ret_4},
		OpcodeRETI:     {reti_1, reti_2, reti_3, reti_4},
		OpcodeRETZ:     {retcc_1, retZ_2, retcc_3, retcc_4, retcc_5},
		OpcodeRETNZ:    {retcc_1, retNZ_2, retcc_3, retcc_4, retcc_5},
		OpcodeRETC:     {retcc_1, retC_2, retcc_3, retcc_4, retcc_5},
		OpcodeRETNC:    {retcc_1, retNC_2, retcc_3, retcc_4, retcc_5},
		OpcodePUSHBC:   {pushBC_1, pushBC_2, push_3, push_4},
		OpcodePUSHDE:   {pushDE_1, pushDE_2, push_3, push_4},
		OpcodePUSHHL:   {pushHL_1, pushHL_2, push_3, push_4},
		OpcodePUSHAF:   {pushAF_1, pushAF_2, push_3, push_4},
		OpcodePOPBC:    {pop1, pop2, popBC_3},
		OpcodePOPDE:    {pop1, pop2, popDE_3},
		OpcodePOPHL:    {pop1, pop2, popHL_3},
		OpcodePOPAF:    {pop1, pop2, popAF_3},
		OpcodeADDHLHL:  {addhlHL_1, addhlHL_2},
		OpcodeADDHLBC:  {addhlBC_1, addhlBC_2},
		OpcodeADDHLDE:  {addhlDE_1, addhlDE_2},
		OpcodeADDHLSP:  {addhlsp_1, addhlsp_2},
		OpcodeLDBCnn:   {ldxxnn_1, ldxxnn_2, ldBCnn_3},
		OpcodeLDDEnn:   {ldxxnn_1, ldxxnn_2, ldDEnn_3},
		OpcodeLDHLnn:   {ldxxnn_1, ldxxnn_2, ldHLnn_3},
		OpcodeLDSPnn:   {ldxxnn_1, ldxxnn_2, ldSPnn_3},
		OpcodeLDHLn:    {ldhln_1, ldhln_2, ldhln_3},
		OpcodeLDHLSPe:  {ldhlspe_1, ldhlspe_2, ldhlspe_3},
		OpcodeLDSPHL:   {ldsphl_1, ldsphl_2},
		OpcodeLDHLAInc: {ldhlAinc_1, ldhlr_2},
		OpcodeLDHLADec: {ldhlAdec_1, ldhlr_2},
		OpcodeLDHLA:    {ldhlA_1, ldhlr_2},
		OpcodeLDHLB:    {ldhlB_1, ldhlr_2},
		OpcodeLDHLC:    {ldhlC_1, ldhlr_2},
		OpcodeLDHLD:    {ldhlD_1, ldhlr_2},
		OpcodeLDHLE:    {ldhlE_1, ldhlr_2},
		OpcodeLDHLH:    {ldhlH_1, ldhlr_2},
		OpcodeLDHLL:    {ldhlL_1, ldhlr_2},
		OpcodeLDBCA:    {ldBCA_1, ldrra_2},
		OpcodeLDDEA:    {ldDEA_1, ldrra_2},
		OpcodeLDHCA:    {ldhca_1, ldhca_2},
		OpcodeLDHAC:    {ldhac_1, ldhac_2},
		OpcodeLDnnSP:   {ldnnsp_1, ldnnsp_2, ldnnsp_3, ldnnsp_4, ldnnsp_5},
		OpcodeLDnnA:    {ldnna_1, ldnna_2, ldnna_3, ldnna_4},
		OpcodeLDAnn:    {ldann_1, ldann_2, ldann_3, ldann_4},
		OpcodeCPn:      {alun_1, CPn_2},
		OpcodeSUBn:     {alun_1, SUBn_2},
		OpcodeORn:      {alun_1, ORn_2},
		OpcodeANDn:     {alun_1, ANDn_2},
		OpcodeADDn:     {alun_1, ADDn_2},
		OpcodeADCn:     {alun_1, ADCn_2},
		OpcodeSBCn:     {alun_1, SBCn_2},
		OpcodeXORn:     {alun_1, XORn_2},
		OpcodeLDHnA:    {ldhna_1, ldhna_2, ldhna_3},
		OpcodeLDHAn:    {ldhan_1, ldhan_2, ldhan_3},
		OpcodeLDADE:    {ldaDE_1, ldarr_2},
		OpcodeLDABC:    {ldaBC_1, ldarr_2},
		OpcodeLDAn:     {ldrn_1, ldAn_2},
		OpcodeLDBn:     {ldrn_1, ldBn_2},
		OpcodeLDCn:     {ldrn_1, ldCn_2},
		OpcodeLDDn:     {ldrn_1, ldDn_2},
		OpcodeLDEn:     {ldrn_1, ldEn_2},
		OpcodeLDHn:     {ldrn_1, ldHn_2},
		OpcodeLDLn:     {ldrn_1, ldLn_2},
		OpcodeCB:       {runCB_1, runCB_2, runCB_3, runCB_4},
		OpcodeRST0x00:  {rst_1, rst_2, rst_3_00, rst_4},
		OpcodeRST0x08:  {rst_1, rst_2, rst_3_08, rst_4},
		OpcodeRST0x10:  {rst_1, rst_2, rst_3_10, rst_4},
		OpcodeRST0x18:  {rst_1, rst_2, rst_3_18, rst_4},
		OpcodeRST0x20:  {rst_1, rst_2, rst_3_20, rst_4},
		OpcodeRST0x28:  {rst_1, rst_2, rst_3_28, rst_4},
		OpcodeRST0x30:  {rst_1, rst_2, rst_3_30, rst_4},
		OpcodeRST0x38:  {rst_1, rst_2, rst_3_38, rst_4},
		OpcodeUndefD3:  {notImplemented},
		OpcodeUndefDB:  {notImplemented},
		OpcodeUndefDD:  {notImplemented},
		OpcodeUndefE3:  {notImplemented},
		OpcodeUndefE4:  {notImplemented},
		OpcodeUndefEB:  {notImplemented},
		OpcodeUndefEC:  {notImplemented},
		OpcodeUndefED:  {notImplemented},
		OpcodeUndefF4:  {notImplemented},
		OpcodeUndefFC:  {notImplemented},
		OpcodeUndefFD:  {notImplemented},
	}
}

func noop(gb *Gameboy) bool {
	return false
}

func endNoop(gb *Gameboy) bool {
	return true
}

func rra(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(RRA(gb.CPU.Regs.A, gb.CPU.Regs.GetFlagC()))
	return true
}

func rla(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(RLA(gb.CPU.Regs.A, gb.CPU.Regs.GetFlagC()))
	return true
}

func rlca(gb *Gameboy) bool { gb.CPU.Regs.SetFlagsAndA(RLCA(gb.CPU.Regs.A)); return true }
func rrca(gb *Gameboy) bool { gb.CPU.Regs.SetFlagsAndA(RRCA(gb.CPU.Regs.A)); return true }

func notImplemented(gb *Gameboy) bool {
	panicf("not implemented opcode %v", gb.CPU.Regs.IR)
	return false
}

func jphl(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.GetHL()))
	return true
}

func halt(gb *Gameboy) bool {
	gb.CPU.Halted = true
	return true
}

func jre_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func jre_2(gb *Gameboy) bool {
	return false
}

func jre_3(gb *Gameboy) bool {
	if gb.CPU.Regs.TempZ&SignBit8 != 0 {
		gb.CPU.SetPC(gb.CPU.Regs.PC - Addr(gb.CPU.Regs.TempZ.SignedAbs()))
	} else {
		gb.CPU.SetPC(gb.CPU.Regs.PC + Addr(gb.CPU.Regs.TempZ))
	}
	return true
}

func jpnn_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func jpnn_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func jpnn_3(gb *Gameboy) bool {
	return false
}

func jpnn_4(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	return true
}

func jrcce_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func jrZe_2(gb *Gameboy) bool  { return jrcce_2(gb, gb.CPU.Regs.GetFlagZ()) }
func jrNZe_2(gb *Gameboy) bool { return jrcce_2(gb, !gb.CPU.Regs.GetFlagZ()) }
func jrCe_2(gb *Gameboy) bool  { return jrcce_2(gb, gb.CPU.Regs.GetFlagC()) }
func jrNCe_2(gb *Gameboy) bool { return jrcce_2(gb, !gb.CPU.Regs.GetFlagC()) }

func jrcce_2(gb *Gameboy, cond bool) bool {
	if cond {
		gb.CPU.LastBranchResult = +1
		newPC := Data16(int16(gb.CPU.Regs.PC) + int16(int8(gb.CPU.Regs.TempZ)))
		gb.CPU.Regs.SetWZ(newPC)
		return false
	}

	gb.CPU.LastBranchResult = -1
	return true
}

func jrcce_3(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	return true
}

func jpccnn_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func jpccnn_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func jpZnn_3(gb *Gameboy) bool  { return jpccnn_3(gb, gb.CPU.Regs.GetFlagZ()) }
func jpNZnn_3(gb *Gameboy) bool { return jpccnn_3(gb, !gb.CPU.Regs.GetFlagZ()) }
func jpCnn_3(gb *Gameboy) bool  { return jpccnn_3(gb, gb.CPU.Regs.GetFlagC()) }
func jpNCnn_3(gb *Gameboy) bool { return jpccnn_3(gb, !gb.CPU.Regs.GetFlagC()) }

func jpccnn_3(gb *Gameboy, cond bool) bool {
	if cond {
		gb.CPU.LastBranchResult = +1
		gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
		return false
	}
	gb.CPU.LastBranchResult = -1
	return true
}

var jpccnn_4 = endNoop

func pushAF_1(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.A)
	return false
}

func pushAF_2(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.F)
	return false
}

func pushBC_1(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.B)
	return false
}

func pushBC_2(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.C)
	return false
}

func pushDE_1(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.D)
	return false
}

func pushDE_2(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.E)
	return false
}

func pushHL_1(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.H)
	return false
}

func pushHL_2(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.L)
	return false
}

var push_3 = noop
var push_4 = endNoop

func pop1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func pop2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func popAF_3(gb *Gameboy) bool {
	gb.CPU.Regs.A = gb.CPU.Regs.TempW
	gb.CPU.Regs.F = gb.CPU.Regs.TempZ & 0xf0
	return true
}

func popBC_3(gb *Gameboy) bool {
	gb.CPU.Regs.B = gb.CPU.Regs.TempW
	gb.CPU.Regs.C = gb.CPU.Regs.TempZ
	return true
}

func popDE_3(gb *Gameboy) bool {
	gb.CPU.Regs.D = gb.CPU.Regs.TempW
	gb.CPU.Regs.E = gb.CPU.Regs.TempZ
	return true
}

func popHL_3(gb *Gameboy) bool {
	gb.CPU.Regs.H = gb.CPU.Regs.TempW
	gb.CPU.Regs.L = gb.CPU.Regs.TempZ
	return true
}

func callnn_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func callnn_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func callnn_3(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	return false
}

func callnn_4(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteData(gb.CPU.Regs.PC.MSB())
	return false
}

func callnn_5(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.PC.LSB())
	return false
}

func callnn_6(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	return true
}

func callccnn_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func callccnn_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func callZnn_3(gb *Gameboy) bool {
	return callccnn_3(gb, gb.CPU.Regs.GetFlagZ())
}
func callNZnn_3(gb *Gameboy) bool {
	return callccnn_3(gb, !gb.CPU.Regs.GetFlagZ())
}
func callCnn_3(gb *Gameboy) bool {
	return callccnn_3(gb, gb.CPU.Regs.GetFlagC())
}
func callNCnn_3(gb *Gameboy) bool {
	return callccnn_3(gb, !gb.CPU.Regs.GetFlagC())
}

func callccnn_3(gb *Gameboy, cond bool) bool {
	if cond {
		gb.CPU.LastBranchResult = +1
		gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	} else {
		gb.CPU.LastBranchResult = -1
		return true
	}
	return false
}

func callccnn_4(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteData(gb.CPU.Regs.PC.MSB())
	return false
}

func callccnn_5(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.PC.LSB())
	return false
}

func callccnn_6(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	return true
}

func ret_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ret_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func ret_3(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	return false
}

var ret_4 = endNoop

func reti_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func reti_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func reti_3(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	// TODO verify if this is the right cycle
	gb.Interrupts.IME = true
	return false
}

var reti_4 = endNoop

var retcc_1 = noop

func retZ_2(gb *Gameboy) bool {
	return retcc_2(gb, gb.CPU.Regs.GetFlagZ())
}

func retNZ_2(gb *Gameboy) bool {
	return retcc_2(gb, !gb.CPU.Regs.GetFlagZ())
}

func retC_2(gb *Gameboy) bool {
	return retcc_2(gb, gb.CPU.Regs.GetFlagC())
}

func retNC_2(gb *Gameboy) bool {
	return retcc_2(gb, !gb.CPU.Regs.GetFlagC())
}

func retcc_2(gb *Gameboy, cond bool) bool {
	if cond {
		gb.CPU.LastBranchResult = +1
		gb.WriteAddress(gb.CPU.Regs.SP)
		gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
		gb.CPU.Regs.TempZ = gb.Data
	} else {
		gb.CPU.LastBranchResult = -1
		return true
	}
	return false
}

func retcc_3(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.CPU.SetSP(gb.CPU.Regs.SP + 1)
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func retcc_4(gb *Gameboy) bool {
	gb.CPU.SetPC(Addr(gb.CPU.Regs.GetWZ()))
	return false
}

var retcc_5 = endNoop

func ldAA(gb *Gameboy) bool { return true }
func ldBA(gb *Gameboy) bool { gb.CPU.Regs.B = gb.CPU.Regs.A; return true }
func ldCA(gb *Gameboy) bool { gb.CPU.Regs.C = gb.CPU.Regs.A; return true }
func ldDA(gb *Gameboy) bool { gb.CPU.Regs.D = gb.CPU.Regs.A; return true }
func ldEA(gb *Gameboy) bool { gb.CPU.Regs.E = gb.CPU.Regs.A; return true }
func ldHA(gb *Gameboy) bool { gb.CPU.Regs.H = gb.CPU.Regs.A; return true }
func ldLA(gb *Gameboy) bool { gb.CPU.Regs.L = gb.CPU.Regs.A; return true }

func ldAB(gb *Gameboy) bool { gb.CPU.Regs.A = gb.CPU.Regs.B; return true }
func ldBB(gb *Gameboy) bool { return true }
func ldCB(gb *Gameboy) bool { gb.CPU.Regs.C = gb.CPU.Regs.B; return true }
func ldDB(gb *Gameboy) bool { gb.CPU.Regs.D = gb.CPU.Regs.B; return true }
func ldEB(gb *Gameboy) bool { gb.CPU.Regs.E = gb.CPU.Regs.B; return true }
func ldHB(gb *Gameboy) bool { gb.CPU.Regs.H = gb.CPU.Regs.B; return true }
func ldLB(gb *Gameboy) bool { gb.CPU.Regs.L = gb.CPU.Regs.B; return true }

func ldAC(gb *Gameboy) bool { gb.CPU.Regs.A = gb.CPU.Regs.C; return true }
func ldBC(gb *Gameboy) bool { gb.CPU.Regs.B = gb.CPU.Regs.C; return true }
func ldCC(gb *Gameboy) bool { return true }
func ldDC(gb *Gameboy) bool { gb.CPU.Regs.D = gb.CPU.Regs.C; return true }
func ldEC(gb *Gameboy) bool { gb.CPU.Regs.E = gb.CPU.Regs.C; return true }
func ldHC(gb *Gameboy) bool { gb.CPU.Regs.H = gb.CPU.Regs.C; return true }
func ldLC(gb *Gameboy) bool { gb.CPU.Regs.L = gb.CPU.Regs.C; return true }

func ldAD(gb *Gameboy) bool { gb.CPU.Regs.A = gb.CPU.Regs.D; return true }
func ldBD(gb *Gameboy) bool { gb.CPU.Regs.B = gb.CPU.Regs.D; return true }
func ldCD(gb *Gameboy) bool { gb.CPU.Regs.C = gb.CPU.Regs.D; return true }
func ldDD(gb *Gameboy) bool { return true }
func ldED(gb *Gameboy) bool { gb.CPU.Regs.E = gb.CPU.Regs.D; return true }
func ldHD(gb *Gameboy) bool { gb.CPU.Regs.H = gb.CPU.Regs.D; return true }
func ldLD(gb *Gameboy) bool { gb.CPU.Regs.L = gb.CPU.Regs.D; return true }

func ldAE(gb *Gameboy) bool { gb.CPU.Regs.A = gb.CPU.Regs.E; return true }
func ldBE(gb *Gameboy) bool { gb.CPU.Regs.B = gb.CPU.Regs.E; return true }
func ldCE(gb *Gameboy) bool { gb.CPU.Regs.C = gb.CPU.Regs.E; return true }
func ldDE(gb *Gameboy) bool { gb.CPU.Regs.D = gb.CPU.Regs.E; return true }
func ldEE(gb *Gameboy) bool { return true }
func ldHE(gb *Gameboy) bool { gb.CPU.Regs.H = gb.CPU.Regs.E; return true }
func ldLE(gb *Gameboy) bool { gb.CPU.Regs.L = gb.CPU.Regs.E; return true }

func ldAH(gb *Gameboy) bool { gb.CPU.Regs.A = gb.CPU.Regs.H; return true }
func ldBH(gb *Gameboy) bool { gb.CPU.Regs.B = gb.CPU.Regs.H; return true }
func ldCH(gb *Gameboy) bool { gb.CPU.Regs.C = gb.CPU.Regs.H; return true }
func ldDH(gb *Gameboy) bool { gb.CPU.Regs.D = gb.CPU.Regs.H; return true }
func ldEH(gb *Gameboy) bool { gb.CPU.Regs.E = gb.CPU.Regs.H; return true }
func ldHH(gb *Gameboy) bool { return true }
func ldLH(gb *Gameboy) bool { gb.CPU.Regs.L = gb.CPU.Regs.H; return true }

func ldAL(gb *Gameboy) bool { gb.CPU.Regs.A = gb.CPU.Regs.L; return true }
func ldBL(gb *Gameboy) bool { gb.CPU.Regs.B = gb.CPU.Regs.L; return true }
func ldCL(gb *Gameboy) bool { gb.CPU.Regs.C = gb.CPU.Regs.L; return true }
func ldDL(gb *Gameboy) bool { gb.CPU.Regs.D = gb.CPU.Regs.L; return true }
func ldEL(gb *Gameboy) bool { gb.CPU.Regs.E = gb.CPU.Regs.L; return true }
func ldHL(gb *Gameboy) bool { gb.CPU.Regs.H = gb.CPU.Regs.L; return true }
func ldLL(gb *Gameboy) bool { return true }

func andreg_A(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.A) }
func andreg_B(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.B) }
func andreg_C(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.C) }
func andreg_D(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.D) }
func andreg_E(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.E) }
func andreg_H(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.H) }
func andreg_L(gb *Gameboy) bool { return andreg(gb, gb.CPU.Regs.L) }

func andreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(AND(gb.CPU.Regs.A, reg))
	return true
}

func xorreg_A(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.A) }
func xorreg_B(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.B) }
func xorreg_C(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.C) }
func xorreg_D(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.D) }
func xorreg_E(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.E) }
func xorreg_H(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.H) }
func xorreg_L(gb *Gameboy) bool { return xorreg(gb, gb.CPU.Regs.L) }

func xorreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(XOR(gb.CPU.Regs.A, reg))
	return true
}

func orreg_A(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.A) }
func orreg_B(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.B) }
func orreg_C(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.C) }
func orreg_D(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.D) }
func orreg_E(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.E) }
func orreg_H(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.H) }
func orreg_L(gb *Gameboy) bool { return orreg(gb, gb.CPU.Regs.L) }

func orreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(OR(gb.CPU.Regs.A, reg))
	return true
}

func addreg_A(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.A) }
func addreg_B(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.B) }
func addreg_C(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.C) }
func addreg_D(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.D) }
func addreg_E(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.E) }
func addreg_H(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.H) }
func addreg_L(gb *Gameboy) bool { return addreg(gb, gb.CPU.Regs.L) }

func addreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, reg, false))
	return true
}

func adcreg_A(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.A) }
func adcreg_B(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.B) }
func adcreg_C(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.C) }
func adcreg_D(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.D) }
func adcreg_E(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.E) }
func adcreg_H(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.H) }
func adcreg_L(gb *Gameboy) bool { return adcreg(gb, gb.CPU.Regs.L) }

func adcreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, reg, gb.CPU.Regs.GetFlagC()))
	return true
}

func inchlind_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func inchlind_2(gb *Gameboy) bool {
	res := ADD(gb.CPU.Regs.TempZ, 1, false)
	// apparently doesn't set C.
	gb.CPU.Regs.SetFlagH(res.H)
	gb.CPU.Regs.SetFlagZ(res.Z())
	gb.CPU.Regs.SetFlagN(res.N)
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.Regs.TempZ = res.Value
	gb.WriteData(gb.CPU.Regs.TempZ)
	return false
}

func inchlind_3(gb *Gameboy) bool {
	return true
}

func dechlind_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func dechlind_2(gb *Gameboy) bool {
	res := SUB(gb.CPU.Regs.TempZ, 1, false)
	// apparently doesn't set C.
	gb.CPU.Regs.SetFlagH(res.H)
	gb.CPU.Regs.SetFlagZ(res.Z())
	gb.CPU.Regs.SetFlagN(res.N)
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.Regs.TempZ = res.Value
	gb.WriteData(gb.CPU.Regs.TempZ)
	return false
}

var dechlind_3 = endNoop

func di(gb *Gameboy) bool {
	gb.Interrupts.SetIMENextCycle = false
	gb.Interrupts.SetIME(gb.Mem, false)
	return true
}

func ei(gb *Gameboy) bool {
	gb.Interrupts.SetIMENextCycle = true
	return true
}

func aluHL_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ORHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(OR(gb.CPU.Regs.A, gb.CPU.Regs.TempZ))
	return true
}

func ANDHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(AND(gb.CPU.Regs.A, gb.CPU.Regs.TempZ))
	return true
}

func XORHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(XOR(gb.CPU.Regs.A, gb.CPU.Regs.TempZ))
	return true
}

func SUBHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, false))
	return true
}

func SBCHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, gb.CPU.Regs.GetFlagC()))
	return true
}

func CPHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlags(SUB(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, false))
	return true
}

func ADDHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, false))
	return true
}

func ADCHL_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, gb.CPU.Regs.GetFlagC()))
	return true
}

func daa(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(DAA(gb.CPU.Regs.A, gb.CPU.Regs.GetFlagC(), gb.CPU.Regs.GetFlagN(), gb.CPU.Regs.GetFlagH()))
	return true
}

func cpl(gb *Gameboy) bool {
	gb.CPU.Regs.A ^= 0xff
	gb.CPU.Regs.SetFlagN(true)
	gb.CPU.Regs.SetFlagH(true)
	return true
}

func ccf(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagC(!gb.CPU.Regs.GetFlagC())
	gb.CPU.Regs.SetFlagN(false)
	gb.CPU.Regs.SetFlagH(false)
	return true
}

func scf(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagC(true)
	gb.CPU.Regs.SetFlagN(false)
	gb.CPU.Regs.SetFlagH(false)
	return true
}

func subreg_A(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.A) }
func subreg_B(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.B) }
func subreg_C(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.C) }
func subreg_D(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.D) }
func subreg_E(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.E) }
func subreg_H(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.H) }
func subreg_L(gb *Gameboy) bool { return subreg(gb, gb.CPU.Regs.L) }

func subreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, reg, false))
	return true
}

func sbcreg_A(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.A) }
func sbcreg_B(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.B) }
func sbcreg_C(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.C) }
func sbcreg_D(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.D) }
func sbcreg_E(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.E) }
func sbcreg_H(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.H) }
func sbcreg_L(gb *Gameboy) bool { return sbcreg(gb, gb.CPU.Regs.L) }

func sbcreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, reg, gb.CPU.Regs.GetFlagC()))
	return true
}

func cpreg_A(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.A) }
func cpreg_B(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.B) }
func cpreg_C(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.C) }
func cpreg_D(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.D) }
func cpreg_E(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.E) }
func cpreg_H(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.H) }
func cpreg_L(gb *Gameboy) bool { return cpreg(gb, gb.CPU.Regs.L) }

func cpreg(gb *Gameboy, reg Data8) bool {
	gb.CPU.Regs.SetFlags(SUB(gb.CPU.Regs.A, reg, false))
	return true
}

func decreg_A(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.A) }
func decreg_B(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.B) }
func decreg_C(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.C) }
func decreg_D(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.D) }
func decreg_E(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.E) }
func decreg_H(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.H) }
func decreg_L(gb *Gameboy) bool { return decreg(gb, &gb.CPU.Regs.L) }

func decreg(gb *Gameboy, reg *Data8) bool {
	result := SUB(*reg, 1, false)
	*reg = result.Value
	gb.CPU.Regs.SetFlagZ(result.Z())
	gb.CPU.Regs.SetFlagH(result.H)
	gb.CPU.Regs.SetFlagN(result.N)
	return true
}

func increg_A(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.A) }
func increg_B(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.B) }
func increg_C(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.C) }
func increg_D(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.D) }
func increg_E(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.E) }
func increg_H(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.H) }
func increg_L(gb *Gameboy) bool { return increg(gb, &gb.CPU.Regs.L) }

func increg(gb *Gameboy, reg *Data8) bool {
	result := ADD(*reg, 1, false)
	*reg = result.Value
	gb.CPU.Regs.SetFlagZ(result.Z())
	gb.CPU.Regs.SetFlagH(result.H)
	gb.CPU.Regs.SetFlagN(result.N)
	return true
}

func INCBC_1(gb *Gameboy) bool { gb.CPU.SetBC(gb.CPU.GetBC() + 1); return true }
func INCDE_1(gb *Gameboy) bool { gb.CPU.SetDE(gb.CPU.GetDE() + 1); return true }
func INCHL_1(gb *Gameboy) bool { gb.CPU.SetHL(gb.CPU.GetHL() + 1); return true }
func INCSP_1(gb *Gameboy) bool { gb.CPU.Regs.SP++; return true }
func DECBC_1(gb *Gameboy) bool { gb.CPU.SetBC(gb.CPU.GetBC() - 1); return true }
func DECDE_1(gb *Gameboy) bool { gb.CPU.SetDE(gb.CPU.GetDE() - 1); return true }
func DECHL_1(gb *Gameboy) bool { gb.CPU.SetHL(gb.CPU.GetHL() - 1); return true }
func DECSP_1(gb *Gameboy) bool { gb.CPU.Regs.SP--; return true }

var iduOp_2 = endNoop

func alun_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ORn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(OR(gb.CPU.Regs.A, gb.CPU.Regs.TempZ))
	return true
}

func ANDn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(AND(gb.CPU.Regs.A, gb.CPU.Regs.TempZ))
	return true
}

func XORn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(XOR(gb.CPU.Regs.A, gb.CPU.Regs.TempZ))
	return true
}

func SUBn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, false))
	return true
}

func SBCn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(SUB(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, gb.CPU.Regs.GetFlagC()))
	return true
}

func CPn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlags(SUB(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, false))
	return true
}

func ADDn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, false))
	return true
}

func ADCn_2(gb *Gameboy) bool {
	gb.CPU.Regs.SetFlagsAndA(ADD(gb.CPU.Regs.A, gb.CPU.Regs.TempZ, gb.CPU.Regs.GetFlagC()))
	return true
}

func ldrn_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldAn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.A) }
func ldBn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.B) }
func ldCn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.C) }
func ldDn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.D) }
func ldEn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.E) }
func ldHn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.H) }
func ldLn_2(gb *Gameboy) bool { return ldrn_2(gb, &gb.CPU.Regs.L) }

func ldrn_2(gb *Gameboy, reg *Data8) bool {
	gb.CPU.IncPC()
	*reg = gb.CPU.Regs.TempZ
	return true
}

func ldrhl_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	return false
}

func ldAhl_2(gb *Gameboy) bool { gb.CPU.Regs.A = gb.Data; return true }
func ldBhl_2(gb *Gameboy) bool { gb.CPU.Regs.B = gb.Data; return true }
func ldChl_2(gb *Gameboy) bool { gb.CPU.Regs.C = gb.Data; return true }
func ldDhl_2(gb *Gameboy) bool { gb.CPU.Regs.D = gb.Data; return true }
func ldEhl_2(gb *Gameboy) bool { gb.CPU.Regs.E = gb.Data; return true }
func ldHhl_2(gb *Gameboy) bool { gb.CPU.Regs.H = gb.Data; return true }
func ldLhl_2(gb *Gameboy) bool { gb.CPU.Regs.L = gb.Data; return true }

func ldahlinc_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.SetHL(gb.CPU.GetHL() + 1)
	return false
}

func ldahlinc_2(gb *Gameboy) bool {
	gb.CPU.Regs.A = gb.Data
	return true
}

func ldahldec_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.SetHL(gb.CPU.GetHL() - 1)
	return false
}

func ldahldec_2(gb *Gameboy) bool {
	gb.CPU.Regs.A = gb.Data
	return true
}

func ldhlAinc_1(gb *Gameboy) bool { return ldhlr_1(gb, gb.CPU.Regs.A, +1) }
func ldhlAdec_1(gb *Gameboy) bool { return ldhlr_1(gb, gb.CPU.Regs.A, -1) }
func ldhlA_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.A, 0) }
func ldhlB_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.B, 0) }
func ldhlC_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.C, 0) }
func ldhlD_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.D, 0) }
func ldhlE_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.E, 0) }
func ldhlH_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.H, 0) }
func ldhlL_1(gb *Gameboy) bool    { return ldhlr_1(gb, gb.CPU.Regs.L, 0) }

func ldhlr_1(gb *Gameboy, reg Data8, inc int16) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.CPU.SetHL(Data16(int16(gb.CPU.GetHL()) + inc))
	gb.WriteData(reg)
	return false
}

var ldhlr_2 = endNoop

func ldBCA_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetBC()))
	gb.WriteData(gb.CPU.Regs.A)
	return false
}
func ldDEA_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetDE()))
	gb.WriteData(gb.CPU.Regs.A)
	return false
}

var ldrra_2 = endNoop

func ldhca_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.C)))
	gb.WriteData(gb.CPU.Regs.A)
	return false
}

func ldhca_2(gb *Gameboy) bool {
	return true
}

func ldhac_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.C)))
	gb.CPU.Regs.A = gb.Data
	return false
}

func ldhac_2(gb *Gameboy) bool {
	return true
}

func ldnnsp_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldnnsp_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func ldnnsp_3(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
	gb.WriteData(gb.CPU.Regs.SP.LSB())
	gb.CPU.Regs.SetWZ(gb.CPU.Regs.GetWZ() + 1)
	return false
}

func ldnnsp_4(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
	gb.WriteData(gb.CPU.Regs.SP.MSB())
	return false
}

func ldnnsp_5(gb *Gameboy) bool {
	return true
}

func ldnna_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldnna_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func ldnna_3(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
	gb.WriteData(gb.CPU.Regs.A)
	return false
}

func ldnna_4(gb *Gameboy) bool {
	return true
}

func ldann_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldann_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func ldann_3(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.Regs.GetWZ()))
	gb.CPU.Regs.A = gb.Data
	return false
}

func ldann_4(gb *Gameboy) bool {
	return true
}

func ldhna_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldhna_2(gb *Gameboy) bool {
	gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.TempZ)))
	gb.CPU.IncPC()
	gb.WriteData(gb.CPU.Regs.A)
	return false
}

func ldhna_3(gb *Gameboy) bool {
	return true
}

func ldhan_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldhan_2(gb *Gameboy) bool {
	gb.WriteAddress(Addr(join16(0xff, gb.CPU.Regs.TempZ)))
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldhan_3(gb *Gameboy) bool {
	gb.CPU.Regs.A = gb.CPU.Regs.TempZ
	return true
}

func ldaBC_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetBC()))
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldaDE_1(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetDE()))
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldarr_2(gb *Gameboy) bool {
	gb.CPU.Regs.A = gb.CPU.Regs.TempZ
	return true
}

func addhlBC_1(gb *Gameboy) bool { return addhlrr_1(gb, gb.CPU.Regs.C) }
func addhlDE_1(gb *Gameboy) bool { return addhlrr_1(gb, gb.CPU.Regs.E) }
func addhlHL_1(gb *Gameboy) bool { return addhlrr_1(gb, gb.CPU.Regs.L) }

func addhlrr_1(gb *Gameboy, lo Data8) bool {
	result := ADD(gb.CPU.Regs.L, lo, false)
	gb.CPU.Regs.L = result.Value
	gb.CPU.Regs.SetFlagC(result.C)
	gb.CPU.Regs.SetFlagH(result.H)
	gb.CPU.Regs.SetFlagN(result.N)
	return false
}

func addhlBC_2(gb *Gameboy) bool { return addhlrr_2(gb, gb.CPU.Regs.B) }
func addhlDE_2(gb *Gameboy) bool { return addhlrr_2(gb, gb.CPU.Regs.D) }
func addhlHL_2(gb *Gameboy) bool { return addhlrr_2(gb, gb.CPU.Regs.H) }

func addhlrr_2(gb *Gameboy, hi Data8) bool {
	result := ADD(gb.CPU.Regs.H, hi, gb.CPU.Regs.GetFlagC())
	gb.CPU.Regs.H = result.Value
	gb.CPU.Regs.SetFlagC(result.C)
	gb.CPU.Regs.SetFlagH(result.H)
	gb.CPU.Regs.SetFlagN(result.N)
	return true
}

func addhlsp_1(gb *Gameboy) bool {
	result := ADD(gb.CPU.Regs.L, gb.CPU.Regs.SP.LSB(), false)
	gb.CPU.Regs.L = result.Value
	gb.CPU.Regs.SetFlagC(result.C)
	gb.CPU.Regs.SetFlagH(result.H)
	gb.CPU.Regs.SetFlagN(result.N)
	return false
}

func addhlsp_2(gb *Gameboy) bool {
	result := ADD(gb.CPU.Regs.H, gb.CPU.Regs.SP.MSB(), gb.CPU.Regs.GetFlagC())
	gb.CPU.Regs.H = result.Value
	gb.CPU.Regs.SetFlagC(result.C)
	gb.CPU.Regs.SetFlagH(result.H)
	gb.CPU.Regs.SetFlagN(result.N)
	return true
}

func addspe_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func addspe_2(gb *Gameboy) bool {
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
	return false
}

func addspe_3(gb *Gameboy) bool {
	res := gb.CPU.Regs.SP.MSB()
	if gb.CPU.Regs.TempW == 1 {
		res++
	} else if gb.CPU.Regs.TempW == 0xff {
		res--
	}
	gb.CPU.Regs.TempW = res
	return false
}

func addspe_4(gb *Gameboy) bool {
	gb.CPU.SetSP(Addr(gb.CPU.Regs.GetWZ()))
	return true
}

func ldxxnn_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldxxnn_2(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempW = gb.Data
	return false
}

func ldBCnn_3(gb *Gameboy) bool {
	gb.CPU.SetBC(gb.CPU.GetWZ())
	return true
}

func ldDEnn_3(gb *Gameboy) bool {
	gb.CPU.SetDE(gb.CPU.GetWZ())
	return true
}

func ldHLnn_3(gb *Gameboy) bool {
	gb.CPU.SetHL(gb.CPU.GetWZ())
	return true
}

func ldSPnn_3(gb *Gameboy) bool {
	gb.CPU.SetSP(Addr(gb.CPU.GetWZ()))
	return true
}

func ldhln_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldhln_2(gb *Gameboy) bool {
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.WriteData(gb.CPU.Regs.TempZ)
	return false
}

var ldhln_3 = endNoop

var ldsphl_1 = noop

func ldsphl_2(gb *Gameboy) bool {
	gb.CPU.Regs.SP = Addr(gb.CPU.GetHL())
	return true
}

func ldhlspe_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.Regs.TempZ = gb.Data
	return false
}

func ldhlspe_2(gb *Gameboy) bool {
	res := ADD(gb.CPU.Regs.SP.LSB(), gb.CPU.Regs.TempZ, false)
	gb.CPU.Regs.L = res.Value
	res.Z0 = true
	gb.CPU.Regs.SetFlags(res)
	return false
}

func ldhlspe_3(gb *Gameboy) bool {
	adj := Data8(0x00)
	if gb.CPU.Regs.TempZ&Bit7 != 0 {
		adj = 0xff
	}
	res := ADD(gb.CPU.Regs.SP.MSB(), adj, gb.CPU.Regs.GetFlagC())
	gb.CPU.Regs.H = res.Value
	return true
}

func rst_1(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.PC.MSB())
	return false
}

func rst_2(gb *Gameboy) bool {
	gb.CPU.SetSP(gb.CPU.Regs.SP - 1)
	gb.WriteAddress(gb.CPU.Regs.SP)
	gb.WriteData(gb.CPU.Regs.PC.LSB())
	return false
}

func rst_3_00(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0000)
	return false
}

func rst_3_08(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0008)
	return false
}

func rst_3_10(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0010)
	return false
}

func rst_3_18(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0018)
	return false
}

func rst_3_20(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0020)
	return false
}

func rst_3_28(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0028)
	return false
}

func rst_3_30(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0030)
	return false
}

func rst_3_38(gb *Gameboy) bool {
	gb.CPU.SetPC(0x0038)
	return false
}

func rst_4(gb *Gameboy) bool {
	return true
}

func NewCBOp(v Data8) CBOp {
	return CBOp{Op: cb((v & 0xf8) >> 3), Target: CBTarget(v & 0x7)}
}

func runCB_1(gb *Gameboy) bool {
	gb.WriteAddress(gb.CPU.Regs.PC)
	gb.CPU.IncPC()
	gb.CPU.CBOp = NewCBOp(gb.Data)
	return false
}

func runCB_2(gb *Gameboy) bool {
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
}

func runCB_3(gb *Gameboy) bool {
	if gb.CPU.CBOp.Op.Is3Cycles() {
		return true
	}
	gb.WriteAddress(Addr(gb.CPU.GetHL()))
	gb.WriteData(gb.CPU.Regs.TempZ)
	return false
}

func runCB_4(gb *Gameboy) bool {
	return true
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
