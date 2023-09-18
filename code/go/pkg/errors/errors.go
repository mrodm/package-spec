package errors

import (
	"fmt"
	"strings"
)

// ValidationError is the interface that every validation error must implement.
type ValidationError interface {
	error

	// Severity returns the validation severity of the error, the higher the pickier. 0 is mandatory.
	Severity() int

	// Code returns a unique identifier of this error.
	Code() string
}

// ValidationPathError is the interface that validation errors related to paths must implement.
type ValidationPathError interface {
	// File returns the file path where the error was raised.
	File() string
}

type ValidationErrors []ValidationError

// Filter filters the validation errors using the function given as a parameter.
func (ve ValidationErrors) Filter(filter func(elem *ValidationError) bool) ValidationErrors {
	var errs ValidationErrors

	for _, item := range ve {
		if filter(&item) {
			errs = append(errs, item)
		}
	}
	return errs
}

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "found 0 validation errors"
	}

	var message strings.Builder
	errorWord := "errors"
	if len(ve) == 1 {
		errorWord = "error"
	}
	fmt.Fprintf(&message, "found %v validation %v:\n", len(ve), errorWord)
	for idx, err := range ve {
		fmt.Fprintf(&message, "%4d. %v\n", idx+1, err)
	}

	return message.String()
}

// Append adds more validation errors.
func (ve *ValidationErrors) Append(moreErrs ValidationErrors) {
	if len(moreErrs) == 0 {
		return
	}

	errs := *ve
	errs = append(errs, moreErrs...)

	*ve = errs
}
