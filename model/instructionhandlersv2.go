package model

type CycleHandler func(cpu *CPU) CycleHandler

// NOP
// Also be used to insert an empty cycle at the end of an instruction

func implNop(cpu *CPU) CycleHandler {
	return nil
}

// LD SP nn

func implLDSPnn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()

	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()

		return func(cpu *CPU) CycleHandler {
			cpu.SetSP(Addr(cpu.Regs.GetWZ()))
			return nil
		}
	}
}

// LD HL nn

func implLDHLnn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()

	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()

		return func(cpu *CPU) CycleHandler {
			cpu.SetHL(cpu.Regs.GetWZ())
			return nil
		}
	}
}

// LD HL, SP+e

func implLDHLSPe(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		res := ADD(cpu.Regs.SP.LSB(), cpu.Regs.TempZ, false)
		cpu.Regs.L = res.Value
		res.Z0 = true
		cpu.Regs.SetFlags(res)
		return func(cpu *CPU) CycleHandler {
			adj := Data8(0x00)
			if cpu.Regs.TempZ&Bit7 != 0 {
				adj = 0xff
			}
			res := ADD(cpu.Regs.SP.MSB(), adj, cpu.Regs.GetFlagC())
			cpu.Regs.H = res.Value
			return nil
		}
	}
}

// LD SP, HL

func implLDSPHL(*CPU) CycleHandler {
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SP = Addr(cpu.GetHL())
		return nil
	}
}

// LD (HL), (PC+n)

func implLDHLn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.store(Addr(cpu.GetHL()), cpu.Regs.TempZ)
		return implNop
	}
}

// LD r, (HL)

func implLDAHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.A) }
func implLDBHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.B) }
func implLDCHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.C) }
func implLDDHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.D) }
func implLDEHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.E) }
func implLDHHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.H) }
func implLDLHL(cpu *CPU) CycleHandler { return implLDrHL(cpu, &cpu.Regs.L) }

func implLDrHL(cpu *CPU, r *Data8) CycleHandler {
	*r = cpu.load(Addr(cpu.GetHL()))
	return implNop
}

// LD A, nn

func implLDAnn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		return func(cpu *CPU) CycleHandler {
			cpu.Regs.A = cpu.load(Addr(cpu.Regs.GetWZ()))
			return implNop
		}
	}
}

// LD nn, A

func implLDnnA(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		return func(cpu *CPU) CycleHandler {
			cpu.store(Addr(cpu.Regs.GetWZ()), cpu.Regs.A)
			return implNop
		}
	}
}

// LD nn, SP

func implLDnnSP(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		return func(cpu *CPU) CycleHandler {
			cpu.store(Addr(cpu.Regs.GetWZ()), cpu.Regs.SP.LSB())
			cpu.Regs.TempZ++
			if cpu.Regs.TempZ == 0 {
				cpu.Regs.TempW++
			}
			return func(cpu *CPU) CycleHandler {
				cpu.store(Addr(cpu.Regs.GetWZ()), cpu.Regs.SP.MSB())
				return implNop
			}
		}
	}
}

// LD (HL), r

func implLDHLA(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.A) }
func implLDHLB(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.B) }
func implLDHLC(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.C) }
func implLDHLD(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.D) }
func implLDHLE(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.E) }
func implLDHLH(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.H) }
func implLDHLL(cpu *CPU) CycleHandler { return implLDHLr(cpu, cpu.Regs.L) }

func implLDHLr(cpu *CPU, r Data8) CycleHandler {
	cpu.store(Addr(cpu.GetHL()), r)
	return implNop
}

// LD (HL+), A

func implLDHLAInc(cpu *CPU) CycleHandler {
	cpu.store(Addr(cpu.GetHL()), cpu.Regs.A)
	cpu.SetHL(cpu.GetHL() + 1)
	return implNop
}

// LD (HL-), A

func implLDHLADec(cpu *CPU) CycleHandler {
	cpu.store(Addr(cpu.GetHL()), cpu.Regs.A)
	cpu.SetHL(cpu.GetHL() - 1)
	return implNop
}

// LD A (C+$FF00)

func implLDHAC(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(Addr(join16(0xff, cpu.Regs.C)))
	return implNop
}

// LD (C+$FF00), A

func implLDHCA(cpu *CPU) CycleHandler {
	cpu.store(Addr(join16(0xff, cpu.Regs.C)), cpu.Regs.A)
	return implNop
}

