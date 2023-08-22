package experiment_apis

import (
	"net/http"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
)

type resource struct {
	app     app.App
	service Service
}

// RegisterRoutes implements Service.
func RegisterRoutes(c fiber.Router, app app.App) {
	resource := &resource{app, NewService(app)}
	c.Get("/exp", resource.getExperiments)
	c.Get("/exp/layers", resource.getLayers)
	c.Get("/exp/layers/override", resource.overrideLayersVariant)
}

func (r *resource) getExperiments(c *fiber.Ctx) error {
	teamName := getTeamName(c)
	if teamName == "" {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Team name not present in the request",
		})
	}

	expData, err := r.service.GetExperimentsData(teamName)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
		return nil
	}

	c.JSON(expData)
	return nil
}

func (r *resource) getLayers(c *fiber.Ctx) error {
	teamName := getTeamName(c)
	if teamName == "" {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Team name not present in the request",
		})
	}

	expData, err := r.service.GetLayerExperimentsData(teamName)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
		return nil
	}

	c.Status(http.StatusOK).JSON(expData)
	return nil
}

func (r *resource) overrideLayersVariant(c *fiber.Ctx) error {
	var reqBody OverrideVariantRequest
	if err := c.BodyParser(&reqBody); err != nil {
		c.Status(400).JSON(&common_structs.HttpError{
			Msg: "Unable to parse the request body",
		})
		return nil
	}

	err := r.service.OverrideVariant(&reqBody)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
		return nil
	}

	c.Status(http.StatusOK)
	return nil
}

func getTeamName(c *fiber.Ctx) string {
	teamName := c.Params("teamName")
	if teamName == "" {
		return ""
	}

	return teamName
}
