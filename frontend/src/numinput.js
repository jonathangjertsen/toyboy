function numberUpdated(id, value) {
    MachineReq.Numbers[id] = value;
    MachineReq.ClickedNumber = id;
    console.log("Set ClickedNumber", id, MachineReq)
    doMachineReq();
}

function rangeUpdated(id, lower, upper) {
    MachineReq.Ranges[id] = { Begin: lower, End: upper }
    MachineReq.ClickedRange = id;
    console.log("Set ClickedRange", id, MachineReq)
    doMachineReq();
}

class InputFields {
    constructor() {
        this.inputs = {};
    }

    initInput(inputElement) {
        const id = inputElement.getAttribute('data-input-id');
        const inputType = inputElement.getAttribute('data-input-type') ?? 'int';
        const min = parseInt(inputElement.getAttribute('data-min')) || -Infinity;
        const max = parseInt(inputElement.getAttribute('data-max')) || Infinity;
        const field = inputElement.querySelector('.numeric-field');
        const hexToggle = inputElement.querySelector('.hex-toggle');
        const confirmBtn = inputElement.querySelector('.confirm-btn');

        let fmt = '';
        if (inputType === 'float') {
            fmt = 'FLOAT';
        } else if ((hexToggle !== null) && (hexToggle.textContent === 'HEX')) {
            fmt = 'HEX';
        } else {
            fmt = 'DEC';
        }
        const state = {
            id: id,
            min: min,
            max: max,
            fmt: fmt,
            lastConfirmedValue: 0,
            element: inputElement,
            field: field,
            hexToggle: hexToggle
        };
        state.lastConfirmedValue = this.getCurrentValue(state);

        this.inputs[id] = state;

        // Toggle hex/dec
        hexToggle?.addEventListener('click', () => {
            if (state.fmt == 'HEX') {
                state.fmt = 'DEC';
            } else {
                state.fmt = 'HEX';
            }
            hexToggle.textContent = state.fmt;
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
            
            switch (state.fmt) {
                case 'HEX':
                    return parseInt(text, 16);
                case 'DEC':
                    return parseInt(text, 10);
                case 'FLOAT':
                    return parseFloat(text);
            }
        } catch (e) {
            return null;
        }
    }

    isInRange(value, state) {
        return value >= state.min && value <= state.max;
    }

    setHexButtonColor(hexToggle, state) {
        if (hexToggle === null) {
            return;
        }
        if (state.fmt === 'HEX') {
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
            if (state.fmt === 'HEX') {
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
        const inputType = rangeElement.getAttribute('data-input-type');
        const min = parseInt(rangeElement.getAttribute('data-min')) || -Infinity;
        const max = parseInt(rangeElement.getAttribute('data-max')) || Infinity;
        const lowerField = rangeElement.querySelector('.lower-field');
        const upperField = rangeElement.querySelector('.upper-field');
        const hexToggle = rangeElement.querySelector('.hex-toggle');
        const confirmBtn = rangeElement.querySelector('.confirm-btn');

        let fmt = '';
        if (inputType === 'float') {
            fmt = 'FLOAT';
        } else if ((hexToggle !== null) && (hexToggle.textContent === 'HEX')) {
            fmt = 'HEX';
        } else {
            fmt = 'DEC';
        }
        const state = {
            id: id,
            min: min,
            max: max,
            fmt: fmt,
            lastConfirmedLower: this.getCurrentValue(lowerField.value, fmt),
            lastConfirmedUpper: this.getCurrentValue(upperField.value, fmt),
            element: rangeElement,
            lowerField: lowerField,
            upperField: upperField,
            hexToggle: hexToggle
        };

        this.ranges[id] = state;

        // Toggle hex/dec
        hexToggle?.addEventListener('click', () => {
            if (state.fmt == 'HEX') {
                state.fmt = 'DEC';
            } else {
                state.fmt = 'HEX';
            }
            hexToggle.textContent = state.fmt
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
            const lower = this.getCurrentValue(state.lowerField, state.fmt);
            const upper = this.getCurrentValue(state.upperField, state.fmt);
            
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
        if (hexToggle === null) {
            return;
        }
        if (state.fmt == 'HEX') {
            hexToggle.classList.add('hex');
            hexToggle.classList.remove('dec');
        } else {
            hexToggle.classList.add('dec');
            hexToggle.classList.remove('hex');
        }
    }

    getCurrentValue(field, fmt) {
        try {
            const text = field.value.trim();
            if (text === '') return null;
            
            switch (fmt) {
                case 'HEX': {
                    return parseInt(text, 16);
                }
                case 'DEC': {
                    return parseInt(text, 10);
                }
                case 'FLOAT': {
                    return parseFloat(text);
                }
            }
        } catch (e) {
            return null;
        }
    }

    isInRange(value, state) {
        return value >= state.min && value <= state.max;
    }

    updateFieldsDisplay(state) {
        const lowerValue = this.getCurrentValue(state.lowerField, state.fmt);
        const upperValue = this.getCurrentValue(state.upperField, state.fmt);
        
        if (lowerValue !== null) {
            if (state.fmt === 'HEX') {
                state.lowerField.value = lowerValue.toString(16).toUpperCase();
            } else {
                state.lowerField.value = lowerValue.toString(10);
            }
        }
        
        if (upperValue !== null) {
            if (state.fmt === 'HEX') {
                state.upperField.value = upperValue.toString(16).toUpperCase();
            } else {
                state.upperField.value = upperValue.toString(10);
            }
        }
        
        this.validateRange(state);
    }

    validateRange(state) {
        const lowerValue = this.getCurrentValue(state.lowerField, state.fmt);
        const upperValue = this.getCurrentValue(state.upperField, state.fmt);
        
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
