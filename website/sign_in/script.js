let authPanel = document.getElementById("sign_in_panel");

authPanel.addEventListener("submit", (e) => {
    e.preventDefault()

    let username = document.getElementById("username")
    let password = document.getElementById("password")
    fetch("/auth/login", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': "Basic " + btoa(username.value + ':' + password.value)
        }
    })
        .then((response) => response.json())
        .then((data) => {
            window.open("/", "_self")
        })
        .catch((error) => {
            console.error('Error:', error);
        });
})