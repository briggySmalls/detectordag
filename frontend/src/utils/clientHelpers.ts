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

// Request account
function requestAccount(router: VueRouter, auth: AuthBundle) {
  logger.debug('Requesting account details');
  clients.accounts.getAccount(auth.accountId, `Bearer ${auth.token}`)
    .then((response) => {
      // Save the account details to the store
      logger.debug('Saving account details');
      store.commit('setAccount', response.data);
    })
    .catch((error) => {
      logger.debug(`Account request error: ${error.text}`);
      // Clear the token (we're assuming that's why we failed)
      storage.clear();
      // Get the user to reauthenticate
      router.push('/login');
    });
}

export default requestAccount;
