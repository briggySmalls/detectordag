<template>
  <b-alert
    v-if="error !== null"
    class="mt-3"
    variant="danger"
    show
  >
    {{ error.message }}
  </b-alert>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class ErrorComponent extends Vue {
  @Prop() private readonly error!: Error;
}
</script>
