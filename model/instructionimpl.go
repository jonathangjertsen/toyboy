package model

import (
	"fmt"
	"slices"
)

type CycleHandler func(cpu *CPU) CycleHandler

// NOP
// Also be used to insert an empty cycle at the end of an instruction

func implNop(cpu *CPU) CycleHandler {
	return nil
}

// LD SP nn

func implLDSPnn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDSPnn_2
}

func implLDSPnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()

	return implLDSPnn_3
}

func implLDSPnn_3(cpu *CPU) CycleHandler {
	cpu.SetSP(cpu.Regs.Temp.Addr())
	return nil
}

// LD HL nn

func implLDHLnn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDHLnn_2
}

func implLDHLnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implLDHLnn_3
}

func implLDHLnn_3(cpu *CPU) CycleHandler {
	cpu.Regs.HL = cpu.Regs.Temp
	return nil
}

// LD HL, SP+e

func implLDHLSPe(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDHLSPe_2
}

func implLDHLSPe_2(cpu *CPU) CycleHandler {
	res := ADD(cpu.Regs.SP.LSB(), cpu.Regs.Temp.LSR, false)
	cpu.Regs.HL.LSR = res.Value
	res.Z0 = true
	cpu.Regs.SetFlags(res)
	return implLDHLSPe_3
}

func implLDHLSPe_3(cpu *CPU) CycleHandler {
	adj := Data8(0x00)
	if cpu.Regs.Temp.LSR&Bit7 != 0 {
		adj = 0xff
	}
	res := ADD(cpu.Regs.SP.MSB(), adj, cpu.Regs.GetFlagC())
	cpu.Regs.HL.MSR = res.Value
	return nil
}

// LD SP, HL

func implLDSPHL(*CPU) CycleHandler {
	return implLDSPHL_2
}

func implLDSPHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SP = cpu.Regs.HL.Addr()
	return nil
}

// LD (HL), (PC+n)

func implLDHLn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDHLn_2
}

func implLDHLn_2(cpu *CPU) CycleHandler {
	cpu.store(cpu.Regs.HL.Addr(), cpu.Regs.Temp.LSR)
	return implNop
}

// LD r, (HL)

func implLDAHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.A) }
func implLDBHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.BC.MSR) }
func implLDCHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.BC.LSR) }
func implLDDHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.DE.MSR) }
func implLDEHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.DE.LSR) }
func implLDHHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.HL.MSR) }
func implLDLHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.HL.LSR) }

func implLDrHL(cpu *CPU, r *Data8) CycleHandler {
	*r = cpu.load(cpu.Regs.HL.Addr())
	return implNop
}

// LD A, nn

func implLDAnn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDAnn_2
}

func implLDAnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implLDAnn_3
}

func implLDAnn_3(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(cpu.Regs.Temp.Addr())
	return implNop
}

// LD nn, A

func implLDnnA(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDnnA_2
}

func implLDnnA_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implLDnnA_3
}

func implLDnnA_3(cpu *CPU) CycleHandler {
	cpu.store(cpu.Regs.Temp.Addr(), cpu.Regs.A)
	return implNop
}

// LD nn, SP

func implLDnnSP(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implLDnnSP_2
}

func implLDnnSP_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implLDnnSP_3
}

func implLDnnSP_3(cpu *CPU) CycleHandler {
	cpu.store(cpu.Regs.Temp.Addr(), cpu.Regs.SP.LSB())
	cpu.Regs.Temp.Inc()
	return implLDnnSP_4
}

func implLDnnSP_4(cpu *CPU) CycleHandler {
	cpu.store(cpu.Regs.Temp.Addr(), cpu.Regs.SP.MSB())
	return implNop
}

// LD (HL), r

func implLDHLA(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.A) }
func implLDHLB(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.BC.MSR) }
func implLDHLC(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.BC.LSR) }
func implLDHLD(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.DE.MSR) }
func implLDHLE(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.DE.LSR) }
func implLDHLH(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.HL.MSR) }
func implLDHLL(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.HL.LSR) }

func implLDHLr(cpu *CPU, r Data8) CycleHandler {
	cpu.store(cpu.Regs.HL.Addr(), r)

	return implNop
}

// LD (HL+), A

func implLDHLAInc(cpu *CPU) CycleHandler {
	cpu.store(cpu.Regs.HL.Addr(), cpu.Regs.A)
	cpu.Regs.HL.Inc()

	return implNop
}

// LD (HL-), A

func implLDHLADec(cpu *CPU) CycleHandler {
	cpu.store(cpu.Regs.HL.Addr(), cpu.Regs.A)
	cpu.Regs.HL.Dec()

	return implNop
}

// LD A (C+$FF00)

