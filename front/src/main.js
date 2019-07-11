import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import VeeValidate from 'vee-validate'
import Buefy from 'buefy'
import '@fortawesome/fontawesome-free/css/all.min.css'

Vue.config.productionTip = false

Vue.use(VeeValidate)
Vue.use(Buefy, {
  defaultIconPack: 'fas'
})

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
