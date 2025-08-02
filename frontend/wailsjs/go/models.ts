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
	    ROMLocation: string;
	    Model: model.Config;
	    PProfURL: string;
	    GUI: ConfigGUI;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Location = source["Location"];
	        this.ROMLocation = source["ROMLocation"];
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
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	export class Range {
	    Begin: number;
	    End: number;
	
	    static createFrom(source: any = {}) {
	        return new Range(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Begin = source["Begin"];
	        this.End = source["End"];
	    }
	}
	export class MachineStateRequest {
	    OpenBoxes: Record<string, boolean>;
	    Numbers: Record<string, number>;
	    Ranges: Record<string, Range>;
	    ClickedNumber: string;
	    ClickedRange: string;
	
	    static createFrom(source: any = {}) {
	        return new MachineStateRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OpenBoxes = source["OpenBoxes"];
	        this.Numbers = source["Numbers"];
	        this.Ranges = this.convertValues(source["Ranges"], Range, true);
	        this.ClickedNumber = source["ClickedNumber"];
	        this.ClickedRange = source["ClickedRange"];
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
	
	export class ConfigDisassembler {
	    Trace: boolean;
	    Enable: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDisassembler(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Trace = source["Trace"];
	        this.Enable = source["Enable"];
	    }
	}
	export class ConfigDebug {
	    RewindSize: number;
	    PanicOnStackUnderflow: boolean;
	    Disassembler: ConfigDisassembler;
	
	    static createFrom(source: any = {}) {
	        return new ConfigDebug(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.RewindSize = source["RewindSize"];
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
	    BootROM: ConfigBootROM;
	    Debug: ConfigDebug;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Clock = this.convertValues(source["Clock"], ConfigClock);
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
	
	
	

}

