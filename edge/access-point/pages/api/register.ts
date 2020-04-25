import { NextApiRequest, NextApiResponse } from 'next'
import { ApiClient, AccountsApi, AuthenticationApi, Credentials, Token, DeviceRegistered, MutableDevice } from '../../lib/client';
import pify from 'pify';
import fs from 'fs';
import path from 'path';


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

function writeFile(filepath: string, content: string): Promise<void> {
  // Ensure the directory is present
  return fs.promises
    .mkdir(path.dirname(filepath), { recursive: true })
    .then(() => fs.writeFile(filepath));
}

// Handle form submission
export default async (req: NextApiRequest, res: NextApiResponse) => {
  // // Pull out the results
  const formData = req.body;
  // Authenticate
  let token: Token = null;
  try {
    const [data, result] = await pify(wrapper.authentication.auth.bind(wrapper.authentication), {multiArgs: true})(
        new Credentials(formData.username, formData.password));
    console.log(result.text)
    token = data;
  } catch (err) {
    handleError(res, err);
    return
  }
  // Register new device
  let registered: DeviceRegistered = null;
  try {
    const [data, result] = await pify(wrapper.accounts.registerDevice.bind(wrapper.accounts), {multiArgs: true})(
      new MutableDevice(formData.deviceName),
      `Bearer ${token.token}`,
      process.env.BALENA_DEVICE_UUID,
      token.accountId);
    console.log(result.text);
    registered = data;
  } catch (err) {
    handleError(res, err);
    return
  }
  // Save certificates to files
  const certsPath = process.env.CERTS_DIR
  const publicPromise = writeFile(path.join(certsPath, 'thing.public.pem'), registered.certificate.publicKey);
  const privatePromise = writeFile(path.join(certsPath, 'thing.private.key'), registered.certificate.privateKey);
  await Promise.all([publicPromise, privatePromise]);
  // Return our success
  res.setHeader('Content-Type', 'application/json');
  res.status(200).json(registered);
}
