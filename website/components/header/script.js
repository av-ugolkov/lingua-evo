import getBrowserFingerprint from '../../tools/get-browser-fingerprint.js';

let btnDictionary = document.getElementById("btnDictionary")
btnDictionary.addEventListener("click", () => {
    fetch("/account/dictionary", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    }).catch((error) => {
        console.error('error:', error);
    })
})

let btnAccount = document.getElementById("btnAccount")
btnAccount.addEventListener("click", () => {
    fetch("/account", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    }).catch((error) => {
        console.error('error:', error);
    })
})

let btnLogout = document.getElementById("btnLogout")
btnLogout.addEventListener("click", () => {
    let token = localStorage.getItem('access_token')
    if (token == null) {
        window.open("/", "_self")
        return
    }

    fetch("/auth/logout", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': "Bearer " + token,
            'Fingerprint': getBrowserFingerprint(),
        }
    }).then((response) => {
        if (response.status == 200) {
            localStorage.removeItem('access_token')

            window.open("/", "_self")
        }
    }).catch((error) => {
        console.error('error:', error);
    })
})