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
        setTimeout(function () {
            refreshToken(token);
        }, 1800000);
        return json.token;
    } else {
        localStorage.removeItem('token');
        console.log("invalid token");
        window.location.replace("/login");
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
        getDevices(token);
    });
    const logoutButton = document.getElementById("logout");
    logoutButton.addEventListener('click', async event => {
        localStorage.removeItem('token');
        window.location.replace("/login");
    });
}

let currentToken = localStorage.getItem('token');
if (currentToken == null) {
    window.location.replace("/login");
}

document.addEventListener('DOMContentLoaded', (event) => {
    //the event occurred
    init();
});

async function getDevices(token) {
    const options = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token
        },
    };
    let response = await fetch("/api/auth/employees", options);
    console.log(response.status);
    const main = document.getElementById("devices");
    while (main.firstChild) {
        main.removeChild(main.firstChild);
    }
    if (response.status === 200) {
        let json = await response.json();
        for (let item of json) {
            console.log(item);
            const a = document.createElement("a"); //clickable
            a.classList.add("no-color"); //no colors for this a tag
            a.href = `/public/devices?id=${item.device.DevID}`; // device info stuff
            const row = document.createElement('div'); //row
            row.classList.add("row", "mb-1");
            const col = document.createElement('div');// 12 column
            col.classList.add("col-12");
            const card = document.createElement('div'); //card
            card.classList.add("card", "bg-bgsecond", "border-left-primary", "shadow", "h-100", "py-2");
            const cardBody = document.createElement('div'); // card body
            cardBody.classList.add("card-body");

            const innerRow = document.createElement('div'); // card inner row
            innerRow.classList.add("row", "no-gutters", "align-items-center");

            const usernameColumn = document.createElement("div"); //1st col
            usernameColumn.classList.add("col-4", "col-md-2");
            const username = document.createElement('div');
            username.classList.add("h5", "font-weight-bold", "mb-1");
            let usernameText = "null";
            if (item.user.username !== "") {
                usernameText = item.user.username;
            }
            username.innerText = usernameText;

            const deviceIdCol = document.createElement("div"); //2nd col
            deviceIdCol.classList.add("col-4", "col-md-2");
            const deviceIDkey = document.createElement('div');
            deviceIDkey.classList.add("h5", "d-inline", "font-weight-bold");
            deviceIDkey.innerText = "ID:";
            const deviceIDVal = document.createElement('div');
            deviceIDVal.classList.add("h5", "d-inline", "font-weight-bold");
            deviceIDVal.innerText = item.device.DevID;

            const deviceVerCol = document.createElement("div"); //3rd col
            deviceVerCol.classList.add("col-4", "col-md-2");
            const deviceVerVal = document.createElement('div');
            deviceVerVal.classList.add("h5", "d-inline", "font-weight-bold");
            deviceVerVal.innerText = `v${item.device.sfwv}`;

            const tempCol = document.createElement('div'); //4th col
            tempCol.classList.add("col-4", "col-md-2");
            const temp = document.createElement('div');
            temp.classList.add("h5", "d-inline", "font-weight-bold");
            temp.innerText = item.device.dhtT + " ";
            const tempIcon = document.createElement('div');
            tempIcon.classList.add("fas", "fa-temperature-low", "fa-1x");
            let hsl1 = map(item.device.dhtT, -10, 50, 180, 0);
            tempIcon.style.color = `hsl(${hsl1},100%,50%)`;
            const humCol = document.createElement('div');

            humCol.classList.add("col-4", "col-md-2"); //5th col
            const hum = document.createElement('div');
            hum.classList.add("h5", "d-inline", "font-weight-bold");
            hum.innerText = item.device.dhtH + "% ";
            const humIcon = document.createElement('div');
            humIcon.classList.add("fas", "fa-water");
            
            const alarmCol = document.createElement('div'); //6th col
            alarmCol.classList.add("col-4", "col-md-2");
            const alarmIcon = document.createElement('div');
            const lostIcon = document.createElement('div');
            if (item.device.LostFlag === 1) {
                lostIcon.classList.add("far", "fa-compass", "fa-2x", "alert-icon-color");
            } else {
                lostIcon.classList.add("far", "fa-compass", "fa-2x", "text-gray-500");
            }
            if (item.device.alarm_state === 1) {
                alarmIcon.classList.add("fas", "fa-volume-up", "fa-2x", "alert-icon-color", "mr-2");
            } else {
                alarmIcon.classList.add("fas", "fa-volume-mute", "fa-2x", "text-gray-500", "mr-2" );
            }


            humCol.append(hum, humIcon);
            tempCol.append(temp, tempIcon);
            deviceIdCol.append(deviceIDkey, deviceIDVal);
            usernameColumn.append(username);
            alarmCol.append(alarmIcon);
            alarmCol.append(lostIcon);
            deviceVerCol.append(deviceVerVal);

            innerRow.append(usernameColumn, deviceVerCol, deviceIdCol, tempCol, humCol);
            innerRow.append(alarmCol);
            cardBody.append(innerRow);
            card.append(cardBody);
            col.append(card);
            row.append(col);
            a.append(row);
            main.append(a);
        }
    } else {
        let message = await response.json();
        alert(message.message);
        window.location.replace("/login");

    }
    setTimeout(function () {
        getDevices(localStorage.getItem('token'));
    }, 30000);
}

function map(x, in_min, in_max, out_min, out_max) {
    return (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min;
}