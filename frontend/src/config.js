const DefaultConfig = {
    "Model": {
        "Clock": {
            "SpeedPercent": 1010
        },
        "ROM": {
            "Location": "assets/cartridges/unbricked.gb"
        },
        "BootROM": {
            "Variant": "DMGBoot",
            "Skip": false
        },
        "Debug": {
            "PanicOnStackUnderflow": false,
            "MaxNOPCount": 10,
            "Disassembler": {
                "Trace": false
            }
        }
    },
    "PProfURL": "localhost:6060",
    "GUI": {
        "Graphics": {
        "Overlay": false,
        "GridColor": {
            "R": 136,
            "G": 136,
            "B": 136,
            "A": 255
        },
        "FillColor": {
            "R": 240,
            "G": 240,
            "B": 240,
            "A": 255
        },
        "DashLen": 4,
        "GridThickness": 1,
        "Font": "Basic"
        },
        "VRAMMem": {
        "Box": {
            "Show": false,
            "Height": 660,
            "Width": 650
        }
        },
        "WRAMMem": {
        "Box": {
            "Show": false,
            "Height": 660,
            "Width": 650
        }
        },
        "ProgMem": {
        "Box": {
            "Show": false,
            "Height": 200,
            "Width": 650
        }
        },
        "HRAMMem": {
        "Box": {
            "Show": false,
            "Height": 200,
            "Width": 650
        }
        },
        "OAMMem": {
        "Box": {
            "Show": false,
            "Height": 220,
            "Width": 650
        }
        },
        "Rewind": {
        "Box": {
            "Show": false,
            "Height": 300,
            "Width": 400
        }
        },
        "Disassembly": {
        "Box": {
            "Show": false,
            "Height": 500,
            "Width": 700
        }
        },
        "Debugger": {
        "Box": {
            "Show": true,
            "Height": 500,
            "Width": 400
        }
        },
        "Timing": {
        "Box": {
            "Show": true,
            "Height": 200,
            "Width": 400
        }
        },
        "Registers": {
        "Box": {
            "Show": false,
            "Height": 400,
            "Width": 200
        }
        },
        "APU": {
        "Box": {
            "Show": false,
            "Height": 400,
            "Width": 200
        }
        },
        "PPU": {
        "Box": {
            "Show": false,
            "Height": 800,
            "Width": 200
        }
        },
        "ViewPort": {
        "Box": {
            "Show": true,
            "Height": 700,
            "Width": 750
        },
        "Graphics": {
            "ShowGrid": false,
            "BlockSize": 8,
            "Scale": 4,
            "ShowAddress": true,
            "StartAddress": 0,
            "BlockIncrement": 8,
            "LineIncrement": 1,
            "DecimalAddress": true,
            "ShowOffsets": true
        }
        },
        "JoyPad": {
        "Box": {
            "Show": true,
            "Height": 350,
            "Width": 450
        },
        "Graphics": {
            "ShowGrid": false,
            "BlockSize": 8,
            "Scale": 8,
            "ShowAddress": false,
            "StartAddress": 0,
            "BlockIncrement": 0,
            "LineIncrement": 0,
            "DecimalAddress": false,
            "ShowOffsets": false
        }
        },
        "TileData": {
        "Box": {
            "Show": false,
            "Height": 400,
            "Width": 500
        },
        "Graphics": {
            "ShowGrid": true,
            "BlockSize": 8,
            "Scale": 2,
            "ShowAddress": true,
            "StartAddress": 32768,
            "BlockIncrement": 16,
            "LineIncrement": 0,
            "DecimalAddress": false,
            "ShowOffsets": false
        }
        },
        "TileMap1": {
        "Box": {
            "Show": false,
            "Height": 650,
            "Width": 650
        },
        "Graphics": {
            "ShowGrid": true,
            "BlockSize": 8,
            "Scale": 2,
            "ShowAddress": true,
            "StartAddress": 38912,
            "BlockIncrement": 1,
            "LineIncrement": 0,
            "DecimalAddress": false,
            "ShowOffsets": false
        }
        },
        "TileMap2": {
        "Box": {
            "Show": false,
            "Height": 650,
            "Width": 650
        },
        "Graphics": {
            "ShowGrid": true,
            "BlockSize": 8,
            "Scale": 2,
            "ShowAddress": true,
            "StartAddress": 39936,
            "BlockIncrement": 1,
            "LineIncrement": 0,
            "DecimalAddress": false,
            "ShowOffsets": false
        }
        },
        "OAMBuffer": {
        "Box": {
            "Show": true,
            "Height": 150,
            "Width": 450
        },
        "Graphics": {
            "ShowGrid": true,
            "BlockSize": 8,
            "Scale": 4,
            "ShowAddress": false,
            "StartAddress": 0,
            "BlockIncrement": 1,
            "LineIncrement": 0,
            "DecimalAddress": false,
            "ShowOffsets": false
        }
        },
        "OAMGraphics": {
        "Box": {
            "Show": true,
            "Height": 200,
            "Width": 450
        },
        "Graphics": {
            "ShowGrid": true,
            "BlockSize": 8,
            "Scale": 4,
            "ShowAddress": false,
            "StartAddress": 65024,
            "BlockIncrement": 4,
            "LineIncrement": 0,
            "DecimalAddress": false,
            "ShowOffsets": false
        }
        },
        "OAMList": {
        "Box": {
            "Show": false,
            "Height": 400,
            "Width": 450
        }
        }
    }
}

function getConfig() {
   return loadObject("config", DefaultConfig);
}

function saveConfig(config) {
    storeObject("config", config);
}
