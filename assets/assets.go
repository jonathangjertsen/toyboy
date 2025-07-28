package assets

import _ "embed"

//go:embed bootrom/dmg_boot_patched.bin
var DMGBoot []byte