func implLDHAC(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(Addr(join16(0xff, cpu.Regs.BC.LSR)))

	return implNop
}

// LD (C+$FF00), A

func implLDHCA(cpu *CPU) CycleHandler {
	cpu.store(Addr(join16(0xff, cpu.Regs.BC.LSR)), cpu.Regs.A)

	return implNop
}

// LD (n+$FF00), A

func implLDHnA(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()

	return implLDHnA_2
}

func implLDHnA_2(cpu *CPU) CycleHandler {
	cpu.store(Addr(join16(0xff, cpu.Regs.Temp.LSR)), cpu.Regs.A)

	return implNop
}

// LD A,(n+$FF00)

func implLDHAn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()

	return implLDHAn_2
}

func implLDHAn_2(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(Addr(join16(0xff, cpu.Regs.Temp.LSR)))

	return implNop
}

// LD r, r'

func implLDAA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.A) }
func implLDBA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.A) }
func implLDCA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.A) }
func implLDDA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.A) }
func implLDEA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.A) }
func implLDHA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.A) }
func implLDLA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.A) }

func implLDAB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.BC.MSR) }
func implLDBB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.BC.MSR) }
func implLDCB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.BC.MSR) }
func implLDDB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.BC.MSR) }
func implLDEB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.BC.MSR) }
func implLDHB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.BC.MSR) }
func implLDLB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.BC.MSR) }

func implLDAC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.BC.LSR) }
func implLDBC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.BC.LSR) }
func implLDCC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.BC.LSR) }
func implLDDC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.BC.LSR) }
func implLDEC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.BC.LSR) }
func implLDHC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.BC.LSR) }
func implLDLC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.BC.LSR) }

func implLDAD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.DE.MSR) }
func implLDBD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.DE.MSR) }
func implLDCD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.DE.MSR) }
func implLDDD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.DE.MSR) }
func implLDED(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.DE.MSR) }
func implLDHD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.DE.MSR) }
func implLDLD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.DE.MSR) }

func implLDAE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.DE.LSR) }
func implLDBE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.DE.LSR) }
func implLDCE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.DE.LSR) }
func implLDDE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.DE.LSR) }
func implLDEE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.DE.LSR) }
func implLDHE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.DE.LSR) }
func implLDLE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.DE.LSR) }

func implLDAH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.HL.MSR) }
func implLDBH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.HL.MSR) }
func implLDCH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.HL.MSR) }
func implLDDH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.HL.MSR) }
func implLDEH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.HL.MSR) }
func implLDHH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.HL.MSR) }
func implLDLH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.HL.MSR) }

func implLDAL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.HL.LSR) }
func implLDBL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.MSR, cpu.Regs.HL.LSR) }
func implLDCL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.BC.LSR, cpu.Regs.HL.LSR) }
func implLDDL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.MSR, cpu.Regs.HL.LSR) }
func implLDEL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.DE.LSR, cpu.Regs.HL.LSR) }
func implLDHL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.MSR, cpu.Regs.HL.LSR) }
func implLDLL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.HL.LSR, cpu.Regs.HL.LSR) }

func implLDrr(dest *Data8, src Data8) CycleHandler {
	*dest = src

	return nil
}

// LD r, n

func implLDAn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.A) }
func implLDBn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.BC.MSR) }
func implLDCn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.BC.LSR) }
func implLDDn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.DE.MSR) }
func implLDEn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.DE.LSR) }
func implLDHn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.HL.MSR) }
func implLDLn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.HL.LSR) }

func implLDrn(cpu *CPU, reg *Data8) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	cpu.Regs.TempPtr[0] = reg

	return implLDrn_2
}

func implLDrn_2(cpu *CPU) CycleHandler {
	*cpu.Regs.TempPtr[0] = cpu.Regs.Temp.LSR
	return nil
}

// LD (rr'), a

func implLDBCA(cpu *CPU) CycleHandler { return implLDrrA(cpu, cpu.Regs.BC.MSR, cpu.Regs.BC.LSR) }
func implLDDEA(cpu *CPU) CycleHandler { return implLDrrA(cpu, cpu.Regs.DE.MSR, cpu.Regs.DE.LSR) }

func implLDrrA(cpu *CPU, msr, lsr Data8) CycleHandler {
	cpu.store(Addr(join16(msr, lsr)), cpu.Regs.A)
	return implNop
}

// LD (rr'), nn

func implLDBCnn(cpu *CPU) CycleHandler { return implLDxxnn(cpu, &cpu.Regs.BC.MSR, &cpu.Regs.BC.LSR) }
func implLDDEnn(cpu *CPU) CycleHandler { return implLDxxnn(cpu, &cpu.Regs.DE.MSR, &cpu.Regs.DE.LSR) }

