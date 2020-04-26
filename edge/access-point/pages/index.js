import { Component, useState } from 'react';
import axios from 'axios';
import WithLoading from '../components/WithLoading';
import Layout from '../components/layout';
import { makeStyles } from '@material-ui/core/styles';
import { useRouter } from 'next/router'
import Alert from '@material-ui/lab/Alert';
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
}));

// Create our form (we will wrap this)
const Form = ({ onSubmit }) => {
  // Create styles for use
  const classes = useStyles();
  // Get the nextjs router
  const router = useRouter()
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
        name="deviceName"
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
export default function Home() {
  // Create styles for use
  const classes = useStyles();
  // Declare isLoading state
  const [isLoading, setIsLoading] = useState(false);
  // Declare isLoading state
  const [error, setError] = useState(null);
  // Callback for submit event
  const handleSubmit = (event) => {
    // Don't actually submit
    event.preventDefault();
    // Get the form data
    const formData = new FormData(event.target);
    let data = {};
    formData.forEach((value, key) => {data[key] = value});
    // Clear any errors
    setError(null);
    // Submit the form data
    axios({
      method: 'post',
      url: '/api/register',
      data: data,
    }).then((response) => {
      // Transition to the success page
      router.push('/success');
    }).catch((error) => {
      // Show the error
      setError(error);
    }).then(() => {
      // Indicate we've finished loading
      setIsLoading(false);
    });
    // Update the state
    setIsLoading(true);
  }
  // Render the component
  return (
    <Layout>
      <p>
        Register your device to get started
      </p>
      {error &&
        <Alert severity="error">{error.message}</Alert>
      }
      <FormWithLoading isLoading={isLoading} onSubmit={handleSubmit} />
    </Layout>
  )
}
