function numberUpdated(id, value) {
    MachineReq.Numbers[id] = value;
    doMachineReq();
}

function rangeUpdated(id, lower, upper) {
    MachineReq.Ranges[id] = { Begin: lower, End: upper }
    doMachineReq();
}

class InputFields {
    constructor() {
        this.inputs = {};
    }

    initInput(inputElement) {
        const id = inputElement.getAttribute('data-input-id');
        const min = parseInt(inputElement.getAttribute('data-min')) || -Infinity;
        const max = parseInt(inputElement.getAttribute('data-max')) || Infinity;
        const field = inputElement.querySelector('.numeric-field');
        const hexToggle = inputElement.querySelector('.hex-toggle');
        const confirmBtn = inputElement.querySelector('.confirm-btn');

        const initHex = hexToggle.textContent === 'HEX'
        const state = {
            id: id,
            min: min,
            max: max,
            isHex: initHex,
            lastConfirmedValue: parseInt(field.value, initHex ? 16 : 10),
            element: inputElement,
            field: field,
            hexToggle: hexToggle
        };

        this.inputs[id] = state;

        // Toggle hex/dec
        hexToggle.addEventListener('click', () => {
            state.isHex = !state.isHex;
            hexToggle.textContent = state.isHex ? 'HEX' : 'DEC';
            this.setHexButtonColor(hexToggle, state);
            this.updateFieldDisplay(state);
        });

        // Validate on input
        field.addEventListener('input', () => {
            this.validateInput(state);
        });

        // Confirm button
        confirmBtn.addEventListener('click', () => {
            const value = this.getCurrentValue(state);
            if (value !== null && this.isInRange(value, state)) {
                state.lastConfirmedValue = value;
                numberUpdated(id, value);
                this.validateInput(state);
            }
        });

        // Initial validation
        this.validateInput(state);
        this.setHexButtonColor(hexToggle, state);
    }

    getCurrentValue(state) {
        try {
            const text = state.field.value.trim();
            if (text === '') return null;
            
            if (state.isHex) {
                return parseInt(text, 16);
            } else {
                return parseInt(text, 10);
            }
        } catch (e) {
            return null;
        }
    }

    isInRange(value, state) {
        return value >= state.min && value <= state.max;
    }

    setHexButtonColor(hexToggle, state) {
        if (state.isHex) {
            hexToggle.classList.add('hex');
            hexToggle.classList.remove('dec');
        } else {
            hexToggle.classList.add('dec');
            hexToggle.classList.remove('hex');
        }
    }
    updateFieldDisplay(state) {
        const currentValue = this.getCurrentValue(state);
        if (currentValue !== null) {
            if (state.isHex) {
                state.field.value = currentValue.toString(16).toUpperCase();
            } else {
                state.field.value = currentValue.toString(10);
            }
        }
        this.validateInput(state);
    }

    validateInput(state) {
        const currentValue = this.getCurrentValue(state);
        
        state.field.classList.remove('changed', 'invalid');
        
        if (state.field.value.trim() == '') {
            state.field.classList.add('empty');
        } else {
            state.field.classList.remove('empty');
        }
        if (currentValue === null || !this.isInRange(currentValue, state)) {
            state.field.classList.add('invalid');
        } else if (currentValue !== state.lastConfirmedValue) {
            state.field.classList.add('changed');
        }
    }
}

// Range Input Component Class
class RangeFields {
    constructor() {
        this.ranges = {};
    }

