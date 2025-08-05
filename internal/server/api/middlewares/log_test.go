package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/mocks"
	"github.com/LekcRg/GophKeeper/internal/server/api/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"resty.dev/v3"
)

func testHandler(t *testing.T, code int, f ...func()) http.HandlerFunc {
	t.Helper()

	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if len(f) > 0 {
			f[0]()
		}

		w.WriteHeader(code)
		_, err := w.Write([]byte("test"))
		require.NoError(t, err)
	})
}

func setupRequestLoggerTest(path string, handler http.Handler) (
	*observer.ObservedLogs, *httptest.Server, func(),
) {
	core, observedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	m := New(&config.Config{}, logger, response.NewResponder(logger), &mocks.MockUserRepo{})

	mux := http.NewServeMux()
	mux.Handle(path, m.RequestLogger(handler))

	server := httptest.NewServer(mux)

	cleanup := func() {
		server.Close()
	}

	return observedLogs, server, cleanup
}

func TestRequestLogger(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		method string
		path   string
		code   int
	}

	tests := []test{
		{
			name:   "GET",
			path:   "/test",
			method: http.MethodGet,
			code:   http.StatusOK,
		},
		{
			name:   "DELETE",
			path:   "/test/1",
			method: http.MethodDelete,
			code:   http.StatusUnauthorized,
		},
		{
			name:   "POST",
			path:   "/",
			method: http.MethodPost,
			code:   http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logsList, server, closeServer := setupRequestLoggerTest(tt.path, testHandler(t, tt.code))
			defer closeServer()

			client := resty.New()
			defer client.Close()

			res, err := client.R().
				Execute(tt.method, server.URL+tt.path)
			require.NoError(t, err)

			assert.Equal(t, tt.code, res.StatusCode())

			allLogs := logsList.All()[0]
			countFound := 0

			for _, val := range allLogs.Context {
				switch val.Key {
				case "method":
					assert.Equal(t, tt.method, val.String)
				case "path":
					assert.Equal(t, tt.path, val.String)
				case "status":
					assert.Equal(t, int64(tt.code), val.Integer)
				default:
					continue
				}

				countFound++
			}

			const expectedFields = 3

			assert.Equal(t, expectedFields, countFound, "info log lower then need")
		})
	}
}

func TestRequestLoggerTime(t *testing.T) {
	t.Parallel()

	timeout := time.Millisecond * 200

	logsList, server, closeServer := setupRequestLoggerTest(
		"/", testHandler(t, 200, func() {
			time.Sleep(timeout)
		}))
	defer closeServer()

	client := resty.New()
	defer client.Close()
	_, err := client.R().
		Get(server.URL)
	require.NoError(t, err)

	allLogs := logsList.All()[0]
	found := false

	for _, val := range allLogs.Context {
		if val.Key == "time" {
			assert.True(t, val.Integer >= int64(timeout),
				"expected %v >= %v", time.Duration(val.Integer), timeout)

			found = true
		}
	}

	assert.True(t, found, "not found time in logger")
}
