package main

import (
	"context"
	"fmt"
	"time"

	"log"
	"math/rand"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitJaeger initializes a Jaeger Tracer
func InitJaeger(service string) (opentracing.Tracer, func()) {
	cfg := config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("cannot initialize Jaeger Tracer: %v", err)
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, func() { closer.Close() }
}

// Database simulates the Database service
func Database(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "DatabaseArea")
	span.SetTag("component", "Database")
	defer span.Finish()

	// logic of app
	time.Sleep(2 * time.Second)

	// Simulate processing time

	Redis(ctx)

	Thirdparty(ctx)
	// Call database function within the same context

}

// Database simulates the Database service
func Redis(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "RedisArea")
	span.SetTag("component", "Redis")
	defer span.Finish()

	// logic of app
	time.Sleep(2 * time.Second) // Simulate processing time

	// Call database function within the same context

}

// Database simulates the Database service
func Thirdparty(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "ThirdPartyArea")
	span.SetTag("component", "Thirdparty")
	defer span.Finish()

	// logic of app
	time.Sleep(2 * time.Second) // Simulate processing time

	// Call database function within the same context

}

// Database simulates the database service
func PaymentService(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, "PaymentArea")
	defer span.Finish()

	fmt.Println("Processing database request...")
	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second) // Simulate database processing time
}

func main() {

	tracer, closer := InitJaeger("TicketPurchasingSystem")
	defer closer()

	// Start a new root span for the main operation
	span := tracer.StartSpan("NewMainRequest")
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// Simulate the video list request
	Database(ctx)

	PaymentService(ctx)

	// Finish the main span
	span.Finish()

	fmt.Println("Request completed")
}
