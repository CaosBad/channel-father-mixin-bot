import Vue from 'vue'
import Vuex from 'vuex'
import storage from './local'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    account: {},
    bot:{}
  },
  getters: {
    getBot: state => {
      return state.bot
    }
  },
  mutations: {
    SaveAccount (state, account) {
      state.account = JSON.parse(account)
      storage.set("account", account)
    },
    SaveBot (state, bot) {
      state.bot = bot
    }
  },
  actions: {
  }
})
