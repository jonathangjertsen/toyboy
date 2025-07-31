package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jonathangjertsen/toyboy/model"
)

type RAMEntry struct {
	Addr model.Addr
	Val  model.Data8
}

type RAM []RAMEntry

func (r *RAM) UnmarshalJSON(data []byte) error {
	var raw [][2]uint64
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	entries := make([]RAMEntry, len(raw))
	for i, pair := range raw {
		entries[i] = RAMEntry{Addr: model.Addr(pair[0]), Val: model.Data8(pair[1])}
	}
	*r = entries
	return nil
}

type Cycle struct {
	Addr *uint16
	Val  *uint8
	RW   string
}

type Cycles []Cycle

func (c *Cycles) UnmarshalJSON(data []byte) error {
	var raw [][]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var result []Cycle
	for _, step := range raw {
		var cy Cycle
		if len(step) > 0 && step[0] != nil {
			if a, ok := step[0].(float64); ok {
				addr := uint16(a)
				cy.Addr = &addr
			}
		}
		if len(step) > 1 && step[1] != nil {
			if v, ok := step[1].(float64); ok {
				val := uint8(v)
				cy.Val = &val
			}
		}
		if len(step) > 2 && step[2] != nil {
			if s, ok := step[2].(string); ok {
				cy.RW = s
			}
		}
		result = append(result, cy)
	}
	*c = result
	return nil
}

type CPUState struct {
	PC  model.Addr  `json:"pc"`
	SP  model.Addr  `json:"sp"`
	A   model.Data8 `json:"a"`
	B   model.Data8 `json:"b"`
	C   model.Data8 `json:"c"`
	D   model.Data8 `json:"d"`
	E   model.Data8 `json:"e"`
	F   model.Data8 `json:"f"`
	H   model.Data8 `json:"h"`
	L   model.Data8 `json:"l"`
	IME uint8       `json:"ime"`
	EI  uint8       `json:"ei"`
	RAM RAM         `json:"ram"`
}

func TestCaseStateFromCPU(cpu *model.CPU) *CPUState {
	return &CPUState{
		PC: cpu.Regs.PC,
		SP: cpu.Regs.SP,
		A:  cpu.Regs.A,
		B:  cpu.Regs.B,
		C:  cpu.Regs.C,
		D:  cpu.Regs.D,
		E:  cpu.Regs.E,
		F:  cpu.Regs.F,
		H:  cpu.Regs.H,
		L:  cpu.Regs.L,
	}
}

func (cpus *CPUState) String() string {
	return fmt.Sprintf(
		"PC=%s SP=%s A=%s B=%s C=%s D=%s E=%s L=%s H=%s Z=%v C=%v H=%v N=%v\n",
		cpus.PC.Hex(),
		cpus.SP.Hex(),
		cpus.A.Hex(),
		cpus.B.Hex(),
		cpus.C.Hex(),
		cpus.D.Hex(),
		cpus.E.Hex(),
		cpus.L.Hex(),
		cpus.H.Hex(),
		cpus.F.Bit(model.FlagBitZ),
		cpus.F.Bit(model.FlagBitC),
		cpus.F.Bit(model.FlagBitH),
		cpus.F.Bit(model.FlagBitN),
	)
}

type TestCase struct {
	Name    string   `json:"name"`
	Initial CPUState `json:"initial"`
	Final   CPUState `json:"final"`
	Cycles  Cycles   `json:"cycles,omitempty"`
}

func (tc *TestCase) String(cpu *model.CPU) string {
	return fmt.Sprintf(
		"%s\n---------\nInitial expect:\n%s---<%d cycles>---\nFinal expect:\n%s-----\nFinal actual:\n%s",
		tc.Name,
		tc.Initial.String(),
		len(tc.Cycles),
		tc.Final.String(),
		TestCaseStateFromCPU(cpu).String(),
	)
}

func MustReadFile(name string) []TestCase {
	content, err := os.ReadFile(fmt.Sprintf("sm83/v1/%s.json", name))
	if err != nil {
		panic(err)
	}
	var tcs []TestCase
	err = json.Unmarshal(content, &tcs)
	if err != nil {
		panic(err)
	}
	return tcs
}

type TestBus struct {
	RAM  [1 << 16]model.Data8
	addr model.Addr
	data model.Data8
}

func (tb *TestBus) BeginCoreDump() func() {
	return func() {}
}

func (tb *TestBus) PushState() (model.Addr, model.Data8) {
	return tb.addr, tb.data
}

func (tb *TestBus) PopState(addr model.Addr, data model.Data8) {
	tb.addr, tb.data = addr, data
}

func (tb *TestBus) Reset() {
	tb.addr = 0
	tb.data = 0
}

func (tb *TestBus) InCoreDump() bool {
	return false
}

func (tb *TestBus) WriteAddress(addr model.Addr) {
	tb.addr = addr
	tb.data = tb.RAM[tb.addr]
}

func (tb *TestBus) ProbeAddress(addr model.Addr) model.Data8 {
	return tb.RAM[addr]
}

