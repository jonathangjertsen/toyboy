with open('assets/bootrom/dmg_boot.bin', 'rb') as f:
    data = bytearray(f.read())

# skip vram zeroing
data[0x04:0x0c] = 0x00

# scroll faster
data[0x6d] = 0x00
data[0x6e] = 0x00
data[0x6f] = 0x00

with open('assets/bootrom/dmg_boot_patched.bin', 'wb') as f:
    f.write(data)
