<template>
  <Topbar :error="error" title="Settings">
    <!-- Email list -->
    <b-form v-if="emails" @submit.prevent="submit">
      <b-form-group
        description="These are the emails we'll use to notify you when you dag spots a change."
      >
        <label for="emails">Notification emails:</label>
        <b-form-tags
          v-model="emails"
          :tag-validator="emailValidator"
          no-outer-focus
          class="mb-2"
        >
          <template
            v-slot="{ tags, inputAttrs, inputHandlers, addTag, removeTag }"
          >
            <b-input-group aria-controls="my-custom-emails-list">
              <input
                v-bind="inputAttrs"
                v-on="inputHandlers"
                placeholder="New email - Press enter to add"
                class="form-control"
              />
              <b-input-group-append>
                <b-button @click="addTag()" variant="primary">Add</b-button>
              </b-input-group-append>
            </b-input-group>
            <b-list-group>
              <b-list-group-item
                v-for="email in tags"
                :key="email"
                class="d-flex justify-content-between align-items-center"
              >
                {{ email }}
                <b-button
                  @click="removeTag(email)"
                  variant="secondary"
                  size="sm"
                  :aria-controls="`my-custom-emails-email_${email.replace(
                    /\s/g,
                    '_'
                  )}_`"
                  ><b-icon-x-circle-fill></b-icon-x-circle-fill
                ></b-button>
              </b-list-group-item>
            </b-list-group>
          </template>
        </b-form-tags>
      </b-form-group>
      <b-button type="submit">Save</b-button>
    </b-form>
    <!-- Loading -->
    <b-spinner v-else-if="loading"></b-spinner>
  </Topbar>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import Topbar from '../layouts/Topbar.vue';

@Component({
  components: {
    Topbar,
  },
})
export default class AccountView extends Vue {
  // Emails to display in the form
  private emails: string[] | null = null;

  // Errors in API requests
  private error: Error | null = null;

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
    this.$clients.accounts
      .getAccount(auth.accountId, `Bearer ${auth.token}`)
      .then((response) => {
        // Save the account details to the store
        this.$logger.debug('Saving account details');
        this.$store.commit('setAccount', response.data);
      })
      .catch((err) => this.$checkUnauthorised(err, (error) => {
        // Record the error
        this.error = error;
        this.$logger.debug(`Account request error: ${error.response}`);
      }));
  }

  // The emails from the store
  private get storedEmails() {
    const { account } = this.$store.state;
    return account !== null ? account.emails : null;
  }

  // Assign emails from the store (when changed)
  @Watch('storedEmails')
  private onPropertyChanged(
    value: string[],
    _: string[], // eslint-disable-line @typescript-eslint/no-unused-vars
  ) {
    this.emails = value;
  }

  // Says if wer are loading device content
  private get loading() {
    return this.emails === null && this.error === null;
  }

  // Submit update to API
  private submit(_: Event) {
    // eslint-disable-line @typescript-eslint/no-unused-vars
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
    this.$clients.accounts
      .updateAccount(auth.accountId, `Bearer ${auth.token}`, {
        emails: this.emails,
      })
      .then((response) => {
        // Save the account details to the store
        this.$logger.debug('Saving account details');
        this.$store.commit('setAccount', response.data);
      })
      .catch((error) => {
        this.$logger.debug(`Account update error: ${error.response}`);
        // Set the error
        this.error = error;
      });
    // Indicate that our emails are updating
    this.emails = null;
  }

  public emailValidator(email: string) {
    // eslint-disable-line class-methods-use-this
    const re = /\S+@\S+\.\S+/;
    return re.test(email);
  }
}
</script>
