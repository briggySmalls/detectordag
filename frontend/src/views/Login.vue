<template>
  <div class="login">
    <h1>This is the login page</h1>
    <form
        @submit="submit"
        action="https://vuejs.org/"
        method="post">
      <div>
        <label for="email">Email</label>
        <input id="email" type="email" name="email" v-model="email" required>
      </div>
      <div>
        <label for="password">Password</label>
        <input id="password" type="password" name="password" v-model="password" required>
      </div>
      <input type="submit" value="Submit">
    </form>
    <div v-if="error">
      {{ error.message }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import {
  ApiClient, AuthenticationApi, Credentials, Token,
} from '../../lib/client';

@Component
export default class Login extends Vue {
  private email = '';

  private password = '';

  private client: AuthenticationApi;

  public error!: Error = null;

  public constructor() {
    // Call super
    super();
    // Create client
    ApiClient.instance.basePath = 'http://localhost:8080/api/v1';
    this.client = new AuthenticationApi();
  }

  public submit(event: Event) { // eslint-disable-line class-methods-use-this
    // Submit a request to the backend
    console.log(`Submitting request, {"email": "${this.email}", "password": "${this.password}"}`);
    // Create the request body
    const creds = new Credentials(this.email, this.password);
    // Submit the request
    this.client.auth(creds, this.handleLogin);
    // Do not actually perform a post action
    event.preventDefault();
  }

  private handleLogin( // eslint-disable-line class-methods-use-this
    error: Error, data: Token, response: any,
  ) {
    if (error) {
      // Assign the error
      this.error = error;
      // Also log it
      console.error(error);
    } else {
      console.log(`API called successfully. Returned data: ${data}`);
      // Record the token and account
      localStorage.setItem('token', data.token);
      localStorage.setItem('accountId', data.accountId);
      // Navigate home
      this.$router.push('/');
    }
  }
}
</script>
