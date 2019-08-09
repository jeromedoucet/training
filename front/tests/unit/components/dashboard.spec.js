import { mount, createLocalVue } from '@vue/test-utils'

import Buefy from 'buefy'
import Vuex from 'vuex'
import VueRouter from 'vue-router'
import store from '@/store'
import Dashboard from '@/components/Dashboard.vue'
import fetchMock from 'fetch-mock'

const localVue = createLocalVue()
localVue.use(Buefy)
localVue.use(VueRouter)
localVue.use(Vuex)

const router = new VueRouter({
  routes: []
})

const createDashboard = propsData => mount(Dashboard, { sync: false, propsData, store, localVue, router })

describe('dashboard', () => {
  beforeEach(async() => {
  })
});
