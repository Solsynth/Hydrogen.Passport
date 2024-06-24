import "virtual:uno.css"

import "./assets/utils.css"

import { createApp } from "vue"
import { createPinia } from "pinia"

import "vuetify/styles"
import { createVuetify } from "vuetify"
import { md3 } from "vuetify/blueprints"
import * as components from "vuetify/components"
import * as labsComponents from "vuetify/labs/components"
import * as directives from "vuetify/directives"

import "@mdi/font/css/materialdesignicons.min.css"
import "@fontsource/roboto/latin.css"
import "@unocss/reset/tailwind.css"

import index from "./index.vue"
import router from "./router"

const app = createApp(index)

app.use(
  createVuetify({
    directives,
    components: {
      ...components,
      ...labsComponents,
    },
    blueprint: md3,
    theme: {
      defaultTheme: "original",
      themes: {
        original: {
          colors: {
            primary: "#4a5099",
            secondary: "#2196f3",
            accent: "#009688",
            error: "#f44336",
            warning: "#ff9800",
            info: "#03a9f4",
            success: "#4caf50",
          },
        },
      },
    },
  }),
)

app.use(createPinia())
app.use(router)

app.mount("#app")