// LD (n+$FF00), A

func implLDHnA(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.store(Addr(join16(0xff, cpu.Regs.TempZ)), cpu.Regs.A)
		return implNop
	}
}

// LD A,(n+$FF00)

func implLDHAn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.A = cpu.load(Addr(join16(0xff, cpu.Regs.TempZ)))
		return implNop
	}
}

// LD r, r'

func implLDAA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.A) }
func implLDBA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.A) }
func implLDCA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.A) }
func implLDDA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.A) }
func implLDEA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.A) }
func implLDHA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.A) }
func implLDLA(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.A) }

func implLDAB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.B) }
func implLDBB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.B) }
func implLDCB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.B) }
func implLDDB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.B) }
func implLDEB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.B) }
func implLDHB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.B) }
func implLDLB(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.B) }

func implLDAC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.C) }
func implLDBC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.C) }
func implLDCC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.C) }
func implLDDC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.C) }
func implLDEC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.C) }
func implLDHC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.C) }
func implLDLC(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.C) }

func implLDAD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.D) }
func implLDBD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.D) }
func implLDCD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.D) }
func implLDDD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.D) }
func implLDED(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.D) }
func implLDHD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.D) }
func implLDLD(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.D) }

func implLDAE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.E) }
func implLDBE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.E) }
func implLDCE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.E) }
func implLDDE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.E) }
func implLDEE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.E) }
func implLDHE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.E) }
func implLDLE(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.E) }

func implLDAH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.H) }
func implLDBH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.H) }
func implLDCH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.H) }
func implLDDH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.H) }
func implLDEH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.H) }
func implLDHH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.H) }
func implLDLH(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.H) }

func implLDAL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.A, cpu.Regs.L) }
func implLDBL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.B, cpu.Regs.L) }
func implLDCL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.C, cpu.Regs.L) }
func implLDDL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.D, cpu.Regs.L) }
func implLDEL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.E, cpu.Regs.L) }
func implLDHL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.H, cpu.Regs.L) }
func implLDLL(cpu *CPU) CycleHandler { return implLDrr(&cpu.Regs.L, cpu.Regs.L) }

func implLDrr(dest *Data8, src Data8) CycleHandler {
	*dest = src
	return nil
}

// LD r, n

func implLDAn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.A) }
func implLDBn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.B) }
func implLDCn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.C) }
func implLDDn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.D) }
func implLDEn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.E) }
func implLDHn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.H) }
func implLDLn(cpu *CPU) CycleHandler { return implLDrn(cpu, &cpu.Regs.L) }

func implLDrn(cpu *CPU, reg *Data8) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()

	return func(cpu *CPU) CycleHandler {
		*reg = cpu.Regs.TempZ
		return nil
	}
}

// LD (rr'), a

func implLDBCA(cpu *CPU) CycleHandler { return implLDrrA(cpu, cpu.Regs.B, cpu.Regs.C) }
func implLDDEA(cpu *CPU) CycleHandler { return implLDrrA(cpu, cpu.Regs.D, cpu.Regs.E) }

func implLDrrA(cpu *CPU, msr, lsr Data8) CycleHandler {
	cpu.store(Addr(join16(msr, lsr)), cpu.Regs.A)
	return implNop
}

// LD (rr'), nn

func implLDBCnn(cpu *CPU) CycleHandler { return implLDxxnn(cpu, &cpu.Regs.B, &cpu.Regs.C) }
func implLDDEnn(cpu *CPU) CycleHandler { return implLDxxnn(cpu, &cpu.Regs.D, &cpu.Regs.E) }

func implLDxxnn(cpu *CPU, msr, lsr *Data8) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()

	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()

		return func(cpu *CPU) CycleHandler {
			*msr = cpu.Regs.TempW
			*lsr = cpu.Regs.TempZ
			return nil
		}
	}
}

// LD A, (rr')

func implLDABC(cpu *CPU) CycleHandler { return implLDArr(cpu, cpu.Regs.B, cpu.Regs.C) }
func implLDADE(cpu *CPU) CycleHandler { return implLDArr(cpu, cpu.Regs.D, cpu.Regs.E) }

func implLDArr(cpu *CPU, msr, lsr Data8) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(join16(msr, lsr)))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.A = cpu.Regs.TempZ
		return nil
	}
}

// LD A, (HL+)

func implLDAHLInc(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(Addr(cpu.GetHL()))
	cpu.Regs.L++
	if cpu.Regs.L == 0 {
		cpu.Regs.H++
	}
	return implNop
}

