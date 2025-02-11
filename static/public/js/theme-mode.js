const DATA_THEME_KEY = 'wis-data-theme';

try {
    let colorScheme = '';
    const mode = localStorage.getItem(DATA_THEME_KEY) || 'system';
    if (mode === 'system') {
        const mql = window.matchMedia('(prefers-color-scheme: dark)');
        if (mql.matches) {
            colorScheme = 'dark';
        } else {
            colorScheme = 'light';
        }
    }
    if (mode === 'light') {
        colorScheme = 'light';
    }
    if (mode === 'dark') {
        colorScheme = 'dark';
    }

    if (colorScheme) {
        document.documentElement.setAttribute('data-theme', colorScheme);
        if (mode !== 'system') {
            localStorage.setItem(DATA_THEME_KEY, colorScheme);
        }

        var toggleBtn = document.getElementById('data-theme-mode-btn');
        toggleBtn.innerText = colorScheme === 'dark' ? 'light_mode' : 'dark_mode';
    }
} catch (e) {
    console.error('Error getting initial color scheme', e);
}

function toggleDataThemeMode(btn_id) {
    var curr = document.documentElement.getAttribute('data-theme') || localStorage.getItem(DATA_THEME_KEY);
    var next = curr === 'dark' ? 'light' : 'dark';

    localStorage.setItem(DATA_THEME_KEY, next);
    document.documentElement.setAttribute('data-theme', next);

    var btn = document.getElementById(btn_id || 'data-theme-mode-btn');
    btn.innerText = next === 'dark' ? 'light_mode' : 'dark_mode';
}