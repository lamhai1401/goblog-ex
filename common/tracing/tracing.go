package tracing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/sirupsen/logrus"
)

// Tracer instance
var tracer opentracing.Tracer

var serviceName = "opentracing-span"

// SetTracer can be used by unit tests to provide a NoopTracer instance. Real users should always
// use the InitTracing func.
func SetTracer(initializedTracer opentracing.Tracer) {
	tracer = initializedTracer
}

// InitTracing connects the calling service to Zipkin and initializes the tracer.
func InitTracing(zipkinURL string, serviceName string) {
	logrus.Infof("Connecting to zipkin server at %v", zipkinURL)
	reporter := zipkinhttp.NewReporter(fmt.Sprintf("%s/api/v1/spans", zipkinURL))

	endpoint, err := zipkin.NewEndpoint(serviceName, "127.0.0.1:0")
	if err != nil {
		logrus.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		logrus.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer = zipkinot.Wrap(nativeTracer)

	logrus.Infof("Successfully started zipkin tracer for service '%v'", serviceName)
}

// AddTracingToReqFromContext adds tracing information to an OUTGOING HTTP request
func AddTracingToReqFromContext(ctx context.Context, req *http.Request) {
	if ctx.Value(serviceName) == nil {
		return
	}
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err := tracer.Inject(
		ctx.Value(serviceName).(opentracing.Span).Context(),
		opentracing.HTTPHeaders,
		carrier)
	if err != nil {
		panic("Unable to inject tracing context: " + err.Error())
	}
}
