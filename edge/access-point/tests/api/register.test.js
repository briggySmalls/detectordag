import { AccountsApi, AuthenticationApi, Credentials, Token, DeviceRegistered, MutableDevice } from '../../lib/client';

import register from "../pages/api/register.js";

// Mock axios for us
jest.mock('client');

describe("With Enzyme", () => {
  it('Walks the happy path', () => {
    // First expect a request to auth
    axios.post.mockResolvedValueOnce()
  });
});
