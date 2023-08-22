package thirdparty

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	HttpRequest = "HttpRequest"
)

var counters = map[string]string{
	HttpRequest: "http_request",
}
var CountersMetric = map[string]*prometheus.CounterVec{}

func init() {
	for name, value := range counters {
		CountersMetric[name] = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: value,
		}, []string{})
		prometheus.MustRegister(CountersMetric[name])
	}

}

func RegisterMetrics(c *fiber.Ctx) error {
	// prometheus.Register()
	return nil
}
