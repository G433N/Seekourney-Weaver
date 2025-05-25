package utils

// PanicOnError panics on non-nil error.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
