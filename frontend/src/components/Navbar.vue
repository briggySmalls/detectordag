<template>
  <b-navbar toggleable="lg">
    <!-- Logo -->
    <b-navbar-brand href="#">
      <img
        id="logo" alt="DetectorDag logo" src="../assets/logo.svg"
        class="d-inline-block">
        DetectorDag
    </b-navbar-brand>
    <!-- Navbar -->
    <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>
    <b-collapse id="nav-collapse" is-nav>
      <b-navbar-nav class="ml-auto">
        <!-- Logged in navbar content -->
        <template v-if="$store.state.account">
          <b-nav-text>{{ username }}</b-nav-text>
          <b-nav-item>
            <b-button size="sm" v-on:click="logout">Logout</b-button>
          </b-nav-item>
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
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Navbar extends Vue {
  private get username() {
    const { account } = this.$store.state;
    return (account) ? account.username : '?';
  }
}
</script>
