<template>
  <div class="review">
    <b-navbar class="justify-content-center">
      <b-navbar-brand href="#">
        <img id="logo" alt="Detectordag logo" src="../assets/logo.svg"
             class="d-inline-block">
        Detectordag
      </b-navbar-brand>
      <div>
        {{ username }}
      </div>
    </b-navbar>
    <h1>Review dags</h1>
    <b-button class="mt-2 mb-2" v-on:click="request" :disabled="isRefreshing">Refresh</b-button>
    <b-card-group deck>
      <DeviceComponent v-for="device in devices" :key="device.deviceId" :device="device" />
    </b-card-group>
    <ErrorComponent :error="error" />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { AccountsApi, Device } from '../../lib/client';
import { storage } from '../utils';
import DeviceComponent from '../components/Device.vue';
import ErrorComponent from '../components/Error.vue';

@Component({
  components: {
    DeviceComponent,
    ErrorComponent,
  },
})
export default class Review extends Vue {
  private error: Error | null = null;

  private devices: Device[] | null = null;

  private client: AccountsApi;

  private isRefreshing = false;

  public constructor() {
    // Call super
    super();
    // Create client
    this.client = new AccountsApi();
  }

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
    const { bundle } = this.storage;
    // Redirect to login if these are not present
    if (bundle == null) {
      this.$logger.debug('Token not available');
      this.$router.push('/login');
      return;
    }
    // Get the devices
    this.$logger.debug('Requesting account\'s devices');
    this.client.getDevices(`Bearer ${bundle.token}`, bundle.accountId, this.handleDevices);
  }

  private handleDevices(error: Error, data: Device[], response: any): any {
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
    this.isRefreshing = false;
  }

  private get username() {
    const account = this.$store.state.account;
    return (account) ? account.username : '?';
  }
}
</script>
