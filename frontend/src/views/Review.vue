<template>
  <Topbar
    title="Your Dags"
    :error="error"
  >
    <template #header>
      <h1 class="my-5">
        Your Dags
        <b-button
          class="ml-2"
          :disabled="loading"
          @click="request"
        >
          Refresh
        </b-button>
      </h1>
    </template>
    <div>
      <!-- Device list -->
      <b-card-group
        v-if="devices"
        deck
      >
        <template v-for="(device, index) in devices">
          <!-- wrap every 2 on sm -->
          <div
            v-if="index % 2 === 0"
            :key="device.deviceId"
            class="w-100 d-none d-sm-block d-md-none"
          />
          <!-- wrap every 3 on md -->
          <div
            v-if="index % 3 === 0"
            :key="device.deviceId"
            class="w-100 d-none d-md-block d-lg-none"
          />
          <!-- wrap every 4 on lg -->
          <div
            v-if="index % 4 === 0"
            :key="device.deviceId"
            class="w-100 d-none d-lg-block d-xl-none"
          />
          <DeviceComponent
            :key="device.deviceId"
            :device="device"
            @errored="errorHandler"
          />
        </template>
      </b-card-group>
      <!-- Loading -->
      <b-spinner v-if="loading" />
    </div>
  </Topbar>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Device } from '../../lib/client';
import DeviceComponent from '../components/Device.vue';
import Topbar from '../layouts/Topbar.vue';

@Component({
  components: {
    DeviceComponent,
    Topbar,
  },
})
export default class Review extends Vue {
  // Errors in API requests
  private error: Error | null = null;

  public created() {
    // Ensure we get some device info
    if (this.devices === null) {
      this.request();
    }
  }

  // Says if wer are loading device content
  private get loading() {
    return this.devices === null && this.error === null;
  }

  private get devices(): Device[] | null {
    return this.$store.state.devices;
  }

  private request() {
    // Clear any existing devices
    this.$store.commit('clearDevices');
    this.error = null;
    // Fetch the token/accountId
    const authBundle = this.$storage.bundle;
    // Redirect to login if these are not present
    if (authBundle == null) {
      this.$logger.debug('Token not available');
      this.$router.push('/login');
      return;
    }
    // Get the devices
    this.$logger.debug("Requesting account's devices");
    this.$clients.accounts
      .getDevices(authBundle.accountId, `Bearer ${authBundle.token}`)
      .then((request) => {
        // Display the requested devices
        this.$store.commit('setDevices', request.data);
      })
      .catch((err) => this.$checkUnauthorised(err, (error) => {
        // Assign the error
        this.error = error;
        this.$logger.debug(`Devices request error: ${error.response}`);
      }));
  }

  private errorHandler(error: Error) {
    this.error = error;
  }
}
</script>
