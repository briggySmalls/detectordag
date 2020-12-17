<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { BootstrapVue } from 'bootstrap-vue';

// Install BootstrapVue
Vue.use(BootstrapVue);

@Component
export default class App extends Vue {}
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
}
</style>
