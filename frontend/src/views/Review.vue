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

  private devices: Device[] | null = null;

  private isRefreshing = false;

  public created() {
    // Make a request upon landing on the page
    this.request();
  }

  private request() {
    // Clear any existing devices
    this.devices = null;
    this.error = null;
    this.isRefreshing = true;
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

  private logout() {
    // Clear the token and account
    storage.clear();
    this.$store.commit('clearAccount');
    // Redirect to the login page
    this.$router.push('/login');
  }

  private handleDevices(error: Error, data: Device[], response: any): any {
    // Indicate the request is finished
    this.isRefreshing = false;
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
    this.devices = data;
  }
}
</script>

<style lang="scss" scoped>
#logo {
  width: 5em;
  height: 5em;
}
</style>
