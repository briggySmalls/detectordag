<template>
  <b-card
    no-body
    border-variant="dark"
    header-border-variant="dark"
  >
    <template #header>
      <b-spinner v-if="isLoading" />
      <EditableText
        v-else
        @edited="updateDeviceName"
      >
        <h4 class="mb-0">
          {{ device.name }}
        </h4>
      </EditableText>
    </template>

    <b-card-body>
      {{/* Define a graphic to illustrate the status */}}
      <div class="status-graphic d-inline-block">
        <div
          class="power-icon-container"
          :class="[deviceStateClass]"
        >
          <img
            v-if="powerState === powerStateEnum.Off"
            alt="Power off"
            src="../assets/power-off.svg"
            class="power-icon"
          >
          <img
            v-else-if="powerState === powerStateEnum.On"
            alt="Power on"
            src="../assets/power-on.svg"
            class="power-icon"
          >
        </div>
        <img
          v-if="connectionState === connectionStateEnum.Disconnected"
          alt="Disconnected"
          src="../assets/no-signal.svg"
          class="connection-icon"
        >
      </div>
      {{/* Add textual descriptions of the status */}}
      <b-card-title class="mt-4">
        {{ deviceStateInfo[deviceState].title }}
      </b-card-title>
      <b-card-text>{{ deviceStateInfo[deviceState].description }}</b-card-text>
    </b-card-body>

    <b-list-group flush>
      <b-list-group-item v-if="powerState === powerStateEnum.Off">
        <!-- Add further detail about losing power -->
        Lost power
        <span
          ref="stateUpdatedTime"
          :datetime="device.state.updated"
        >
          {{ device.state.updated }}
        </span>
      </b-list-group-item>
      <b-list-group-item v-if="connectionState === connectionStateEnum.Disconnected">
        <!-- Add further detail about the dag losing connection -->
        Lost connection
        <span
          ref="connectionUpdatedTime"
          :datetime="device.connection.updated"
        >
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
import EditableText from './EditableText.vue';

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

enum ConnectionState {
  Connected = 1,
  Disconnected,
}

interface StateData {
  class: string;
  title: string;
  description: string;
}

@Component({
  components: {
    EditableText,
  },
})
export default class Device extends Vue {
  // Declare some enums so we can use them in the template
  private readonly deviceStateEnum: typeof DeviceState = DeviceState;

  private readonly powerStateEnum: typeof PowerState = PowerState;

  private readonly connectionStateEnum: typeof ConnectionState = ConnectionState;

  @Prop() private device!: DeviceModel;

  private isLoading = false;

  private readonly deviceStateInfo: Record<DeviceState, StateData> = {
    [DeviceState.On]: {
      class: 'on',
      title: 'On',
      description: 'The power is on!',
    },
    [DeviceState.Off]: {
      class: 'off',
      title: 'Off',
      description: 'Your dag says that the power is off',
    },
    [DeviceState.WasOn]: {
      class: 'was-on',
      title: 'Was On',
      description: "We've lost contact with your dag. The power was on the last we heard...",
    },
    [DeviceState.WasOff]: {
      class: 'was-off',
      title: 'Was Off',
      description:
        'Your dag noticed the power go, and then we lost contact. It may have run out of battery.',
    },
  };

  private mounted() {
    // Gether together the time elements
    // Note: some may be hidden (e.g. on and connected device)
    const els = [
      this.$refs.stateUpdatedTime as HTMLElement,
      this.$refs.connectionUpdatedTime as HTMLElement,
    ].filter((x) => x !== undefined);
    // Render them
    if (els.length) {
      render(els);
    }
  }

  private updateDeviceName(name: string) {
    this.$logger.debug(`need to update to ${name}`);
    // Submit a request to set the device name
    const auth = this.$storage.bundle;
    // Redirect to login if these are not present
    if (auth == null) {
      this.$logger.debug('Token not available');
      this.$router.push('/login');
      return;
    }
    this.$clients.devices.updateDevice(`Bearer ${auth.token}`, this.device.deviceId, { name })
      .then((response) => {
        // Submit the new device info to the store
        this.$store.commit('setDevice', response.data);
        // We've bound to this so it will update automatically!
        this.isLoading = false;
      })
      .catch((error) => {
        this.$logger.debug(`Ah bugger: ${error}`);
      });
    // Indicate we are saving the name
    this.isLoading = true;
  }

  private get deviceStatus(): string {
    switch (this.deviceState) {
      case DeviceState.On:
        return 'On';
      case DeviceState.Off:
        return 'Off';
      case DeviceState.WasOn:
        return 'Was On';
      case DeviceState.WasOff:
        return 'Was Off';
      default:
        throw new Error(`Unexpected state: ${this.deviceState}`);
    }
  }

  private get deviceState(): DeviceState {
    if (this.connectionState === ConnectionState.Connected) {
      return this.powerState === PowerState.On ? DeviceState.On : DeviceState.Off;
    }
    return this.powerState === PowerState.On ? DeviceState.WasOn : DeviceState.WasOff;
  }

  private get deviceStateClass(): string {
    return this.deviceStateInfo[this.deviceState].class;
  }

  private get deviceStateText(): string {
    return this.deviceStateInfo[this.deviceState].title;
  }

  private get connectionState(): ConnectionState {
    switch (this.device.connection.status) {
      case 'connected':
        return ConnectionState.Connected;
      case 'disconnected':
        return ConnectionState.Disconnected;
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
@import '~bootstrap/scss/_functions.scss';
@import '~bootstrap/scss/_variables.scss';

.status-graphic {
  position: relative;
  padding: 0em;

  .power-icon-container {
    border-radius: 50%;
    width: 10em;
    height: 10em;

    &.on,
    &.was-on {
      background-color: theme-color('success');
    }
    &.off,
    &.was-off {
      background-color: theme-color('danger');
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

    border-radius: 1em;
  }
}
</style>
