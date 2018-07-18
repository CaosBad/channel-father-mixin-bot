var storageObj = {
    set: function (key, item) {
        localStorage.setItem(key, JSON.stringify(item))
    },
    get: function (key) {
        var value = localStorage.getItem(key) || '[]';
        value = JSON.parse(value);
        return value
    },
    clear: function (key) {
        localStorage.removeItem(key)
    }
}

export default storageObj