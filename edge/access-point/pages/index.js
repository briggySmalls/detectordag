import { Component, useState } from 'react';
import axios from 'axios';
import WithLoading from '../components/WithLoading';
import RegistrationForm from '../components/registrationForm';
import Layout from '../components/layout';
import { makeStyles } from '@material-ui/core/styles';
import { useRouter } from 'next/router'
import Alert from '@material-ui/lab/Alert';

// Wrap the form in a loader
const FormWithLoading = WithLoading(RegistrationForm);

// Create the actual homepage
export default function Home() {
  // Declare isLoading state
  const [isLoading, setIsLoading] = useState(false);
  // Declare error state
  const [error, setError] = useState(null);
  // Get the nextjs router
  const router = useRouter()
  // Callback for submit event
  const handleSubmit = (data) => {
    // Clear any errors
    setError(null);
    // Submit the form data
    axios.post(
      '/api/register',
      data,
    ).then((response) => {
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
