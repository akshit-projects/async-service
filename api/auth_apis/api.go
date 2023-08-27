package auth_apis

import (
	"fmt"
	"net/http"
	"time"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type resource struct {
	app     app.App
	service Service
}

func RegisterRoutes(c fiber.Router, app app.App) {
	resource := &resource{app, NewService(app)}
	c.Post("/login", resource.loginHandler)
}

func (r *resource) loginHandler(c *fiber.Ctx) error {
	var request AuthRequest
	cfg := r.app.GetConfig()
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Invalid request body",
		})
	}

	// Validate the JWT is valid
	claims, err := r.service.LoginUser(request.IdToken)
	if err != nil {
		return c.Status(http.StatusForbidden).JSON(&common_structs.HttpError{
			Msg: "Invalid request",
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	tc := token.Claims.(jwt.MapClaims)
	expiry := time.Now().Add(10 * time.Minute).Unix()
	tc["exp"] = expiry
	tc["authorized"] = true
	tc["user"] = claims.Email
	if tokenString, err := token.SignedString([]byte(cfg.JWTSecret)); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Unable to authenticate the user.",
		})
	} else {
		return c.Status(200).JSON(ClientToken{
			tokenString,
			fmt.Sprintf("%d", expiry),
		})
	}

}
