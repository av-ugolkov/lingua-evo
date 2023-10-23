
let lableRandom = document.getElementById("lable-random")
let interval = setInterval(function () {
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
            lableRandom.textContent = data["text"]
        })
        .catch((error) => {
            lableRandom.textContent = error
            stopInterval()
        });
}, 2000);


function stopInterval() {
    clearInterval(interval)
}
