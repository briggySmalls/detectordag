<template>
  <Topbar title="Your Dags" :error="error" >
    <b-button
      class="mt-2 mb-2 d-inline-block"
      v-on:click="request"
      :disabled="loading">
      Refresh
    </b-button>
    <div>
      <!-- Device list -->
      <b-card-group v-if="devices" deck>
        <DeviceComponent
          v-for="device in devices"
          :key="device.deviceId"
          :device="device" />
      </b-card-group>
      <!-- Loading -->
      <b-spinner v-if="loading"></b-spinner>
    </div>
  </Topbar>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Device } from '../../lib/client';
import { storage } from '../utils';
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
    return (this.devices === null) && (this.error === null);
  }

  private get devices(): Device[] | null {
    return this.$store.state.devices;
  }

  private request() {
    // Clear any existing devices
    this.$store.commit('clearDevices');
    this.error = null;
    // Fetch the token/accountId
    const authBundle = storage.bundle;
    // Redirect to login if these are not present
    if (authBundle == null) {
      this.$logger.debug('Token not available');
      this.$router.push('/login');
      return;
    }
    // Get the devices
    this.$logger.debug('Requesting account\'s devices');
    this.$clients.accounts.getDevices(authBundle.accountId, `Bearer ${authBundle.token}`)
      .then((request) => {
        // Display the requested devices
        this.$store.commit('setDevices', request.data);
      })
      .catch((error) => {
        // Assign the error
        this.error = error;
        // Also log it
        this.$logger.debug(error.response);
        // If we have authorization issues, redirect to login
        this.$router.push('/login');
      });
  }
}
</script>
