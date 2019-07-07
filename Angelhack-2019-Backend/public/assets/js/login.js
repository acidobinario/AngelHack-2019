async function refreshToken(token) {
    const options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
        },
    };
    let response = await fetch("/api/refresh_token", options);
    let json = await response.json();
    if (json.code === 200) {
        console.log("token refreshed!, new token is " + json.token);
        localStorage.setItem('token', json.token);
        window.location.replace("/index");
        return json.token;
    } else {
        localStorage.removeItem('token');
        console.log("invalid token");
        return null;
    }
}


function init() {
    let token = localStorage.getItem('token');
    if (token != null) {
        console.log('token is not null');
        //TODO: check if the token is valid, if it is not valid , continue
        let x = refreshToken(token).then((newToken) => {
            if (newToken !== null) {
                console.log("setting new token -> " + newToken);
                console.log("redirecting to index...");
            }
        });
        console.log(x);

    }



    const button = document.getElementById("submit");

    button.addEventListener('click', async event => {
        let login = document.getElementById("inputUsername").value;
        let password = document.getElementById("inputPassword").value;

        const data = {"username": login, "password": password};
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        };
        const response = await fetch('/api/login', options);
        const json = await response.json();
        if (json.code === 200) {
            console.log("logged in!!, JWT -> " + json.token);
            //TODO: Save the token in local storage and redirect to index webpage
            localStorage.setItem('token', json.token);
            console.log("redirecting to index...");
            window.location.replace("/index")
        } else {
            alert(json.message);
        }
        console.log(json);
    });
    // Execute a function when the user releases a key on the keyboard
    document.getElementById("inputPassword").addEventListener("keyup", function (event) {
        // Number 13 is the "Enter" key on the keyboard
        if (event.keyCode === 13) {
            // Cancel the default action, if needed
            event.preventDefault();
            // Trigger the button element with a click
            document.getElementById("submit").click();
        }
    });
}

init();
