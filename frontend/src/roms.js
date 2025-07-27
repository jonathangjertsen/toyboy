const b64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

function b64Encode(bytes) {
    const chars = new Array(Math.ceil(bytes.length / 3) * 4);
    let idx = 0;

    // Convert next 3 octets at a time (to 4 sextets)
    let i = 0;
    for (i = 0; i < bytes.length - 2; i += 3) {
        // sextets[0] = octets[0][0:6]
        chars[idx] = b64[bytes[i] >> 2];
        // sextets[1] = octets[0][7:8], octets[1][0:4]
        chars[idx+1] = b64[((bytes[i] & 0x3) << 4) | (bytes[i + 1] >> 4)];
        // sextets[2] = octets[1][0:4], octets[2][0:2]
        chars[idx+2] = b64[((bytes[i + 1] & 0xf) << 2) | (bytes[i + 2] >> 6)];
        // sextets[3] = octets[2][3:8]
        chars[idx+3] = b64[bytes[i + 2] & 0x3f];
        
        idx += 4;
    }

    // Convert tail
    if (i === bytes.length - 1) {
        // sextets[0] = octets[0][0:6]
        chars[idx] = b64[bytes[i] >> 2];
        // sextets[1] = octets[0][7:8], 4 padding zeros
        chars[idx+1] = b64[(bytes[i] & 0x3) << 4];
        // sextets[2] = sextets[3] = padding
        chars[idx+2] = '=';
        chars[idx+3] = '=';
    } else if (i === bytes.length - 2) {
        // sextets[0] = octets[0][0:6]
        chars[idx] = b64[bytes[i] >> 2];
        // sextets[1] = octets[0][7:8], octets[1][0:4]
        chars[idx] = b64[((bytes[i] & 0x3) << 4) | (bytes[i + 1] >> 4)];
        // sextets[2] = octets[1][0:4], 2 padding zeros
        chars[idx] = b64[(bytes[i + 1] & 0xf) << 2];
        // sextets[3] = padding
        chars[idx] = '=';
    }
    return chars.join('');
}

// Load roms from local storage, and if something wrong is stored in there, clean it up
function getROMs() {
    return loadObject("roms", {});
}

function setROMs(roms) {
    storeObject("roms", roms);

    // Update UI
    listROMs(roms);
}

async function handleFileUpload(event) {
    let roms = getROMs();

    const files = [...event.target.files];
    for (const file of files) {
        const buffer = await file.arrayBuffer();
        const bytes = new Uint8Array(buffer);
        const romB64 = b64Encode(bytes)
        roms[file.name] = {
            Name: file.name,
            Size: bytes.length,
            Data: romB64,
            Time: Date.now()
        }
    }

    setROMs(roms);
}
