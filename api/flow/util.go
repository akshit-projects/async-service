package flow_apis

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	FLOW_DB_NAME         = "test"
	FLOW_COLLECTION_NAME = "flow"
)

func addFlow(app app.App, flow *Flow) (*string, error) {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database(FLOW_DB_NAME).Collection(FLOW_COLLECTION_NAME)

	data, err := bson.Marshal(*flow)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Unable to save data")
	}

	result, err := coll.InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Unable to save flow data")
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return &id, nil
}

func getFlows(app app.App, filter *common_structs.APIFilter) ([]Flow, error) {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database(FLOW_DB_NAME).Collection(FLOW_COLLECTION_NAME)

	ctx := context.Background()
	// get all the records
	mFilter := bson.M{}
	for key, value := range filter.Filters {
		mFilter[key] = value
	}
	cursor, err := coll.Find(ctx, mFilter, &options.FindOptions{
		Limit: &filter.Limit,
		Skip:  &filter.Skip,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Unable to get flows from DB")
	}

	defer cursor.Close(ctx)
	flows := []Flow{}
	for cursor.Next(ctx) {
		var flow Flow
		err := cursor.Decode(&flow)
		if err != nil {
			log.Println("Decode error. Unable to decode data.", err)
			continue
		}

		for idx := range flow.Steps {
			// convert mongo format to required format
			step := &flow.Steps[idx]
			kv := step.Meta.(primitive.D)
			mp := make(map[string]interface{})
			for k, v := range kv.Map() {
				mp[k] = v
			}
			step.Meta = mp
		}
		flows = append(flows, flow)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Unable to get flows", err.Error())
		return nil, errors.New("An unknown error occurred")
	}

	return flows, nil
}

func getFlow(app app.App, id string) (*Flow, error) {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database(FLOW_DB_NAME).Collection(FLOW_COLLECTION_NAME)

	ctx := context.Background()
	// get all the records
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("Invalid flow id passed")
	}

	result := coll.FindOne(ctx, bson.D{
		{Key: "_id", Value: oid},
	})

	var flow Flow
	err = result.Decode(&flow)
	if err == mongo.ErrNoDocuments {
		return nil, fiber.ErrNotFound
	} else if err != nil {
		return nil, errors.New("Unable to get flow data")
	}

	for idx := range flow.Steps {
		// convert mongo format to required format
		step := &flow.Steps[idx]
		kv := step.Meta.(primitive.D)
		mp := make(map[string]interface{})
		for k, v := range kv.Map() {
			mp[k] = v
		}
		step.Meta = mp
	}

	return &flow, nil
}

func updateFlow(app app.App, flow *Flow) error {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database(FLOW_DB_NAME).Collection(FLOW_COLLECTION_NAME)

	oid, err := primitive.ObjectIDFromHex(flow.Id)
	if err != nil {
		return errors.New("Invalid flow id passed")
	}

	flow.Id = ""
	data, err := bson.Marshal(*flow)
	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Unable to update data")
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	result := coll.FindOneAndReplace(context.Background(), filter, data)
	if result.Err() != nil {
		return errors.New("Unable to update the flow data")
	}

	return nil
}
