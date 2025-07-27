package req

import (
	"resty.dev/v3"
)

type ResError struct {
	Errors map[string]string
}

func (e *ResError) Error() string {
	errText := "response errors: "

	for key, val := range e.Errors {
		if val == "" {
			continue
		}

		errText += key + ":" + val + " "
	}

	return errText
}

type Request struct {
	client *resty.Client
}

func New() *Request {
	return &Request{
		client: resty.New(),
	}
}

func (r *Request) Close() error {
	return r.client.Close()
}
