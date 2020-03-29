import { Request, Response } from 'superagent';

export as namespace detectordag

export class ApiClient {
  basePath: string;
}

export class AuthenticationApi {
  constructor(client: ApiClient);
  auth(body: Credentials, callback: (error: Error, data: Token, response: Response) => void): Request;
}

export class AccountsApi {
  constructor(client: ApiClient);
  getDevices(authorization: string, accountId: string, callback: (error: Error, data: Device[], response: Response) => void): Request;
  getAccount(authorization: string, accountId: string, callback: (error: Error, data: Account, response: Response) => void): Request;
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
