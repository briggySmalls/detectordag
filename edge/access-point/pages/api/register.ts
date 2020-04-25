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

// Handle form submission
export default async (req: NextApiRequest, res: NextApiResponse) => {
  // Pull out the results
  const formData = req.body;
  // Authenticate
  try {
    let result = await util.promisify(wrapper.authentication.auth.bind(wrapper.authentication))(new Credentials(formData.username, formData.password));
    const token = JSON.parse(result.text);
    res.setHeader('Content-Type', 'application/json');
    res.status(200).json(result);
  } catch (err) {
    console.log(err);
    res.setHeader('Content-Type', 'application/json');
    res.status(500).json({error: err});
  }
  // let result = await client.accounts.registerDevice(new , `Bearer ${auth.token}`, , auth.accountId);
  // const registered = JSON.parse(result.text);

  // // Create an async executor
  // const run = async () => {
  //   // First authorise
  //   // Now save the credentials to file
  //   console.log(token);
  // };

  // Run the executor
  // const promise = new Promise(run).then((result) => {
  //   res.setHeader('Content-Type', 'application/json');
  //   res.status(200).json(result);
  // }).catch((error) => {
  //   res.setHeader('Content-Type', 'application/json');
  //   res.status(500).json(error);
  // });
}
