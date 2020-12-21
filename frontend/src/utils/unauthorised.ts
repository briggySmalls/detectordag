import Vue from 'vue';
import { AxiosError } from 'axios';

// Add the clients to Vue
Vue.prototype.$unauthorised = function unauthorised(error: AxiosError) {
  if (error.response === undefined || error.response.status !== 403) {
    return false;
  }
  // The auth bundle is invalid, so clear it
  this.$storage.clear();
  // Redirect to login
  this.$router.push('/login');
  return true;
};

// Update the type hinting for all Vue instances
declare module 'vue/types/vue' {
  interface Vue {
    $unauthorised(error: AxiosError): boolean;
  }
}
