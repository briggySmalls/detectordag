import Head from 'next/head'
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';

// Create some CSS styles for the page
const useStyles = makeStyles((theme) => ({
  form: {
    display: 'flex',
    flexDirection: 'column',
  },
  formItem: {
    margin: theme.spacing(1),
  },
  root: {
    padding: '5rem 0',
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  }
}));

export default function Home() {
  const classes = useStyles();
  return (
    <Container>
      <Head>
        <title>Detector Dag Device</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={classes.root}>
        <Typography variant="h2" component="h1" gutterBottom>
          Your dag is here to serve!
        </Typography>

        <p>
          Register your device to get started
        </p>

        <form className={classes.form}>
          <TextField
            id="username-input"
            label="Username"
            type="email"
            autoComplete="current-password"
            className={classes.formItem}
            required
          />
          <TextField
            id="password-input"
            label="Password"
            type="password"
            autoComplete="current-password"
            className={classes.formItem}
            required
          />
          <TextField
            id="device-name-input"
            label="Desired device name"
            type="text"
            className={classes.formItem}
            required
          />
          <Button
            variant="contained"
            type="submit"
            className={classes.formItem}
          >
            Submit
          </Button>
        </form>
      </main>
    </Container>
  )
}
