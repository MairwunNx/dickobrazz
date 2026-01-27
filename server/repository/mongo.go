package repository

import "go.mongodb.org/mongo-driver/mongo"

type MongoRepository struct {
	db *mongo.Database
}

func NewMongoRepository(client *mongo.Client, databaseName string) *MongoRepository {
	return &MongoRepository{db: client.Database(databaseName)}
}

func (r *MongoRepository) Database() *mongo.Database {
	return r.db
}
