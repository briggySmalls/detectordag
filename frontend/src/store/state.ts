import { Device, Account } from '../../lib/client';

class State {
  account: Account | null = null;

  devices: Device[] | null = null;
}

export default State;
