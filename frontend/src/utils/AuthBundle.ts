export default class AuthBundle {
  readonly accountId: string;

  readonly token: string;

  constructor(accountId: string, token: string) {
    this.accountId = accountId;
    this.token = token;
  }
}
