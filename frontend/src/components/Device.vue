<template>
  <b-card :title="device.deviceId">
    <b-list-group v-for="(value, key) in device.state" :key="key" flush>
      <b-list-group-item>
        {{ key }}: {{ value }}
      </b-list-group-item>
    </b-list-group>
    <footer>
      Updated: {{ device.updated }}
    </footer>
  </b-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Device as DeviceModel } from '../../lib/client';

@Component
export default class Device extends Vue {
  @Prop() private device!: DeviceModel;
}
</script>
