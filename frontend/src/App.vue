<template>
  <div id="app">
    <b-container>
      <router-view/>
    </b-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { BootstrapVue } from 'bootstrap-vue';
import { ApiClient } from '../lib/client';

// Install BootstrapVue
Vue.use(BootstrapVue);

@Component
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
