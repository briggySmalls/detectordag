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
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component
export default class Login extends Vue {
  private email = '';

  private password = '';

  public submit(event: Event) { // eslint-disable-line class-methods-use-this
    // Submit a request to the backend
    console.log(`Submitting request, {"email": "${this.email}", "password": "${this.password}"}`);
    // Do not actually perform a post action
    event.preventDefault();
  }
}
</script>
