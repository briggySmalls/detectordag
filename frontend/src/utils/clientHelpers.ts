import { Response } from 'superagent';
import logger from './logger';
import store from '../store';
import storage from './Storage';
import router from '../router/router';
import { Account } from '../../lib/client';

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

export default {
  handleAccountResponse,
};
