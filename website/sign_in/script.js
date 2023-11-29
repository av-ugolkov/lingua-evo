import getBrowserFingerprint from './get-browser-fingerprint.js';

let authPanel = document.getElementById("sign_in_panel");

authPanel.addEventListener("submit", async (e) => {
    e.preventDefault()

    let username = document.getElementById("username")
    let password = document.getElementById("password")

    fetch("/auth/login", {
        method: "post",
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': "Basic " + btoa(username.value + ':' + password.value),
            'Fingerprint': getBrowserFingerprint()
        }
    })
        .then((response) => response.json())
        .then((data) => {
            let token = data['access_token'];
            sessionStorage.setItem('access_token', token);

            window.open("/", "_self");
        })
        .catch((error) => {
            console.error('error:', error);
        });
})

function _getFingerprint() {
    return new Promise((resolve, reject) => {
        async function getHash() {
            const options = {
                excludes: {
                    plugins: true,
                    localStorage: true,
                    adBlock: true,
                    screenResolution: true,
                    availableScreenResolution: true,
                    enumerateDevices: true,
                    pixelRatio: true,
                    doNotTrack: true,
                    preprocessor: (key, value) => {
                        if (key === 'userAgent') {
                            const parser = new UAParser(value)
                            // return customized user agent (without browser version)
                            return `${parser.getOS().name} :: ${parser.getBrowser().name} :: ${parser.getEngine().name}`
                        }
                        return value
                    }
                }
            }

            try {
                const components = await Fingerprint2.getPromise(options)
                const values = components.map(component => component.value)
                console.log('fingerprint hash components', components)

                return String(Fingerprint2.x64hash128(values.join(''), 31))
            } catch (e) {
                reject(e)
            }
        }

        if (window.requestIdleCallback) {
            console.log('get fp hash @ requestIdleCallback')
            requestIdleCallback(async () => resolve(await getHash()))
        } else {
            console.log('get fp hash @ setTimeout')
            setTimeout(async () => resolve(await getHash()), 500)
        }
    })
}