export as namespace detectordag

export class ApiClient {
}

export class AuthenticationApi {
  constructor();
  auth(body: Credentials, callback: (error: Error, data: any, response: any) => any): any;
}

export class Credentials {
  constructor(username: string, password: string);
}
