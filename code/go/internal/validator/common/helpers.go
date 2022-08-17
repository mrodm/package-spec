package common

import (
	"errors"
	"os"
	"strconv"
)

func IsDefinedWarningsAsErrors() bool {
	warningsAsErrors := false
	warningsAsErrorsStr, found := os.LookupEnv("PACKAGE_SPEC_WARNINGS_AS_ERRORS")
	if found {
		warningsAsErrors, err = strconv.ParseBool(warningsAsErrorsStr)
		if err != nil {
			return errors.New("invalid value for PACKAGE_SPEC_WARNINGS_AS_ERRORS")
		}
	}
	return warningsAsErrors
}
