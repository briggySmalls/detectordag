<template>
  <div id="app">
    <b-navbar class="justify-content-center">
      <b-navbar-brand href="#">
        <img id="logo" alt="Detectordag logo" src="./assets/logo.svg"
             class="d-inline-block">
        Detectordag
      </b-navbar-brand>
    </b-navbar>
    <router-view/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { BootstrapVue } from 'bootstrap-vue';
import { ApiClient } from '../lib/client';

// Install BootstrapVue
Vue.use(BootstrapVue);

@Component({})
export default class App extends Vue {
  constructor() {
    super();
    // Configure client endpoint
    ApiClient.instance.basePath = 'http://localhost:8080/api/v1';
  }
}
</script>

<style lang="scss">
// Import global styles
@import './app.scss';

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;

  #logo {
    max-width: 400px;
  }
}
</style>
