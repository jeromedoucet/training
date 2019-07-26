const state = {
  authenticated: false
}

const mutations = {
  login (state) {
    state.authenticated = true
  },

  logout (state) {
    state.authenticated = false
  }
}

export default {
  namespaced: true,
  state,
  mutations
}
