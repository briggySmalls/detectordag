import VueRouter from 'vue-router';
import { Response } from 'superagent';
import logger from './logger';
import store from '../store';
import storage from './Storage';
import AuthBundle from './AuthBundle';
import clients from './clients';
import { Account } from '../../lib/client';

// Handy type to represent an Account API response
type handleAccountResponse = (error: Error, data: Account, response: Response) => void;

// Save the account details to the store
function handleAccountResponseFactory(router: VueRouter): handleAccountResponse {
  return (error: Error, data: Account, response: Response) => {
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
  };
}

// Request account
function requestAccount(router: VueRouter, auth: AuthBundle) {
  logger.debug('Requesting account details');
  clients.accounts.getAccount(`Bearer ${auth.token}`, auth.accountId, handleAccountResponseFactory(router));
}

export {
  requestAccount,
  handleAccountResponseFactory,
};