func implLDxxnn(cpu *CPU, msr, lsr *Data8) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	cpu.Regs.TempPtr[0] = msr
	cpu.Regs.TempPtr[1] = lsr

	return implLDxxnn_2
}

func implLDxxnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()

	return implLDxxnn_3
}

func implLDxxnn_3(cpu *CPU) CycleHandler {
	*cpu.Regs.TempPtr[0] = cpu.Regs.Temp.MSR
	*cpu.Regs.TempPtr[1] = cpu.Regs.Temp.LSR
	return nil
}

// LD A, (rr')

func implLDABC(cpu *CPU) CycleHandler { return implLDArr(cpu, cpu.Regs.BC.MSR, cpu.Regs.BC.LSR) }
func implLDADE(cpu *CPU) CycleHandler { return implLDArr(cpu, cpu.Regs.DE.MSR, cpu.Regs.DE.LSR) }

func implLDArr(cpu *CPU, msr, lsr Data8) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(Addr(join16(msr, lsr)))
	return implLDArr_2
}

func implLDArr_2(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.Regs.Temp.LSR
	return nil
}

// LD A, (HL+)

func implLDAHLInc(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(cpu.Regs.HL.Addr())
	cpu.Regs.HL.Inc()
	return implNop
}

// LD A, (HL-)

func implLDAHLDec(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(cpu.Regs.HL.Addr())
	cpu.Regs.HL.Dec()
	return implNop
}

// JR PC+e

func implJRe(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implJRe_2
}

func implJRe_2(cpu *CPU) CycleHandler {
	newPC := Data16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.Temp.LSR)))
	cpu.Regs.Temp.Set(newPC)
	return implJRe_3
}

func implJRe_3(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return nil
}

// JR cc, PC+e

func implJRZe(cpu *CPU) CycleHandler {
	return implJRcce(cpu, cpu.Regs.GetFlagZ())
}

func implJRNZe(cpu *CPU) CycleHandler {
	return implJRcce(cpu, !cpu.Regs.GetFlagZ())
}

func implJRCe(cpu *CPU) CycleHandler {
	return implJRcce(cpu, cpu.Regs.GetFlagC())
}

func implJRNCe(cpu *CPU) CycleHandler {
	return implJRcce(cpu, !cpu.Regs.GetFlagC())
}

func implJRcce(cpu *CPU, cond bool) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()

	if !cond {
		cpu.lastBranchResult = -1
		return implNop
	}

	cpu.lastBranchResult = +1
	return implJRcce_2
}

func implJRcce_2(cpu *CPU) CycleHandler {
	newPC := Data16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.Temp.LSR)))
	cpu.Regs.Temp.Set(newPC)
	return implJRcce_3
}

func implJRcce_3(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return nil
}

// JP HL

func implJPHL(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.HL.Addr())
	return nil
}

// JP nn

func implJPnn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implJPnn_2
}

func implJPnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implJPnn_3
}

func implJPnn_3(cpu *CPU) CycleHandler {
	cpu.lastBranchResult = +1
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return implNop
}

// JP cc, nn

func implJPZnn(cpu *CPU) CycleHandler {
	return implJPccnn(cpu, cpu.Regs.GetFlagZ())
}

func implJPNZnn(cpu *CPU) CycleHandler {
	return implJPccnn(cpu, !cpu.Regs.GetFlagZ())
}

func implJPCnn(cpu *CPU) CycleHandler {
	return implJPccnn(cpu, cpu.Regs.GetFlagC())
}

func implJPNCnn(cpu *CPU) CycleHandler {
	return implJPccnn(cpu, !cpu.Regs.GetFlagC())
}

func implJPccnn(cpu *CPU, cond bool) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	cpu.Regs.TempCond = cond
	return implJPccnn_2
}

func implJPccnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	if !cpu.Regs.TempCond {
		cpu.lastBranchResult = -1
		return nil
	}
	cpu.lastBranchResult = +1
	return implJPccnn_3
}

func implJPccnn_3(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return implNop
}

// CALL nn

func implCALLnn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implCALLnn_2
}

func implCALLnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implCALLnn_3
}

func implCALLnn_3(*CPU) CycleHandler {
	return implCALLnn_4
}

func implCALLnn_4(cpu *CPU) CycleHandler {
	cpu.push((cpu.Regs.PC + 1).MSB())
	return implCALLnn_5
}

func implCALLnn_5(cpu *CPU) CycleHandler {
	cpu.push((cpu.Regs.PC + 1).LSB())
	return implCALLnn_6
}

func implCALLnn_6(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return nil
}

