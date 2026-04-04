import { createApp } from 'vue'
import App from './App.vue'
import router from './router.js'

import './assets/css/variables.css'
import './assets/css/base.css'
import './assets/css/layout.css'
import './assets/css/components.css'

const app = createApp(App)
app.use(router)
app.mount('#app')
