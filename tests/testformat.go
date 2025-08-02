package tests

import (
	"bytes"
	"encoding/gob"
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
	Addr uint16
	Val  uint8
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
				cy.Addr = addr
			}
		}
		if len(step) > 1 && step[1] != nil {
			if v, ok := step[1].(float64); ok {
				val := uint8(v)
				cy.Val = val
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

func (cs Cycles) String() string {
	out := ""
	for _, c := range cs {
		out += fmt.Sprintf("@%04x %02x %s\n", c.Addr, c.Val, c.RW)
	}
	return out
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

func TestCaseStateFromCPU(gb *model.Gameboy) *CPUState {
	return &CPUState{
		PC: gb.CPU.Regs.PC,
		SP: gb.CPU.Regs.SP,
		A:  gb.CPU.Regs.A,
		B:  gb.CPU.Regs.B,
		C:  gb.CPU.Regs.C,
		D:  gb.CPU.Regs.D,
		E:  gb.CPU.Regs.E,
		F:  gb.CPU.Regs.F,
		H:  gb.CPU.Regs.H,
		L:  gb.CPU.Regs.L,
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

func (tc *TestCase) String(gb *model.Gameboy) string {
	return fmt.Sprintf(
		"%s\n---------\nInitial:\n%s---\n%s---\nFinal expect:\n%s-----\nFinal actual:\n%s",
		tc.Name,
		tc.Initial.String(),
		tc.Cycles.String(),
		tc.Final.String(),
		TestCaseStateFromCPU(gb).String(),
	)
}

func MustReadFile(name string) []TestCase {
	var tcs []TestCase

	gobfname := fmt.Sprintf("sm83/v1/%s.gob", name)
	gobf, err := os.Open(gobfname)
	if err == nil {
		defer gobf.Close()
		dec := gob.NewDecoder(gobf)
		if err := dec.Decode(&tcs); err == nil {
			return tcs
		}
	}
	content, err := os.ReadFile(fmt.Sprintf("sm83/v1/%s.json", name))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &tcs)
	if err != nil {
		panic(err)
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&tcs); err == nil {
		os.WriteFile(gobfname, buf.Bytes(), 0o666)
	}
	return tcs
}

func Run(t *testing.T, tcs []TestCase, opcode model.Opcode, gb *model.Gameboy, audio model.Audio) {
	t.Helper()
	for i, tc := range tcs {
		RunOne(t, i, tc, opcode, gb, audio)
		if t.Failed() {
			break
		}
	}
}

func RunOne(t *testing.T, i int, tc TestCase, opcode model.Opcode, gb *model.Gameboy, audio model.Audio) {
	t.Helper()

	clock := model.NewClock()

	fs := model.FrameSync{Ch: make(chan func(*model.ViewPort), 1)}

	config := model.DefaultConfig
	config.Debug.Disassembler.Enable = false
	config.BootROM.Variant = "None"
	config.Debug.RewindSize = 16
	gb.Init(&config, clock)

	gb.CPU.Regs.A = tc.Initial.A
	gb.CPU.Regs.B = tc.Initial.B
	gb.CPU.Regs.C = tc.Initial.C
	gb.CPU.Regs.D = tc.Initial.D
	gb.CPU.Regs.E = tc.Initial.E
	gb.CPU.Regs.F = tc.Initial.F
	gb.CPU.Regs.H = tc.Initial.H
	gb.CPU.Regs.L = tc.Initial.L
	gb.CPU.Regs.PC = tc.Initial.PC
	gb.CPU.Regs.SP = tc.Initial.SP
	for _, entry := range tc.Initial.RAM {
		gb.Mem[entry.Addr] = model.Data8(entry.Val)
	}
	gb.PureRAM = true

	// Emulate a fetch
	gb.CPU.Regs.IR = model.Opcode(gb.Mem[tc.Initial.PC])
	gb.CPU.Regs.PC++

	clock.MCycle(len(tc.Cycles), gb, audio, &fs)

	defer func() {
		if t.Failed() {
			gb.CPU.Dump(gb)
		}
	}()

	if have, want := gb.CPU.Regs.A, tc.Final.A; have != want {
		t.Fatalf("Test %d Register A have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.B, tc.Final.B; have != want {
		t.Fatalf("Test %d Register B have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.C, tc.Final.C; have != want {
		t.Fatalf("Test %d Register C have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.D, tc.Final.D; have != want {
		t.Fatalf("Test %d Register D have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.E, tc.Final.E; have != want {
		t.Fatalf("Test %d Register E have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.F.Bit(model.FlagBitZ), tc.Final.F.Bit(model.FlagBitZ); have != want {
		t.Fatalf("Test %d Flag Z have %v want %v. Full test: %s", i, have, want, tc.String(gb))
	}
	if have, want := gb.CPU.Regs.F.Bit(model.FlagBitC), tc.Final.F.Bit(model.FlagBitC); have != want {
		t.Fatalf("Test %d Flag C have %v want %v. Full test: %s", i, have, want, tc.String(gb))
	}
	if have, want := gb.CPU.Regs.F.Bit(model.FlagBitH), tc.Final.F.Bit(model.FlagBitH); have != want {
		t.Fatalf("Test %d Flag H have %v want %v. Full test: %s", i, have, want, tc.String(gb))
	}
	if have, want := gb.CPU.Regs.F.Bit(model.FlagBitN), tc.Final.F.Bit(model.FlagBitN); have != want {
		t.Fatalf("Test %d Flag N have %v want %v. Full test: %s", i, have, want, tc.String(gb))
	}
	if have, want := gb.CPU.Regs.F, tc.Final.F; have != want {
		t.Fatalf("Test %d Register F have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.H, tc.Final.H; have != want {
		t.Fatalf("Test %d Register H have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.L, tc.Final.L; have != want {
		t.Fatalf("Test %d Register L have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.PC, tc.Final.PC+1; have != want {
		t.Fatalf("Test %d Register PC have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}
	if have, want := gb.CPU.Regs.SP, tc.Final.SP; have != want {
		t.Fatalf("Test %d Register SP have %s want %s. Full test: %s", i, have.Hex(), want.Hex(), tc.String(gb))
	}

	for _, entry := range tc.Final.RAM {
		if have, want := gb.Mem[entry.Addr], model.Data8(entry.Val); have != want {
			t.Fatalf("Test %d Addr %s have %s want %s. Full test: %s", i, entry.Addr.Hex(), have.Hex(), want.Hex(), tc.String(gb))
		}
	}
}
