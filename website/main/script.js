window.onload = function () {
    let token = sessionStorage.getItem('access_token')
    if (token == null) {
        return
    }
    fetch("/get-account-data?access_token=" + sessionStorage.getItem('access_token'), {
        method: "get",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },

    }).then((response) => response.json())
        .then((data) => {
            let isLogin = data["IsLogin"]

            if (isLogin) {
                document.getElementById("dictionary_btn").innerHTML = `<div class="border"><button id=\"btnDictionary" type="button" class="accountBtn">Dictionary</button></div>`;
                document.getElementById("account_panel").innerHTML = `<button id="btnAccount" type="button" class="accountBtn" value="account">` + data["Name"] + `</button>`
            } else {
                document.getElementById("account_panel").innerHTML = `<div class="border">
            <button id="btnSignup" type="button" class="accountBtn" value="signup">Sign Up</button>
            |
            <button id="btnLogin" type="button" class="accountBtn" value="login">Login</button>
          </div>`
            }
        }).catch((error) => {
            console.log(error)
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