import Vue from 'vue';
import { ApiClient, AccountsApi, AuthenticationApi } from '../../lib/client';

// Define a wrapper for the different clients
class ClientWrapper {
  public readonly accounts: AccountsApi;

  public readonly authentication: AuthenticationApi;

  public constructor(client: ApiClient) {
    this.accounts = new AccountsApi(client);
    this.authentication = new AuthenticationApi(client);
  }
}

// Create the underlying client
const client = new ApiClient();
client.basePath = 'http://localhost:8080/api/v1';

// Create an instance of our wrapper
const wrapper = new ClientWrapper(client);

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
