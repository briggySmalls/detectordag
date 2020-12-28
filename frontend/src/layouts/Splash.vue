<template>
  <b-container
    id="splash"
    class="mt-5"
    fluid="sm"
  >
    <!-- Logo -->
    <img
      id="logo"
      class="my-1"
      alt="detector dag logo"
      src="../assets/logo.svg"
    >
    <!-- Title -->
    <h1 class="my-2">
      {{ title }}
    </h1>
    <!-- Content -->
    <slot />
    <!-- Error -->
    <ErrorComponent :error="error" />
  </b-container>
</template>

<script lang="ts">
import { Component, Vue, Prop } from 'vue-property-decorator';
import ErrorComponent from '../components/Error.vue';

@Component({
  components: {
    ErrorComponent,
  },
})
export default class Splash extends Vue {
  @Prop() private readonly error!: Error;

  @Prop() private readonly title!: string;
}
</script>

<style lang="scss" scoped>
@import '../app.scss';

@include media-breakpoint-up(sm) {
  #content-container {
    max-width: map-get($grid-breakpoints, sm);
  }
}
</style>
