package form

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LekcRg/GophKeeper/internal/client/req"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Errors struct {
	fields      map[string]string
	knownFields map[string]bool // поля, которые должны отображаться как field errors
	Message     string
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

func (e *Errors) HandleError(err error, key string) {
	errText := err.Error()
	if e.knownFields[key] {
		e.setFieldError(key, errText)
		return
	}

	e.setMessage(errText)
}

func (e *Errors) HandleAPIError(err error) {
	var (
		listErrs map[string]string
		valErr   validation.Errors
	)

	var resErr *req.ResError
	//nolint:gocritic // errors switch
	if errors.As(err, &resErr) {
		listErrs = resErr.Errors
	} else if errors.As(err, &valErr) {
		listErrs = make(map[string]string, len(valErr))
		for key, val := range valErr {
			listErrs[key] = val.Error()
		}
	} else {
		e.setMessage(err.Error())

		return
	}

	var unknownErrors []string

	for key, val := range listErrs {
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
