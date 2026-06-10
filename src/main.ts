import { createApp } from "vue"
import { createPinia } from "pinia"
import vuetify from "@/plugins/vuetify"
import "@mdi/font/css/materialdesignicons.min.css"
import App from "./App.vue"

const app = createApp(App)
app.use(createPinia()).use(vuetify).mount("#app")
