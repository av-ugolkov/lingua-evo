let btnDictionary = document.getElementById("btnDictionary")
btnDictionary.addEventListener("click", () => {
    fetch("/account/dictionary", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})

let btnAccount = document.getElementById("btnAccount")
btnAccount.addEventListener("click", () => {
    fetch("/account", {
        method: "get",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})

let btnLogout = document.getElementById("btnLogout")
btnLogout.addEventListener("click", () => {
    fetch("/logout", {
        method: "post",
    }).then((data) => {
        window.open(data["url"], "_self")
        console.log(data);
    })
})