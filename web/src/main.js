import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import Buefy from 'buefy'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueMoment from 'vue-moment'


import './registerServiceWorker'

Vue.use(VueAxios, axios)
Vue.use(Buefy)
Vue.use(VueMoment)
Vue.axios.defaults.baseURL = 'https://channel-api.otcxin.one/'
Vue.axios.defaults.headers.post['Content-Type'] = 'application/json'

Vue.config.productionTip = false

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')