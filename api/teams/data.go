package teams_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/akshitbansal-1/async-testing/be/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// Get team details from DB
func GetTeams(ctx context.Context, app app.App) ([]Team, error) {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database("test").Collection("teams")
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var teams []Team
	if err = cursor.All(ctx, &teams); err != nil {
		return nil, err
	}

	for _, result := range teams {
		cursor.Decode(&result)
		output, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s\n", output)
	}

	return teams, nil
}

// Add a member for the team
func AddMember(ctx context.Context, app app.App, teamId string, userEmail string) error {
	dbClient := app.GetMongoClient()
	coll := dbClient.Database("test").Collection("teams")
	result := coll.FindOne(ctx, bson.D{
		{Key: "id", Value: teamId},
	})
	if result.Err() != nil {
		return result.Err()
	}

	var team Team
	result.Decode(&result)
	if utils.Contains(team.Members, userEmail) {
		return errors.New("User id already present in the team")
	}

	team.Members = append(team.Members, userEmail)
	data, _ := bson.Marshal(team)
	updateResult, err := coll.UpdateOne(ctx, bson.D{
		{Key: "id", Value: teamId},
	}, data)
	if err != nil || updateResult.ModifiedCount != 1 {
		return errors.New("Unable to update members in the team")
	}

	return nil
}
