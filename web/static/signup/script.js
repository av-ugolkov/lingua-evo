console.log("Script Auth")

let authPanel = document.getElementById("signupPanel");

authPanel.addEventListener("submit", (e) => {
    e.preventDefault()

    let responseStatus = 404
    let email = document.getElementById("email")
    let username = document.getElementById("username")
    let password = document.getElementById("password")
    fetch("signup", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email: email.value, username: username.value, password: password.value })
    })
        .then((response) => response.json())
        .then((data) => {
            if (responseStatus == 201) {
                console.log('Success:', data)
                window.open(data["url"], "_self")
            }
        })
        .catch((err) => {
            console.error('Error:', err)
        });
})