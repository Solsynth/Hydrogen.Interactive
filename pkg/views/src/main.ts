import "virtual:uno.css"

import { createApp } from "vue"
import { createPinia } from "pinia"

import "vuetify/styles"
import { createVuetify } from "vuetify"
import * as components from "vuetify/components"
import * as directives from "vuetify/directives"

import "@mdi/font/css/materialdesignicons.min.css"

import index from "./index.vue"
import router from "./router"

const app = createApp(index)

app.use(
  createVuetify({
    components,
    directives,
    theme: {
      themes: {
        light: {
          primary: "#4a5099",
          secondary: "#2196f3",
          accent: "#009688",
          error: "#f44336",
          warning: "#ff9800",
          info: "#03a9f4",
          success: "#4caf50"
        }
      }
    }
  })
)

app.use(createPinia())
app.use(router)

app.mount("#app")
