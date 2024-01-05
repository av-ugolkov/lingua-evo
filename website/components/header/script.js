let bntSignIn = document.getElementById("btnSignUp")
bntSignIn.addEventListener("click", () => {
    fetch("/signup", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})

let bntLogin = document.getElementById("btnSignIn")
bntLogin.addEventListener("click", () => {
    fetch("/signin", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})