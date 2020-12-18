<template>
  <Splash title='Page not found!' :error='error' />
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import Splash from '../layouts/Splash.vue';

@Component({
  components: {
    Splash,
  },
})
export default class NotFound extends Vue {
  private readonly error = null;
}
</script>