func (tb *TestBus) ProbeRange(begin, end model.Addr) []model.Data8 {
	return tb.RAM[begin : end+1]
}

func (tb *TestBus) WriteData(data model.Data8) {
	tb.data = data
	tb.RAM[tb.addr] = data
}

func (tb *TestBus) GetAddress() model.Addr {
	return tb.addr
}

func (tb *TestBus) GetData() model.Data8 {
	return tb.data
}

func (tb *TestBus) GetCounters(model.Addr) (uint64, uint64) {
	return 0, 0
}

func (tb *TestBus) GetPeripheral(any) {
	panic("getPeripheral not implemented")
}

func Run(t *testing.T, tcs []TestCase, opcode model.Opcode) {
	t.Helper()
	for i, tc := range tcs {
		RunOne(t, i, tc, opcode)
		if t.Failed() {
			break
		}
	}
}

func RunOne(t *testing.T, i int, tc TestCase, opcode model.Opcode) {
	t.Helper()

	audio, devnull := model.AudioStub()
	defer func() { close(devnull) }()
	testBus := TestBus{}
	ints := model.NewInterrupts(testBus.RAM[:])

	clock := model.NewRealtimeClock(model.DefaultConfig.Clock, audio, ints)
	defer func() { clock.Stop() }()

	apu := model.NewAPU(clock, &model.DefaultConfig, testBus.RAM[:])
	model.NewTimer(clock, testBus.RAM[:], apu, ints)
	model.NewPPU(clock, ints, testBus.RAM[:], &model.DefaultConfig, nil)

	config := model.DefaultConfig
	config.Debug.PanicOnStackUnderflow = false
	cpu := model.NewCPU(clock, nil, &testBus, &config, nil)
	cpu.Reset()
	cpu.Regs.A = tc.Initial.A
	cpu.Regs.B = tc.Initial.B
	cpu.Regs.C = tc.Initial.C
	cpu.Regs.D = tc.Initial.D
	cpu.Regs.E = tc.Initial.E
	cpu.Regs.F = tc.Initial.F
	cpu.Regs.H = tc.Initial.H
	cpu.Regs.L = tc.Initial.L
	cpu.Regs.PC = tc.Initial.PC
	cpu.Regs.SP = tc.Initial.SP
	for _, entry := range tc.Initial.RAM {
		testBus.RAM[entry.Addr] = model.Data8(entry.Val)
	}
	cpu.Regs.IR = model.Opcode(testBus.RAM[cpu.Regs.PC])

	clock.MCycle(len(tc.Cycles)+1, ints)

	if have, want := cpu.Regs.A, tc.Final.A; have != want {
		t.Fatalf("Test %d Register A have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.B, tc.Final.B; have != want {
		t.Fatalf("Test %d Register B have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.C, tc.Final.C; have != want {
		t.Fatalf("Test %d Register C have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.D, tc.Final.D; have != want {
		t.Fatalf("Test %d Register D have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.E, tc.Final.E; have != want {
		t.Fatalf("Test %d Register E have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.F.Bit(model.FlagBitZ), tc.Final.F.Bit(model.FlagBitZ); have != want {
		t.Fatalf("Test %d Flag Z have %v want %v. Full test: %s", i, have, want, tc.String(cpu))
	}
	if have, want := cpu.Regs.F.Bit(model.FlagBitC), tc.Final.F.Bit(model.FlagBitC); have != want {
		t.Fatalf("Test %d Flag C have %v want %v. Full test: %s", i, have, want, tc.String(cpu))
	}
	if have, want := cpu.Regs.F.Bit(model.FlagBitH), tc.Final.F.Bit(model.FlagBitH); have != want {
		t.Fatalf("Test %d Flag H have %v want %v. Full test: %s", i, have, want, tc.String(cpu))
	}
	if have, want := cpu.Regs.F.Bit(model.FlagBitN), tc.Final.F.Bit(model.FlagBitN); have != want {
		t.Fatalf("Test %d Flag N have %v want %v. Full test: %s", i, have, want, tc.String(cpu))
	}
	if have, want := cpu.Regs.F, tc.Final.F; have != want {
		t.Fatalf("Test %d Register F have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.H, tc.Final.H; have != want {
		t.Fatalf("Test %d Register H have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.L, tc.Final.L; have != want {
		t.Fatalf("Test %d Register L have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.PC, tc.Final.PC+1; have != want {
		t.Fatalf("Test %d Register PC have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}
	if have, want := cpu.Regs.SP, tc.Final.SP; have != want {
		t.Fatalf("Test %d Register SP have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(cpu))
	}

	for _, entry := range tc.Final.RAM {
		if have, want := testBus.RAM[entry.Addr], model.Data8(entry.Val); have != want {
			t.Fatalf("Test %d Addr %s have %s want %s. Full test: %s", i, entry.Addr.Hex(), have.Hex(), want.Hex(), tc.String(cpu))
		}
	}
}
