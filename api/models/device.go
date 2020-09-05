package models

import (
	"context"

	"github.com/Loptt/home-automation-system/api/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Device represents a device in a controller.
type Device struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Pin  int32              `bson:"pin,omitempty" json:"pin"`
	Name string             `bson:"name,omitempty" json:"name"`
}

// DeviceController interacts with the device collection in the DB.
type DeviceController struct {
	collection *mongo.Collection
}

// GetAll retrieves all deivces.
func (dc *DeviceController) GetAll(ctx context.Context) ([]Device, error) {
	cursor, err := dc.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0)

	for cursor.Next(ctx) {
		var device Device
		if err := cursor.Decode(&device); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// GetByID retrieves a device given a MongoDB ID.
func (dc *DeviceController) GetByID(ctx context.Context, id primitive.ObjectID) (*Device, error) {
	var device Device
	deviceResult := dc.collection.FindOne(ctx, bson.M{"_id": id})
	if deviceResult.Err() != nil {
		return nil, deviceResult.Err()
	}

	deviceResult.Decode(&device)

	return &device, nil
}

// Create a new device in the db.
func (dc *DeviceController) Create(ctx context.Context, device *Device) error {
	_, err := dc.collection.InsertOne(ctx, *device)
	if err != nil {
		return err
	}

	return nil
}

// Delete a device given an ID
func (dc *DeviceController) Delete(ctx context.Context, id primitive.ObjectID) error {
	deleteResult := dc.collection.FindOneAndDelete(ctx, bson.M{"_id": id})
	if deleteResult.Err() != nil {
		return deleteResult.Err()
	}

	return nil
}

// NewDeviceController creates a DeviceController with the correct collection.
func NewDeviceController() (*DeviceController, error) {
	db, err := db.Database()
	if err != nil {
		return nil, err
	}
	return &DeviceController{collection: db.Collection("devices")}, nil
}
