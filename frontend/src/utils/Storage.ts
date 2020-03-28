import AuthBundle from './AuthBundle';

// Helper class for managing local storage
class Storage {
  readonly fieldName = 'authBundle';

  public save(auth: AuthBundle) {
    localStorage.setItem(this.fieldName, JSON.stringify(auth));
  }

  public clear() {
    localStorage.remove(this.fieldName);
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

// Export our storage instance
export default new Storage();
