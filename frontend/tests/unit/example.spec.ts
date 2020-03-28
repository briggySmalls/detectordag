import { expect } from 'chai';
import { shallowMount } from '@vue/test-utils';
import ErrorComponent from '@/components/Error.vue';

describe('Error.vue', () => {
  it('renders props.error when passed', () => {
    const message = 'error message';
    const error = new Error(message);
    const wrapper = shallowMount(ErrorComponent, {
      propsData: { error },
    });
    expect(wrapper.text()).to.include(message);
  });
});
