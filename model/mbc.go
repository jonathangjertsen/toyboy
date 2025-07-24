package model

//go:generate go-enum --flag --nocomments

// ENUM(None, 1, 2, 3, MMM01, 5, 6, 7, PocketCamera, BandaiTAMA5, HuC3, HuC1)
type MBCID uint8

const (
	ROMBankSize = 16 * 1024
	RAMBankSize = 8 * 1024
)

type MBCFeatures struct {
	ID        MBCID
	RAM       bool
	Battery   bool
	RTC       bool
	Rumble    bool
	NROMBanks int
	NRAMBanks int
}

func (mbcf *MBCFeatures) TotalROMSize() int {
	return ROMBankSize * mbcf.NROMBanks
}

func (mbcf *MBCFeatures) TotalRAMSize() int {
	return RAMBankSize * mbcf.NRAMBanks
}

func GetMBCFeatures(code, romsiz, ramsiz uint8) MBCFeatures {
	var mbc MBCFeatures
	switch code {
	case 0x00:
		mbc.ID = MBCIDNone
	case 0x01:
		mbc.ID = MBCID1
	case 0x02:
		mbc.ID = MBCID1
		mbc.RAM = true
	case 0x03:
		mbc.ID = MBCID1
		mbc.RAM = true
		mbc.Battery = true
	case 0x05:
		mbc.ID = MBCID2
	case 0x06:
		mbc.ID = MBCID2
		mbc.RAM = true
		mbc.Battery = true
	case 0x08:
		mbc.ID = MBCIDNone
		mbc.RAM = true
	case 0x09:
		mbc.ID = MBCIDNone
		mbc.RAM = true
		mbc.Battery = true
	case 0x0b:
		mbc.ID = MBCIDMMM01
	case 0x0c:
		mbc.ID = MBCIDMMM01
		mbc.RAM = true
	case 0x0d:
		mbc.ID = MBCIDMMM01
		mbc.RAM = true
		mbc.Battery = true
	case 0x0f:
		mbc.ID = MBCID3
		mbc.RTC = true
		mbc.Battery = true
	case 0x10:
		mbc.ID = MBCID3
		mbc.RTC = true
		mbc.RAM = true
		mbc.Battery = true
	case 0x11:
		mbc.ID = MBCID3
	case 0x12:
		mbc.ID = MBCID3
		mbc.RAM = true
	case 0x13:
		mbc.ID = MBCID3
		mbc.RAM = true
		mbc.Battery = true
	case 0x19:
		mbc.ID = MBCID5
	case 0x1a:
		mbc.ID = MBCID5
		mbc.RAM = true
	case 0x1b:
		mbc.ID = MBCID5
		mbc.RAM = true
		mbc.Battery = true
	case 0x1c:
		mbc.ID = MBCID5
		mbc.Rumble = true
	case 0x1d:
		mbc.ID = MBCID5
		mbc.Rumble = true
		mbc.RAM = true
	case 0x1e:
		mbc.ID = MBCID5
		mbc.Rumble = true
		mbc.RAM = true
		mbc.Battery = true
	case 0x20:
		mbc.ID = MBCID6
	case 0x22:
		mbc.ID = MBCID7
		mbc.Rumble = true
		mbc.RAM = true
		mbc.Battery = true
	case 0xfc:
		mbc.ID = MBCIDPocketCamera
	case 0xfd:
		mbc.ID = MBCIDBandaiTAMA5
	case 0xfe:
		mbc.ID = MBCIDHuC3
		mbc.RTC = true
	case 0xff:
		mbc.ID = MBCIDHuC1
		mbc.RAM = true
		mbc.Battery = true
	}

	if romsiz <= 0x08 {
		mbc.NROMBanks = 1 << romsiz
	} else {
		panicf("rom size byte 0x%02x not supported", romsiz)
	}

	if mbc.RAM {
		switch ramsiz {
		case 2:
			mbc.NRAMBanks = 1
		case 3:
			mbc.NRAMBanks = 4
		case 4:
			mbc.NRAMBanks = 16
		case 5:
			mbc.NRAMBanks = 8
		}
	}
	return mbc
}
