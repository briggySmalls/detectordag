import { Request, Response } from 'superagent';

export as namespace detectordag

export class ApiClient {
  basePath: string;
}

type Callback<T> = (error: Error, data: T, response: Response) => void;

export class AuthenticationApi {
  constructor(client: ApiClient);
  auth(body: Credentials, callback: Callback<Token>): Request;
}

export class AccountsApi {
  constructor(client: ApiClient);
  getDevices(authorization: string, accountId: string, callback: Callback<Device[]>): Request;
  getAccount(authorization: string, accountId: string, callback: Callback<Account>): Request;
  updateAccount(body: Emails, authorization: string, accountId: string, callback: Callback<Account>): Request;
}

export class Credentials {
  constructor(username: string, password: string);
}

export class Token {
  token: string;
  accountId: string;
}

export class Device {
  name: string;
  deviceId: string;
  state: DeviceState;
  updated: Date;
}

export class DeviceState {
  power: boolean;
}

export class Account {
  username: string;
  emails: string[];
}

export class Emails {
  constructor(emails: string[]);
}