// CALL cc nn
func implCALLZnn(cpu *CPU) CycleHandler  { return implCALLccnn(cpu, cpu.Regs.GetFlagZ()) }
func implCALLNZnn(cpu *CPU) CycleHandler { return implCALLccnn(cpu, !cpu.Regs.GetFlagZ()) }
func implCALLCnn(cpu *CPU) CycleHandler  { return implCALLccnn(cpu, cpu.Regs.GetFlagC()) }
func implCALLNCnn(cpu *CPU) CycleHandler { return implCALLccnn(cpu, !cpu.Regs.GetFlagC()) }

func implCALLccnn(cpu *CPU, cond bool) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	cpu.Regs.TempCond = cond
	return implCALLccnn_2
}

func implCALLccnn_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.fetchByteAtPC()
	return implCALLccnn_3
}

func implCALLccnn_3(cpu *CPU) CycleHandler {
	if !cpu.Regs.TempCond {
		return nil
	}
	return implCALLccnn_4
}

func implCALLccnn_4(cpu *CPU) CycleHandler {
	cpu.push((cpu.Regs.PC + 1).MSB())
	return implCALLccnn_5
}

func implCALLccnn_5(cpu *CPU) CycleHandler {
	cpu.push((cpu.Regs.PC + 1).LSB())
	return implCALLccnn_6
}

func implCALLccnn_6(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return nil
}

// RST n

func implRST0x00(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x00) }
func implRST0x08(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x08) }
func implRST0x10(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x10) }
func implRST0x18(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x18) }
func implRST0x20(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x20) }
func implRST0x28(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x28) }
func implRST0x30(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x30) }
func implRST0x38(cpu *CPU) CycleHandler { return implRSTn(cpu, 0x38) }

func implRSTn(cpu *CPU, vec Data8) CycleHandler {
	cpu.push((cpu.Regs.PC + 1).MSB())
	cpu.Regs.Temp.LSR = vec
	cpu.Regs.Temp.MSR = 0
	return implRSTn_2
}

func implRSTn_2(cpu *CPU) CycleHandler {
	cpu.push((cpu.Regs.PC + 1).LSB())
	return implRSTn_3
}

func implRSTn_3(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return nil
}

// RET

func implRET(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.pop()
	return implRET_2
}

func implRET_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.pop()
	return implRET_3
}

func implRET_3(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	return implNop
}

// RETI

func implRETI(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.pop()
	return implRETI_2
}

func implRETI_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.pop()
	return implRETI_3
}

func implRETI_3(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr())
	// TODO verify if this is the right cycle
	if cpu.Interrupts != nil {
		cpu.Interrupts.IME = true
	}
	return implNop
}

// RET cc

func implRETZ(cpu *CPU) CycleHandler {
	return implRETcc(cpu, cpu.Regs.GetFlagZ())
}

func implRETNZ(cpu *CPU) CycleHandler {
	return implRETcc(cpu, !cpu.Regs.GetFlagZ())
}

func implRETC(cpu *CPU) CycleHandler {
	return implRETcc(cpu, cpu.Regs.GetFlagC())
}

func implRETNC(cpu *CPU) CycleHandler {
	return implRETcc(cpu, !cpu.Regs.GetFlagC())
}

func implRETcc(cpu *CPU, cond bool) CycleHandler {
	if !cond {
		cpu.lastBranchResult = -1
		return implNop
	}
	cpu.lastBranchResult = +1
	return implRETcc_2
}

func implRETcc_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.pop()
	return implRETcc_3
}

func implRETcc_3(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.pop()
	return implRETcc_4
}

func implRETcc_4(cpu *CPU) CycleHandler {
	cpu.SetPC(cpu.Regs.Temp.Addr() - 1)
	return implNop
}

// PUSH rr'

func implPUSHAF(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.A, cpu.Regs.F) }
func implPUSHBC(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.BC.MSR, cpu.Regs.BC.LSR) }
func implPUSHDE(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.DE.MSR, cpu.Regs.DE.LSR) }
func implPUSHHL(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.HL.MSR, cpu.Regs.HL.LSR) }

func implPUSHrr(cpu *CPU, msb, lsb Data8) CycleHandler {
	cpu.push(msb)
	cpu.Regs.Temp.LSR = lsb
	return implPUSHrr_2
}

func implPUSHrr_2(cpu *CPU) CycleHandler {
	cpu.push(cpu.Regs.Temp.LSR)
	return implPUSHrr_3
}

func implPUSHrr_3(*CPU) CycleHandler {
	return implNop
}

// POP rr'

func implPOPAF(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.A, &cpu.Regs.F) }
func implPOPBC(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.BC.MSR, &cpu.Regs.BC.LSR) }
func implPOPDE(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.DE.MSR, &cpu.Regs.DE.LSR) }
func implPOPHL(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.HL.MSR, &cpu.Regs.HL.LSR) }

