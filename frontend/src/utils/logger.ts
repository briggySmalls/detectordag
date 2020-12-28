import Vue from 'vue';
import winston from 'winston';
import 'setimmediate'; // Polyfill for winston browser console

// Create the logger
const logger = winston.createLogger({
  level: process.env.NODE_ENV !== 'production' ? 'debug' : 'error',
  format: winston.format.simple(),
  transports: [new winston.transports.Console()],
});

// Add the logger to Vue
// Add a logger to all vue instances
Vue.prototype.$logger = logger;

// Update the type hinting for all Vue instances
declare module 'vue/types/vue' {
  interface Vue {
    $logger: winston.Logger;
  }
}

// Export the logger if people want to use it directly
export default logger;