// LD A, (HL-)

func implLDAHLDec(cpu *CPU) CycleHandler {
	cpu.Regs.A = cpu.load(Addr(cpu.GetHL()))
	cpu.Regs.L--
	if cpu.Regs.L == 0xFF {
		cpu.Regs.H--
	}
	return implNop
}

// JR PC+e

func implJRe(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		newPC := Data16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.TempZ)))
		cpu.Regs.SetWZ(newPC)
		return func(cpu *CPU) CycleHandler {
			cpu.SetPC(Addr(cpu.Regs.GetWZ()))
			return nil
		}
	}
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
	cpu.Regs.TempZ = cpu.fetchByteAtPC()

	if cond {
		cpu.lastBranchResult = +1
		return func(cpu *CPU) CycleHandler {
			newPC := Data16(int16(cpu.Regs.PC) + int16(int8(cpu.Regs.TempZ)))
			cpu.Regs.SetWZ(newPC)
			return func(cpu *CPU) CycleHandler {
				cpu.SetPC(Addr(cpu.Regs.GetWZ()))
				return nil
			}
		}
	}
	cpu.lastBranchResult = -1
	return implNop
}

// JP HL

func implJPHL(cpu *CPU) CycleHandler {
	cpu.SetPC(Addr(cpu.GetHL() - 1))
	return nil
}

// JP nn

func implJPnn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		return func(cpu *CPU) CycleHandler {
			cpu.lastBranchResult = +1
			cpu.SetPC(Addr(cpu.Regs.GetWZ() - 1))
			return implNop
		}
	}
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
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		if !cond {
			cpu.lastBranchResult = -1
			return implNop
		}
		return func(cpu *CPU) CycleHandler {
			cpu.lastBranchResult = +1
			cpu.SetPC(Addr(cpu.Regs.GetWZ() - 1))
			return implNop
		}
	}
}

// CALL nn

func implCALLnn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		return func(cpu *CPU) CycleHandler {
			return func(cpu *CPU) CycleHandler {
				cpu.push((cpu.Regs.PC + 1).MSB())
				return func(cpu *CPU) CycleHandler {
					cpu.push((cpu.Regs.PC + 1).LSB())
					return func(cpu *CPU) CycleHandler {
						cpu.SetPC(Addr(cpu.Regs.GetWZ()))
						return nil
					}
				}
			}
		}
	}
}

// CALL cc nn

func implCALLZnn(cpu *CPU) CycleHandler  { return implCALLccnn(cpu, cpu.Regs.GetFlagZ()) }
func implCALLNZnn(cpu *CPU) CycleHandler { return implCALLccnn(cpu, !cpu.Regs.GetFlagZ()) }
func implCALLCnn(cpu *CPU) CycleHandler  { return implCALLccnn(cpu, cpu.Regs.GetFlagC()) }
func implCALLNCnn(cpu *CPU) CycleHandler { return implCALLccnn(cpu, !cpu.Regs.GetFlagC()) }

func implCALLccnn(cpu *CPU, cond bool) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.fetchByteAtPC()
		return func(cpu *CPU) CycleHandler {
			if !cond {
				return nil
			}
			return func(cpu *CPU) CycleHandler {
				cpu.push((cpu.Regs.PC + 1).MSB())
				return func(cpu *CPU) CycleHandler {
					cpu.push((cpu.Regs.PC + 1).LSB())
					return func(cpu *CPU) CycleHandler {
						cpu.SetPC(Addr(cpu.Regs.GetWZ() - 1))
						return nil
					}
				}
			}
		}
	}
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
	return func(cpu *CPU) CycleHandler {
		cpu.push((cpu.Regs.PC + 1).LSB())
		return func(cpu *CPU) CycleHandler {
			cpu.SetPC(Addr(join16(0x00, vec) - 1))
			return implNop
		}
	}
}

// RET

func implRET(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.pop()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.pop()
		return func(cpu *CPU) CycleHandler {
			cpu.SetPC(Addr(cpu.Regs.GetWZ() - 1))
			return implNop
		}
	}
}

// RETI

