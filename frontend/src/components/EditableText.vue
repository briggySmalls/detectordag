<template>
    <div>
        <b-form v-if="isEditing" inline>
            <b-input-group>
                <label class="sr-only" for="inline-form-input-name">Device name</label>
                <b-form-input
                    id="inline-form-input-name"
                    class="mb-2 mr-sm-2 mb-sm-0"
                    placeholder="Device name"
                    type="text"
                    v-model="value"
                ></b-form-input>
                <b-input-group-append>
                    <b-button
                        @click="submit"
                        variant="primary">
                        Set
                    </b-button>
                </b-input-group-append>
            </b-input-group>
        </b-form>
        <template v-else>
            <slot/>
            <b-icon-x-circle-fill @click="edit()"/>
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
