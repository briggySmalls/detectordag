<template>
  <b-navbar :sticky="true" toggleable="lg" type="dark" variant="info">
    <!-- Logo -->
    <b-navbar-brand to="/review" href="#" :active="$route.name === 'Review'">
      <img
        id="logo" alt="detector dag logo" src="../assets/logo.svg"
        class="d-inline-block">
        detector dag
    </b-navbar-brand>
    <!-- Navbar -->
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav>
        <b-nav-item to="/review" :active="$route.name ==='Review'">review</b-nav-item>
        <b-nav-item to="/account" :active="$route.name ==='Account'">settings</b-nav-item>
      </b-navbar-nav>
      <b-navbar-nav class="ml-auto">
        <b-nav-item v-on:click="logout">Sign out</b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component
export default class Navbar extends Vue {
  // Provide the username from the store
  private get username() {
    const { account } = this.$store.state;
    return (account) ? account.username : '?';
  }

  // Log the user out
  private logout() {
    // Clear the token and account
    this.$storage.clear();
    this.$store.commit('clearAccount');
    // Redirect to the login page
    this.$router.push('/login');
  }
}
</script>

<style lang="scss" scoped>
#logo {
  width: 3em;
  height: 3em;
}
nav {
  font-weight: 800;

  .nav-icon {
    width: 3em;
    height: 3em;
  }
}
</style>
