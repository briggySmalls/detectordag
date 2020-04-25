import { Component, useState } from 'react';
import Head from 'next/head'
import WithLoading from '../components/WithLoading';
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

// Create our form (we will wrap this)
const Form = ({ onSubmit }) => {
  // Create styles for use
  const classes = useStyles();
  // Render the form
  return (
    <form className={classes.form} onSubmit={onSubmit}>
      <TextField
        id="username-input"
        name="username"
        label="Username"
        type="email"
        autoComplete="current-password"
        className={classes.formItem}
        required
      />
      <TextField
        id="password-input"
        name="password"
        label="Password"
        type="password"
        autoComplete="current-password"
        className={classes.formItem}
        required
      />
      <TextField
        id="device-name-input"
        name="device-name"
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
  )
}

// Wrap the form in a loader
const FormWithLoading = WithLoading(Form);

// Create the actual homepage
function Home() {
  // Create styles for use
  const classes = useStyles();
  // Declare isLoading state
  const [isLoading, setIsLoading] = useState(false);
  // Callback for submit event
  const handleSubmit = (event) => {
    // Don't actually submit
    event.preventDefault();
    // Get the form data
    const formData = new FormData(event.target);
    let data = {};
    formData.forEach((value, key) => {data[key] = value});
    // Submit the form data
    console.log(data);
    // axios({
    //   method: 'post',
    //   url: '/api/register',
    //   data: data,
    // });
    // Update the state
    setIsLoading(false);
  }
  // Render the component
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

        <FormWithLoading isLoading={isLoading} onSubmit={handleSubmit} />
      </main>
    </Container>
  )
}

export default Home;
