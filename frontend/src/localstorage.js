function loadObject(key, defaultValue) {
    let obj = defaultValue;
    const storedStr = localStorage.getItem(key);
    if (storedStr !== null) {
        try {
            let stored = JSON.parse(storedStr)
            const isProbablyObject = typeof stored == 'object' && !Array.isArray(stored) && stored !== null;
            if (isProbablyObject) {
                obj = stored;
            } else {
                localStorage.removeItem(key);
            }
        } catch (_) {
            localStorage.removeItem(key);
        }
    }
    return obj;
}

function storeObject(key, value) {
    localStorage.setItem(key, JSON.stringify(value));
}