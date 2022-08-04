import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index'
import naive from './plugins/naive'
import store from './store'

createApp(App).use(router).use(naive).use(store).mount('#app')
