
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