func implPOPrr(cpu *CPU, msb, lsb *Data8) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.pop()

	cpu.Regs.TempPtr[0] = msb
	cpu.Regs.TempPtr[1] = lsb

	return implPOPrr_2
}

func implPOPrr_2(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.MSR = cpu.pop()

	return implPOPrr_3
}

func implPOPrr_3(cpu *CPU) CycleHandler {
	msb := cpu.Regs.TempPtr[0]
	lsb := cpu.Regs.TempPtr[1]

	*msb = cpu.Regs.Temp.MSR
	if lsb == &cpu.Regs.F {
		*lsb = cpu.Regs.Temp.LSR & 0xf0
	} else {
		*lsb = cpu.Regs.Temp.LSR
	}

	return nil
}

// RLA

func implRLA(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(RLA(cpu.Regs.A, cpu.Regs.GetFlagC()))
	return nil
}

func implRRA(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(RRA(cpu.Regs.A, cpu.Regs.GetFlagC()))
	return nil
}

func implRLCA(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(RLCA(cpu.Regs.A))
	return nil
}

func implRRCA(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(RRCA(cpu.Regs.A))
	return nil
}

// AND A, r

func implANDA(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.A) }
func implANDB(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.BC.MSR) }
func implANDC(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.BC.LSR) }
func implANDD(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.DE.MSR) }
func implANDE(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.DE.LSR) }
func implANDH(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.HL.MSR) }
func implANDL(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.HL.LSR) }

func implANDr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, r))
	return nil
}

// AND A, n

func implANDn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implANDn_2
}

func implANDn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, cpu.Regs.Temp.LSR))
	return nil
}

// AND A, (HL)

func implANDHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implANDHL_2
}

func implANDHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, cpu.Regs.Temp.LSR))
	return nil
}

// CP A, r

func implCPA(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.A) }
func implCPB(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.BC.MSR) }
func implCPC(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.BC.LSR) }
func implCPD(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.DE.MSR) }
func implCPE(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.DE.LSR) }
func implCPH(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.HL.MSR) }
func implCPL(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.HL.LSR) }

func implCPr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlags(SUB(cpu.Regs.A, r, false))
	return nil
}

// CP A, n

func implCPn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implCPn_2
}

func implCPn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlags(SUB(cpu.Regs.A, cpu.Regs.Temp.LSR, false))
	return nil
}

// CP A, (HL)

func implCPHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implCPHL_2
}

func implCPHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlags(SUB(cpu.Regs.A, cpu.Regs.Temp.LSR, false))
	return nil
}

// SUB A, r

func implSUBA(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.A) }
func implSUBB(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.BC.MSR) }
func implSUBC(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.BC.LSR) }
func implSUBD(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.DE.MSR) }
func implSUBE(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.DE.LSR) }
func implSUBH(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.HL.MSR) }
func implSUBL(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.HL.LSR) }

func implSUBr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, r, false))
	return nil
}

// SUB A, n

func implSUBn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implSUBn_2
}

func implSUBn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.Temp.LSR, false))
	return nil
}

// SUB A, (HL)

func implSUBHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implSUBHL_2
}

func implSUBHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.Temp.LSR, false))
	return nil
}

// SBC A, r

func implSBCA(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.A) }
func implSBCB(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.BC.MSR) }
func implSBCC(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.BC.LSR) }
func implSBCD(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.DE.MSR) }
func implSBCE(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.DE.LSR) }
func implSBCH(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.HL.MSR) }
func implSBCL(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.HL.LSR) }

func implSBCr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, r, cpu.Regs.GetFlagC()))
	return nil
}

// SBC A, n

func implSBCn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implSBCn_2
}

func implSBCn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.Temp.LSR, cpu.Regs.GetFlagC()))
	return nil
}

// SBC A, (HL)

func implSBCHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implSBCHL_2
}

func implSBCHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.Temp.LSR, cpu.Regs.GetFlagC()))
	return nil
}

// ADD A, r

func implADDA(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.A) }
func implADDB(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.BC.MSR) }
func implADDC(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.BC.LSR) }
func implADDD(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.DE.MSR) }
func implADDE(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.DE.LSR) }
func implADDH(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.HL.MSR) }
func implADDL(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.HL.LSR) }

func implADDr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, r, false))
	return nil
}

// ADD A, n

func implADDn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implADDn_2
}

func implADDn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.Temp.LSR, false))
	return nil
}

// ADD A, (HL)

func implADDHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implADDHL_2
}

func implADDHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.Temp.LSR, false))
	return nil
}

// ADD HL, rr'

