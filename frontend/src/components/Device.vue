<template>
  <b-card :title="device.name">
    <b-card-text>
      Status: {{ powerStatus }}
    </b-card-text>
    <template v-slot:footer>
      Updated: {{ device.state.updated.toLocaleString() }}
    </template>
  </b-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Device as DeviceModel } from '../../lib/client';

@Component
export default class Device extends Vue {
  @Prop() private device!: DeviceModel;

  private get powerStatus(): string {
    return (this.device.state.power) ? '⚡️ On' : '❗️ Off';
  }
}
</script>
