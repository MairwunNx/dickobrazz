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
	users        *mongo.Collection
	usersOnce    sync.Once
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

func CollectionUsers(db *mongo.Client) *mongo.Collection {
	usersOnce.Do(func() {
		users = db.Database("dickbot_db").Collection("users")
	})
	return users
}
