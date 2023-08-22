package app

import (
	"github.com/akshitbansal-1/async-testing/be/config"
	thirdparty "github.com/akshitbansal-1/async-testing/be/third_party"
	"go.mongodb.org/mongo-driver/mongo"
)


type App interface {
	GetMongoClient() *mongo.Client
	GetCacheClient() thirdparty.CacheClient
	GetConfig() *config.Configuration
}

type app struct {
	mongoClient *mongo.Client
	cacheClient thirdparty.CacheClient
	config *config.Configuration
}

func (app *app) GetMongoClient() *mongo.Client {
	return app.mongoClient
}

func (app *app) GetCacheClient() thirdparty.CacheClient {
	return app.cacheClient
}

func (app *app) GetConfig() *config.Configuration{
	return app.config
}

func NewApp(config *config.Configuration) App {
	app := &app{}

	app.mongoClient = thirdparty.NewMongoClient(config.MongoConnectionString)
	app.cacheClient = thirdparty.NewCacheClient(config.RedisConfiguration)
	app.config = config


	return app
}