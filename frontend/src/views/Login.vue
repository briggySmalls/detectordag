<template>
  <b-container id="login">
    <h1>Login</h1>
    <b-form @submit="submit">
      <b-form-group
        id="email"
        label="Email:"
        label-for="email">
        <b-form-input
          id="email" type="email" name="email"
          v-model="email"
          trim required>
        </b-form-input>
      </b-form-group>
      <b-form-group
        id="password"
        label="Password:"
        label-for="password">
        <b-form-input
          id="password" type="password" name="password"
          v-model="password"
          trim required>
        </b-form-input>
      </b-form-group>
      <b-button type="submit" >Submit</b-button>
    </b-form>
    <ErrorComponent :error="error" />
  </b-container>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import {
  AuthenticationApi, AccountsApi, Credentials, Token,
} from '../../lib/client';
import { Storage, AuthBundle } from '../utils';
import ErrorComponent from '../components/Error.vue';

@Component({
  components: {
    ErrorComponent,
  },
})

// The Login route is responsible for ensuring a user's token and account are stored
export default class Login extends Vue {
  private email = '';

  private password = '';

  private authClient: AuthenticationApi;

  private accountsClient: AccountsApi;

  private storage: Storage;

  public error: Error | null = null;

  public constructor() {
    // Call super
    super();
    // Create clients
    this.authClient = new AuthenticationApi();
    this.accountsClient = new AccountsApi();
    // Create storage helper
    this.storage = new Storage();
  }

  public created() {
    // First we check if the user already has a token
    const { bundle } = this.storage;
    if (bundle == null) {
      // There is no token, so display login
      this.$logger.debug('No preexisting token found');
      return;
    }
    // Check if the token is valid
    this.$logger.debug('Found token, requesting accounts');
    this.accountsClient.getAccount(`Bearer ${bundle.token}`, bundle.accountId, this.handleAccount);
  }

  public submit(event: Event) {
    // Create the request body
    const creds = new Credentials(this.email, this.password);
    // Submit the request
    this.authClient.auth(creds, this.handleLogin);
    // Do not actually perform a post action
    event.preventDefault();
  }

  private handleAccount(error: Error, data: Account, response: any) {
    // Handle any errors
    if (error) {
      // Assign the error
      this.error = error;
      // Also log it
      console.error(response.text);
      return;
    }
    // We have got the account, job done!
    this.finishUp(data);
  }

  private handleLogin(error: Error, data: Token, response: any) {
    if (error) {
      // Assign the error
      this.error = error;
      // Also log it
      console.error(response.text);
      return;
    }
    // Record the token and account in local storage
    const bundle = new AuthBundle(data.accountId, data.token);
    this.storage.save(bundle);
    // Now we have a valid token, let's get the account details
    this.requestAccount(bundle);
  }

  // Request the account and configure callback
  private requestAccount(bundle: Token) {
    this.accountsClient.getAccount(`Bearer ${bundle.token}`, bundle.accountId, this.handleAccount);
  }

  // Save the account and move on
  private finishUp(account: Account) {
    // We have got the account, so save it
    this.$store.commit('setAccount', account);
    // Redirect to the dashboard
    this.$router.push('/dashboard');
  }
}
</script>

<style lang="scss" scoped>
#login {
  max-width: 30em;
}
</style>
