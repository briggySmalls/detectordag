import Vue from 'vue';
import { AxiosError } from 'axios';

// Add the clients to Vue
Vue.prototype.$checkUnauthorised = function checkUnauthorised(
  error: AxiosError,
  callback: (error: AxiosError) => void,
) {
  if (error.response === undefined || error.response.status !== 403) {
    callback(error);
    return;
  }
  // The auth bundle is invalid, so clear it
  this.$storage.clear();
  // Redirect to login
  this.$router.push('/login');
};

// Update the type hinting for all Vue instances
declare module 'vue/types/vue' {
  interface Vue {
    $checkUnauthorised(error: AxiosError, callback: (error: AxiosError) => void): boolean;
  }
}
