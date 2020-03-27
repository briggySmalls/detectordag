<template>
  <b-container id="login">
    <h1>Login</h1>
    <b-form
      @submit="submit"
      action="https://vuejs.org/"
      method="post">
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
    <div v-if="error">
      {{ error.message }}
    </div>
  </b-container>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { AuthenticationApi, Credentials, Token } from '../../lib/client';

@Component
export default class Login extends Vue {
  private email = '';

  private password = '';

  private client: AuthenticationApi;

  public error: Error | null = null;

  public constructor() {
    // Call super
    super();
    // Create client
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
      this.error = new Error(response.text);
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

<style lang="scss" scoped>
#login {
  max-width: 30em;
}
</style>