func implRETI(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.pop()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.pop()
		return func(cpu *CPU) CycleHandler {
			cpu.SetPC(Addr(cpu.Regs.GetWZ() - 1))
			// TODO verify if this is the right cycle
			if cpu.Interrupts != nil {
				cpu.Interrupts.IME = true
			}
			return implNop
		}
	}
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
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempZ = cpu.pop()
		return func(cpu *CPU) CycleHandler {
			cpu.Regs.TempW = cpu.pop()
			return func(cpu *CPU) CycleHandler {
				cpu.SetPC(Addr(cpu.Regs.GetWZ()) - 1)
				return implNop
			}
		}
	}
}

// PUSH rr'

func implPUSHAF(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.A, cpu.Regs.F) }
func implPUSHBC(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.B, cpu.Regs.C) }
func implPUSHDE(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.D, cpu.Regs.E) }
func implPUSHHL(cpu *CPU) CycleHandler { return implPUSHrr(cpu, cpu.Regs.H, cpu.Regs.L) }

func implPUSHrr(cpu *CPU, msb, lsb Data8) CycleHandler {
	cpu.push(msb)
	return func(cpu *CPU) CycleHandler {
		cpu.push(lsb)
		return func(cpu *CPU) CycleHandler {
			return implNop
		}
	}
}

// POP rr'

func implPOPAF(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.A, &cpu.Regs.F) }
func implPOPBC(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.B, &cpu.Regs.C) }
func implPOPDE(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.D, &cpu.Regs.E) }
func implPOPHL(cpu *CPU) CycleHandler { return implPOPrr(cpu, &cpu.Regs.H, &cpu.Regs.L) }

func implPOPrr(cpu *CPU, msb, lsb *Data8) CycleHandler {
	cpu.Regs.TempZ = cpu.pop()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.TempW = cpu.pop()
		return func(cpu *CPU) CycleHandler {
			*msb = cpu.Regs.TempW
			if lsb == &cpu.Regs.F {
				*lsb = cpu.Regs.TempZ & 0xf0
			} else {
				*lsb = cpu.Regs.TempZ
			}
			return nil
		}
	}
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
func implANDB(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.B) }
func implANDC(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.C) }
func implANDD(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.D) }
func implANDE(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.E) }
func implANDH(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.H) }
func implANDL(cpu *CPU) CycleHandler { return implANDr(cpu, cpu.Regs.L) }

func implANDr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, r))
	return nil
}

// AND A, n

func implANDn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, cpu.Regs.TempZ))
		return nil
	}
}

// AND A, (HL)

func implANDHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(AND(cpu.Regs.A, cpu.Regs.TempZ))
		return nil
	}
}

// CP A, r

func implCPA(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.A) }
func implCPB(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.B) }
func implCPC(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.C) }
func implCPD(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.D) }
func implCPE(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.E) }
func implCPH(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.H) }
func implCPL(cpu *CPU) CycleHandler { return implCPr(cpu, cpu.Regs.L) }

func implCPr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlags(SUB(cpu.Regs.A, r, false))
	return nil
}

// CP A, n

func implCPn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlags(SUB(cpu.Regs.A, cpu.Regs.TempZ, false))
		return nil
	}
}

// CP A, (HL)

func implCPHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlags(SUB(cpu.Regs.A, cpu.Regs.TempZ, false))
		return nil
	}
}

// SUB A, r

func implSUBA(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.A) }
func implSUBB(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.B) }
func implSUBC(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.C) }
func implSUBD(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.D) }
func implSUBE(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.E) }
func implSUBH(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.H) }
func implSUBL(cpu *CPU) CycleHandler { return implSUBr(cpu, cpu.Regs.L) }

func implSUBr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, r, false))
	return nil
}

// SUB A, n

func implSUBn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.TempZ, false))
		return nil
	}
}

// SUB A, (HL)

func implSUBHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.TempZ, false))
		return nil
	}
}

// SBC A, r

func implSBCA(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.A) }
func implSBCB(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.B) }
func implSBCC(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.C) }
func implSBCD(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.D) }
func implSBCE(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.E) }
func implSBCH(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.H) }
func implSBCL(cpu *CPU) CycleHandler { return implSBCr(cpu, cpu.Regs.L) }

func implSBCr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, r, cpu.Regs.GetFlagC()))
	return nil
}

// SBC A, n

func implSBCn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.TempZ, cpu.Regs.GetFlagC()))
		return nil
	}
}

// SBC A, (HL)

func implSBCHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(SUB(cpu.Regs.A, cpu.Regs.TempZ, cpu.Regs.GetFlagC()))
		return nil
	}
}

// ADD A, r

