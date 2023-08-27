let intervalCheck;
export const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('token_exp');
};

export const login = (data) => {
    localStorage.setItem('token', data.token);
    localStorage.setItem('token_exp', data.tokenExp || Date.now() + 6000000);
}

export const checkLoginState = () => {
    const expiry = localStorage.getItem('token_exp');
    const token = localStorage.getItem('token');
    if (!expiry || !token) {
        logout();
        return false;
    }
    const expiryNum = parseInt(expiry) * 1000;
    if (Date.now() < new Date(expiryNum)) {
        if (intervalCheck) {
            clearInterval(intervalCheck);
        }
        intervalCheck = setInterval(() => {
            logout();
            fireLogout();
            clearInterval(intervalCheck);
        }, expiryNum - Date.now());
        return true;
    }
    logout();
}

export const fireLogout = () => {
    let event = new Event('userLogout');
    dispatchEvent(event);
}

export const checkAndFireLogout = (status) => {
    if (status === 429) {
        fireLogout();
    }
}