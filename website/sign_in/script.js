import getBrowserFingerprint from '../tools/get-browser-fingerprint.js';

let authPanel = document.getElementById("sign_in_panel");

authPanel.addEventListener("submit", async (e) => {
    e.preventDefault()

    let username = document.getElementById("username")
    let password = document.getElementById("password")

    fetch("/auth/login", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': "Basic " + btoa(username.value + ':' + password.value),
            'Fingerprint': getBrowserFingerprint()
        }
    })
        .then((response) => response.json())
        .then((data) => {
            let token = data['access_token'];
            sessionStorage.setItem('access_token', token);

            window.open("/", "_self");
        })
        .catch((error) => {
            console.error('error:', error);
        });
})