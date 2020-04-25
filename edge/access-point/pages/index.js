import Head from 'next/head'
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';

export default function Home() {
  return (
    <Container>
      <Head>
        <title>Detector Dag Device</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <Typography variant="h2" component="h1" gutterBottom>
          Your dag is here to serve!
        </Typography>

        <p className="description">
          Register your device to get started
        </p>

        <form>
          <TextField
            id="username-input"
            label="Username"
            type="email"
            autoComplete="current-password"
            required
          />
          <TextField
            id="password-input"
            label="Password"
            type="password"
            autoComplete="current-password"
            required
          />
          <TextField
            id="device-name-input"
            label="Desired device name"
            type="text"
            required
          />
        </form>
      </main>

      <style jsx>{`
        main {
          padding: 5rem 0;
          flex: 1;
          display: flex;
          flex-direction: column;
          justify-content: center;
          align-items: center;
        }
      `}</style>

      <style jsx global>{`
        html,
        body {
          padding: 0;
          margin: 0;
          font-family: -apple-system, BlinkMacSystemFont, Segoe UI, Roboto,
            Oxygen, Ubuntu, Cantarell, Fira Sans, Droid Sans, Helvetica Neue,
            sans-serif;
        }

        * {
          box-sizing: border-box;
        }
      `}</style>
    </Container>
  )
}
