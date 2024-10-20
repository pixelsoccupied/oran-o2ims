package middleware

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/openshift-kni/oran-o2ims/internal/service/alarms/api/generated"
	"net/http"
	"os"
	"time"
)

type Middleware = func(http.Handler) http.Handler

// CreateMwStack a simple helper function to call
// middlewares in a chain instead of wrapping them one at a time.
func CreateMwStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			cur := middlewares[i]
			next = cur(next)
		}
		return next
	}
}

// LogDuration log time taken to complete a request.
// TODO: This is just get started with middleware but should be replaced with something that's more suitable for production i.e OpenTelemetry
// https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/examples/prometheus/main.go
func LogDuration() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			next.ServeHTTP(w, r)
			fmt.Printf("%q took %s\n", r.RequestURI, time.Since(startTime))
		})
	}
}

// AlarmsOapiValidation to validate all incoming requests as specified in the spec
func AlarmsOapiValidation() Middleware {
	// This also validates the spec
	swagger, err := generated.GetSwagger()
	if err != nil {
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc, // No auth needed even when we have something in spec
		},
	})
}
