import { mount, createLocalVue } from '@vue/test-utils'

import Buefy from 'buefy'
import VueRouter from 'vue-router'
import Vuex from 'vuex'
import store from '@/store'
import App from '@/App.vue'
import fetchMock from 'fetch-mock'

const localVue = createLocalVue()
localVue.use(Buefy)
localVue.use(VueRouter)
localVue.use(Vuex)

const router = new VueRouter({
  routes: []
})

const createApp = propsData => mount(App, { sync: false, propsData, store, localVue, router })

describe('app', () => {
  beforeEach(async() => {
    await store.commit('user/logout')
  })

  it("Should allow login or sign in when not authenticated", () => {
    const cmp = createApp()

    expect(cmp.find("#login-link").isVisible()).toBe(true)
    expect(cmp.find("#sign-in-link").isVisible()).toBe(true)
    expect(cmp.find("#logout-link").isVisible()).toBe(false)
  });

  it("Should allow logout when authenticated", async () => {
    await store.commit('user/login')
    const cmp = createApp()

    expect(store.state.user.authenticated).toBe(true)
    expect(cmp.find("#login-link").isVisible()).toBe(false)
    expect(cmp.find("#sign-in-link").isVisible()).toBe(false)
    expect(cmp.find("#logout-link").isVisible()).toBe(true)
  });

  it("Should logout", async () => {
    await store.commit('user/login')
    const cmp = createApp()

    const evt = { preventDefault: jest.fn() }
    fetchMock.postOnce((url, opt) => {
      return (
        url === '/app/public/logout'
      )
    }, {})

    await cmp.vm.onLogout(evt)

    expect(store.state.user.authenticated).toBe(false)
    expect(cmp.find("#login-link").isVisible()).toBe(true)
    expect(cmp.find("#sign-in-link").isVisible()).toBe(true)
    expect(cmp.find("#logout-link").isVisible()).toBe(false)
  });
});
