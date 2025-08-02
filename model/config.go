package model

type Config struct {
	Clock   ConfigClock
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

type ConfigDebug struct {
	RewindSize            int
	PanicOnStackUnderflow bool
	Disassembler          ConfigDisassembler
}

type ConfigDisassembler struct {
	Trace  bool
	Enable bool
}

var DefaultConfig = Config{
	Clock: ConfigClock{
		SpeedPercent: 100.0,
	},
	BootROM: ConfigBootROM{
		Skip:    false,
		Variant: "DMGBoot",
	},
	Debug: ConfigDebug{
		PanicOnStackUnderflow: true,
		Disassembler: ConfigDisassembler{
			Trace:  false,
			Enable: true,
		},
	},
}
