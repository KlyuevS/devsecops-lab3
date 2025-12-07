package mongo

import (
	"context"
	"go-news/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Хранилище данных.
type Store struct {
	db *mongo.Client
}

const (
	dbName         = "go-news"
	collectionName = "posts"
)

// Конструктор объекта хранилища.
func New(connectionString string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	// не отключаем клиент сразу — вернём подключение в Store
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: client,
	}
	return &s, err
}

func (s *Store) Tasks() ([]storage.Task, error) {
	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var posts []storage.Task
	for cur.Next(context.Background()) {
		var p storage.Task
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

func (s *Store) AddTask(p storage.Task) error {
	collection := s.db.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return err
}
func (s *Store) UpdateTask(p storage.Task) error {
	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{{Key: "ID", Value: p.ID}}
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "ResponsibleID", Value: p.ResponsibleID},
			{Key: "ResponsibleName", Value: p.ResponsibleName},
			{Key: "Context", Value: p.Context},
			{Key: "DueDate", Value: p.DueDate},
			{Key: "AssignedAt", Value: p.AssignedAt},
		},
	}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return err
}
func (s *Store) DeleteTask(p storage.Task) error {
	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{{Key: "ID", Value: p.ID}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return err
}
