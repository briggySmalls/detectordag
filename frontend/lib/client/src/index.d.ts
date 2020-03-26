export as namespace detectordag

export class ApiClient {
  instance: ApiClient;
}

export class AuthenticationApi {
  constructor();
  auth(body: Credentials, callback: (error: Error, data: Token, response: any) => any): any;
}

export class AccountsApi {
  constructor();
  getDevices(authorization: string, accountId: string, callback: (error: Error, data: Device[], response: any) => any): any;
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
  power: bool;
}
