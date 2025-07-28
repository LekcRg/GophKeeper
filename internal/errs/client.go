package errs

import "errors"

var ErrMustContainHTTP = errors.New("must contain http:// or https://")
