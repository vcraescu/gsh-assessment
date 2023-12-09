const {createApp, ref} = Vue

createApp({
    setup() {
        const quantity = ref(250)
        const rows = ref()
        const error = ref()

        const calculate = function () {
            error.value = ""
            rows.value = null

            const xmlHttp = new XMLHttpRequest();
            xmlHttp.open("POST", "/orders", false);

            xmlHttp.send(JSON.stringify({
                quantity: quantity.value,
            }));

            if (xmlHttp.status >= 500) {
                error.value = xmlHttp.statusText

                return
            }

            const resp = JSON.parse(xmlHttp.responseText)

            if (xmlHttp.status >= 400) {
                if (resp.error) {
                    error.value = resp.error
                } else {
                    error.value = "Error occurred"
                }

                return
            }

            rows.value = resp.data.rows

            return false
        }

        return {
            quantity,
            calculate,
            rows,
            error,
        }
    }
}).mount('#app')
