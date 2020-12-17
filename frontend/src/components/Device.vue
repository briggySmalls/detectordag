<template>
  <b-card :title="device.name">
    <b-card-body>
      <b-card-sub-title class="mb-2">{{ deviceStatus }}</b-card-sub-title>
    </b-card-body>
    <b-list-group flush>
      <b-list-group-item v-if="isOff">
        Turned off at: {{ device.state.updated.toLocaleString() }}
      </b-list-group-item>
      <b-list-group-item v-if="isDisconnected">
        Last seen at: {{ device.connection.updated.toLocaleString() }}
      </b-list-group-item>
    </b-list-group>
  </b-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { Device as DeviceModel } from '../../lib/client';

enum DeviceState {
  On = 1,
  Off,
  WasOn,
  WasOff,
}

enum PowerState {
  On = 1,
  Off,
}

enum ConnectedState {
  Connected = 1,
  Disconnected,
}

@Component
export default class Device extends Vue {
  @Prop() private device!: DeviceModel;

  private get deviceStatus(): string {
    switch (this.state) {
      case DeviceState.On:
        return 'On';
      case DeviceState.Off:
        return 'Off';
      case DeviceState.WasOn:
        return 'Was On';
      case DeviceState.WasOff:
        return 'Was Off';
      default:
        throw new Error(`Unexpected state: ${this.state}`);
    }
  }

  private get state(): DeviceState {
    if (this.connectionStatue === ConnectedState.Connected) {
      return (this.powerState === PowerState.On) ? DeviceState.On : DeviceState.Off;
    }
    return (this.powerState === PowerState.On) ? DeviceState.WasOn : DeviceState.WasOff;
  }

  private get isOff(): bool {
    return this.powerState === PowerState.Off;
  }

  private get isDisconnected(): bool {
    return this.connectionState === ConnectedState.Disconnected;
  }

  private get connectionState(): ConnectedState {
    switch (this.device.connection.status) {
      case 'connected':
        return ConnectedState.Connected;
      case 'disconnected':
        return ConnectedState.Disconnected;
      default:
        throw new Error(`Unexpected connection state: "${this.device.connection.status}"`);
    }
  }

  private get powerState(): PowerState {
    switch (this.device.state.power) {
      case 'on':
        return PowerState.On;
      case 'off':
        return PowerState.Off;
      default:
        throw new Error(`Unexpected device power: "${this.device.state.power}"`);
    }
  }
}
</script>
