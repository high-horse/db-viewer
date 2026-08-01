import { createApp } from 'vue'
import App from './App.vue'
import { Quasar, Dark } from "quasar"

import "./style.css"
import "quasar/src/css/index.sass"
import "@quasar/extras/material-icons/material-icons.css"

const app = createApp(App)

app.use(Quasar, {
  plugins: {
    Dark
  }
})

const savedPreference = localStorage.getItem("dark-mode")
const isDark = savedPreference === null ? true : savedPreference === "true"

Dark.set(isDark)

// CRITICAL FOR TAILWIND: Sync the HTML element class on startup
document.documentElement.classList.toggle("dark", isDark)

app.mount("#app")
