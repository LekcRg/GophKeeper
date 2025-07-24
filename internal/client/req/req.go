package req

import "resty.dev/v3"

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