func implADDHLHL(cpu *CPU) CycleHandler { return implADDHLrr(cpu, cpu.Regs.HL.MSR, cpu.Regs.HL.LSR) }
func implADDHLBC(cpu *CPU) CycleHandler { return implADDHLrr(cpu, cpu.Regs.BC.MSR, cpu.Regs.BC.LSR) }
func implADDHLDE(cpu *CPU) CycleHandler { return implADDHLrr(cpu, cpu.Regs.DE.MSR, cpu.Regs.DE.LSR) }
func implADDHLSP(cpu *CPU) CycleHandler {
	return implADDHLrr(cpu, cpu.Regs.SP.MSB(), cpu.Regs.SP.LSB())
}

func implADDHLrr(cpu *CPU, msr, lsr Data8) CycleHandler {
	result := ADD(cpu.Regs.HL.LSR, lsr, false)
	cpu.Regs.Temp.LSR = msr
	cpu.Regs.HL.LSR = result.Value
	cpu.Regs.SetFlagC(result.C)
	cpu.Regs.SetFlagH(result.H)
	cpu.Regs.SetFlagN(result.N)
	return implADDHLrr_2
}

func implADDHLrr_2(cpu *CPU) CycleHandler {
	result := ADD(cpu.Regs.HL.MSR, cpu.Regs.Temp.LSR, cpu.Regs.GetFlagC())
	cpu.Regs.HL.MSR = result.Value
	cpu.Regs.SetFlagC(result.C)
	cpu.Regs.SetFlagH(result.H)
	cpu.Regs.SetFlagN(result.N)
	return nil
}

// ADD SP, PC+e

func implADDSPe(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implADDSPe_2
}

func implADDSPe_2(cpu *CPU) CycleHandler {
	zSign := cpu.Regs.Temp.LSR&Bit7 != 0
	result := ADD(cpu.Regs.SP.LSB(), cpu.Regs.Temp.LSR, false)
	cpu.Regs.Temp.LSR = result.Value
	cpu.Regs.Temp.MSR = 0
	cpu.Regs.SetFlags(result)
	cpu.Regs.SetFlagZ(false)
	if c := cpu.Regs.GetFlagC(); c && !zSign {
		cpu.Regs.Temp.MSR = 1
	} else if !c && zSign {
		cpu.Regs.Temp.MSR = 0xff
	}
	return implADDSPe_3
}

func implADDSPe_3(cpu *CPU) CycleHandler {
	res := cpu.Regs.SP.MSB()
	if cpu.Regs.Temp.MSR == 1 {
		res++
	} else if cpu.Regs.Temp.MSR == 0xff {
		res--
	}
	cpu.Regs.Temp.MSR = res
	return implADDSPe_4
}

func implADDSPe_4(cpu *CPU) CycleHandler {
	cpu.SetSP(cpu.Regs.Temp.Addr())
	return nil
}

// ADC A, r

func implADCA(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.A) }
func implADCB(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.BC.MSR) }
func implADCC(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.BC.LSR) }
func implADCD(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.DE.MSR) }
func implADCE(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.DE.LSR) }
func implADCH(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.HL.MSR) }
func implADCL(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.HL.LSR) }

func implADCr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, r, cpu.Regs.GetFlagC()))
	return nil
}

// ADC A, n

func implADCn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implADCn_2
}

func implADCn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.Temp.LSR, cpu.Regs.GetFlagC()))
	return nil
}

// ADC A, (HL)

func implADCHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implADCHL_2
}

func implADCHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.Temp.LSR, cpu.Regs.GetFlagC()))
	return nil
}

// OR A, r

func implORA(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.A) }
func implORB(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.BC.MSR) }
func implORC(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.BC.LSR) }
func implORD(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.DE.MSR) }
func implORE(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.DE.LSR) }
func implORH(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.HL.MSR) }
func implORL(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.HL.LSR) }

func implORr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, r))
	return nil
}

// OR A, n

func implORn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implORn_2
}

func implORn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, cpu.Regs.Temp.LSR))
	return nil
}

// OR A, (HL)

func implORHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implORHL_2
}

func implORHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, cpu.Regs.Temp.LSR))
	return nil
}

// XOR A, r

func implXORA(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.A) }
func implXORB(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.BC.MSR) }
func implXORC(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.BC.LSR) }
func implXORD(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.DE.MSR) }
func implXORE(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.DE.LSR) }
func implXORH(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.HL.MSR) }
func implXORL(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.HL.LSR) }

func implXORr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, r))
	return nil
}

// XOR A, n

func implXORn(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()
	return implXORn_2
}

func implXORn_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, cpu.Regs.Temp.LSR))
	return nil
}

// XOR A, (HL)

func implXORHL(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implXORHL_2
}

