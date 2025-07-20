package mongodb

import (
	"context"
	"fmt"
	"gocrudapp/model"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoClient struct {
	Client mongo.Collection
}

func (c MongoClient) CreateUser(user model.User) (string, error) {
	res, err := c.Client.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(bson.ObjectID).Hex(), nil
}

func (c MongoClient) GetUserByID(id string) (model.User, error) {
	docId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return model.User{}, fmt.Errorf("invalid ID format: %v", err)
	}

	var user  model.User
	filter := bson.D{{Key: "_id", Value: docId}}

	err = c.Client.FindOne(context.Background(), filter).Decode(&user)

	if err == mongo.ErrNoDocuments {
		return model.User{}, fmt.Errorf("user not found")
	} else if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (c MongoClient) GetAllUsers() ([]model.User, error) {
	filter := bson.D{}

	cur, err := c.Client.Find(context.Background(), filter)

	if err != nil {
		return []model.User{}, err
	}

	defer cur.Close(context.Background())

	var users []model.User
	for cur.Next(context.Background()) {
		var user model.User
		err := cur.Decode(&user)

		if err != nil {
			slog.Error("error decoding user", slog.String("error", err.Error()))
			continue // skip this user if there's an error decoding
		} else {
			users = append(users, user)
		}
	}
	fmt.Printf("Found %d users\n", len(users))
	return users, nil
}

func (c MongoClient) UpdateUserAgeByID(id string, age int) (int, error) {
	docId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: docId}}
	updateStmt := bson.D{{Key: "$set", Value: bson.D{{Key: "age", Value: age}}}}

	res, err := c.Client.UpdateOne(context.Background(), filter, updateStmt)

	if err != nil {
		return 0, fmt.Errorf("error updating user age: %v", err)
	}

	return int(res.ModifiedCount), nil
}

func (c MongoClient) DeleteUserByID(id string) (int, error) {
	docId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.D{{Key: "_id", Value: docId}}

	res, err := c.Client.DeleteOne(context.Background(), filter)

	if res.DeletedCount == 0 {
		return 0, fmt.Errorf("no user found with ID: %s", id)
	} else if err != nil {
		return 0, fmt.Errorf("error deleting user: %v", err)
	}
	
	slog.Info("User deleted successfully", slog.String("id", id))

	return int(res.DeletedCount), nil
}

func (c MongoClient) DeleteAllUsers() (int, error) {
	filter := bson.D{}
	res, err := c.Client.DeleteMany(context.Background(), filter)

	if err != nil {
		return 0, fmt.Errorf("error deleting users: %v", err)
	}

	return int(res.DeletedCount), nil
}