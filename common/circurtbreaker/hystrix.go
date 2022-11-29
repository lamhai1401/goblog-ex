package circurtbreaker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/lamhai1401/goblog-ex/common/tracing"
	"github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

// Client to do http requests with
var Client *http.Client

// RETRIES is the number of retries to do in the retrier.
var retries = 3

// CallUsingCircuitBreaker performs a HTTP call inside a circuit breaker.
func CallUsingCircuitBreaker(ctx context.Context, breakerName string, url string, method string) ([]byte, error) {
	output := make(chan []byte, 1)
	errors := hystrix.Go(breakerName, func() error {
		req, _ := http.NewRequest(method, url, nil)
		tracing.AddTracingToReqFromContext(ctx, req)
		logrus.Debugf("Call in breaker")
		err := callWithRetries(req, output)
		return err // For hystrix, forward the err from the retrier. It's nil if OK.
	}, func(err error) error {
		logrus.Errorf("In fallback function for breaker %v, error: %v", breakerName, err.Error())
		circuit, _, _ := hystrix.GetCircuit(breakerName)
		logrus.Errorf("Circuit state is: %v", circuit.IsOpen())
		return err
	})

	select {
	case out := <-output:
		logrus.Debugf("Call in breaker %v successful", breakerName)
		return out, nil

	case err := <-errors:
		logrus.Debugf("Got error on channel in breaker %v. Msg: %v", breakerName, err.Error())
		return nil, err
	}
}

func callWithRetries(req *http.Request, output chan []byte) error {
	r := retrier.New(retrier.ConstantBackoff(retries, 100*time.Millisecond), nil)
	attempt := 0
	err := r.Run(func() error {
		attempt++
		resp, err := Client.Do(req)
		if err == nil && resp.StatusCode < 299 {
			responseBody, err := io.ReadAll(resp.Body)
			if err == nil {
				output <- responseBody
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("status was %v", resp.StatusCode)
		}
		logrus.Errorf("Retrier failed, attempt %v", attempt)
		return err
	})
	return err
}
