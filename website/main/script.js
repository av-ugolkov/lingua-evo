console.log('index')

let lableRandom = document.getElementById("lable-random")
setInterval(function () {
    fetch("/word/get_random", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ language_code: 'en-GB' })
    })
        .then((response) => response.json())
        .then((data) => {
            console.log('Success:', data["text"]);
            lableRandom.textContent = data["text"]
        })
        .catch((error) => {
            console.error('Error:', error);
        });
}, 2000);
