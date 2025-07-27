async function waitFrame() {
    await new Promise(requestAnimationFrame);
}
const FPS =  document.getElementById("fps")
const Keys = ['w', 'a', 's', 'd', 'k', 'l', 'n', 'm'];
const KeyState = Object.fromEntries(Keys.map(k => [k, false]));
Frame = null;

async function run() {
    config = await window.go.main.App.GetConfig();
    console.log(config);

    await window.go.main.App.StartGB()
    const ws = new WebSocket('ws://localhost:8081/vp');
    ws.binaryType = 'arraybuffer';

    ws.onmessage = function(event) {
        Frame = new Uint8Array(event.data);
    };
}

run();

async function updateSpeed() {
    FPS.innerHTML = await window.go.main.App.GetSpeedPct();
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
    }
    await window.go.main.App.SetKeyState(KeyState);
}

document.addEventListener('keydown', (e) => {
    setKeyState(e.key.toLowerCase(), true);
});

document.addEventListener('keyup', (e) => {
    setKeyState(e.key.toLowerCase(), false);
});
