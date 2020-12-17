<template>
  <b-card :title="device.name">
    <div class="status-graphic d-inline-block">
      <div class="power-icon-container" v-bind:class="[deviceStateClass]">
        <img v-if="powerState === powerStateEnum.Off"
          alt="Power off" src="../assets/power-off.svg" class="power-icon">
        <img v-else-if="powerState === powerStateEnum.On"
          alt="Power on" src="../assets/power-on.svg" class="power-icon">
      </div>
      <img v-if="connectionState === connectedStateEnum.Disconnected"
        alt="Disconnected" src="../assets/no-signal.svg" class="connection-icon">
    </div>
    <b-card-body>
      <b-card-sub-title class="mb-2">{{ deviceStateText }}</b-card-sub-title>
    </b-card-body>
    <b-list-group flush>
      <b-list-group-item v-if="powerState === powerStateEnum.Off">
        Lost power
        <span class="time" :datetime="device.state.updated">
          {{ device.state.updated }}
        </span>
      </b-list-group-item>
      <b-list-group-item v-if="connectionState === connectedStateEnum.Disconnected">
        Lost connection
        <span class="time" :datetime="device.connection.updated">
          {{ device.connection.updated }}
        </span>
      </b-list-group-item>
    </b-list-group>
  </b-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';
import { render } from 'timeago.js';
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

interface StateData {
  class: string;
  text: string;
}

@Component
export default class Device extends Vue {
  // Declare some enums so we can use them in the template
  private const deviceStateEnum: typeof DeviceState = DeviceState

  private const powerStateEnum: typeof PowerState = PowerState

  private const connectedStateEnum: typeof ConnectedState = ConnectedState

  @Prop() private device!: DeviceModel;

  private const deviceStateInfo: Record<DeviceState, StateData> = {
    [DeviceState.On]: { class: 'on', text: 'On' },
    [DeviceState.Off]: { class: 'off', text: 'Off' },
    [DeviceState.WasOn]: { class: 'was-on', text: 'Was On' },
    [DeviceState.WasOff]: { class: 'was-off', text: 'Was Off' },
  };

  private mounted() {
    // Render times nicely
    render(this.$el.querySelectorAll('.time'));
  }

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

  private get deviceState(): DeviceState {
    if (this.connectionStatue === ConnectedState.Connected) {
      return (this.powerState === PowerState.On) ? DeviceState.On : DeviceState.Off;
    }
    return (this.powerState === PowerState.On) ? DeviceState.WasOn : DeviceState.WasOff;
  }

  private get deviceStateClass(): string {
    return this.deviceStateInfo[this.deviceState].class;
  }

  private get deviceStateText(): string {
    return this.deviceStateInfo[this.deviceState].text;
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


<style lang="scss" scoped>
@import "node_modules/bootstrap/scss/bootstrap";

.status-graphic {
  position: relative;
  padding: 1em;

  .power-icon-container {
    border-radius: 50%;
    width: 10em;
    height: 10em;

    &.on, &.was-on {
      background-color: theme-color("success");
    }
    &.off, &.was-off {
      background-color: theme-color("danger");
    }
  }

  img.power-icon {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  img.connection-icon {
    position: absolute;
    right: 0;
    top: 0;
    padding: 0.8em;
    width: 3em;
    height: 3em;
    background-color: theme-color("warning");
    border-radius: 1em;
  }
}
</style>
