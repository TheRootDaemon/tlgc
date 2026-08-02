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
