package errors

import (
	"fmt"

	"cuelang.org/go/cue/errors"
)

func Describe(desc string, err error) error {
	return fmt.Errorf("%s: %s", desc, errors.Details(err, nil))
}
