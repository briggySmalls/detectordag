import logger from './logger';
import store from '../store';
import storage from './storage';
import router from '../router';
import { Response } from 'superagent';
import { Account, Response } from '../../lib/client';


// Save the account details to the store
function handleAccountResponse(error: Error, data: Account, response: Response) {
  // Handle errors
  if (error) {
    logger.debug(`Account request error: ${response.text}`);
    // Clear the token (we're assuming that's why we failed)
    storage.clear();
    // Get the user to reauthenticate
    router.push('/login');
    return;
  }
  // Save the account details to the store
  logger.debug('Saving account details');
  store.commit('setAccount', data);
}

export {
  handleAccountResponse,
}
