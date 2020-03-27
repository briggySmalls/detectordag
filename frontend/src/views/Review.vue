<template>
  <div class="review">
    <h1>This is the review page</h1>
    <b-button v-on:click="request" :disabled="isRefreshing">Refresh</b-button>
    <b-card-group deck>
      <DeviceComponent v-for="device in devices" :key="device.deviceId" :device="device" />
    </b-card-group>
    <b-alert variant="danger" v-if="error">
      {{ error.message }}
    </b-alert>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { AccountsApi, Device } from '../../lib/client';
import { Storage } from '../utils';
import DeviceComponent from '../components/Device.vue';

@Component({
  components: {
    DeviceComponent,
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
    console.log(`Making request for devices on account ${bundle.accountId}`);
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
    } else {
      console.log(data);
      // Display the requested devices
      this.devices = data;
    }
    this.isRefreshing = false;
  }
}
</script>
