package application

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Cock struct {
	ID          string    `bson:"_id"`
	Size        int32     `bson:"size"`
	Nickname    string    `bson:"nickname"`
	UserID      int64     `bson:"user_id"`
	RequestedAt time.Time `bson:"requested_at"`
}

type UserCockRace struct { // Для аггрегаций только
	UserID    int64  `bson:"_id"`
	Nickname  string `bson:"nickname"`
	TotalSize int32  `bson:"total_size"`
}

func InitializeMongoConnection(ctx context.Context, log *Logger) *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.F("MONGODB_URI does not have value, set it in .env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAppName("Dickobrazz").SetTimeout(10*time.Second))
	if err != nil {
		log.F("Failed to connect to MongoDB", InnerError, err)
	}

	log.I("Successfully connected to MongoDB!")
	return client
}

func (app *Application) SaveCockToMongo(log *Logger, cock *Cock) {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	if _, err := collection.InsertOne(app.ctx, cock); err != nil {
		log.E("Failed to save cock to MongoDB", InnerError, err)
	} else {
		log.I("Successfully saved cock to MongoDB")
	}
}

func (app *Application) AggregateCockSizes(log *Logger) []UserCockRace {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$user_id"},
				{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
			}},
		},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$limit", Value: 13}},
	}

	cursor, err := collection.Aggregate(app.ctx, pipeline)
	if err != nil {
		log.E("Failed to aggregate cock sizes", InnerError, err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil { // Кааааак же похуй...
			log.E("Failed to close mongo cursor", InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse aggregation results", InnerError, err)
	}

	log.I("Successfully aggregated cock sizes")
	return results
}

func (app *Application) GetUserAggregatedCock(log *Logger, userID int64) *UserCockRace {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
	}

	cursor, err := collection.Aggregate(app.ctx, pipeline)
	if err != nil {
		log.E("Failed to aggregate user cock sizes", InnerError, err)
		return nil
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", InnerError, err)
		}
	}(cursor, app.ctx)

	var result UserCockRace
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode aggregation result", InnerError, err)
			return nil
		}
		return &result
	}

	log.I("No cocks found for user")
	return nil
}

//pipeline := `[
//	{"$match": { "user_name": "mairwunnx" }},
//	{"$group": { "_id": "$brand", "count": { "$sum": 1 } }},
//	{"$project": { "brand": "$_id", "_id": 0, "count": 1 }}
//]`
