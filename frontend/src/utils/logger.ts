import winston from 'winston';
import 'setimmediate'; // Polyfill for winston browser console


export default winston.createLogger({
  level: (process.env.NODE_ENV !== 'production') ? 'debug' : 'error',
  format: winston.format.simple(),
  transports: [
    new winston.transports.Console(),
  ],
});
