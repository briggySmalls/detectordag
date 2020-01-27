package shared

func WrapError(err error, msg string) {
	return fmt.Errorf("%s: %w", msg, err)
}
