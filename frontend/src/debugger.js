const RunBtn = document.getElementById("run-btn");
const PauseBtn = document.getElementById("pause-btn");
const StepBtn = document.getElementById("step-btn");

RunBtn.addEventListener('click', () => {
    runBtn()
})
PauseBtn.addEventListener('click', () => {
    pauseBtn()
})
StepBtn.addEventListener('click', () => {
    stepBtn()
})

async function runBtn() {
    await window.go.main.App.Start();
}

async function pauseBtn() {
    await window.go.main.App.Pause();
}

async function stepBtn() {
    await window.go.main.App.Step();
}
