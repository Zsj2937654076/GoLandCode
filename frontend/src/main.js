import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import './assets/styles/main.css'

// Create Vue application instance
const app = createApp(App)

// Use plugins
app.use(store)
app.use(router)
app.use(ElementPlus)

// Mount the app
app.mount('#app') 