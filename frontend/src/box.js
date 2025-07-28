// Global state management system
class Boxes {
    constructor() {
        this.states = {};
        this.callbacks = [];
    }

    initBox(id, isCollapsed) {
        this.states[id] = isCollapsed
        this.changed();
    }

    toggle(id) {
        return this.set(id, !this.get(id))
    }

    get(id) {
        return this.states[id];
    }

    set(id, open) {
        this.states[id] = open;
        this.changed();
        return this.states[id];
    }

    all() {
        return { ...this.states };
    }

    collapsed() {
        return Object.keys(this.states).filter(id => !this.states[id]);
    }

    expanded() {
        return Object.keys(this.states).filter(id => this.states[id]);
    }

    onStateChange(callback) {
        this.callbacks.push(callback);
    }

    changed() {
        this.callbacks.forEach(callback => callback(this.all()));
    }
}

const GlobalBoxes = new Boxes();

function initializeBox(boxElement) {
    const id = boxElement.getAttribute('data-box-id');
    const collapsed = boxElement.classList.contains('collapsed');

    GlobalBoxes.initBox(id, !collapsed);

    const header = boxElement.querySelector('.box-header');
    const button = boxElement.querySelector('.collapse-button');

    const toggleHandler = (e) => {
        e.preventDefault();
        if (GlobalBoxes.toggle(id)) {
            boxElement.classList.remove('collapsed');
        } else {
            boxElement.classList.add('collapsed');
        }
    };

    header.addEventListener('click', toggleHandler);
    button.addEventListener('click', (e) => {
        e.stopPropagation();
        toggleHandler(e);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.box').forEach(initializeBox);
    MachineReq.OpenBoxes = GlobalBoxes.all();
    doMachineReq();

    GlobalBoxes.onStateChange((states) => {
        MachineReq.OpenBoxes = states;
        doMachineReq();
    });
});


function loadBoxes(backendStates) {
    Object.keys(backendStates).forEach(boxId => {
        const boxElement = document.querySelector(`[data-box-id="${boxId}"]`);
        if (boxElement) {
            GlobalBoxes.set(boxId, backendStates[boxId]);
            if (backendStates[boxId]) {
                boxElement.classList.remove('collapsed');
            } else {
                boxElement.classList.add('collapsed');
            }
        }
    });
}
