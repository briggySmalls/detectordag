<template>
  <Splash id="login" title="detector dag" :error="error">
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
import { AuthBundle } from '../utils';
import Splash from '../layouts/Splash.vue';

@Component({
  components: {
    Splash,
  },
})
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
    this.$clients.authentication.auth({ username: this.email, password: this.password })
      .then((response) => {
        // Record the token and account in local storage
        this.$logger.debug('Auth response received');
        // Build and save the authorization data into a handy bundle
        const bundle = new AuthBundle(response.data.accountId, response.data.token);
        this.$storage.save(bundle);
        // Redirect to review
        this.$router.push('review');
      })
      .catch((error) => {
        // Log the real response
        // See https://github.com/swagger-api/swagger-codegen/issues/2602
        this.$logger.debug(error.response);
        // Assign the error
        this.error = error;
      })
      .then(() => {
        // Indicate we've finished-up
        this.isRequesting = false;
      });
    // Do not actually perform a post action
    event.preventDefault();
  }
}
</script>

<style lang="scss" scoped>
#login {
  max-width: 20em;
}
</style>
-
