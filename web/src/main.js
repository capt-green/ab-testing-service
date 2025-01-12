import { createApp } from 'vue'
import App from './App.vue'
import routes from './router'
import { createPinia } from 'pinia'
import './style.css'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(routes)
app.mount('#app')
