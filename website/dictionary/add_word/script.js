console.log("Script Add Word")

let addWord = document.getElementById("add_word");

addWord.addEventListener("submit", (e) => {
    e.preventDefault()

    let native_word = document.getElementById("native_word")
    let native_lang = document.getElementById("native_lang")
    let tran_word = document.getElementById("tran_word")
    let tran_lang = document.getElementById("tran_lang")
    let example = document.getElementById("example")
    let pronunciation = document.getElementById("pronunciation")

    fetch("add_word", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            native_word: native_word.value,
            native_lang: native_lang.value,
            tran_word: tran_word.value,
            tran_lang: tran_lang.value,
            example: example.value,
            pronunciation: pronunciation.value
        })
    })
        .then((response) => {
            response.json();
            console.log('Success response');
        })
        .then((data) => {
            console.log('Success:', data);
            window.open(data["url"], "_self")
        })
        .catch((error) => {
            console.error('Error:', error);
        });
})