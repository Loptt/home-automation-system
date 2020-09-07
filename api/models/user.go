package models

import (
	"context"
	"net/http"

	"github.com/Loptt/home-automation-system/api/db"
	"github.com/Loptt/home-automation-system/api/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user in a controller.
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Email    string             `bson:"email,omitempty" json:"email"`
	Password string             `bson:"password,omitempty" json:"password"`
}

// UserController interacts with the user collection in the DB.
type UserController struct {
	collection *mongo.Collection
}

// GetAll retrieves all users.
func (uc *UserController) GetAll(ctx context.Context) ([]User, error) {
	cursor, err := uc.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetByID retrieves a user given a MongoDB ID.
func (uc *UserController) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	userResult := uc.collection.FindOne(ctx, bson.M{"_id": id})
	if userResult.Err() != nil {
		if userResult.Err() == mongo.ErrNoDocuments {
			return nil, errors.NewServerError("No users found", http.StatusNotFound)
		}
		return nil, userResult.Err()
	}

	userResult.Decode(&user)

	return &user, nil
}

// Create a new user in the db.
func (uc *UserController) Create(ctx context.Context, user *User) error {
	_, err := uc.collection.InsertOne(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}

// Update an existing user in the db.
func (uc *UserController) Update(ctx context.Context, id primitive.ObjectID, user *User) error {
	updateResult := uc.collection.FindOneAndReplace(ctx, bson.M{"_id": id}, *user)
	if updateResult.Err() != nil {
		if updateResult.Err() == mongo.ErrNoDocuments {
			return errors.NewServerError("No users found", http.StatusNotFound)
		}
		return updateResult.Err()
	}

	return nil
}

// Delete a device given an ID
func (uc *UserController) Delete(ctx context.Context, id primitive.ObjectID) error {
	deleteResult := uc.collection.FindOneAndDelete(ctx, bson.M{"_id": id})
	if deleteResult.Err() != nil {
		if deleteResult.Err() == mongo.ErrNoDocuments {
			return errors.NewServerError("No devices found", http.StatusNotFound)
		}
		return deleteResult.Err()
	}

	return nil
}

// NewUserController creates a UserController with the correct collection.
func NewUserController() (*UserController, error) {
	db, err := db.Database()
	if err != nil {
		return nil, err
	}
	return &UserController{collection: db.Collection("users")}, nil
}
