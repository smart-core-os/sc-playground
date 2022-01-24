package errs

// First returns the first non-nil error from errs or nil.
func First(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
