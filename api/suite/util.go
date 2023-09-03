package suite_apis

import (
	"context"
	"errors"
	"fmt"
	"log"

	flow_apis "github.com/akshitbansal-1/async-testing/be/api/flow"
	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/common_structs"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	SUITE_DB_NAME         = "test"
	SUITE_COLLECTION_NAME = "suite"
)

func getSuites(app app.App, filter *common_structs.APIFilter) ([]Suite, error) {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database(SUITE_DB_NAME).Collection(SUITE_COLLECTION_NAME)

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
		return nil, errors.New("Unable to get suites from DB")
	}

	defer cursor.Close(ctx)
	suites := []Suite{}
	for cursor.Next(ctx) {
		var suite Suite
		err := cursor.Decode(&suite)
		if err != nil {
			log.Println("Decode error. Unable to decode data.", err)
			continue
		}
		// TODO update if flows are deleted
		suites = append(suites, suite)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Unable to get flows", err.Error())
		return nil, errors.New("An unknown error occurred")
	}

	return suites, nil
}

func addSuite(app app.App, s *Suite) (*string, error) {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database(SUITE_DB_NAME).Collection(SUITE_COLLECTION_NAME)
	data, err := bson.Marshal(*s)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Unable to save suite data")
	}

	if err = checkFlowsExistence(app, s.FlowIds); err != nil {
		return nil, err
	}

	result, err := coll.InsertOne(context.Background(), data)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("Unable to save suite data")
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return &id, nil
}

func checkFlowsExistence(app app.App, ids []string) error {
	dbClient := app.GetMongoClient()
	fCol := dbClient.Database(flow_apis.FLOW_DB_NAME).
		Collection(flow_apis.FLOW_COLLECTION_NAME)

	objectIds := make([]primitive.ObjectID, len(ids))

	for i, oid := range ids {
		id, err := primitive.ObjectIDFromHex(oid)
		if err != nil {
			return errors.New("Invalid object id passed -> " + oid)
		}
		objectIds[i] = id
	}
	count, err := fCol.CountDocuments(context.Background(), bson.M{
		"_id": bson.M{
			"$in": objectIds,
		},
	})
	if err != nil {
		return errors.New("Unable to save suite data")
	}

	if int(count) != len(ids) {
		return fiber.ErrNotFound
	}

	return nil
}
