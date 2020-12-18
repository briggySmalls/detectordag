import Vue from 'vue';
import AuthBundle from './AuthBundle';

// Helper class for managing local storage
class Storage {
  readonly fieldName = 'authBundle';

  public save(auth: AuthBundle) {
    localStorage.setItem(this.fieldName, JSON.stringify(auth));
  }

  public clear() {
    localStorage.removeItem(this.fieldName);
  }

  public get bundle(): AuthBundle | null {
    // Try to get the bundle
    const bundle = localStorage.getItem(this.fieldName);
    if (!bundle) {
      return null;
    }
    return JSON.parse(bundle);
  }
}

// Create a singleton storage instance
const storage = new Storage();

// Add the clients to Vue
Vue.prototype.$storage = storage;

// Update the type hinting for all Vue instances
declare module 'vue/types/vue' {
  interface Vue {
    $storage: Storage;
  }
}

// Export our storage instance
export default storage;
