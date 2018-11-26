package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"google.golang.org/genproto/googleapis/api/monitoredres"

	"cloud.google.com/go/logging"
)

type (
	traceIDKey     string
	typeRequestKey string
)

var (
	projectID                 = os.Getenv("GOOGLE_CLOUD_PROJECT")
	traceID    traceIDKey     = "trace_id"
	requestKey typeRequestKey = "request"
	resource                  = &monitoredres.MonitoredResource{
		Labels: map[string]string{
			"module_id":  os.Getenv("GAE_SERVICE"),
			"project_id": projectID,
			"version_id": os.Getenv("GAE_VERSION"),
		},
		Type: "gae_app",
	}
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(r)

		lg, err := newLogger(ctx)
		if err != nil {
			panic(err)
		}
		defer lg.close()

		lg.debug(ctx, "this is debug")
		lg.info(ctx, "this is info")
		lg.warn(ctx, "this is warn")
		lg.error(ctx, "this is error")
		lg.critical(ctx, "this is critical")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func newContext(r *http.Request) context.Context {
	trace := r.Header.Get("X-Cloud-Trace-Context")
	ctx := context.WithValue(r.Context(), traceID, trace)
	return context.WithValue(ctx, requestKey, r)
}

type logger struct {
	client *logging.Client
	logger *logging.Logger
}

func newLogger(ctx context.Context) (*logger, error) {
	client, err := logging.NewClient(ctx, fmt.Sprintf("projects/%s", projectID))
	if err != nil {
		return nil, err
	}
	return &logger{client, client.Logger("app_logs")}, nil
}

func (l logger) info(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, logging.Info, format, args...)
}

func (l logger) error(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, logging.Error, format, args...)
}

func (l logger) warn(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, logging.Warning, format, args...)
}

func (l logger) critical(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, logging.Critical, format, args...)
}

func (l logger) debug(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, logging.Debug, format, args...)
}

func (l logger) log(ctx context.Context, severity logging.Severity, format string, args ...interface{}) {
	trace := ctx.Value(traceID).(string)
	l.logger.Log(logging.Entry{
		Trace:    fmt.Sprintf("projects/%s/traces/%s", projectID, strings.Split(trace, "/")[0]),
		Severity: severity,
		Payload:  fmt.Sprintf(format, args...),
		Resource: resource,
	})
}

func (l logger) close() error {
	return l.client.Close()
}
