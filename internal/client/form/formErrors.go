package form

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/req"
)

type Errors struct {
	fields      map[string]string
	Message     string
	knownFields map[string]bool // поля, которые должны отображаться как field errors
}

func NewErrors(knownFields []string) *Errors {
	known := make(map[string]bool)
	for _, field := range knownFields {
		known[field] = true
	}

	return &Errors{
		fields:      make(map[string]string),
		knownFields: known,
	}
}

func (e *Errors) Clear() {
	e.fields = make(map[string]string)
	e.Message = ""
}

func (e *Errors) setFieldError(field, message string) {
	e.fields[field] = message
}

func (e *Errors) setMessage(message string) {
	e.Message = message
}

func (e *Errors) GetFieldError(field string) string {
	return e.fields[field]
}

func (e *Errors) HandleAPIError(err error) {
	var resErr *req.ResError
	if !errors.As(err, &resErr) || resErr == nil || resErr.Errors == nil {
		e.setMessage(err.Error())
		return
	}

	var unknownErrors []string

	for key, val := range resErr.Errors {
		if val == "" {
			continue
		}

		if e.knownFields[key] {
			e.setFieldError(key, val)
		} else {
			unknownErrors = append(unknownErrors, fmt.Sprintf("%s: %s", key, val))
		}
	}

	if len(unknownErrors) > 0 {
		e.setMessage(strings.Join(unknownErrors, "; "))
	}
}
