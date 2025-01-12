package database

import "go.mongodb.org/mongo-driver/mongo"

var cocks *mongo.Collection

func CollectionCocks(db *mongo.Client) *mongo.Collection {
	if cocks == nil {
		cocks = db.Database("dickbot_db").Collection("cocks")
	}
	return cocks
}
