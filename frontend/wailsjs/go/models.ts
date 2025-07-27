export namespace color {
	
	export class RGBA {
	    R: number;
	    G: number;
	    B: number;
	    A: number;
	
	    static createFrom(source: any = {}) {
	        return new RGBA(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.R = source["R"];
	        this.G = source["G"];
	        this.B = source["B"];
	        this.A = source["A"];
	    }
	}

}

export namespace main {
	
	export class ConfigOAMList {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigOAMList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigOAMGraphics {
	    Box: ConfigBox;
	    Graphics: ConfigGraphics;
	
	    static createFrom(source: any = {}) {
	        return new ConfigOAMGraphics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigOAMBuffer {
	    Box: ConfigBox;
	    Graphics: ConfigGraphics;
	
	    static createFrom(source: any = {}) {
	        return new ConfigOAMBuffer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigTileMap {
	    Box: ConfigBox;
	    Graphics: ConfigGraphics;
	
	    static createFrom(source: any = {}) {
	        return new ConfigTileMap(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigTileData {
	    Box: ConfigBox;
	    Graphics: ConfigGraphics;
	
	    static createFrom(source: any = {}) {
	        return new ConfigTileData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigJoyPad {
	    Box: ConfigBox;
	    Graphics: ConfigGraphics;
	
	    static createFrom(source: any = {}) {
	        return new ConfigJoyPad(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigGraphics {
	    ShowGrid: boolean;
	    BlockSize: number;
	    Scale: number;
	    ShowAddress: boolean;
	    StartAddress: number;
	    BlockIncrement: number;
	    LineIncrement: number;
	    DecimalAddress: boolean;
	    ShowOffsets: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConfigGraphics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ShowGrid = source["ShowGrid"];
	        this.BlockSize = source["BlockSize"];
	        this.Scale = source["Scale"];
	        this.ShowAddress = source["ShowAddress"];
	        this.StartAddress = source["StartAddress"];
	        this.BlockIncrement = source["BlockIncrement"];
	        this.LineIncrement = source["LineIncrement"];
	        this.DecimalAddress = source["DecimalAddress"];
	        this.ShowOffsets = source["ShowOffsets"];
	    }
	}
	export class ConfigViewPort {
	    Box: ConfigBox;
	    Graphics: ConfigGraphics;
	
	    static createFrom(source: any = {}) {
	        return new ConfigViewPort(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigPPU {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigPPU(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigAPU {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigAPU(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigRegisters {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigRegisters(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigTiming {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigTiming(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigDebugger {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDebugger(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigDisassembly {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDisassembly(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigRewind {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigRewind(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigOAM {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigOAM(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigHRAM {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigHRAM(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigProgram {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigProgram(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigWRAM {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigWRAM(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigBox {
	    Show: boolean;
	    Height: number;
	    Width: number;
	
	    static createFrom(source: any = {}) {
	        return new ConfigBox(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Show = source["Show"];
	        this.Height = source["Height"];
	        this.Width = source["Width"];
	    }
	}
	export class ConfigVRAM {
	    Box: ConfigBox;
	
	    static createFrom(source: any = {}) {
	        return new ConfigVRAM(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Box = this.convertValues(source["Box"], ConfigBox);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigGraphicsGlobal {
	    Overlay: boolean;
	    GridColor: color.RGBA;
	    FillColor: color.RGBA;
	    DashLen: number;
	    GridThickness: number;
	    Font: string;
	
	    static createFrom(source: any = {}) {
	        return new ConfigGraphicsGlobal(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Overlay = source["Overlay"];
	        this.GridColor = this.convertValues(source["GridColor"], color.RGBA);
	        this.FillColor = this.convertValues(source["FillColor"], color.RGBA);
	        this.DashLen = source["DashLen"];
	        this.GridThickness = source["GridThickness"];
	        this.Font = source["Font"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigGUI {
	    Graphics: ConfigGraphicsGlobal;
	    VRAMMem: ConfigVRAM;
	    WRAMMem: ConfigWRAM;
	    ProgMem: ConfigProgram;
	    HRAMMem: ConfigHRAM;
	    OAMMem: ConfigOAM;
	    Rewind: ConfigRewind;
	    Disassembly: ConfigDisassembly;
	    Debugger: ConfigDebugger;
	    Timing: ConfigTiming;
	    Registers: ConfigRegisters;
	    APU: ConfigAPU;
	    PPU: ConfigPPU;
	    ViewPort: ConfigViewPort;
	    JoyPad: ConfigJoyPad;
	    TileData: ConfigTileData;
	    TileMap1: ConfigTileMap;
	    TileMap2: ConfigTileMap;
	    OAMBuffer: ConfigOAMBuffer;
	    OAMGraphics: ConfigOAMGraphics;
	    OAMList: ConfigOAMList;
	
	    static createFrom(source: any = {}) {
	        return new ConfigGUI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Graphics = this.convertValues(source["Graphics"], ConfigGraphicsGlobal);
	        this.VRAMMem = this.convertValues(source["VRAMMem"], ConfigVRAM);
	        this.WRAMMem = this.convertValues(source["WRAMMem"], ConfigWRAM);
	        this.ProgMem = this.convertValues(source["ProgMem"], ConfigProgram);
	        this.HRAMMem = this.convertValues(source["HRAMMem"], ConfigHRAM);
	        this.OAMMem = this.convertValues(source["OAMMem"], ConfigOAM);
	        this.Rewind = this.convertValues(source["Rewind"], ConfigRewind);
	        this.Disassembly = this.convertValues(source["Disassembly"], ConfigDisassembly);
	        this.Debugger = this.convertValues(source["Debugger"], ConfigDebugger);
	        this.Timing = this.convertValues(source["Timing"], ConfigTiming);
	        this.Registers = this.convertValues(source["Registers"], ConfigRegisters);
	        this.APU = this.convertValues(source["APU"], ConfigAPU);
	        this.PPU = this.convertValues(source["PPU"], ConfigPPU);
	        this.ViewPort = this.convertValues(source["ViewPort"], ConfigViewPort);
	        this.JoyPad = this.convertValues(source["JoyPad"], ConfigJoyPad);
	        this.TileData = this.convertValues(source["TileData"], ConfigTileData);
	        this.TileMap1 = this.convertValues(source["TileMap1"], ConfigTileMap);
	        this.TileMap2 = this.convertValues(source["TileMap2"], ConfigTileMap);
	        this.OAMBuffer = this.convertValues(source["OAMBuffer"], ConfigOAMBuffer);
	        this.OAMGraphics = this.convertValues(source["OAMGraphics"], ConfigOAMGraphics);
	        this.OAMList = this.convertValues(source["OAMList"], ConfigOAMList);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Config {
	    Location: string;
	    Model: model.Config;
	    PProfURL: string;
	    GUI: ConfigGUI;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Location = source["Location"];
	        this.Model = this.convertValues(source["Model"], model.Config);
	        this.PProfURL = source["PProfURL"];
	        this.GUI = this.convertValues(source["GUI"], ConfigGUI);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

}

export namespace model {
	
	export class Mixer {
	    RegChannelPan: number;
	    RegMasterVolumeVINPan: number;
	
	    static createFrom(source: any = {}) {
	        return new Mixer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RegChannelPan = source["RegChannelPan"];
	        this.RegMasterVolumeVINPan = source["RegMasterVolumeVINPan"];
	    }
	}
	export class NoiseChannel {
	    // Go type: PeriodCounter
	    PeriodCounter: any;
	    // Go type: LengthTimer
	    LengthTimer: any;
	    // Go type: Envelope
	    Envelope: any;
	    RegLengthTimer: number;
	    RegVolumeEnvelope: number;
	    RegRNG: number;
	    RegCtl: number;
	
	    static createFrom(source: any = {}) {
	        return new NoiseChannel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PeriodCounter = this.convertValues(source["PeriodCounter"], null);
	        this.LengthTimer = this.convertValues(source["LengthTimer"], null);
	        this.Envelope = this.convertValues(source["Envelope"], null);
	        this.RegLengthTimer = source["RegLengthTimer"];
	        this.RegVolumeEnvelope = source["RegVolumeEnvelope"];
	        this.RegRNG = source["RegRNG"];
	        this.RegCtl = source["RegCtl"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class WaveChannel {
	    // Go type: PeriodCounter
	    PeriodCounter: any;
	    // Go type: LengthTimer
	    LengthTimer: any;
	    RegDACEn: number;
	    RegLengthTimer: number;
	    RegOutputLevel: number;
	    RegPeriodLow: number;
	    RegPeriodHighCtl: number;
	
	    static createFrom(source: any = {}) {
	        return new WaveChannel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PeriodCounter = this.convertValues(source["PeriodCounter"], null);
	        this.LengthTimer = this.convertValues(source["LengthTimer"], null);
	        this.RegDACEn = source["RegDACEn"];
	        this.RegLengthTimer = source["RegLengthTimer"];
	        this.RegOutputLevel = source["RegOutputLevel"];
	        this.RegPeriodLow = source["RegPeriodLow"];
	        this.RegPeriodHighCtl = source["RegPeriodHighCtl"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PulseChannel {
	    // Go type: PeriodCounter
	    PeriodCounter: any;
	    // Go type: LengthTimer
	    LengthTimer: any;
	    // Go type: DutyGenerator
	    DutyGenerator: any;
	    // Go type: Envelope
	    Envelope: any;
	    RegLengthDuty: number;
	    RegVolumeEnvelope: number;
	    RegPeriodLow: number;
	    RegPeriodHighCtl: number;
	
	    static createFrom(source: any = {}) {
	        return new PulseChannel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PeriodCounter = this.convertValues(source["PeriodCounter"], null);
	        this.LengthTimer = this.convertValues(source["LengthTimer"], null);
	        this.DutyGenerator = this.convertValues(source["DutyGenerator"], null);
	        this.Envelope = this.convertValues(source["Envelope"], null);
	        this.RegLengthDuty = source["RegLengthDuty"];
	        this.RegVolumeEnvelope = source["RegVolumeEnvelope"];
	        this.RegPeriodLow = source["RegPeriodLow"];
	        this.RegPeriodHighCtl = source["RegPeriodHighCtl"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Envelope {
	
	
	    static createFrom(source: any = {}) {
	        return new Envelope(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class DutyGenerator {
	
	
	    static createFrom(source: any = {}) {
	        return new DutyGenerator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class LengthTimer {
	
	
	    static createFrom(source: any = {}) {
	        return new LengthTimer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class PeriodCounter {
	
	
	    static createFrom(source: any = {}) {
	        return new PeriodCounter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Sweep {
	    RegSweep: number;
	
	    static createFrom(source: any = {}) {
	        return new Sweep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RegSweep = source["RegSweep"];
	    }
	}
	export class PulseChannelWithSweep {
	    Sweep: Sweep;
	    // Go type: PeriodCounter
	    PeriodCounter: any;
	    // Go type: LengthTimer
	    LengthTimer: any;
	    // Go type: DutyGenerator
	    DutyGenerator: any;
	    // Go type: Envelope
	    Envelope: any;
	    RegLengthDuty: number;
	    RegVolumeEnvelope: number;
	    RegPeriodLow: number;
	    RegPeriodHighCtl: number;
	
	    static createFrom(source: any = {}) {
	        return new PulseChannelWithSweep(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Sweep = this.convertValues(source["Sweep"], Sweep);
	        this.PeriodCounter = this.convertValues(source["PeriodCounter"], null);
	        this.LengthTimer = this.convertValues(source["LengthTimer"], null);
	        this.DutyGenerator = this.convertValues(source["DutyGenerator"], null);
	        this.Envelope = this.convertValues(source["Envelope"], null);
	        this.RegLengthDuty = source["RegLengthDuty"];
	        this.RegVolumeEnvelope = source["RegVolumeEnvelope"];
	        this.RegPeriodLow = source["RegPeriodLow"];
	        this.RegPeriodHighCtl = source["RegPeriodHighCtl"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class APU {
	    Offset: number;
	    Data: number[];
	    MasterCtl: number;
	    Pulse1: PulseChannelWithSweep;
	    Pulse2: PulseChannel;
	    Wave: WaveChannel;
	    Noise: NoiseChannel;
	    Mixer: Mixer;
	    DIVAPU: number;
	
	    static createFrom(source: any = {}) {
	        return new APU(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Offset = source["Offset"];
	        this.Data = source["Data"];
	        this.MasterCtl = source["MasterCtl"];
	        this.Pulse1 = this.convertValues(source["Pulse1"], PulseChannelWithSweep);
	        this.Pulse2 = this.convertValues(source["Pulse2"], PulseChannel);
	        this.Wave = this.convertValues(source["Wave"], WaveChannel);
	        this.Noise = this.convertValues(source["Noise"], NoiseChannel);
	        this.Mixer = this.convertValues(source["Mixer"], Mixer);
	        this.DIVAPU = source["DIVAPU"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Pixel {
	    ColorIDXBGPriority: number;
	    Palette: number;
	
	    static createFrom(source: any = {}) {
	        return new Pixel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ColorIDXBGPriority = source["ColorIDXBGPriority"];
	        this.Palette = source["Palette"];
	    }
	}
	export class FIFO {
	    Slots: Pixel[];
	    Level: number;
	    ShiftPos: number;
	    PushPos: number;
	
	    static createFrom(source: any = {}) {
	        return new FIFO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Slots = this.convertValues(source["Slots"], Pixel);
	        this.Level = source["Level"];
	        this.ShiftPos = source["ShiftPos"];
	        this.PushPos = source["PushPos"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Shifter {
	    Discard: number;
	    Suspended: boolean;
	    X: number;
	    LastShifted: number;
	    PPU?: PPU;
	
	    static createFrom(source: any = {}) {
	        return new Shifter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Discard = source["Discard"];
	        this.Suspended = source["Suspended"];
	        this.X = source["X"];
	        this.LastShifted = source["LastShifted"];
	        this.PPU = this.convertValues(source["PPU"], PPU);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SpriteFetcher {
	    Cycle: number;
	    State: number;
	    X: number;
	    TileIndex: number;
	    TileLSBAddr: number;
	    TileLSB: number;
	    TileMSB: number;
	    Suspended: boolean;
	    PPU?: PPU;
	    SpriteIDX: number;
	    DoneX: number;
	
	    static createFrom(source: any = {}) {
	        return new SpriteFetcher(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Cycle = source["Cycle"];
	        this.State = source["State"];
	        this.X = source["X"];
	        this.TileIndex = source["TileIndex"];
	        this.TileLSBAddr = source["TileLSBAddr"];
	        this.TileLSB = source["TileLSB"];
	        this.TileMSB = source["TileMSB"];
	        this.Suspended = source["Suspended"];
	        this.PPU = this.convertValues(source["PPU"], PPU);
	        this.SpriteIDX = source["SpriteIDX"];
	        this.DoneX = source["DoneX"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Object {
	    X: number;
	    Y: number;
	    TileIndex: number;
	    Attributes: number;
	
	    static createFrom(source: any = {}) {
	        return new Object(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	        this.TileIndex = source["TileIndex"];
	        this.Attributes = source["Attributes"];
	    }
	}
	export class OAMBuffer {
	    Buffer: Object[];
	    Level: number;
	
	    static createFrom(source: any = {}) {
	        return new OAMBuffer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Buffer = this.convertValues(source["Buffer"], Object);
	        this.Level = source["Level"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class UserMessage {
	    // Go type: time
	    Time: any;
	    Message: string;
	    Warn: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserMessage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Time = this.convertValues(source["Time"], null);
	        this.Message = source["Message"];
	        this.Warn = source["Warn"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ClockRT {
	    Cycle: number;
	    // Go type: atomic
	    Running: any;
	
	    static createFrom(source: any = {}) {
	        return new ClockRT(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Cycle = source["Cycle"];
	        this.Running = this.convertValues(source["Running"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DisInstruction {
	    Raw: number[];
	    Address: number;
	    Opcode: number;
	    Visited: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DisInstruction(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Raw = source["Raw"];
	        this.Address = source["Address"];
	        this.Opcode = source["Opcode"];
	        this.Visited = source["Visited"];
	    }
	}
	export class Block {
	    CanExplore: boolean;
	    Begin: number;
	    Source: number[];
	    Decoded: DisInstruction[];
	
	    static createFrom(source: any = {}) {
	        return new Block(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.CanExplore = source["CanExplore"];
	        this.Begin = source["Begin"];
	        this.Source = source["Source"];
	        this.Decoded = this.convertValues(source["Decoded"], DisInstruction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Debug {
	    Config?: ConfigDisassembler;
	    // Go type: Block
	    Program: any;
	    // Go type: Block
	    HRAM: any;
	    // Go type: Block
	    WRAM: any;
	    PC: number;
	    // Go type: atomic
	    BreakX: any;
	    // Go type: atomic
	    BreakY: any;
	    // Go type: atomic
	    BreakPC: any;
	    // Go type: atomic
	    BreakIR: any;
	    // Go type: ClockRT
	    CLK?: any;
	    Warnings: Record<string, UserMessage>;
	
	    static createFrom(source: any = {}) {
	        return new Debug(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Config = this.convertValues(source["Config"], ConfigDisassembler);
	        this.Program = this.convertValues(source["Program"], null);
	        this.HRAM = this.convertValues(source["HRAM"], null);
	        this.WRAM = this.convertValues(source["WRAM"], null);
	        this.PC = source["PC"];
	        this.BreakX = this.convertValues(source["BreakX"], null);
	        this.BreakY = this.convertValues(source["BreakY"], null);
	        this.BreakPC = this.convertValues(source["BreakPC"], null);
	        this.BreakIR = this.convertValues(source["BreakIR"], null);
	        this.CLK = this.convertValues(source["CLK"], null);
	        this.Warnings = this.convertValues(source["Warnings"], UserMessage, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Timer {
	    Mem: MemoryRegion;
	    APU?: APU;
	    Interrupts?: Interrupts;
	    DIV: number;
	
	    static createFrom(source: any = {}) {
	        return new Timer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mem = this.convertValues(source["Mem"], MemoryRegion);
	        this.APU = this.convertValues(source["APU"], APU);
	        this.Interrupts = this.convertValues(source["Interrupts"], Interrupts);
	        this.DIV = source["DIV"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Prohibited {
	    FEA0toFEFF: MemoryRegion;
	    FF71toFF7F: MemoryRegion;
	
	    static createFrom(source: any = {}) {
	        return new Prohibited(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FEA0toFEFF = this.convertValues(source["FEA0toFEFF"], MemoryRegion);
	        this.FF71toFF7F = this.convertValues(source["FF71toFF7F"], MemoryRegion);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Serial {
	    Mem: MemoryRegion;
	
	    static createFrom(source: any = {}) {
	        return new Serial(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Mem = this.convertValues(source["Mem"], MemoryRegion);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Joypad {
	    Interrupts?: Interrupts;
	    Written: MemoryRegion;
	    Action: number;
	    Direction: number;
	
	    static createFrom(source: any = {}) {
	        return new Joypad(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Interrupts = this.convertValues(source["Interrupts"], Interrupts);
	        this.Written = this.convertValues(source["Written"], MemoryRegion);
	        this.Action = source["Action"];
	        this.Direction = source["Direction"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MBCFeatures {
	    ID: number;
	    RAM: boolean;
	    Battery: boolean;
	    RTC: boolean;
	    Rumble: boolean;
	    NROMBanks: number;
	    NRAMBanks: number;
	
	    static createFrom(source: any = {}) {
	        return new MBCFeatures(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.RAM = source["RAM"];
	        this.Battery = source["Battery"];
	        this.RTC = source["RTC"];
	        this.Rumble = source["Rumble"];
	        this.NROMBanks = source["NROMBanks"];
	        this.NRAMBanks = source["NRAMBanks"];
	    }
	}
	export class Cartridge {
	    ROM: number[][];
	    RAM: number[][];
	    CurrROMBank0: MemoryRegion;
	    CurrROMBankN: MemoryRegion;
	    CurrRAMBank: MemoryRegion;
	    MBCFeatures: MBCFeatures;
	    ExtRAMEnabled: boolean;
	    BankNo1: number;
	    BankNo2: number;
	    BankModeSel: number;
	    SelectedRAMBank: number;
	    SelectedROMBank: number;
	
	    static createFrom(source: any = {}) {
	        return new Cartridge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ROM = source["ROM"];
	        this.RAM = source["RAM"];
	        this.CurrROMBank0 = this.convertValues(source["CurrROMBank0"], MemoryRegion);
	        this.CurrROMBankN = this.convertValues(source["CurrROMBankN"], MemoryRegion);
	        this.CurrRAMBank = this.convertValues(source["CurrRAMBank"], MemoryRegion);
	        this.MBCFeatures = this.convertValues(source["MBCFeatures"], MBCFeatures);
	        this.ExtRAMEnabled = source["ExtRAMEnabled"];
	        this.BankNo1 = source["BankNo1"];
	        this.BankNo2 = source["BankNo2"];
	        this.BankModeSel = source["BankModeSel"];
	        this.SelectedRAMBank = source["SelectedRAMBank"];
	        this.SelectedROMBank = source["SelectedROMBank"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BootROMLock {
	    Offset: number;
	    Data: number[];
	    BootOff: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BootROMLock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Offset = source["Offset"];
	        this.Data = source["Data"];
	        this.BootOff = source["BootOff"];
	    }
	}
	export class ConfigDisassembler {
	    Trace: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDisassembler(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Trace = source["Trace"];
	    }
	}
	export class ConfigDebug {
	    PanicOnStackUnderflow: boolean;
	    Disassembler: ConfigDisassembler;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDebug(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PanicOnStackUnderflow = source["PanicOnStackUnderflow"];
	        this.Disassembler = this.convertValues(source["Disassembler"], ConfigDisassembler);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigBootROM {
	    Variant: string;
	    Skip: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConfigBootROM(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Variant = source["Variant"];
	        this.Skip = source["Skip"];
	    }
	}
	export class ConfigROM {
	    Location: string;
	
	    static createFrom(source: any = {}) {
	        return new ConfigROM(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Location = source["Location"];
	    }
	}
	export class ConfigClock {
	    SpeedPercent: number;
	
	    static createFrom(source: any = {}) {
	        return new ConfigClock(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SpeedPercent = source["SpeedPercent"];
	    }
	}
	export class Config {
	    Clock: ConfigClock;
	    ROM: ConfigROM;
	    BootROM: ConfigBootROM;
	    Debug: ConfigDebug;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Clock = this.convertValues(source["Clock"], ConfigClock);
	        this.ROM = this.convertValues(source["ROM"], ConfigROM);
	        this.BootROM = this.convertValues(source["BootROM"], ConfigBootROM);
	        this.Debug = this.convertValues(source["Debug"], ConfigDebug);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Bus {
	    Data: number;
	    Address: number;
	    Config?: Config;
	    BootROMLock?: BootROMLock;
	    BootROM?: MemoryRegion;
	    VRAM?: MemoryRegion;
	    HRAM?: MemoryRegion;
	    WRAM?: MemoryRegion;
	    APU?: APU;
	    OAM?: MemoryRegion;
	    PPU?: PPU;
	    Cartridge?: Cartridge;
	    Joypad?: Joypad;
	    Interrupts?: Interrupts;
	    Serial?: Serial;
	    Prohibited?: Prohibited;
	    Timer?: Timer;
	
	    static createFrom(source: any = {}) {
	        return new Bus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = source["Data"];
	        this.Address = source["Address"];
	        this.Config = this.convertValues(source["Config"], Config);
	        this.BootROMLock = this.convertValues(source["BootROMLock"], BootROMLock);
	        this.BootROM = this.convertValues(source["BootROM"], MemoryRegion);
	        this.VRAM = this.convertValues(source["VRAM"], MemoryRegion);
	        this.HRAM = this.convertValues(source["HRAM"], MemoryRegion);
	        this.WRAM = this.convertValues(source["WRAM"], MemoryRegion);
	        this.APU = this.convertValues(source["APU"], APU);
	        this.OAM = this.convertValues(source["OAM"], MemoryRegion);
	        this.PPU = this.convertValues(source["PPU"], PPU);
	        this.Cartridge = this.convertValues(source["Cartridge"], Cartridge);
	        this.Joypad = this.convertValues(source["Joypad"], Joypad);
	        this.Interrupts = this.convertValues(source["Interrupts"], Interrupts);
	        this.Serial = this.convertValues(source["Serial"], Serial);
	        this.Prohibited = this.convertValues(source["Prohibited"], Prohibited);
	        this.Timer = this.convertValues(source["Timer"], Timer);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DMA {
	    Reg: number;
	    Source: number;
	    Dest: number;
	    Bus?: Bus;
	
	    static createFrom(source: any = {}) {
	        return new DMA(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Reg = source["Reg"];
	        this.Source = source["Source"];
	        this.Dest = source["Dest"];
	        this.Bus = this.convertValues(source["Bus"], Bus);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MemoryRegion {
	    Offset: number;
	    Data: number[];
	
	    static createFrom(source: any = {}) {
	        return new MemoryRegion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Offset = source["Offset"];
	        this.Data = source["Data"];
	    }
	}
	export class Interrupts {
	    MemIF: MemoryRegion;
	    MemIE: MemoryRegion;
	    IME: boolean;
	    PendingInterrupt: number;
	
	    static createFrom(source: any = {}) {
	        return new Interrupts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.MemIF = this.convertValues(source["MemIF"], MemoryRegion);
	        this.MemIE = this.convertValues(source["MemIE"], MemoryRegion);
	        this.IME = source["IME"];
	        this.PendingInterrupt = source["PendingInterrupt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Stat {
	    Reg: number;
	    Interrupts?: Interrupts;
	
	    static createFrom(source: any = {}) {
	        return new Stat(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Reg = source["Reg"];
	        this.Interrupts = this.convertValues(source["Interrupts"], Interrupts);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PPU {
	    RegLCDC: number;
	    Stat: Stat;
	    RegSCY: number;
	    RegSCX: number;
	    RegWY: number;
	    RegWX: number;
	    RegLY: number;
	    RegLYC: number;
	    RegBGP: number;
	    RegOBP0: number;
	    RegOBP1: number;
	    DMA: DMA;
	    Bus?: Bus;
	    Interrupts?: Interrupts;
	    MemoryRegion: MemoryRegion;
	    Debug?: Debug;
	    FrameCount: number;
	    Config?: Config;
	    Mode: number;
	    OAMScanCycle: number;
	    OAMBuffer: OAMBuffer;
	    PixelDrawCycle: number;
	    BackgroundFetcher: BackgroundFetcher;
	    SpriteFetcher: SpriteFetcher;
	    Shifter: Shifter;
	    BackgroundFIFO: FIFO;
	    SpriteFIFO: FIFO;
	    BGPalette: number;
	    OBJPalette0: number;
	    OBJPalette1: number;
	    HBlankRemainingCycles: number;
	    VBlankLineRemainingCycles: number;
	    FBBackground: number[][];
	    FBWindow: number[][];
	    FBViewport: number[][];
	
	    static createFrom(source: any = {}) {
	        return new PPU(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RegLCDC = source["RegLCDC"];
	        this.Stat = this.convertValues(source["Stat"], Stat);
	        this.RegSCY = source["RegSCY"];
	        this.RegSCX = source["RegSCX"];
	        this.RegWY = source["RegWY"];
	        this.RegWX = source["RegWX"];
	        this.RegLY = source["RegLY"];
	        this.RegLYC = source["RegLYC"];
	        this.RegBGP = source["RegBGP"];
	        this.RegOBP0 = source["RegOBP0"];
	        this.RegOBP1 = source["RegOBP1"];
	        this.DMA = this.convertValues(source["DMA"], DMA);
	        this.Bus = this.convertValues(source["Bus"], Bus);
	        this.Interrupts = this.convertValues(source["Interrupts"], Interrupts);
	        this.MemoryRegion = this.convertValues(source["MemoryRegion"], MemoryRegion);
	        this.Debug = this.convertValues(source["Debug"], Debug);
	        this.FrameCount = source["FrameCount"];
	        this.Config = this.convertValues(source["Config"], Config);
	        this.Mode = source["Mode"];
	        this.OAMScanCycle = source["OAMScanCycle"];
	        this.OAMBuffer = this.convertValues(source["OAMBuffer"], OAMBuffer);
	        this.PixelDrawCycle = source["PixelDrawCycle"];
	        this.BackgroundFetcher = this.convertValues(source["BackgroundFetcher"], BackgroundFetcher);
	        this.SpriteFetcher = this.convertValues(source["SpriteFetcher"], SpriteFetcher);
	        this.Shifter = this.convertValues(source["Shifter"], Shifter);
	        this.BackgroundFIFO = this.convertValues(source["BackgroundFIFO"], FIFO);
	        this.SpriteFIFO = this.convertValues(source["SpriteFIFO"], FIFO);
	        this.BGPalette = source["BGPalette"];
	        this.OBJPalette0 = source["OBJPalette0"];
	        this.OBJPalette1 = source["OBJPalette1"];
	        this.HBlankRemainingCycles = source["HBlankRemainingCycles"];
	        this.VBlankLineRemainingCycles = source["VBlankLineRemainingCycles"];
	        this.FBBackground = source["FBBackground"];
	        this.FBWindow = source["FBWindow"];
	        this.FBViewport = source["FBViewport"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BackgroundFetcher {
	    Cycle: number;
	    State: number;
	    X: number;
	    TileIndex: number;
	    TileLSBAddr: number;
	    TileLSB: number;
	    TileMSB: number;
	    Suspended: boolean;
	    PPU?: PPU;
	    TileIndexAddr: number;
	    TileOffsetX: number;
	    TileOffsetY: number;
	    WindowYReached: boolean;
	    WindowFetching: boolean;
	    WindowLineCounter: number;
	    WindowPixelRenderedThisScanline: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BackgroundFetcher(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Cycle = source["Cycle"];
	        this.State = source["State"];
	        this.X = source["X"];
	        this.TileIndex = source["TileIndex"];
	        this.TileLSBAddr = source["TileLSBAddr"];
	        this.TileLSB = source["TileLSB"];
	        this.TileMSB = source["TileMSB"];
	        this.Suspended = source["Suspended"];
	        this.PPU = this.convertValues(source["PPU"], PPU);
	        this.TileIndexAddr = source["TileIndexAddr"];
	        this.TileOffsetX = source["TileOffsetX"];
	        this.TileOffsetY = source["TileOffsetY"];
	        this.WindowYReached = source["WindowYReached"];
	        this.WindowFetching = source["WindowFetching"];
	        this.WindowLineCounter = source["WindowLineCounter"];
	        this.WindowPixelRenderedThisScanline = source["WindowPixelRenderedThisScanline"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	export class CodeSection {
	    Instructions: DisInstruction[];
	
	    static createFrom(source: any = {}) {
	        return new CodeSection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Instructions = this.convertValues(source["Instructions"], DisInstruction);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	export class DataSection {
	    Raw: number[];
	    Address: number;
	
	    static createFrom(source: any = {}) {
	        return new DataSection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Raw = source["Raw"];
	        this.Address = source["Address"];
	    }
	}
	export class Disassembly {
	    PC: number;
	    Code: CodeSection[];
	    Data: DataSection[];
	
	    static createFrom(source: any = {}) {
	        return new Disassembly(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.PC = source["PC"];
	        this.Code = this.convertValues(source["Code"], CodeSection);
	        this.Data = this.convertValues(source["Data"], DataSection);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExecLogEntry {
	    Instruction: DisInstruction;
	    BranchResult: number;
	
	    static createFrom(source: any = {}) {
	        return new ExecLogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Instruction = this.convertValues(source["Instruction"], DisInstruction);
	        this.BranchResult = source["BranchResult"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Rewind {
	    Buffer: ExecLogEntry[];
	    Idx: number;
	    Full: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Rewind(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Buffer = this.convertValues(source["Buffer"], ExecLogEntry);
	        this.Idx = source["Idx"];
	        this.Full = source["Full"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PixelFIFODump {
	    Slots: Pixel[];
	    Level: number;
	
	    static createFrom(source: any = {}) {
	        return new PixelFIFODump(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Slots = this.convertValues(source["Slots"], Pixel);
	        this.Level = source["Level"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PPUDump {
	    Registers: number[];
	    BGFIFO: PixelFIFODump;
	    SpriteFIFO: PixelFIFODump;
	    LastShifted: number;
	    OAMScanCycle: number;
	    PixelDrawCycle: number;
	    HBlankRemainingCycles: number;
	    VBlankLineRemainingCycles: number;
	    PixelShifter: Shifter;
	    BackgroundFetcher: BackgroundFetcher;
	    SpriteFetcher: SpriteFetcher;
	    OAMBuffer: OAMBuffer;
	
	    static createFrom(source: any = {}) {
	        return new PPUDump(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Registers = source["Registers"];
	        this.BGFIFO = this.convertValues(source["BGFIFO"], PixelFIFODump);
	        this.SpriteFIFO = this.convertValues(source["SpriteFIFO"], PixelFIFODump);
	        this.LastShifted = source["LastShifted"];
	        this.OAMScanCycle = source["OAMScanCycle"];
	        this.PixelDrawCycle = source["PixelDrawCycle"];
	        this.HBlankRemainingCycles = source["HBlankRemainingCycles"];
	        this.VBlankLineRemainingCycles = source["VBlankLineRemainingCycles"];
	        this.PixelShifter = this.convertValues(source["PixelShifter"], Shifter);
	        this.BackgroundFetcher = this.convertValues(source["BackgroundFetcher"], BackgroundFetcher);
	        this.SpriteFetcher = this.convertValues(source["SpriteFetcher"], SpriteFetcher);
	        this.OAMBuffer = this.convertValues(source["OAMBuffer"], OAMBuffer);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RegisterFile {
	    A: number;
	    F: number;
	    B: number;
	    C: number;
	    D: number;
	    E: number;
	    H: number;
	    L: number;
	    PC: number;
	    SP: number;
	    IR: number;
	    TempZ: number;
	    TempW: number;
	
	    static createFrom(source: any = {}) {
	        return new RegisterFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.A = source["A"];
	        this.F = source["F"];
	        this.B = source["B"];
	        this.C = source["C"];
	        this.D = source["D"];
	        this.E = source["E"];
	        this.H = source["H"];
	        this.L = source["L"];
	        this.PC = source["PC"];
	        this.SP = source["SP"];
	        this.IR = source["IR"];
	        this.TempZ = source["TempZ"];
	        this.TempW = source["TempW"];
	    }
	}
	export class CoreDump {
	    Cycle: number;
	    Regs: RegisterFile;
	    ProgramStart: number;
	    ProgramEnd: number;
	    Program: number[];
	    HRAM: number[];
	    OAM: number[];
	    VRAM: number[];
	    APU: number[];
	    PPU: PPUDump;
	    Rewind?: Rewind;
	    Disassembly?: Disassembly;
	
	    static createFrom(source: any = {}) {
	        return new CoreDump(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Cycle = source["Cycle"];
	        this.Regs = this.convertValues(source["Regs"], RegisterFile);
	        this.ProgramStart = source["ProgramStart"];
	        this.ProgramEnd = source["ProgramEnd"];
	        this.Program = source["Program"];
	        this.HRAM = source["HRAM"];
	        this.OAM = source["OAM"];
	        this.VRAM = source["VRAM"];
	        this.APU = source["APU"];
	        this.PPU = this.convertValues(source["PPU"], PPUDump);
	        this.Rewind = this.convertValues(source["Rewind"], Rewind);
	        this.Disassembly = this.convertValues(source["Disassembly"], Disassembly);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	

}

