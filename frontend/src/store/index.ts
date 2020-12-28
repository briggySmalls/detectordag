import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    account: null,
    devices: null,
  },
  mutations: {
    setAccount(state, newAccount) {
      state.account = newAccount;
    },
    clearAccount(state) {
      state.account = null;
    },
    setDevices(state, devices) {
      state.devices = devices;
    },
    clearDevices(state) {
      state.devices = null;
    },
  },
  actions: {},
  modules: {},
});
