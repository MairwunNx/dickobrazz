package database

import "go.mongodb.org/mongo-driver/mongo"

var cocks *mongo.Collection
var achievements *mongo.Collection

func CollectionCocks(db *mongo.Client) *mongo.Collection {
	if cocks == nil {
		cocks = db.Database("dickbot_db").Collection("cocks")
	}
	return cocks
}

func CollectionAchievements(db *mongo.Client) *mongo.Collection {
	if achievements == nil {
		achievements = db.Database("dickbot_db").Collection("achievements")
	}
	return achievements
}
