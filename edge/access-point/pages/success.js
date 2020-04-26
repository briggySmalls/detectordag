import Layout from '../components/layout';
import Alert from '@material-ui/lab/Alert';

// Create the actual homepage
export default function Success() {
  // Render the component
  return (
    <Layout>
      <Alert severity="success">Your dag has been registered successfully!</Alert>
    </Layout>
  )
}
