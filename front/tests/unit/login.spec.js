import { mount, createLocalVue } from '@vue/test-utils'

import fetchMock from 'fetch-mock'
import Buefy from 'buefy'
import VeeValidate from 'vee-validate';
import VueRouter from 'vue-router'
import Login from '@/views/Login.vue'

const localVue = createLocalVue()
localVue.use(VeeValidate);
localVue.use(Buefy)
localVue.use(VueRouter)

const router = new VueRouter({
  routes: []
})

const createLogin = propsData => mount(Login, { sync: false, propsData, localVue, router })

describe('login', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    localStorage.setItem('authenticated', '0')
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
    expect(localStorage.getItem('authenticated')).toBe('1')

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
    expect(localStorage.getItem('authenticated')).toBe('0')
    expect(cmp.vm.errorMessage).toBe('Wrong credentials.')

  });
});
