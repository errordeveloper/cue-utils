package errors

import (
	"fmt"

	"cuelang.org/go/cue/errors"
)

func Describe(desc string, err error) error {
	msg := errors.Details(err, &errors.Config{})
	return fmt.Errorf("%s: %s", desc, msg)
}
