import Vue from 'vue';
import Vuex from 'vuex';
import { Device } from '../../lib/client';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    account: null,
    devices: Array <Device>(),
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
      state.devices = [];
    },
    setDevice(state, device) {
      const index = state.devices.findIndex((x) => x.deviceId === device.deviceId);
      if (index === -1) {
        state.devices.push(device);
      }
      state.devices[index] = device;
    },
  },
  actions: {
  },
  modules: {
  },
});
