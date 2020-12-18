<template>
  <Topbar :error="error" title="Settings">
    <!-- Email list -->
    <b-form v-if="emails" @submit.prevent="submit">
      <b-form-group>
        <label for="tags">Notification emails:</label>
        <b-form-tags
          input-id="tags"
          v-model="emails"
          placeholder="Add email"
          name=""
          class="mb-2">
        </b-form-tags>
      </b-form-group>
      <b-button type="submit">Submit</b-button>
    </b-form>
    <!-- Loading -->
    <b-spinner v-else></b-spinner>
  </Topbar>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import Topbar from '../layouts/Topbar.vue';
import { storage } from '../utils';
import requestAccount from '../utils/clientHelpers';

@Component({
  components: {
    Topbar,
  },
})
export default class AccountView extends Vue {
  // Emails to display in the form
  private emails: string[] | null = null;

  public created() {
    this.emails = this.storedEmails;
  }

  // The emails from the store
  private get storedEmails() {
    const { account } = this.$store.state;
    return (account !== null) ? account.emails : null;
  }

  // Assign emails from the store (when changed)
  @Watch('storedEmails')
  private onPropertyChanged(
    value: string[], _: string[], // eslint-disable-line @typescript-eslint/no-unused-vars
  ) {
    this.emails = value;
  }

  // Submit update to API
  private submit(_: Event) { // eslint-disable-line @typescript-eslint/no-unused-vars
    this.$logger.debug('Emails submitted');
    // Get auth token
    const auth = storage.bundle;
    // Redirect to login if these are not present
    if (auth == null) {
      this.$logger.debug('Token not available');
      this.$router.push('/login');
      return;
    }
    // Submit account update
    this.$logger.debug(`Submitting emails for update: ${this.emails}`);
    if (this.emails == null) {
      // TODO: Handle this error case better
      return;
    }
    // Request the account
    requestAccount(this.$router, auth);
    // Indicate that our emails are updating
    this.emails = null;
  }
}
</script>
