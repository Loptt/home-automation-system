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

// Device represents a device in a controller.
type Device struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Pin    int32              `bson:"pin,omitempty" json:"pin"`
	Name   string             `bson:"name,omitempty" json:"name"`
	User   primitive.ObjectID `bson:"user,omitempty" json:"user"`
	Status string             `bson:"status,omitempty" json:"status"`
}

// DeviceEvents represents a device with populated events.
type DeviceEvents struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Pin    int32              `bson:"pin,omitempty" json:"pin"`
	Name   string             `bson:"name,omitempty" json:"name"`
	User   primitive.ObjectID `bson:"user,omitempty" json:"user"`
	Status string             `bson:"status,omitempty" json:"status"`
	Events []Event            `bson:"events,omitempty" json:"events"`
}

// SetDevice initializes a DeviceEvents with a given Device object.
func (de *DeviceEvents) SetDevice(d Device) {
	de.ID = d.ID
	de.Pin = d.Pin
	de.Name = d.Name
	de.User = d.User
	de.Status = d.Status
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
		if deviceResult.Err() == mongo.ErrNoDocuments {
			return nil, errors.NewServerError("No devices found", http.StatusNotFound)
		}
		return nil, deviceResult.Err()
	}

	deviceResult.Decode(&device)

	return &device, nil
}

// GetByUser retrieves a devices given a user ID.
func (dc *DeviceController) GetByUser(ctx context.Context, user primitive.ObjectID) ([]Device, error) {
	cursor, err := dc.collection.Find(ctx, bson.M{"user": user})
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

// GetByUserWithEvents retrieves a devices given a user ID with all the associated events.
func (dc *DeviceController) GetByUserWithEvents(ctx context.Context, user primitive.ObjectID) ([]DeviceEvents, error) {
	cursor, err := dc.collection.Find(ctx, bson.M{"user": user})
	if err != nil {
		return nil, err
	}

	ec, err := NewEventController()
	if err != nil {
		return nil, err
	}

	devices := make([]DeviceEvents, 0)

	for cursor.Next(ctx) {
		var device Device
		var deviceEvent DeviceEvents
		if err := cursor.Decode(&device); err != nil {
			return nil, err
		}

		deviceEvent.SetDevice(device)

		events, err := ec.GetByDevice(ctx, device.ID)
		if err != nil {
			return nil, err
		}

		deviceEvent.Events = events

		devices = append(devices, deviceEvent)
	}

	return devices, nil
}

// Create a new device in the db.
func (dc *DeviceController) Create(ctx context.Context, device *Device) error {
	_, err := dc.collection.InsertOne(ctx, *device)
	if err != nil {
		return err
	}

	return nil
}

// Update an existing device in the db.
func (dc *DeviceController) Update(ctx context.Context, id primitive.ObjectID, device *Device) error {
	updateResult := dc.collection.FindOneAndReplace(ctx, bson.M{"_id": id}, *device)
	if updateResult.Err() != nil {
		if updateResult.Err() == mongo.ErrNoDocuments {
			return errors.NewServerError("No devices found", http.StatusNotFound)
		}
		return updateResult.Err()
	}

	return nil
}

// Delete a device given an ID
func (dc *DeviceController) Delete(ctx context.Context, id primitive.ObjectID) error {
	deleteResult := dc.collection.FindOneAndDelete(ctx, bson.M{"_id": id})
	if deleteResult.Err() != nil {
		if deleteResult.Err() == mongo.ErrNoDocuments {
			return errors.NewServerError("No devices found", http.StatusNotFound)
		}
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