func implXORHL_2(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, cpu.Regs.Temp.LSR))
	return nil
}

// INC r

func implINCA(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.A) }
func implINCB(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.BC.MSR) }
func implINCC(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.BC.LSR) }
func implINCD(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.DE.MSR) }
func implINCE(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.DE.LSR) }
func implINCH(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.HL.MSR) }
func implINCL(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.HL.LSR) }

func implINCr(cpu *CPU, r *Data8) CycleHandler {
	result := ADD(*r, 1, false)
	*r = result.Value
	cpu.Regs.SetFlagZ(result.Z())
	cpu.Regs.SetFlagH(result.H)
	cpu.Regs.SetFlagN(result.N)
	return nil
}

// INC rr'

func implINCBC(cpu *CPU) CycleHandler {
	cpu.Regs.BC.Inc()
	return implNop
}

func implINCDE(cpu *CPU) CycleHandler {
	cpu.Regs.DE.Inc()
	return implNop
}

func implINCHL(cpu *CPU) CycleHandler {
	cpu.Regs.HL.Inc()
	return implNop
}

func implINCSP(cpu *CPU) CycleHandler {
	cpu.Regs.SP++
	return implNop
}

// INC (HL)

func implINCHLInd(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implINCHLInd_2
}

func implINCHLInd_2(cpu *CPU) CycleHandler {
	res := ADD(cpu.Regs.Temp.LSR, 1, false)
	// apparently doesn't set C.
	cpu.Regs.SetFlagH(res.H)
	cpu.Regs.SetFlagZ(res.Z())
	cpu.Regs.SetFlagN(res.N)
	cpu.Bus.WriteAddress(cpu.Regs.HL.Addr())
	cpu.Regs.Temp.LSR = res.Value
	cpu.Bus.WriteData(cpu.Regs.Temp.LSR)
	return implNop
}

// DEC r

func implDECA(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.A) }
func implDECB(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.BC.MSR) }
func implDECC(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.BC.LSR) }
func implDECD(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.DE.MSR) }
func implDECE(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.DE.LSR) }
func implDECH(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.HL.MSR) }
func implDECL(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.HL.LSR) }

func implDECr(cpu *CPU, r *Data8) CycleHandler {
	result := SUB(*r, 1, false)
	*r = result.Value
	cpu.Regs.SetFlagZ(result.Z())
	cpu.Regs.SetFlagH(result.H)
	cpu.Regs.SetFlagN(result.N)
	return nil
}

// DEC rr'

func implDECBC(cpu *CPU) CycleHandler {
	cpu.Regs.BC.Dec()
	return implNop
}

func implDECDE(cpu *CPU) CycleHandler {
	cpu.Regs.DE.Dec()
	return implNop
}

func implDECHL(cpu *CPU) CycleHandler {
	cpu.Regs.HL.Dec()
	return implNop
}

func implDECSP(cpu *CPU) CycleHandler {
	cpu.Regs.SP--
	return implNop
}

// DEC (HL)

func implDECHLInd(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.load(cpu.Regs.HL.Addr())
	return implDECHLInd_2
}

func implDECHLInd_2(cpu *CPU) CycleHandler {
	res := SUB(cpu.Regs.Temp.LSR, 1, false)
	// apparently doesn't set C.
	cpu.Regs.SetFlagH(res.H)
	cpu.Regs.SetFlagZ(res.Z())
	cpu.Regs.SetFlagN(res.N)
	cpu.Bus.WriteAddress(cpu.Regs.HL.Addr())
	cpu.Regs.Temp.LSR = res.Value
	cpu.Bus.WriteData(cpu.Regs.Temp.LSR)
	return implNop
}

// DAA

func implDAA(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagsAndA(DAA(cpu.Regs.A, cpu.Regs.GetFlagC(), cpu.Regs.GetFlagN(), cpu.Regs.GetFlagH()))
	return nil
}

// CPL

func implCPLaka2f(cpu *CPU) CycleHandler {
	cpu.Regs.A ^= 0xff
	cpu.Regs.SetFlagN(true)
	cpu.Regs.SetFlagH(true)
	return nil
}

// CCF

func implCCF(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagC(!cpu.Regs.GetFlagC())
	cpu.Regs.SetFlagN(false)
	cpu.Regs.SetFlagH(false)
	return nil
}

// SCF

func implSCF(cpu *CPU) CycleHandler {
	cpu.Regs.SetFlagC(true)
	cpu.Regs.SetFlagN(false)
	cpu.Regs.SetFlagH(false)
	return nil
}

// DI

func implDI(cpu *CPU) CycleHandler {
	if cpu.Interrupts != nil {
		cpu.Interrupts.setIMENextCycle = false
		cpu.Interrupts.SetIME(false)
	}
	return nil
}

