
const elem = document.getElementById.bind(document);

const romInput = elem('romInput');
const romInputButton = elem('romInputButton');
const romInfo = elem('fileInfo');

romInputButton.addEventListener('click', () => romInput.click());

// Listen for file selection
romInput.addEventListener('change', handleFileUpload);
romInputButton.addEventListener('click', () => romInput.click());


function listROMs(roms) {
    if (roms == null || roms.length === 0) {
        fileInfo.innerHTML = '<p>No files uploaded yet</p>';
        return;
    }
    
    let html = '<h3>Uploaded ROMs:</h3>';
    for (const name in roms) {
        let rom = roms[name];
        html += `
            <div class="uploadedROM">
                <strong>${rom.Name}</strong> (${rom.Size/1024} kB)
                <div class="actions">
                    <button class="action-btn play-btn" data-name="${rom.Name}">▶️ Play</button>
                    <button class="action-btn delete-btn" data-name="${rom.Name}">🗑️ Delete</button>
                </div>
            </div>
        `;
    };
    romInfo.innerHTML = html;
}

romInfo.addEventListener('click', function(e) {
    let name = e.target.dataset.name;
    if (e.target.classList.contains("play-btn")) {
        config = getConfig();
        config.Model.ROM.Location = name;
    }
    if (e.target.classList.contains("delete-btn")) {
        roms = getROMs();
        delete roms[name];
        setROMs(roms);
    }
}, false);

// On load
listROMs(getROMs());
