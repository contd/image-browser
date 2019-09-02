import Vue from 'vue'
import axios from 'axios'
import BootstrapVue from 'bootstrap-vue'
import App from './App.vue'
import store from './store'
import '@/assets/customized.scss'

Vue.use(BootstrapVue)
Vue.prototype.$BASEURL = 'http://localhost:6969/api'
Vue.prototype.$axios = axios
Vue.config.productionTip = false

new Vue({
  store,
  render: h => h(App)
}).$mount("#app")
