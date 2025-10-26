package utils

import "log"

// failOnError logs the provided message and panics if err is non-nil.
// This is a helper used in this package for simplicity.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
