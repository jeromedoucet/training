import { mount, createLocalVue } from '@vue/test-utils'

import Buefy from 'buefy'
import VueRouter from 'vue-router'
import App from '@/App.vue'

const localVue = createLocalVue()
localVue.use(Buefy)
localVue.use(VueRouter)

const router = new VueRouter({
  routes: []
})

const createApp = propsData => mount(App, { sync: false, propsData, localVue, router })

describe('app', () => {
  beforeEach(() => {
    localStorage.setItem('authenticated', '0')
  })

  it("Should allow login or sign in when not authenticated", async () => {
    const cmp = createApp()

    expect(cmp.find("#login-link").isVisible()).toBe(true)
    expect(cmp.find("#sign-in-link").isVisible()).toBe(true)
    expect(cmp.find("#logout-link").isVisible()).toBe(false)
  });

  it("Should allow logout when authenticated", async () => {
    localStorage.setItem('authenticated', '1')

  });
});
