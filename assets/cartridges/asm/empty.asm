; from https://gbdev.io/gb-asm-tutorial/part1/hello_world.html

INCLUDE "hardware.inc"

SECTION "Header", ROM0[$100]
    jp EntryPoint

    NINTENDO_LOGO

    ds $150 - @, 0

EntryPoint:
    jp EntryPoint
