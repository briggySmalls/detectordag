import Vue from 'vue';
import Vuex from 'vuex';
import State from './state';

Vue.use(Vuex);

export default new Vuex.Store({
  state: new State(),
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
    setDevice(state, device) {
      // We'd expect the devices to already be here
      if (state.devices === null) {
        throw new Error('Cannot setDevice before setDevices');
      }
      const index = state.devices.findIndex((x) => x.deviceId === device.deviceId);
      // Use splice so changes to array are noticed
      state.devices.splice(index, 1, device);
    },
  },
  actions: {},
  modules: {},
});
