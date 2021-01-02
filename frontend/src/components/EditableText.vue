<template>
  <div>
    <b-form
      v-if="isEditing"
      inline
      class="justify-content-center"
    >
      <b-input-group>
        <label
          class="sr-only"
          for="inline-form-input-name"
        >Device name</label>
        <b-form-input
          id="inline-form-input-name"
          v-model="value"
          class="mb-2 mr-sm-2 mb-sm-0"
          placeholder="Device name"
          type="text"
        />
        <b-input-group-append>
          <b-button
            variant="primary"
            @click="submit"
          >
            Set
          </b-button>
        </b-input-group-append>
      </b-input-group>
    </b-form>
    <template v-else>
      <div class="d-inline-block"><slot /></div>
      <div class="d-inline-block"><b-icon-pencil-square class="edit-icon" @click="edit()" /></div>
    </template>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component
export default class EditableText extends Vue {
  private isEditing = false;

  private value = '';

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
    this.$emit('edited', this.value);
  }
}
</script>

<style lang="scss" scoped>
.edit-icon {
  margin-left: 1em;
  cursor: pointer;
}
</style>
