import { AppProps } from 'next/app'
import 'typeface-roboto'; // Import roboto because material-ui likes it

// This default export is required in a new `pages/_app.js` file.
export default function DetectorDag({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />
}
