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
          class="mb-2 mr-sm-2 mb-sm-0"
          type="text"
          :placeholder="value"
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
      <div class="d-inline-block"><slot>{{ value }}</slot></div>
      <div class="d-inline-block"><b-icon-pencil-square class="edit-icon" @click="edit()" /></div>
    </template>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Vue } from 'vue-property-decorator';

@Component
export default class EditableText extends Vue {
  private isEditing = false;

  @Prop() private value!: string;

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
