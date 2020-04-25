import CircularProgress from '@material-ui/core/CircularProgress';

function WithLoading(Component) {
  return function WihLoadingComponent({ isLoading, ...props }) {
    // Short-circuit with the provided component
    if (!isLoading) {
        return (<Component {...props} />);
    }
    // Indicate that the component is loading
    return (<CircularProgress />);
  }
}

export default WithLoading;
