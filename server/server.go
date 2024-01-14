package server

import (
	"fmt"
	"time"

	auth_apis "github.com/akshitbansal-1/async-testing/be/api/auth"
	experiment_apis "github.com/akshitbansal-1/async-testing/be/api/experiment"
	flow_apis "github.com/akshitbansal-1/async-testing/be/api/flow"
	teams_api "github.com/akshitbansal-1/async-testing/be/api/teams"
	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/repository/cache"
	thirdparty "github.com/akshitbansal-1/async-testing/be/third_party"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewServer(app app.App) {
	server := fiber.New()

	server.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001", // Replace with your frontend's URL
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// add rate limiter
	addRateLimiter(server, app)
	server.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	server.Use(logger.New(logger.Config{
		DisableColors: true,
	}))
	server.Get("/metrics", thirdparty.RegisterMetrics)
	registerApis(server, app)
	server.Listen(":3000")
}

func addRateLimiter(server *fiber.App, app app.App) {
	server.Use(limiter.New(limiter.Config{
		Max: 20,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("rl:%s", c.IP())
		},
		LimiterMiddleware: limiter.SlidingWindow{},
		Expiration:        1 * time.Minute,
		Storage:           cache.NewCustomRateLimiterStorage(app.GetCacheClient()),
	}))
}

func registerApis(server *fiber.App, app app.App) {
	v1Group := server.Group("/api/v1")
	teams_api.RegisterRoutes(v1Group, app)
	experiment_apis.RegisterRoutes(v1Group, app)
	flow_apis.RegisterRoutes(v1Group, app)
	auth_apis.RegisterRoutes(v1Group, app)
}
