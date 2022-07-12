import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementUI from 'element-ui';
import axios from 'axios';
import CodeEditor from 'bin-code-editor';
import Contextmenu from "vue-contextmenujs"

Vue.config.productionTip = false
Vue.use(ElementUI)
Vue.use(CodeEditor)
Vue.use(Contextmenu)
axios.defaults.baseURL = "http://127.0.0.1:10580"
Vue.prototype.$axios = axios

new Vue({
    router,
    store,
    render: function(h) { return h(App) }
}).$mount('#app')