// EI

func implEI(cpu *CPU) CycleHandler {
	if cpu.Interrupts != nil {
		cpu.Interrupts.setIMENextCycle = true
	}
	return nil
}

// HALT

func implHALT(cpu *CPU) CycleHandler {
	cpu.halted = true
	return nil
}

// Undefined instructions

func implUndefD3(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefDB(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefDD(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefE3(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefE4(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefEB(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefEC(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefED(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefF4(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefFC(cpu *CPU) CycleHandler { return implUndefined(cpu) }
func implUndefFD(cpu *CPU) CycleHandler { return implUndefined(cpu) }

func implUndefined(cpu *CPU) CycleHandler {
	panicf("hit undefined instruction %s (0x%02x)", cpu.Regs.IR, uint8(cpu.Regs.IR))

	// Real GB does something like this; system stops
	return implUndefined
}

// STOP

func implSTOP(cpu *CPU) CycleHandler {
	fmt.Printf("WARNING: Emulating STOP with HALT")
	return implHALT(cpu)
}

// CB

func implCB(cpu *CPU) CycleHandler {
	cpu.Regs.Temp.LSR = cpu.fetchByteAtPC()

	return implCB_2
}

func implCB_2(cpu *CPU) CycleHandler {
	cbOp := NewCBOp(cpu.Regs.Temp.LSR)

	if cbOp.Target == CBTargetIndirectHL {
		cpu.Bus.WriteAddress(cpu.Regs.HL.Addr())
	}
	var val Data8
	switch cbOp.Target {
	case CBTargetB:
		val = cpu.Regs.BC.MSR
	case CBTargetC:
		val = cpu.Regs.BC.LSR
	case CBTargetD:
		val = cpu.Regs.DE.MSR
	case CBTargetE:
		val = cpu.Regs.DE.LSR
	case CBTargetH:
		val = cpu.Regs.HL.MSR
	case CBTargetL:
		val = cpu.Regs.HL.LSR
	case CBTargetIndirectHL:
		val = cpu.Bus.GetData()
	case CBTargetA:
		val = cpu.Regs.A
	default:
		panic("unknown CBOp target")
	}
	val = cpu.doCBOp(val, cbOp.Op)
	switch cbOp.Target {
	case CBTargetB:
		cpu.Regs.BC.MSR = val
	case CBTargetC:
		cpu.Regs.BC.LSR = val
	case CBTargetD:
		cpu.Regs.DE.MSR = val
	case CBTargetE:
		cpu.Regs.DE.LSR = val
	case CBTargetH:
		cpu.Regs.HL.MSR = val
	case CBTargetL:
		cpu.Regs.HL.LSR = val
	case CBTargetIndirectHL:
		cpu.Regs.Temp.LSR = val
		if cbOp.Op.Is3Cycles() {
			return implNop
		}
		return implCB_3
	case CBTargetA:
		cpu.Regs.A = val
	default:
		panic("unknown CBOp target")
	}
	return nil
}

func implCB_3(cpu *CPU) CycleHandler {
	cpu.Bus.WriteAddress(cpu.Regs.HL.Addr())
	cpu.Bus.WriteData(cpu.Regs.Temp.LSR)
	return implNop
}

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
func NewCBOp(v Data8) CBOp {
	return CBOp{Op: cb((v & 0xf8) >> 3), Target: CBTarget(v & 0x7)}
}

func (cpu *CPU) doCBOp(val Data8, cbOp cb) Data8 {
	switch cbOp {
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
		panicf("unknown op = %+v", cbOp)
	}
	return val
}

func (cpu *CPU) cbbit(val, mask Data8) {
	cpu.Regs.SetFlagZ(val&mask == 0)
	cpu.Regs.SetFlagN(false)
	cpu.Regs.SetFlagH(true)
}

// Helpers

func (cpu *CPU) push(v Data8) {
	cpu.SetSP(cpu.Regs.SP - 1)
	cpu.store(cpu.Regs.SP, v)
}

func (cpu *CPU) pop() Data8 {
	v := cpu.load(cpu.Regs.SP)
	cpu.SetSP(cpu.Regs.SP + 1)
	return v
}

func (cpu *CPU) load(addr Addr) Data8 {
	cpu.Bus.WriteAddress(addr)
	return cpu.Bus.GetData()
}

func (cpu *CPU) fetchByteAtPC() Data8 {
	data := cpu.load(cpu.Regs.PC)
	cpu.IncPC()
	return data
}

func (cpu *CPU) store(addr Addr, v Data8) {
	cpu.Bus.WriteAddress(addr)
	cpu.Bus.WriteData(v)
}
