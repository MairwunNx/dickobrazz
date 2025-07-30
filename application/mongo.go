package application

import (
	"context"
	"dickobot/application/database"
	"dickobot/application/logging"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func InitializeMongoConnection(ctx context.Context, log *logging.Logger) *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.F("MONGODB_URI does not have value, set it in .env file")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAppName("Dickobrazz").SetTimeout(10*time.Second))
	if err != nil {
		log.F("Failed to connect to MongoDB", logging.InnerError, err)
	}

	log.I("Successfully connected to MongoDB!")
	return client
}

func (app *Application) SaveCockToMongo(log *logging.Logger, cock *Cock) {
	collection := database.CollectionCocks(app.db)

	if _, err := collection.InsertOne(app.ctx, cock); err != nil {
		log.E("Failed to save cock to MongoDB", logging.InnerError, err)
	} else {
		log.I("Successfully saved cock to MongoDB")
	}
}

func (app *Application) AggregateCockSizes(log *logging.Logger) []UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineTopUsersBySize())
	if err != nil {
		log.E("Failed to aggregate cock sizes", logging.InnerError, err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil { // Кааааак же похуй...
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse aggregation results", logging.InnerError, err)
	}

	log.I("Successfully aggregated cock sizes")
	return results
}

func (app *Application) GetUserAggregatedCock(log *logging.Logger, userID int64) *UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineUserTotalSize(userID))
	if err != nil {
		log.E("Failed to aggregate user cock sizes", logging.InnerError, err)
		return nil
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result UserCockRace
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode aggregation result", logging.InnerError, err)
			return nil
		}
		return &result
	}

	log.I("No cocks found for user")
	return nil
}