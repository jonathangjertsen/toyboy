package model

type Config struct {
	Clock   ConfigClock
	ROM     ConfigROM
	BootROM ConfigBootROM
	Debug   ConfigDebug
}

type ConfigBootROM struct {
	Variant string
	Skip    bool
}

type ConfigClock struct {
	SpeedPercent float64
}

type ConfigROM struct {
	Location string
}

type ConfigDebug struct {
	PanicOnStackUnderflow bool
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
		Skip:    false,
		Variant: "DMGBoot",
	},
	Debug: ConfigDebug{
		PanicOnStackUnderflow: true,
		Disassembler: ConfigDisassembler{
			Trace: true,
		},
	},
}
