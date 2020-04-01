<template>
  <Splash :title='title' :error='error' />
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
  private readonly title = 'Page not found!';

  private readonly error = null;
}
</script>
