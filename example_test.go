package fibermonitor_test

import (
	"log"
	"time"

	"github.com/gofiber/fiber"
	fibermonitor "github.com/labbsr0x/fiber-monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type FakeDependencyChecker struct{}

func (m *FakeDependencyChecker) GetDependencyName() string {
	return "fake-dependency"
}

func (m *FakeDependencyChecker) Check() fibermonitor.DependencyStatus {
	return fibermonitor.DOWN
}

func YourHandler(c *fiber.Ctx) {
	c.Status(fiber.StatusOK).Send("fiber-monitor!\n")
}

func main() {
	// Creates fiber-monitor instance
	monitor, err := fibermonitor.New("v1.0.0", fibermonitor.DefaultErrorMessageKey, fibermonitor.DefaultBuckets)
	if err != nil {
		panic(err)
	}

	dependencyChecker := &FakeDependencyChecker{}
	monitor.AddDependencyChecker(dependencyChecker, time.Second*30)

	app := fiber.New()

	// Register fiber-monitor middleware
	app.Use(monitor.Prometheus())
	// Register metrics endpoint
	app.Get("/metrics", func(c *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())(c.Fasthttp)
	})
	// Routes consist of a path and a handler function.
	app.Get("/", YourHandler)

	// Bind to a port and pass our router in
	log.Fatal(app.Listen(3000))
}
