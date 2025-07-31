package req

import (
	"github.com/LekcRg/GophKeeper/internal/config"
	"resty.dev/v3"
)

type ResError struct {
	Errors map[string]string
}

const minErrStatus = 299

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
	config *config.ClientConfig
}

func New(cfg *config.ClientConfig) *Request {
	return &Request{
		client: resty.New(),
		config: cfg,
	}
}

func (r *Request) Close() error {
	return r.client.Close()
}