    initRange(rangeElement) {
        const id = rangeElement.getAttribute('data-range-id');
        const min = parseInt(rangeElement.getAttribute('data-min')) || -Infinity;
        const max = parseInt(rangeElement.getAttribute('data-max')) || Infinity;
        const lowerField = rangeElement.querySelector('.lower-field');
        const upperField = rangeElement.querySelector('.upper-field');
        const hexToggle = rangeElement.querySelector('.hex-toggle');
        const confirmBtn = rangeElement.querySelector('.confirm-btn');

        const initHex = hexToggle.textContent === 'HEX';
        const state = {
            id: id,
            min: min,
            max: max,
            isHex: initHex,
            lastConfirmedLower: parseInt(lowerField.value, initHex ? 16 : 10),
            lastConfirmedUpper: parseInt(upperField.value, initHex ? 16 : 10),
            element: rangeElement,
            lowerField: lowerField,
            upperField: upperField,
            hexToggle: hexToggle
        };

        this.ranges[id] = state;

        // Toggle hex/dec
        hexToggle.addEventListener('click', () => {
            state.isHex = !state.isHex;
            hexToggle.textContent = state.isHex ? 'HEX' : 'DEC';
            this.setHexButtonColor(hexToggle, state);
            this.updateFieldsDisplay(state);
        });

        // Validate on input
        lowerField.addEventListener('input', () => {
            this.validateRange(state);
        });
        upperField.addEventListener('input', () => {
            this.validateRange(state);
        });

        // Confirm button
        confirmBtn.addEventListener('click', () => {
            const lower = this.getCurrentValue(state.lowerField, state.isHex);
            const upper = this.getCurrentValue(state.upperField, state.isHex);
            
            if (lower !== null && upper !== null && 
                this.isInRange(lower, state) && this.isInRange(upper, state) &&
                lower <= upper) {
                state.lastConfirmedLower = lower;
                state.lastConfirmedUpper = upper;
                rangeUpdated(id, lower, upper);
                this.validateRange(state);
            }
        });

        // Initial validation
        this.validateRange(state);
        this.setHexButtonColor(hexToggle, state);
    }

    setHexButtonColor(hexToggle, state) {
        if (state.isHex) {
            hexToggle.classList.add('hex');
            hexToggle.classList.remove('dec');
        } else {
            hexToggle.classList.add('dec');
            hexToggle.classList.remove('hex');
        }
    }

    getCurrentValue(field, isHex) {
        try {
            const text = field.value.trim();
            if (text === '') return null;
            
            if (isHex) {
                return parseInt(text, 16);
            } else {
                return parseInt(text, 10);
            }
        } catch (e) {
            return null;
        }
    }

    isInRange(value, state) {
        return value >= state.min && value <= state.max;
    }

    updateFieldsDisplay(state) {
        const lowerValue = this.getCurrentValue(state.lowerField, !state.isHex);
        const upperValue = this.getCurrentValue(state.upperField, !state.isHex);
        
        if (lowerValue !== null) {
            if (state.isHex) {
                state.lowerField.value = lowerValue.toString(16).toUpperCase();
            } else {
                state.lowerField.value = lowerValue.toString(10);
            }
        }
        
        if (upperValue !== null) {
            if (state.isHex) {
                state.upperField.value = upperValue.toString(16).toUpperCase();
            } else {
                state.upperField.value = upperValue.toString(10);
            }
        }
        
        this.validateRange(state);
    }

    validateRange(state) {
        const lowerValue = this.getCurrentValue(state.lowerField, state.isHex);
        const upperValue = this.getCurrentValue(state.upperField, state.isHex);
        
        // Clear all classes
        state.lowerField.classList.remove('changed', 'invalid');
        state.upperField.classList.remove('changed', 'invalid');
        
        // Validate lower field
        if (lowerValue === null || !this.isInRange(lowerValue, state)) {
            state.lowerField.classList.add('invalid');
        } else if (lowerValue !== state.lastConfirmedLower) {
            state.lowerField.classList.add('changed');
        }
        
        // Validate upper field
        if (upperValue === null || !this.isInRange(upperValue, state)) {
            state.upperField.classList.add('invalid');
        } else if (upperValue !== state.lastConfirmedUpper) {
            state.upperField.classList.add('changed');
        }
        
        // Check range validity (lower <= upper)
        if (lowerValue !== null && upperValue !== null && lowerValue > upperValue) {
            state.lowerField.classList.add('invalid');
            state.upperField.classList.add('invalid');
        }
    }
}

const GlobalInputFields = new InputFields();
const GlobalRangeFields = new RangeFields();

document.addEventListener('DOMContentLoaded', () => {
    const numericInputs = document.querySelectorAll('.numeric-input');
    numericInputs.forEach(input => GlobalInputFields.initInput(input));

    const rangeInputs = document.querySelectorAll('.range-input');
    console.log("rangeInputs", rangeInputs);
    rangeInputs.forEach(range => GlobalRangeFields.initRange(range));
});
