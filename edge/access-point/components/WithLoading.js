import CircularProgress from '@material-ui/core/CircularProgress';
import PropTypes from "prop-types";

export default function WithLoading(Component) {
  const WithLoadingComponent = ({ isLoading, ...props }) => {
    // Short-circuit with the provided component
    if (!isLoading) {
        return (<Component {...props} />);
    }
    // Indicate that the component is loading
    return (<CircularProgress />);
  }
  WithLoadingComponent.propTypes = {
    isLoading: PropTypes.bool.isRequired,
    props: PropTypes.any,
  };
}

WithLoading.propTypes = {
    Component: PropTypes.element.isRequired,
};
