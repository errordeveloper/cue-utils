package errors

import (
	"bytes"
	"fmt"
	"io"

	"cuelang.org/go/cue/errors"
)

func Describe(desc string, err error) error {
	w := bytes.NewBuffer(nil)

	for _, err := range errors.Errors(err) {
		errors.Print(w, err, &errors.Config{Cwd: "/"})
	}

	msg, err := io.ReadAll(w)
	if err != nil {
		return err
	}
	return fmt.Errorf("%s: %s", desc, msg)
}
