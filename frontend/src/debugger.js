const RunBtn = document.getElementById("run-btn");
const PauseBtn = document.getElementById("pause-btn");
const StepBtn = document.getElementById("step-btn");
const LoadBtn = document.getElementById("load-btn");
const SaveBtn = document.getElementById("save-btn");
const ExecLogBtn = document.getElementById("execlog-btn");

RunBtn.addEventListener('click', () => {
    runBtn()
})
PauseBtn.addEventListener('click', () => {
    pauseBtn()
})
StepBtn.addEventListener('click', () => {
    stepBtn()
})
SaveBtn.addEventListener('click', () => {
    saveBtn()
})
LoadBtn.addEventListener('click', () => {
    loadBtn()
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

async function loadBtn() {
    await window.go.main.App.Load();
}

async function saveBtn() {
    await window.go.main.App.Save();
}
