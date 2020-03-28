<template>
  <div class="review">
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
import { Storage } from '../utils';
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

  private storage: Storage;

  private isRefreshing = false;

  public constructor() {
    // Call super
    super();
    // Create client
    this.client = new AccountsApi();
    // Create storage helper
    this.storage = new Storage();
  }

  public created() {
    // Make a request immediately
    this.request();
  }

  public request() {
    // Clear any existing devices
    this.devices = null;
    this.error = null;
    this.isRefreshing = true;
    // Fetch the token/accountId
    const { bundle } = this.storage;
    // Redirect to login if these are not present
    if (bundle == null) {
      this.$router.push('/login');
      return;
    }
    // Get the devices
    this.client.getDevices(`Bearer ${bundle.token}`, bundle.accountId, this.handleDevices);
  }

  public handleDevices(error: Error, data: Device[], response: any): any {
    if (error) {
      // Assign the error
      this.error = error;
      // Also log it
      console.error(response.text);
      // If we have authorization issues, redirect to login
      this.$router.push('/login');
      return;
    }
    // Display the requested devices
    this.devices = data;
    this.isRefreshing = false;
  }
}
</script>
