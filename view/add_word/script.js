console.log("Script Add Word")

let addWord = document.getElementById("add_word");

addWord.addEventListener("submit", (e)=>{
    e.preventDefault()

    let orig_word=document.getElementById("orig_word")
    let orig_lang=document.getElementById("orig_lang")
    let tran_word=document.getElementById("tran_word")
    let tran_lang=document.getElementById("tran_lang")
    fetch("add_word",{
        method:"post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({orig_word: orig_word.value, orig_lang:orig_lang.value, tran_word: tran_word.value, tran_lang: tran_lang.value})
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