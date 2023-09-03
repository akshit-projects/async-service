package suite_apis

import (
	"math"
	"net/http"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

const MAX_SUITE_LIMIT = 10

type resource struct {
	app     app.App
	service Service
}

// RegisterRoutes implements Service.
func RegisterRoutes(c fiber.Router, app app.App) {
	resource := &resource{app, NewService(app)}
	c.Post("/suite", resource.addSuite)
	c.Get("/suite", resource.getSuites)
}

func (r *resource) getSuites(c *fiber.Ctx) error {
	filter := getFilter(c)
	suites, err := r.service.GetSuites(filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: "Unable to get flows",
		})
	}

	return c.JSON(suites)

}

func getFilter(c *fiber.Ctx) *common_structs.APIFilter {
	search := c.Query("search", "")
	searchFilter := map[string]interface{}{
		"name": bson.M{
			"$regex": "(?i).*" + search + ".*",
		},
	}
	limit := int64(math.Max(float64(c.QueryInt("limit", MAX_SUITE_LIMIT)), MAX_SUITE_LIMIT))
	return &common_structs.APIFilter{
		Filters: searchFilter,
		Limit:   limit,
	}
}

func (r *resource) addSuite(c *fiber.Ctx) error {
	var suite Suite
	if err := c.BodyParser(&suite); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Invalid request body",
		})
	}

	if sId, err := r.service.AddSuite(&suite); err != nil {
		return c.Status(err.Status).JSON(&common_structs.HttpError{
			Msg: err.Msg,
		})
	} else {
		return c.Status(http.StatusCreated).JSON(sId)
	}
}
