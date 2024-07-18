package mongodb

import (
	"Task/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskMongoDb struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewTask() (*TaskMongoDb, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	collection := client.Database("Tasks").Collection("task")
	return &TaskMongoDb{client: client, collection: collection}, nil
}

func (t *TaskMongoDb) StoreNewTask(task models.Tasks) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := t.collection.InsertOne(ctx, task)
	return err
}

func (t *TaskMongoDb) StoreGetTasks() ([]*models.Tasks, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tasks []*models.Tasks
	cursor, err := t.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TaskMongoDb) StoreUpdateTask(id primitive.ObjectID, tasks models.Tasks) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := t.collection.UpdateByID(ctx, id, bson.M{"$set": tasks})
	return err
}

func (t *TaskMongoDb) StoreDeleteTask(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := t.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err 
}
