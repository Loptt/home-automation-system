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

// Configuration represents a single user global config parameters
type Configuration struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	User            primitive.ObjectID `bson:"user,omitempty" json:"user"`
	SystemStatus    string             `bson:"systemStatus,omitempty" json:"systemStatus"`
	RainPercentage  int                `bson:"rainPercentage,omitempty" json:"rainPercentage"`
	DefaultDuration int                `bson:"defaultDuration,omitempty" json:"defaultDuration"`
	Update          bool               `bson:"update,omitempty" json:"update"`
}

// ConfigurationController interacts with the configuration collection in the DB.
type ConfigurationController struct {
	collection *mongo.Collection
}

// GetAll retrieves all deivces.
func (dc *ConfigurationController) GetAll(ctx context.Context) ([]Configuration, error) {
	cursor, err := dc.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	configurations := make([]Configuration, 0)

	for cursor.Next(ctx) {
		var configuration Configuration
		if err := cursor.Decode(&configuration); err != nil {
			return nil, err
		}
		configurations = append(configurations, configuration)
	}

	return configurations, nil
}

// GetByID retrieves a configuration given a MongoDB ID.
func (dc *ConfigurationController) GetByID(ctx context.Context, id primitive.ObjectID) (*Configuration, error) {
	var configuration Configuration
	configurationResult := dc.collection.FindOne(ctx, bson.M{"_id": id})
	if configurationResult.Err() != nil {
		if configurationResult.Err() == mongo.ErrNoDocuments {
			return nil, errors.NewServerError("No configurations found", http.StatusNotFound)
		}
		return nil, configurationResult.Err()
	}

	configurationResult.Decode(&configuration)

	return &configuration, nil
}

// GetByUser retrieves a configuration given a user ID.
func (dc *ConfigurationController) GetByUser(ctx context.Context, user primitive.ObjectID) (*Configuration, error) {
	var configuration Configuration
	configurationResult := dc.collection.FindOne(ctx, bson.M{"user": user})
	if configurationResult.Err() != nil {
		if configurationResult.Err() == mongo.ErrNoDocuments {
			return nil, errors.NewServerError("No configurations found for user", http.StatusNotFound)
		}
		return nil, configurationResult.Err()
	}

	configurationResult.Decode(&configuration)

	return &configuration, nil
}

// Create a new configuration in the db.
func (dc *ConfigurationController) Create(ctx context.Context, configuration *Configuration) error {
	_, err := dc.collection.InsertOne(ctx, *configuration)
	if err != nil {
		return err
	}

	return nil
}

// Update an existing configuration in the db.
func (dc *ConfigurationController) Update(ctx context.Context, id primitive.ObjectID, configuration *Configuration) error {
	updateResult := dc.collection.FindOneAndReplace(ctx, bson.M{"_id": id}, *configuration)
	if updateResult.Err() != nil {
		if updateResult.Err() == mongo.ErrNoDocuments {
			return errors.NewServerError("No configurations found", http.StatusNotFound)
		}
		return updateResult.Err()
	}

	return nil
}

// Delete a configuration given an ID
func (dc *ConfigurationController) Delete(ctx context.Context, id primitive.ObjectID) error {
	deleteResult := dc.collection.FindOneAndDelete(ctx, bson.M{"_id": id})
	if deleteResult.Err() != nil {
		if deleteResult.Err() == mongo.ErrNoDocuments {
			return errors.NewServerError("No configurations found", http.StatusNotFound)
		}
		return deleteResult.Err()
	}

	return nil
}

// NewConfigurationController creates a ConfigurationController with the correct collection.
func NewConfigurationController() (*ConfigurationController, error) {
	db, err := db.Database()
	if err != nil {
		return nil, err
	}
	return &ConfigurationController{collection: db.Collection("configurations")}, nil
}
