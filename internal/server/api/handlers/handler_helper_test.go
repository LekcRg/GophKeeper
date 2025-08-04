package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
)

type serveHTTPOpts struct {
	body   io.Reader
	method string
	target string
}

func serveHTTPWithCtx(ctx context.Context, handler http.HandlerFunc, opts serveHTTPOpts) *httptest.ResponseRecorder {
	method := "POST"
	if opts.method != "" {
		method = opts.method
	}

	target := "/"
	if opts.target != "" {
		target = opts.target
	}

	var body io.Reader = nil
	if opts.body != nil {
		body = opts.body
	}

	req := httptest.NewRequest(method, target, body)
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	h := http.HandlerFunc(handler)
	h.ServeHTTP(res, req)

	return res
}
