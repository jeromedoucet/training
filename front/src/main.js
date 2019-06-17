import Vue from 'vue'
import App from './App.vue'
import router from './router'
import Buefy from 'buefy'
import 'buefy/dist/buefy.css'
import 'bulmaswatch/darkly/bulmaswatch.min.css'
import '@fortawesome/fontawesome-free/css/all.min.css'

Vue.config.productionTip = false

Vue.use(Buefy, {
  defaultIconPack: 'fas'
})

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')
