package database

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	cocks        *mongo.Collection
	cocksOnce    sync.Once
	achievements *mongo.Collection
	achOnce      sync.Once
)

func CollectionCocks(db *mongo.Client) *mongo.Collection {
	cocksOnce.Do(func() {
		cocks = db.Database("dickbot_db").Collection("cocks")
	})
	return cocks
}

func CollectionAchievements(db *mongo.Client) *mongo.Collection {
	achOnce.Do(func() {
		achievements = db.Database("dickbot_db").Collection("achievements")
	})
	return achievements
}
