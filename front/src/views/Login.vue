<template>
  <div class="container">
    <div class="columns is-mobile">
      <div class="column is-half is-offset-one-quarter">
        <div class="card">

          <div class="card-content">
            <div class="level">
              <span class="level-item"/>
                <p class=" title is-3">Login</p>
              <span class="level-item"/>
            </div>

            <b-field
              horizontal
              label="Login"
              :type="{ 'is-danger': errors.has('login') }"
              :message="errors.first('login')">
              <b-input name="login" v-validate="'required'" icon="running" v-model="form.login"></b-input>
            </b-field>

            <b-field
              horizontal
              label="Password"
              :type="{ 'is-danger': errors.has('password') }"
              :message="[{
                    'The password field is required' : errors.firstByRule('password', 'required')
                }]"
              >
              <b-input name="password" v-validate="'required'" icon="lock" type="password" v-model="form.password" password-reveal></b-input>
            </b-field>

            <b-field horizontal>
              <p class="control">
              <button class="button is-primary" @click="onSubmit" :disabled="!isFormValid">
                Login
              </button>
              </p>
            </b-field>

            <b-message type="is-danger" v-if="!!errorMessage">
              {{errorMessage}}
            </b-message>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { login } from '@/requests/user'

export default {
  props: {
  },
  data () {
    return {
      form: {
        login: '',
        password: ''
      },
      errorMessage: null
    }
  },
  computed: {
    isFormValid () {
      return !Object.keys(this.fields).some(key => this.fields[key].invalid)
    }
  },
  methods: {
    onSubmit: async function (evt) {
      evt.preventDefault()
      try {
        await login(this.form)
        localStorage.setItem('authenticated', '1')
        this.$router.push('/')
      } catch (error) {
        this.errorMessage = error.message
      }
    }
  }
}
</script>
