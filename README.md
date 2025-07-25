# toyboy gameboy emulator

toyboy is a gameboy emulator.

As indicated by the name, it is a toy project. 
Contrary to the name, it is written by a fully grown man.

## Architecture / goals

toyboy is written in go.

It tries to be cycle-accurate, everything is triggered from the clock.

## Status

- Runs the bootrom
- Runs the Hello world ROM from GB ASM tutorial, working on the Unbricked ROM
- Handles most instructions
- PPU mostly works

Currently looks like this

![](screenshot.png)

## Resources used

https://gbdev.io/resources.html

https://gekkio.fi/files/gb-docs/gbctr.pdf

https://github.com/AntonioND/giibiiadvance/blob/master/docs/TCAGBD.pdf

https://hacktix.github.io/GBEDG/ppu/

https://gbdev.io/gb-asm-tutorial/ (for ROMs)