func implADDA(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.A) }
func implADDB(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.B) }
func implADDC(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.C) }
func implADDD(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.D) }
func implADDE(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.E) }
func implADDH(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.H) }
func implADDL(cpu *CPU) CycleHandler { return implADDr(cpu, cpu.Regs.L) }

func implADDr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, r, false))
	return nil
}

// ADD A, n

func implADDn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.TempZ, false))
		return nil
	}
}

// ADD A, (HL)

func implADDHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.TempZ, false))
		return nil
	}
}

// ADD HL, rr'

func implADDHLHL(cpu *CPU) CycleHandler { return implADDHLrr(cpu, cpu.Regs.H, cpu.Regs.L) }
func implADDHLBC(cpu *CPU) CycleHandler { return implADDHLrr(cpu, cpu.Regs.B, cpu.Regs.C) }
func implADDHLDE(cpu *CPU) CycleHandler { return implADDHLrr(cpu, cpu.Regs.D, cpu.Regs.E) }
func implADDHLSP(cpu *CPU) CycleHandler {
	return implADDHLrr(cpu, cpu.Regs.SP.MSB(), cpu.Regs.SP.LSB())
}

func implADDHLrr(cpu *CPU, msr, lsr Data8) CycleHandler {
	result := ADD(cpu.Regs.L, lsr, false)
	cpu.Regs.L = result.Value
	cpu.Regs.SetFlagC(result.C)
	cpu.Regs.SetFlagH(result.H)
	cpu.Regs.SetFlagN(result.N)
	return func(cpu *CPU) CycleHandler {
		result := ADD(cpu.Regs.H, msr, cpu.Regs.GetFlagC())
		cpu.Regs.H = result.Value
		cpu.Regs.SetFlagC(result.C)
		cpu.Regs.SetFlagH(result.H)
		cpu.Regs.SetFlagN(result.N)
		return nil
	}
}

// ADD SP, PC+e

func implADDSPe(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
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
		return func(cpu *CPU) CycleHandler {
			res := cpu.Regs.SP.MSB()
			if cpu.Regs.TempW == 1 {
				res++
			} else if cpu.Regs.TempW == 0xff {
				res--
			}
			cpu.Regs.TempW = res
			return func(cpu *CPU) CycleHandler {
				cpu.SetSP(Addr(cpu.Regs.GetWZ()))
				return nil
			}
		}
	}
}

// ADC A, r

func implADCA(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.A) }
func implADCB(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.B) }
func implADCC(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.C) }
func implADCD(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.D) }
func implADCE(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.E) }
func implADCH(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.H) }
func implADCL(cpu *CPU) CycleHandler { return implADCr(cpu, cpu.Regs.L) }

func implADCr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, r, cpu.Regs.GetFlagC()))
	return nil
}

// ADC A, n

func implADCn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.TempZ, cpu.Regs.GetFlagC()))
		return nil
	}
}

// ADC A, (HL)

func implADCHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(ADD(cpu.Regs.A, cpu.Regs.TempZ, cpu.Regs.GetFlagC()))
		return nil
	}
}

// OR A, r

func implORA(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.A) }
func implORB(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.B) }
func implORC(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.C) }
func implORD(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.D) }
func implORE(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.E) }
func implORH(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.H) }
func implORL(cpu *CPU) CycleHandler { return implORr(cpu, cpu.Regs.L) }

func implORr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, r))
	return nil
}

// OR A, n

func implORn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, cpu.Regs.TempZ))
		return nil
	}
}

// OR A, (HL)

func implORHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(OR(cpu.Regs.A, cpu.Regs.TempZ))
		return nil
	}
}

// XOR A, r

func implXORA(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.A) }
func implXORB(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.B) }
func implXORC(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.C) }
func implXORD(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.D) }
func implXORE(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.E) }
func implXORH(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.H) }
func implXORL(cpu *CPU) CycleHandler { return implXORr(cpu, cpu.Regs.L) }

func implXORr(cpu *CPU, r Data8) CycleHandler {
	cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, r))
	return nil
}

// XOR A, n

func implXORn(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, cpu.Regs.TempZ))
		return nil
	}
}

// XOR A, (HL)

func implXORHL(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		cpu.Regs.SetFlagsAndA(XOR(cpu.Regs.A, cpu.Regs.TempZ))
		return nil
	}
}

// INC r

