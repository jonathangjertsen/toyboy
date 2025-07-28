with open('assets/cartridges/tetris.gb', 'rb') as f:
    data = bytearray(f.read())
print(f"original byte at 0x38: {data[0x038B]}")
data[0x038B] = 0x01

with open('assets/cartridges/tetris.gb', 'wb') as f:
    f.write(data)
