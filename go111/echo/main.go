package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"

	"google.golang.org/genproto/googleapis/api/monitoredres"

	"cloud.google.com/go/logging"

	"github.com/labstack/echo"
)

var (
	projectID         = os.Getenv("GOOGLE_CLOUD_PROJECT")
	logClient         = fmt.Sprintf("projects/%s", projectID)
	logID             = "example-logger"
	monitoredResource = &monitoredres.MonitoredResource{
		Labels: map[string]string{
			"project_id": projectID,
			"module_id":  os.Getenv("GAE_SERVICE"),
			"version_id": os.Getenv("GAE_VERSION"),
		},
		Type: "gae_app",
	}
)

func logout(c echo.Context, logger *logging.Logger, severity logging.Severity, payload interface{}) {
	// リクエストヘッダからトレースID情報の抽出を行うAppEngine環境でのみ利用可
	val := c.Request().Header.Get("X-Cloud-Trace-Context")
	id := strings.SplitN(val, "/", 2)[0]
	logger.Log(logging.Entry{
		Trace:    fmt.Sprintf("projects/%s/traces/%s", projectID, id),
		Severity: severity,
		Payload:  payload,
		Resource: monitoredResource,
	})
}

func customRecover(logger *logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10)
					length := runtime.Stack(stack, false)
					logout(c, logger, logging.Critical, string(stack[:length]))
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

func loggerFlusher(logger *logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			next(c)

			go func() {
				if err := logger.Flush(); err != nil {
					fmt.Println(err)
				}
			}()
			return nil
		}
	}
}

func main() {
	ctx := context.Background()
	client, err := logging.NewClient(ctx, logClient)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	logger := client.Logger(logID)
	e.Use(customRecover(logger))
	e.GET("/", func(c echo.Context) error {
		r := rand.Intn(10)
		if r > 7 {
			panic("panic occur")
		}

		logout(c, logger, logging.Info, echo.Map{
			"label":   "header",
			"headers": c.Request().Header,
		})
		envs := make(map[string]interface{})
		for _, env := range os.Environ() {
			val := strings.Split(env, "=")
			envs[val[0]] = val[1]
		}
		logout(c, logger, logging.Info, echo.Map{
			"label": "env_vars",
			"envs":  envs,
		})
		logout(c, logger, logging.Info, fmt.Sprintf("Value: %v", r))
		return c.JSON(200, echo.Map{
			"message": "HelloWorld",
			"value":   r,
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, e)
}
