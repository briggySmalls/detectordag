export as namespace detectordag

export class ApiClient {
  instance: ApiClient;
}

export class AuthenticationApi {
  constructor();
  auth(body: Credentials, callback: (error: Error, data: Token, response: any) => any): any;
}

export class Credentials {
  constructor(username: string, password: string);
}

export class Token {
  token: string;
  accountId: string;
}
