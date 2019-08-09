<template>
  <div class="home">
    <dashboard v-if="authenticated" />
  </div>
</template>

<script>

import Dashboard from '@/components/Dashboard'

export default {
  computed: {
    authenticated () {
      return !!this.$store.state.user.authenticated
    }
  },
  components: {
    Dashboard
  }
}
</script>
