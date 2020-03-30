<template>
  <div class="review">
    <!-- Navbar -->
    <NavbarComponent />
    <!-- Main page -->
    <h1>Review Dags</h1>
    <b-button
      class="mt-2 mb-2 d-inline-block"
      v-on:click="request"
      :disabled="isRefreshing">
      Refresh
    </b-button>
    <b-container>
      <!-- Device list -->
      <b-card-group v-if="devices" deck>
        <DeviceComponent
          v-for="device in devices"
          :key="device.deviceId"
          :device="device" />
      </b-card-group>
      <!-- Loading -->
      <b-spinner v-else></b-spinner>
    </b-container>
    <ErrorComponent :error="error" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { Device } from '../../lib/client';
import { storage } from '../utils';
import DeviceComponent from '../components/Device.vue';
import NavbarComponent from '../components/Navbar.vue';
import ErrorComponent from '../components/Error.vue';

@Component({
  components: {
    DeviceComponent,
    ErrorComponent,
    NavbarComponent,
  },
})
export default class Review extends Vue {
  private error: Error | null = null;

  public created() {
    // Ensure we get some device info
    if (this.devices === null) {
      this.request();
    }
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
    this.$clients.accounts.getDevices(`Bearer ${authBundle.token}`, authBundle.accountId, this.handleDevices);
  }

  private handleDevices(error: Error, data: Device[], response: any): any {
    // Indicate the request is finished
    if (error) {
      // Assign the error
      this.error = error;
      // Also log it
      this.$logger.debug(response.text);
      // If we have authorization issues, redirect to login
      this.$router.push('/login');
      return;
    }
    // Display the requested devices
    this.$store.commit('setDevices', data);
  }

  private get devices(): Device[] | null {
    return this.$store.state.devices;
  }

  private get isRefreshing() {
    return this.$store.state.devices === null;
  }
}
</script>
