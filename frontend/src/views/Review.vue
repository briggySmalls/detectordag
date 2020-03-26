<template>
  <div class="review">
    <h1>This is the review page</h1>
    <div class="device" v-for="device in devices" :key="device.deviceId">
      <h2>Device {{ device.name }} - ({{ device.deviceId }})</h2>
      <h3>Updated {{ device.updated }}</h3>
      <ul>
        <li v-for="(value, key) in device.state" :key="key">
          {{ key }}: {{ value }}
        </li>
      </ul>
    </div>
    <div v-if="error">
      {{ error.message }}
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import { AccountsApi, Device } from '../../lib/client';

@Component
export default class Review extends Vue {
  private error: Error | null = null;

  private devices: Device[] | null = null;

  private client: AccountsApi;

  public constructor() {
    // Call super
    super();
    // Create client
    this.client = new AccountsApi();
  }

  public created() {
    // Fetch the token/accountId
    const token = localStorage.getItem('token');
    const accountId = localStorage.getItem('accountId');
    // Get the devices
    console.log(`Making request for devices on account ${accountId}`);
    this.client.getDevices(token, accountId);
  }

  public handleDevices(error: Error, data: Device[], response: any): any {
    if (error) {
      // Assign the error
      this.error = error;
      // Also log it
      console.error(error);
    } else {
      console.log(`API called successfully. Returned data: ${data}`);
      // Record the token and account
    }
  }
}
</script>
