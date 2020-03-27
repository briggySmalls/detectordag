<template>
  <b-alert class="mt-3" variant="danger" :show="error">
    {{ error.message }}
  </b-alert>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class Device extends Vue {
  @Prop() private error: Error | null;
}
</script>
