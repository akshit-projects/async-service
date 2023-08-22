package teams_api

import (
	"net/http"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
)

type resource struct {
	app app.App
	service Service
}

// RegisterRoutes implements Service.
func RegisterRoutes(c fiber.Router, app app.App) {
	resource := &resource{app, NewService(app)}
	c.Get("/teams", resource.getTeams)
	c.Post("/teams/member", resource.addMember)
}

func (r *resource) getTeams(c *fiber.Ctx) error {
	teams, err := r.service.GetTeams(c.Context())

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Unable to get teams",
		})
	}

	return c.JSON(teams)
}

func (r *resource) addMember(c *fiber.Ctx) error {
	var req AddMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Invalid request body, teamId and userEmail required",
		})
	}

	err := r.service.AddMember(c.Context(), req.TeamId, req.UserEmail)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
	}

	return nil
}