import { makeStyles } from '@material-ui/core/styles';
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
export default function Form ({ onSubmit }) {
  // Create styles for use
  const classes = useStyles();
  // Create a handler
  const innerSubmit = (event) => {
    // Do not actually post a submission
    event.preventDefault();
    // Now call outer handler
    const formData = new FormData(event.target);
    let data = {};
    formData.forEach((value, key) => {data[key] = value});
    onSubmit(data);
  }
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
