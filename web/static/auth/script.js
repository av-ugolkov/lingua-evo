console.log("Script Auth")

let authPanel = document.getElementById("authPanel");

authPanel.addEventListener("submit", (e)=>{
    e.preventDefault()

    let username=document.getElementById("username")
    let password=document.getElementById("password")
    fetch("auth",{
        method:"post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({username: username.value, password: password.value})
    })
        .then((response) => response.json())
        .then((data) => {
            console.log('Success:', data);
            window.open(data["url"],"_self")
        })
        .catch((error) => {
            console.error('Error:', error);
        });
})