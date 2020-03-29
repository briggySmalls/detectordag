<template>
  <b-container id="login" class="mt-5">
    <img id="logo" class="my-1" alt="Detectordag logo" src="../assets/logo.svg">
    <h1 class="my-2">DetectorDag</h1>
    <b-form v-if="!isRequesting" @submit="submit">
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
    <b-spinner v-else label="Spinning"></b-spinner>
    <ErrorComponent :error="error" />
  </b-container>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Credentials, Token } from '../../lib/client';
import { storage, AuthBundle } from '../utils';
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

  public error: Error | null = null;

  private isRequesting = false;

  public submit(event: Event) {
    this.$logger.debug('Login submitted');
    // Request authentication
    this.isRequesting = true;
    this.error = null;
    this.$clients.authentication.auth(new Credentials(this.email, this.password), this.handleLogin);
    // Do not actually perform a post action
    event.preventDefault();
  }

  private handleLogin(error: Error, data: Token, response: any) {
    // Handle any errors
    this.isRequesting = false;
    if (error) {
      // Log the real response
      // See https://github.com/swagger-api/swagger-codegen/issues/2602
      this.$logger.debug(response.text);
      // Assign the error
      this.error = error;
      return;
    }
    // Record the token and account in local storage
    this.$logger.debug('Account data received');
    const bundle = new AuthBundle(data.accountId, data.token);
    storage.save(bundle);
    // Redirect to review
    this.$router.push('review');
  }
}
</script>

<style lang="scss" scoped>
#login {
  max-width: 20em;
}
#logo {
  max-width: 30em;
}
</style>
-
