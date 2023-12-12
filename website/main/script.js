import getBrowserFingerprint from '../tools/get-browser-fingerprint.js';

window.onload = async function () {
    let token = localStorage.getItem('access_token')
    if (token == null) {
        return
    }

    let payload = JSON.parse(atob(token.split(".")[1]));
    let exp = payload["exp"]
    if (Date.now() > exp * 1000) {
        await fetch("/auth/refresh", {
            method: "post",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Fingerprint': getBrowserFingerprint(),
            }
        })
            .then(response => {
                return response.json()
            })
            .then(data => {
                token = data['access_token'];
                localStorage.setItem('access_token', token);
            })
            .catch(error => {
                localStorage.removeItem('access_token')
                console.error('error:', error);
            })
    }

    await fetch("/get-account-panel", {
        method: "get",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Fingerprint': getBrowserFingerprint(),
            'Access-Token': token
        },
    })
        .then(async (response) => {
            if (response.status == 200) {
                document.getElementById("right-side").innerHTML = await response.text()
            }
        })
        .catch((error) => {
            console.error(error)
        })
}


let lableRandom = document.getElementById("random-field")
let interval = setInterval(function () {
    fetch("/word/get_random", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ language_code: 'en' })
    })
        .then((response) => {
            return response.json()
        })
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


// Close the dropdown if the user clicks outside of it
window.onclick = function (event) {
    if (!event.target.matches('.accountBtn')) {
        var dropdowns = document.getElementsByClassName("dropdown-content");
        var i;
        for (i = 0; i < dropdowns.length; i++) {
            var openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
}