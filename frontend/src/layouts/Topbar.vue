<template>
  <div id="topbar">
    <!-- Display a navbar when logged in -->
    <NavbarComponent />
    <b-container id="content-container" fluid="sm">
      <!-- Title -->
      <slot name="header">
        <h1 class="mt-5">{{ title }}</h1>
      </slot>
      <!-- Content -->
      <slot />
      <!-- Error -->
      <ErrorComponent :error="error" />
    </b-container>
    <footer id="footer">Made with ❤ by <a href="https://sambriggs.dev">sam briggs</a></footer>
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
@import '../app.scss';

$footer-padding: 100px;

#topbar {
  min-height: 100vh;
  position: relative;

  #content-container {
    padding-bottom: $footer-padding;
  }

  footer#footer {
    // Make the footer sit at the bottom of the page
    margin-top: 5em;
    position: absolute;
    width: 100%;
    bottom: 0;
  }
}

@include media-breakpoint-up(sm) {
  #content-container {
    max-width: map-get($grid-breakpoints, sm);
  }
}
</style>
