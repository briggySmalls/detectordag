<template>
  <Splash id="login" :title="title" :error="error">
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
  </Splash>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Token } from '../../lib/client';
import { storage, AuthBundle } from '../utils';
import Splash from '../layouts/Splash.vue';

@Component({
  components: {
    Splash,
  },
})
export default class Login extends Vue {
  // Page title
  private readonly title = 'DetectorDag'

  private email = '';

  private password = '';

  public error: Error | null = null;

  private isRequesting = false;

  public submit(event: Event) {
    this.$logger.debug('Login submitted');
    // Request authentication
    this.isRequesting = true;
    this.error = null;
    this.$clients.authentication.auth(
      { username: this.email, password: this.password },
      this.handleLogin,
    );
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
    this.$logger.debug('Auth response received');

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
</style>
-
