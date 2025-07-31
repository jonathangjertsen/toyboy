package assets

import _ "embed"

//go:embed bootrom/dmg_boot.bin
var DMGBoot []byte

//go:embed bootrom/dmg_boot_patched.bin
var DMGBootPatched []byte
