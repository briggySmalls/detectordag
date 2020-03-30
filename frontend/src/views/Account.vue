<template>
  <div class="review">
    <!-- Navbar -->
    <NavbarComponent />
    <!-- Main page -->
    <h1>Account Details</h1>
    <b-container>
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
    </b-container>
    <ErrorComponent :error="error" />
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import NavbarComponent from '../components/Navbar.vue';
import ErrorComponent from '../components/Error.vue';
import { Emails, Account } from '../../lib/client';
import { storage } from '../utils';
import { handleAccountResponse } from '../utils/clientHelpers';

@Component({
  components: {
    ErrorComponent,
    NavbarComponent,
  },
})
export default class Review extends Vue {
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
  private onPropertyChanged(value: string[], oldvalue: string[]) {
    this.emails = value;
  }

  // Submit update to API
  private submit(event: Event) {
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
    this.$clients.accounts.updateAccount(
      new Emails(this.emails), `Bearer ${auth.token}`, auth.accountId, handleAccountResponse,
    );
    // Indicate that our emails are updating
    this.emails = null;
  }
}
</script>
