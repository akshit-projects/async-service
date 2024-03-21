package flow_apis

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/akshitbansal-1/async-testing/be/utils"
	"github.com/akshitbansal-1/async-testing/lib/manager"
	"github.com/akshitbansal-1/async-testing/lib/structs"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	MAX_FLOW_LIMIT = 6
)

type resource struct {
	app     app.App
	service Service
}

// RegisterRoutes implements Service.
func RegisterRoutes(c fiber.Router, app app.App) {
	resource := &resource{app, NewService(app)}
	c.Get("/flow", resource.getFlows)
	c.Get("/flow/run", websocket.New(resource.runFlow))
	c.Get("/flow/:id", resource.getFlow)
	c.Post("/flow", resource.addFlow)
	c.Put("/flow", resource.updateFlow)
	c.Post("/flow/validate", resource.validateSteps)
}

func (r *resource) updateFlow(c *fiber.Ctx) error {
	var flow *structs.Flow
	var err error
	if flow, err = getFlowObject(c); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Unable to parse the request body",
		})
	}

	existingFlow, err := r.service.GetFlow(flow.Id)
	if flow.Id == "" || err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
				Msg: "Flow not found for the id",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
	}

	flow.CreatedAt = existingFlow.CreatedAt
	if err = r.service.UpdateFlow(flow); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
	}

	return c.SendStatus(http.StatusOK)
}

func (r *resource) getFlow(c *fiber.Ctx) error {
	flowId := c.Params("id")
	flow, err := r.service.GetFlow(flowId)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return c.Status(http.StatusNotFound).JSON(&common_structs.HttpError{
				Msg: "Flow not found for the id",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
	}

	return c.JSON(flow)
}

func (r *resource) getFlows(c *fiber.Ctx) error {
	filter, err := getFilter(c)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.APIError{
			Msg: err.Error(),
		})
	}
	flows, err := r.service.GetFlows(filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&common_structs.HttpError{
			Msg: "Unable to get flows",
		})
	}

	return c.JSON(flows)
}

func getFilter(c *fiber.Ctx) (*common_structs.APIFilter, error) {
	search := c.Query("search", "")
	searchFilter := map[string]interface{}{}

	if search != "" {
		searchFilter["name"] = bson.M{
			"$regex": "(?i).*" + search + ".*",
		}
	}

	if ids := c.Query("ids", ""); ids != "" {
		fIds := strings.Split(ids, ",")
		if oids, err := utils.MakeObjectIds(fIds); err != nil {
			return nil, errors.New("Invalid flow ids")
		} else {
			searchFilter["_id"] = bson.M{
				"$in": oids,
			}
		}
	}

	limit := int64(math.Max(float64(c.QueryInt("limit", MAX_FLOW_LIMIT)), MAX_FLOW_LIMIT))
	return &common_structs.APIFilter{
		Filters: searchFilter,
		Limit:   limit,
	}, nil
}

func (r *resource) addFlow(c *fiber.Ctx) error {
	var flow *structs.Flow
	var err error
	if flow, err = getFlowObject(c); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Unable to parse the request body",
		})
	}

	uid, err := r.service.AddFlow(flow)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
	}

	return c.Status(http.StatusCreated).
		JSON(uid)
}

func (r *resource) validateSteps(c *fiber.Ctx) error {
	var flow *structs.Flow
	var err error
	if flow, err = getFlowObject(c); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: "Unable to parse the request body",
		})
	}

	if err := r.service.ValidateSteps(flow); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&common_structs.HttpError{
			Msg: err.Error(),
		})
	}

	return c.SendStatus(http.StatusOK)
}

func (r *resource) runFlow(conn *websocket.Conn) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		data, _ := utils.ToBytes[structs.StepResponse](structs.StepResponse{
			Name:   "",
			Status: structs.STEP_ERROR,
			Response: &structs.StepError{
				Error: "Unable to parse request body",
			},
			Id: "",
		})
		conn.WriteMessage(websocket.TextMessage, data)
		return
	}
	var flow *structs.Flow = &structs.Flow{}
	if err := json.Unmarshal(msg, flow); err != nil {
		data, _ := utils.ToBytes[structs.StepResponse](structs.StepResponse{
			Name:   "",
			Status: structs.STEP_ERROR,
			Response: &structs.StepError{
				Error: "Unable to get request body",
			},
			Id: "",
		})
		conn.WriteMessage(websocket.TextMessage, data)
		return
	}
	ch := manager.NewChannel[*structs.ExecutionStatusUpdate](0)
	go r.service.RunFlow(ch, flow)

	for resp := range ch.GetChannel() {
		data, _ := utils.ToBytes[structs.ExecutionStatusUpdate](*resp)
		conn.WriteMessage(websocket.TextMessage, data)
	}

	conn.Close()
}

func getFlowObject(c *fiber.Ctx) (*structs.Flow, error) {
	var flow structs.Flow
	if err := c.BodyParser(&flow); err != nil {
		return nil, errors.New("")
	}

	return &flow, nil
}
