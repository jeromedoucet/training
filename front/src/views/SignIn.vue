<template>
  <div class="container">
    <div class="columns is-mobile">
      <div class="column is-half is-offset-one-quarter">
        <div class="card">

          <div class="card-content">
            <div class="level">
              <span class="level-item"/>
                <p class=" title is-3">Create your account !</p>
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
              label="First name"
              :type="{ 'is-danger': errors.has('firstName') }"
              :message="errors.first('firstName')"
              >
              <b-input name="firstName" v-validate="'required'" icon="user" v-model="form.firstName"></b-input>
            </b-field>

            <b-field
              horizontal
              label="Last name"
              :type="{ 'is-danger': errors.has('lastName') }"
              :message="errors.first('lastName')"
              >
              <b-input name="lastName" v-validate="'required'" icon="user" v-model="form.lastName"></b-input>
            </b-field>

            <b-field
              horizontal
              label="email"
              :type="{ 'is-danger': errors.has('email') }"
              :message="errors.first('email')"
              >
              <b-input name="email" v-validate="'required|email'" icon="at" type="email" v-model="form.email"></b-input>
            </b-field>

            <b-field
              horizontal
              label="Password"
              :type="{ 'is-danger': errors.has('password') }"
              :message="[{
                    'The password field is required' : errors.firstByRule('password', 'required'),
                    'The password must have 10 characters at least ' : errors.firstByRule('password', 'min')
                }]"
              >
              <b-input name="password" v-validate="'required|min:10'" icon="lock" type="password" v-model="form.password" password-reveal></b-input>
            </b-field>

            <b-field horizontal>
              <p class="control">
              <button class="button is-primary" @click="onSubmit" :disabled="!isFormValid">
                Register
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
import { createUser } from '@/requests/user'

export default {
  props: {
  },
  data () {
    return {
      form: {
        login: '',
        firstName: '',
        lastName: '',
        email: '',
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
        await createUser(this.form)
        localStorage.setItem('authenticated', '1')
        this.$router.push('/')
      } catch (error) {
        this.errorMessage = error.message
      }
    }
  }
}
</script>
