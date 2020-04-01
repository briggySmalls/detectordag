<template>
  <div id="topbar">
    <!-- Display a navbar when logged in -->
    <NavbarComponent />
    <b-container id= "content-container" fluid="sm">
      <!-- Title -->
      <h1>{{ title }}</h1>
      <!-- Content -->
      <slot />
      <!-- Error -->
      <ErrorComponent :error="error" />
    </b-container>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import ErrorComponent from '../components/Error.vue';
import NavbarComponent from '../components/Navbar.vue';

@Component({
  components: {
    ErrorComponent,
    NavbarComponent,
  },
})
export default class Topbar extends Vue {
  @Prop() private readonly error!: Error;

  @Prop() private readonly title!: string;
}
</script>

<style lang="scss" scoped>
@import "~bootstrap/scss/functions";
@import "~bootstrap/scss/variables";
@import "~bootstrap/scss/mixins";

@include media-breakpoint-up(sm) {
  #content-container {
    max-width: map-get($grid-breakpoints, sm);
  }
}
</style>
