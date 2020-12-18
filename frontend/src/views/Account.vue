<template>
  <Topbar :error="error" title="Settings">
    <!-- Email list -->
    <b-form v-if="emails" @submit.prevent="submit">
      <b-form-group
        description="These are the emails we'll use to notify you when you dag spots a change.">
        <label for="emails">Notification emails:</label>
        <b-form-tags v-model="emails" no-outer-focus class="mb-2">
          <template v-slot="{ tags, inputAttrs, inputHandlers, addTag, removeTag }">
            <b-input-group aria-controls="my-custom-tags-list">
              <input
                v-bind="inputAttrs"
                v-on="inputHandlers"
                placeholder="New tag - Press enter to add"
                class="form-control">
              <b-input-group-append>
                <b-button @click="addTag()" variant="primary">Add</b-button>
              </b-input-group-append>
            </b-input-group>
            <ul
              id="my-custom-tags-list"
              class="list-unstyled d-inline-flex flex-wrap mb-0"
              aria-live="polite"
              aria-atomic="false"
              aria-relevant="additions removals"
            >
              <!-- Always use the tag value as the :key, not the index! -->
              <!-- Otherwise screen readers will not read the tag
                  additions and removals correctly -->
              <b-card
                v-for="tag in tags"
                :key="tag"
                :id="`my-custom-tags-tag_${tag.replace(/\s/g, '_')}_`"
                tag="li"
                class="mt-1 mr-1"
                body-class="py-1 pr-2 text-nowrap"
              >
                <strong>{{ tag }}</strong>
                <b-button
                  @click="removeTag(tag)"
                  variant="link"
                  size="sm"
                  :aria-controls="`my-custom-tags-tag_${tag.replace(/\s/g, '_')}_`"
                >remove</b-button>
              </b-card>
            </ul>
          </template>
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
    // Check if we already have the account info
    this.emails = null;
    if (this.storedEmails !== null) {
      // Just copy them over then
      this.emails = this.storedEmails;
      return;
    }
    // Check we have a valid login
    const auth = this.$storage.bundle;
    if (auth == null) {
      this.$logger.debug('Token not available');
      this.$router.push('/login');
      return;
    }
    // Request the accounts to render them
    this.$clients.accounts.getAccount(auth.accountId, `Bearer ${auth.token}`)
      .then((response) => {
      // Save the account details to the store
        this.$logger.debug('Saving account details');
        this.$store.commit('setAccount', response.data);
      })
      .catch((error) => {
        this.$logger.debug(`Account request error: ${error.response}`);
        // Clear the token (we're assuming that's why we failed)
        this.$storage.clear();
        // Get the user to reauthenticate
        this.$router.push('/login');
      });
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
    const auth = this.$storage.bundle;
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
    this.$clients.accounts.updateAccount(auth.accountId, `Bearer ${auth.token}`, this.emails);
    // Indicate that our emails are updating
    this.emails = null;
  }
}
</script>
