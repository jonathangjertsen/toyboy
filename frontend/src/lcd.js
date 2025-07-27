const Canvas = document.getElementById("lcd");
const Renderer = new WebGLFrameRenderer(Canvas);

async function drawFrame(event) {
    Renderer.renderFrame(new Uint8Array(event.data), 160, 144);
    FPS.innerHTML = await window.go.main.App.GetSpeedPct();
}

