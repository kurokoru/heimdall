package plugins

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kurokoru/heimdall/v7"
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	logger zerolog.Logger
}

func NewZerologLogger(log zerolog.Logger) heimdall.Plugin {
	return &zerologLogger{
		logger: log,
	}
}

func (zl *zerologLogger) OnRequestStart(req *http.Request) {
	// ctx := context.WithValue(req.Context(), reqTime, time.Now())
	// *req = *(req.WithContext(ctx))
}

func (zl *zerologLogger) OnRequestEnd(req *http.Request, res *http.Response) {
	reqDuration := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()
	statusCode := res.StatusCode
	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = ioutil.ReadAll(req.Body)
	}
	resBody, _ := ioutil.ReadAll(res.Body)

	zl.logger.Info().
		Str("uri", url).
		Str("method", method).
		Int("http_status", statusCode).
		RawJSON("request_param", reqBody).
		RawJSON("response_message", resBody).
		Dur("elapsed_time", reqDuration).
		Msg("")
}

func (zl *zerologLogger) OnError(req *http.Request, err error) {
	reqDuration := getRequestDuration(req.Context()) / time.Millisecond
	method := req.Method
	url := req.URL.String()
	zl.logger.Error().
		Err(err).
		Str("uri", url).
		Str("method", method).
		Dur("elapsed", reqDuration)
}
