<template>
  <div>
    <b-form
      v-if="isEditing"
      inline
      class="justify-content-center"
    >
      <label
        class="sr-only"
        for="inline-form-input-name"
      >Device name</label>
      <b-input-group>
        <b-form-input
          id="inline-form-input-name"
          v-model="input"
          type="text"
          :placeholder="value"
        />
        <b-input-group-append>
          <b-button
            type="submit"
            variant="primary"
            @click="submit"
          >
            Save
          </b-button>
          <b-button
            type="cancel"
            variant="secondary"
            @click="cancel"
          >
            Cancel
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </b-form>
    <template v-else>
      <div class="d-inline-block">
        <slot>{{ value }}</slot>
      </div>
      <div class="d-inline-block">
        <b-icon-pencil-square
          class="edit-icon"
          @click="edit()"
        />
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class EditableText extends Vue {
  private isEditing = false;

  @Prop() private value!: string;

  private input = '';

  private edit() {
    if (this.isEditing) {
      return;
    }
    this.$logger.debug('triggered!');
    this.isEditing = true;
  }

  private submit() {
    this.$logger.debug('submitted');
    this.isEditing = false;
    // Get the current value
    this.$emit('edited', this.input);
  }

  private cancel() {
    this.$logger.debug('cancelled');
    // Clear anything typed
    this.input = '';
    this.isEditing = false;
  }
}
</script>

<style lang="scss" scoped>
.edit-icon {
  margin-left: 1em;
  cursor: pointer;
}
</style>
