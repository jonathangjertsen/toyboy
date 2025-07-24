package model

type Config struct {
	Clock   ConfigClock
	ROM     ConfigROM
	BootROM ConfigBootROM
	Debug   ConfigDebug
}

type ConfigBootROM struct {
	Location string
	Skip     bool
}

type ConfigClock struct {
	SpeedPercent float64
}

type ConfigROM struct {
	Location string
}

type ConfigDebug struct {
	PanicOnStackUnderflow bool
	MaxNOPCount           int
	Disassembler          ConfigDisassembler
}

type ConfigDisassembler struct {
	Trace bool
}

var DefaultConfig = Config{
	Clock: ConfigClock{
		SpeedPercent: 100.0,
	},
	ROM: ConfigROM{
		"assets/cartridges/01-special.gb",
	},
	BootROM: ConfigBootROM{
		Skip:     false,
		Location: "assets/bootrom/dmg_boot.bin",
	},
	Debug: ConfigDebug{
		PanicOnStackUnderflow: true,
		MaxNOPCount:           10,
		Disassembler: ConfigDisassembler{
			Trace: true,
		},
	},
}
