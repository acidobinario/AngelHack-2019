async function refreshToken(token) {
    const options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
        },
    };
    let response = await fetch("/api/v3/refresh_token", options);
    let json = await response.json();
    if (json.code === 200) {
        console.log("token refreshed!, new token is " + json.token);
        localStorage.setItem('token', json.token);
        setTimeout(function () {
            refreshToken(token);
        }, 1800000);
        return json.token;
    } else {
        localStorage.removeItem('token');
        console.log("invalid token");
        window.location.replace("/public/login");
        return null;
    }
}

function init() {
    let Thetoken;
    refreshToken(localStorage.getItem('token')).then((token) => {
        let dtoken = jwt_decode(token);
        console.log("decoded token: ");
        console.log(dtoken);
        document.getElementById("username").textContent = dtoken.identity;
        document.getElementById("userImage").src = `http://tinygraphs.com/isogrids/${dtoken.identity}?theme=seascape&numcolors=4&size=220&fmt=svg`;
    });
    const logoutButton = document.getElementById("logout");
    logoutButton.addEventListener('click', async event => {
        localStorage.removeItem('token');
        window.location.replace("/public/login");
    });
}

let currentToken = localStorage.getItem('token');
if (currentToken == null) {
    window.location.replace("/public/login");
}

document.addEventListener('DOMContentLoaded', (event) => {
    //the event occurred
    init();
});