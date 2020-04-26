import { makeStyles } from '@material-ui/core/styles';
import Head from 'next/head';
import Container from '@material-ui/core/Container';
import Typography from '@material-ui/core/Typography';

// Create some CSS styles for the page
const useStyles = makeStyles((theme) => ({
  root: {
    padding: '5rem 0',
    flex: 1,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
  }
}));

export default function Layout({ children }) {
    // Create classes
    const classes = useStyles();
    // Create layout
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
            { children }
          </main>
        </Container>
    )
}
