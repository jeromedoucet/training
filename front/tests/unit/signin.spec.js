import { mount, createLocalVue } from '@vue/test-utils'

import fetchMock from 'fetch-mock'
import Buefy from 'buefy'
import VeeValidate from 'vee-validate';
import VueRouter from 'vue-router'
import SignIn from '@/views/SignIn.vue'

const localVue = createLocalVue()
localVue.use(VeeValidate);
localVue.use(Buefy)
localVue.use(VueRouter)

const router = new VueRouter({
  routes: []
})

const createSignIn = propsData => mount(SignIn, { propsData, localVue, router })

// Endpoint de test d'authentication
//
// cas retour erreur serveur => affichage d'une erreur
//
// cas validation formulaire en échec
//
// cas on est authentifié => ne pas pouvoir aller dessus (e2e ?)

describe('Sign in', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    localStorage.setItem('authenticated', '0')
  })

  it("Should create a user and redirect to '/'", async () => {
    // given
    const cmp = createSignIn()
    cmp.setData({ form: { login: 'jerdct', firstName: 'Jérôme', lastName: 'Doucet', email: 'jerdct@gmail.com', password: 'toto' } })
    cmp.vm.$router.push = jest.fn()
    const evt = { preventDefault: jest.fn() }

    fetchMock.postOnce((url, opt) => {
      return (
        url === '/app/public/users' &&
        JSON.parse(opt.body).login === 'jerdct' &&
        JSON.parse(opt.body).firstName === 'Jérôme' &&
        JSON.parse(opt.body).lastName === 'Doucet' &&
        JSON.parse(opt.body).email === 'jerdct@gmail.com' &&
        JSON.parse(opt.body).password === 'toto'
      )
    }, {})

    // when
    await cmp.vm.onSubmit(evt)

    // then
    expect(cmp.vm.$router.push.mock.calls.length).toBe(1)
    expect(cmp.vm.$router.push.mock.calls[0][0]).toBe('/')
    expect(localStorage.getItem('authenticated')).toBe('1')
  })

  it('Should print an error message when issue happend', async () => {
    // given
    const cmp = createSignIn()
    cmp.setData({ form: { login: 'jerdct', firstName: 'Jérôme', lastName: 'Doucet', email: 'jerdct@gmail.com', password: 'toto' } })
    cmp.vm.$router.push = jest.fn()
    const evt = { preventDefault: jest.fn() }

    fetchMock.postOnce((url, opt) => {
      return (
        url === '/app/public/users' &&
        JSON.parse(opt.body).login === 'jerdct' &&
        JSON.parse(opt.body).firstName === 'Jérôme' &&
        JSON.parse(opt.body).lastName === 'Doucet' &&
        JSON.parse(opt.body).email === 'jerdct@gmail.com' &&
        JSON.parse(opt.body).password === 'toto'
      )
    }, {
      status: 409,
      body: { message: 'Another user already exist with this identifier' }
    })

    // when
    await cmp.vm.onSubmit(evt)

    // then
    expect(cmp.vm.$router.push.mock.calls.length).toBe(0)
    expect(localStorage.getItem('authenticated')).toBe('0')
    expect(cmp.vm.errorMessage).toBe('Another user already exist with this identifier')
  })
})
