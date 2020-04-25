import { NextApiRequest, NextApiResponse } from 'next'
import { ApiClient, AccountsApi, AuthenticationApi, Credentials, Token, DeviceRegistered } from '../../lib/client';
import util from 'util';


// Define a wrapper for the different clients
class ClientWrapper {
  public readonly accounts: AccountsApi;

  public readonly authentication: AuthenticationApi;

  public constructor(client: ApiClient) {
    this.accounts = new AccountsApi(client);
    this.authentication = new AuthenticationApi(client);
  }
}

// Create the underlying client
const client = new ApiClient();
client.basePath = `${process.env.API_BASEPATH}/v1`;

// Create an instance of our wrapper
const wrapper = new ClientWrapper(client);

function handleError(res: NextApiResponse, err: Error) {
  console.log(err);
  res.setHeader('Content-Type', 'application/json');
  res.status(500).json({error: err.message});
}

// Handle form submission
export default async (req: NextApiRequest, res: NextApiResponse) => {
  // Pull out the results
  const formData = req.body;
  // Authenticate
  let token: Token = null;
  try {
    let result = await util.promisify((callback) => {
      wrapper.authentication.auth(
        new Credentials(formData.username, formData.password),
        (err, ...results) => callback(err, results));
    })()
    console.log(result.text)
    token = JSON.parse(result.text);
  } catch (err) {
    handleError(res, err);
    return
  }
  // Register new device
  let registered: DeviceRegistered = null;
  try {
    let result = await util.promisify((callback) => {
      wrapper.accounts.registerDevice(
        `Bearer ${token.token}`,
        process.env.BALENA_DEVICE_UUID,
        token.accountId,
        callback);
    })();
    console.log(result.text);
    registered = JSON.parse(result.text);
  } catch (err) {
    handleError(res, err);
    return
  }
  // Return our success
  res.setHeader('Content-Type', 'application/json');
  res.status(200).json(registered);
}
