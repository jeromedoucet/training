import { mount, createLocalVue } from '@vue/test-utils'

import fetchMock from 'fetch-mock'
import Buefy from 'buefy'
import VeeValidate from 'vee-validate';
import VueRouter from 'vue-router'
import Vuex from 'vuex'
import store from '@/store'
import Login from '@/views/Login.vue'

const localVue = createLocalVue()
localVue.use(VeeValidate);
localVue.use(Buefy)
localVue.use(VueRouter)
localVue.use(Vuex)

const router = new VueRouter({
  routes: []
})

const createLogin = propsData => mount(Login, { sync: false, propsData, store, localVue, router })

describe('login', () => {
  beforeEach(async () => {
    await store.commit('user/logout')
    jest.clearAllMocks()
  })

  it("Should log a user in and redirect to '/'", async () => {
    const cmp = createLogin()
    cmp.setData({ form: { login: 'jerdct', password: 'toto' } })
    cmp.vm.$router.push = jest.fn()
    const evt = { preventDefault: jest.fn() }

    fetchMock.postOnce((url, opt) => {
      return (
        url === '/app/public/login' &&
        JSON.parse(opt.body).login === 'jerdct' &&
        JSON.parse(opt.body).password === 'toto'
      )
    }, {})

    // when
    await cmp.vm.onSubmit(evt)

    // then
    expect(cmp.vm.$router.push.mock.calls.length).toBe(1)
    expect(cmp.vm.$router.push.mock.calls[0][0]).toBe('/')
    expect(store.state.user.authenticated).toBe(true)
  });

  it("Should print an error message when issue happend", async () => {
    const cmp = createLogin()
    cmp.setData({ form: { login: 'jerdct', password: 'toto' } })
    cmp.vm.$router.push = jest.fn()
    const evt = { preventDefault: jest.fn() }

    fetchMock.postOnce((url, opt) => {
      return (
        url === '/app/public/login' &&
        JSON.parse(opt.body).login === 'jerdct' &&
        JSON.parse(opt.body).password === 'toto'
      )
    }, {
      status: 409,
      body: { message: 'Wrong credentials.' }
    })

    // when
    await cmp.vm.onSubmit(evt)

    // then
    expect(cmp.vm.$router.push.mock.calls.length).toBe(0)
    expect(store.state.user.authenticated).toBe(false)
    expect(cmp.vm.errorMessage).toBe('Wrong credentials.')
  });
});
