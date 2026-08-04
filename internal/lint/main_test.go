package lint

// errorCode returns the code of the first reported error,
// or "" if the result is clean.
// Rules report at most one error, always at line 0.
func errorCode(r *Result) string {
	if len(r.Errors) == 0 {
		return ""
	}
	return r.Errors[0].Code
}

// errorCodes returns the codes of all reported errors,
// in report order.
// It returns nil if the result is clean.
func errorCodes(r *Result) []string {
	if len(r.Errors) == 0 {
		return nil
	}
	codes := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		codes[i] = e.Code
	}
	return codes
}
