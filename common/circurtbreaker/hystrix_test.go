package circurtbreaker

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000",
	})
	logrus.SetLevel(logrus.DebugLevel)
	Client = &http.Client{}
}

func TestCallUsingResilience(t *testing.T) {
	Convey("Test call circurt breaker", t, func() {
		defer gock.Off()

		Convey("When all failed", func() {
			buildGockMatcherTimes(500, 4)
			hystrix.Flush()
			bytes, err := CallUsingCircuitBreaker(context.TODO(), "TEST", "http://quotes-service", "GET")
			So(err, ShouldNotBeNil)
			So(bytes, ShouldBeNil)
		})

		Convey("Will got some resp after two failed", func() {
			retries = 3
			buildGockMatcherTimes(500, 2)
			body := []byte("Some response")
			buildGockMatcherWithBody(200, string(body))
			hystrix.Flush()

			bytes, err := CallUsingCircuitBreaker(context.TODO(), "TEST", "http://quotes-service", "GET")
			So(err, ShouldBeNil)
			So(bytes, ShouldNotBeNil)
			So(string(body), ShouldEqual, string(bytes))
		})
	})
}

func TestCallHystrixOpensAfterThresholdPassed(t *testing.T) {
	Convey("Test call circurt breaker open after threshold passed", t, func() {
		defer gock.Off()

		// 3 failed
		for a := 0; a < 3; a++ {
			buildGockMatcher(500)
		}

		// 3 success
		for a := 0; a < 3; a++ {
			buildGockMatcherWithBody(200, "")
		}

		hystrix.Flush()

		retries = 0
		hystrix.ConfigureCommand("TEST", hystrix.CommandConfig{
			RequestVolumeThreshold: 5,    // at least 5 request error will get this
			SleepWindow:            1000, // time wait to check recovery
			ErrorPercentThreshold:  50,   // 50% error resp
		})

		// total 6 request and 50% will get error resp
		for a := 0; a < 6; a++ {
			_, _ = CallUsingCircuitBreaker(context.TODO(), "TEST", "http://quotes-service", "GET")
		}

		cb, _, _ := hystrix.GetCircuit("TEST")
		So(cb.IsOpen(), ShouldBeTrue)

		// wait for to cirtcurt closed
		time.Sleep(1100 * time.Millisecond)

		for a := 0; a < 3; a++ {
			buildGockMatcherWithBody(200, "")
		}

		for a := 0; a < 3; a++ {
			_, _ = CallUsingCircuitBreaker(context.TODO(), "TEST", "http://quotes-service", "GET")
		}

		So(cb.IsOpen(), ShouldBeFalse)
	})
}

func buildGockMatcherTimes(status int, times int) {
	for a := 0; a < times; a++ {
		buildGockMatcher(status)
	}
}

func buildGockMatcherWithBody(status int, body string) {
	gock.New("http://quotes-service").
		Reply(status).BodyString(body)
}

func buildGockMatcher(status int) {
	buildGockMatcherWithBody(status, "")
}
