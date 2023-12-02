import getBrowserFingerprint from '../tools/get-browser-fingerprint.js';


window.onload = function () {
    let token = sessionStorage.getItem('access_token')
    if (token == null) {
        return
    }

    fetch("/get-account-data", {
        method: "get",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Fingerprint': getBrowserFingerprint(),
            'Access-Token': token
        },

    }).then((response) => response.json())
        .then((data) => {
            document.getElementById("account_panel").innerHTML = data
        }).catch((error) => {
            console.error(error)
        })
}


let lableRandom = document.getElementById("lable-random")
let interval = setInterval(function () {
    fetch("/word/get_random", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ language_code: 'en' })
    })
        .then((response) => response.json())
        .then((data) => {
            lableRandom.textContent = data["text"]
        })
        .catch((error) => {
            lableRandom.textContent = error
            stopInterval()
        });
}, 60000);


function stopInterval() {
    clearInterval(interval)
}

let bntSignIn = document.getElementById("btnSignup")
bntSignIn.addEventListener("click", () => {
    fetch("/signup", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})

let bntLogin = document.getElementById("btnLogin")
bntLogin.addEventListener("click", () => {
    fetch("/login", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})