import Vue from 'vue';
import { AccountsApi, AuthenticationApi } from '../../lib/client';

// Define a wrapper for the different clients
class ClientWrapper {
  public readonly accounts: AccountsApi;

  public readonly authentication: AuthenticationApi;

  public constructor() {
    const bPath = process.env.VUE_APP_API_BASEPATH || `${process.env.BASE_URL}/api/v1`;
    this.accounts = new AccountsApi(undefined, bPath);
    this.authentication = new AuthenticationApi(undefined, bPath);
  }
}

// Create an instance of our wrapper
const wrapper = new ClientWrapper();

// Add the clients to Vue
Vue.prototype.$clients = wrapper;

// Update the type hinting for all Vue instances
declare module 'vue/types/vue' {
  interface Vue {
    $clients: ClientWrapper;
  }
}

// Export the wrapper if people want to use it directly
export default wrapper;
