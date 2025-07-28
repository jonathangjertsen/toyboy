async function waitFrame() {
    await new Promise(requestAnimationFrame);
}
const FPS =  document.getElementById("fps")
const Keys = ['w', 'a', 's', 'd', 'k', 'l', 'n', 'm'];
const KeyState = Object.fromEntries(Keys.map(k => [k, false]));
const ButtonsArray = [
    document.getElementById("btn-up"),
    document.getElementById("btn-left"),
    document.getElementById("btn-down"),
    document.getElementById("btn-right"),
    document.getElementById("btn-b"),
    document.getElementById("btn-a"),
    document.getElementById("btn-select"),
    document.getElementById("btn-start"),
];
const Buttons = Object.fromEntries(Keys.map((k, i) => [k, ButtonsArray[i]]));
let Frame = null;
const CPURegistersText = document.getElementById("cpu-registers-text");
const PPURegistersText = document.getElementById("ppu-registers-text");
const APURegistersText = document.getElementById("apu-registers-text");
const WRAMText = document.getElementById("wram-text");
const HRAMText = document.getElementById("hram-text");
const OAMText = document.getElementById("oam-text");
const ClockText = document.getElementById("clock-text");
const DisassemblyText = document.getElementById("disassembly-text");

async function run() {
    config = await window.go.main.App.GetConfig();
    console.log(config);

    await window.go.main.App.StartGB();
    await window.go.main.App.StartWebSocketServer();

    const ws = new WebSocket('ws://localhost:8081/data');
    ws.binaryType = 'arraybuffer';

    let dataID = "";
    const decoder = new TextDecoder();
    ws.onmessage = function(event) {
        console.log(dataID);
        if (typeof event.data === "string") {
            dataID = event.data;
            return;
        }
        const data = new Uint8Array(event.data);
        switch (dataID) {
            case "Viewport": {
                Frame = data;
                break;
            }
            case "CPURegisters": {
                CPURegistersText.innerText = decoder.decode(data);
                break;
            }
            case "PPURegisters": {
                PPURegistersText.innerText = decoder.decode(data);
                break;
            }
            case "APURegisters": {
                APURegistersText.innerText = decoder.decode(data);
                break;
            }
            case "Disassembly": {
                DisassemblyText.innerText = decoder.decode(data);
                break;
            }
            case "HRAM": {
                HRAMText.innerText = decoder.decode(data);
                break;
            }
            case "WRAM": {
                WRAMText.innerText = decoder.decode(data);
                break;
            }
            case "OAM": {
                OAMText.innerText = decoder.decode(data);
                break;
            }
            case "CPUState": {
                const cpuState = data[0];
                console.log("CPU state", cpuState, typeof cpuState);
                if (cpuState === 0) {
                    RunBtn.disabled = false;
                    PauseBtn.disabled = true;
                } else {
                    RunBtn.disabled = true;
                    PauseBtn.disabled = false;
                }
                break;
            }
            case "Clock": {
                ClockText.innerText = decoder.decode(data);
                break;
            }
        }
        dataID = "";
    };
}

run();

async function updateSpeed() {
    requestAnimationFrame(updateSpeed);
}
updateSpeed();

function renderLoop() {
    if (Frame !== null) {
        Renderer.renderFrame(Frame, 160, 144);
    }
    requestAnimationFrame(renderLoop);
}
renderLoop();

async function setKeyState(k, state) {
    if (k in KeyState) {
        KeyState[k] = state
        if (state) {
            Buttons[k].classList.add("pressed");
        } else {
            Buttons[k].classList.remove("pressed");
        }
    }
    await window.go.main.App.SetKeyState(KeyState);
}

document.addEventListener('keydown', (e) => {
    setKeyState(e.key.toLowerCase(), true);
});

document.addEventListener('keyup', (e) => {
    setKeyState(e.key.toLowerCase(), false);
});

for (const [k, button] of Object.entries(Buttons)) {
    button.addEventListener('mousedown', (e) => {
        setKeyState(k, true)
    })
    button.addEventListener('mouseup', (e) => {
        setKeyState(k, false)
    })
}

let MachineReq = {
    OpenBoxes: {},
    Numbers: {},
    Ranges: {
        WRAM: {Begin:49152, End:49216},
        Disassembly: {Begin: 0, End: 256},
    },
}

async function doMachineReq() {
    console.log("Sending machine request", MachineReq);
    await window.go.main.App.MachineStateRequest(MachineReq);
    MachineReq.ClickedNumber = "";
    MachineReq.ClickedRange = "";
}