func implINCA(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.A) }
func implINCB(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.B) }
func implINCC(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.C) }
func implINCD(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.D) }
func implINCE(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.E) }
func implINCH(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.H) }
func implINCL(cpu *CPU) CycleHandler { return implINCr(cpu, &cpu.Regs.L) }

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
	cpu.Regs.C++
	if cpu.Regs.C == 0 {
		cpu.Regs.B++
	}
	return implNop
}

func implINCDE(cpu *CPU) CycleHandler {
	cpu.Regs.E++
	if cpu.Regs.E == 0 {
		cpu.Regs.D++
	}
	return implNop
}

func implINCHL(cpu *CPU) CycleHandler {
	cpu.Regs.L++
	if cpu.Regs.L == 0 {
		cpu.Regs.H++
	}
	return implNop
}

func implINCSP(cpu *CPU) CycleHandler {
	cpu.Regs.SP++
	return implNop
}

// INC (HL)

func implINCHLInd(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		res := ADD(cpu.Regs.TempZ, 1, false)
		// apparently doesn't set C.
		cpu.Regs.SetFlagH(res.H)
		cpu.Regs.SetFlagZ(res.Z())
		cpu.Regs.SetFlagN(res.N)
		cpu.writeAddressBus(Addr(cpu.GetHL()))
		cpu.Regs.TempZ = res.Value
		cpu.Bus.WriteData(cpu.Regs.TempZ)
		return implNop
	}
}

// DEC r

func implDECA(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.A) }
func implDECB(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.B) }
func implDECC(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.C) }
func implDECD(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.D) }
func implDECE(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.E) }
func implDECH(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.H) }
func implDECL(cpu *CPU) CycleHandler { return implDECr(cpu, &cpu.Regs.L) }

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
	cpu.Regs.C--
	if cpu.Regs.C == 0xFF {
		cpu.Regs.B--
	}
	return implNop
}

func implDECDE(cpu *CPU) CycleHandler {
	cpu.Regs.E--
	if cpu.Regs.E == 0xFF {
		cpu.Regs.D--
	}
	return implNop
}

func implDECHL(cpu *CPU) CycleHandler {
	cpu.Regs.L--
	if cpu.Regs.L == 0xFF {
		cpu.Regs.H--
	}
	return implNop
}

func implDECSP(cpu *CPU) CycleHandler {
	cpu.Regs.SP--
	return implNop
}

// DEC (HL)

func implDECHLInd(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.load(Addr(cpu.GetHL()))
	return func(cpu *CPU) CycleHandler {
		res := SUB(cpu.Regs.TempZ, 1, false)
		// apparently doesn't set C.
		cpu.Regs.SetFlagH(res.H)
		cpu.Regs.SetFlagZ(res.Z())
		cpu.Regs.SetFlagN(res.N)
		cpu.writeAddressBus(Addr(cpu.GetHL()))
		cpu.Regs.TempZ = res.Value
		cpu.Bus.WriteData(cpu.Regs.TempZ)
		return implNop
	}
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

func implUndefined(cpu *CPU) CycleHandler {
	panicf("hit undefined instruction %s (0x%02x)", cpu.Regs.IR, uint8(cpu.Regs.IR))

	// Real GB does something like this; system stops
	return implUndefined
}

// CB

func implCB(cpu *CPU) CycleHandler {
	cpu.Regs.TempZ = cpu.fetchByteAtPC()

	return implCB_2
}

func implCB_2(cpu *CPU) CycleHandler {
	cbOp := NewCBOp(cpu.Regs.TempZ)

	if cbOp.Target == CBTargetIndirectHL {
		cpu.writeAddressBus(Addr(cpu.GetHL()))
	}
	var val Data8
	switch cbOp.Target {
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
	val = cpu.doCBOp(val, cbOp.Op)
	switch cbOp.Target {
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
		if cbOp.Op.Is3Cycles() {
			return implNop
		}
		return func(cpu *CPU) CycleHandler {
			cpu.writeAddressBus(Addr(cpu.GetHL()))
			cpu.Bus.WriteData(cpu.Regs.TempZ)
			return implNop
		}
	case CBTargetA:
		cpu.Regs.A = val
	default:
		panic("unknown CBOp target")
	}
	return nil
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
	cpu.writeAddressBus(addr)
	return cpu.Bus.GetData()
}

func (cpu *CPU) fetchByteAtPC() Data8 {
	data := cpu.load(cpu.Regs.PC)
	cpu.IncPC()
	return data
}

func (cpu *CPU) store(addr Addr, v Data8) {
	cpu.writeAddressBus(addr)
	cpu.Bus.WriteData(v)
}
