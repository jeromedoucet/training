<template>
  <div>
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <router-link class="r-link" to="/">
          <a class="navbar-item" >
            <img src="@/assets/training_logo.png">
          </a>
        </router-link>
      </div>
      <div class="navbar-end">
        <router-link v-show="!authenticated" id="login-link" class="r-link" to="/login"><a class="navbar-item">Login</a></router-link>
        <router-link v-show="!authenticated" id="sign-in-link" class="r-link" to="/sign-in"><a class="navbar-item">Sign in</a></router-link>
        <a id="logout-link" v-show="authenticated" @click="onLogout" class="navbar-item">Logout</a>
      </div>
    </nav>
    <div>
      <section class="section">
          <router-view/>
      </section>
    </div>
  </div>
</template>
<script>

import { logout } from '@/requests/user'

export default {
  computed: {
    authenticated () {
      return !!this.$store.state.user.authenticated
    }
  },
  methods: {
    onLogout: async function (evt) {
      evt.preventDefault()
      try {
        await logout()
        await this.$store.commit('user/logout')
      } catch (error) {
      }
    }
  }
}

</script>

<style lang="scss">
// Bulma + Bulmaswatch
@import "node_modules/bulmaswatch/darkly/variables";
@import "node_modules/bulma/bulma";
@import "node_modules/bulmaswatch/darkly/overrides";

// Buefy
@import "node_modules/buefy/src/scss/buefy";

.r-link {
  display: inherit;
}
</style>
