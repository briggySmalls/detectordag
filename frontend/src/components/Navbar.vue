<template>
  <b-navbar sticky="true" toggleable="lg" type="dark" variant="info">
    <!-- Logo -->
    <b-navbar-brand href="#">
      <img
        id="logo" alt="detector dag logo" src="../assets/logo.svg"
        class="d-inline-block">
        detector dag
    </b-navbar-brand>
    <!-- Navbar -->
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav class="ml-auto">
        <!-- Logged in navbar content -->
        <template v-if="$store.state.account">
          <b-nav-item-dropdown :text="username" right>
            <b-dropdown-item to="/account" href="#">Settings</b-dropdown-item>
            <b-dropdown-item v-on:click="logout" href="#">Sign out</b-dropdown-item>
          </b-nav-item-dropdown>
        </template>
        <!-- Loading -->
        <template v-else>
          <b-spinner></b-spinner>
        </template>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { storage } from '../utils';

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
    storage.clear();
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
}
</style>
