package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PipelineActiveUsersSince(since time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "requested_at", Value: bson.D{{Key: "$gte", Value: since}}}}}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$user_id"}}}},
		{{Key: "$count", Value: "total"}},
	}
}

func PipelineSizeDistribution() mongo.Pipeline {
	return mongo.Pipeline{
		{{
			Key: "$bucket",
			Value: bson.D{
				{Key: "groupBy", Value: "$size"},
				{Key: "boundaries", Value: bson.A{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 62}},
				{Key: "default", Value: "other"},
				{Key: "output", Value: bson.D{{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}},
			},
		}},
	}
}
