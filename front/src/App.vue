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
        <router-link class="r-link" to="/login"><a class="navbar-item">Login</a></router-link>
        <router-link class="r-link" to="/sign-in"><a class="navbar-item">Sign in</a></router-link>
      </div>
    </nav>
    <div>
      <section class="section">
          <router-view/>
      </section>
    </div>
  </div>
</template>

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
