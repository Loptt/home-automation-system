package models

import (
	"context"
	"net/http"
	"time"

	"github.com/Loptt/home-automation-system/api/db"
	"github.com/Loptt/home-automation-system/api/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type eventTime struct {
	Hour   int `bson:"hour,omitempty" json:"hour"`
	Minute int `bson:"minute,omitempty" json:"minute"`
}

// eventRepetition represents how an event should be repeated. If times is greater than 0,
// the event will be repeated N times before being deleted. If date is anything later than
// Unix epoch, the event will be executed as long as the current date is less that the
// specified date.
// If both are set to their default values, then the event will be executed indefinitely.
type eventRepetition struct {
	Times   int       `bson:"times,omitempty" json:"times"`
	Date    time.Time `bson:"date,omitempty" json:"date"`
	Current int       `bson:"current,omitempty" json:"current"`
}

// Event represents a trigger event for a device
type Event struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name       string             `bson:"name,omitempty" json:"name"`
	Days       []int              `bson:"days,omitempty" json:"days"`
	Time       eventTime          `bson:"time,omitempty" json:"time"`
	Duration   int                `bson:"duration,omitempty" json:"duration"`
	Repetition eventRepetition    `bson:"repetition,omitempty" json:"repetition"`
	Device     primitive.ObjectID `bson:"device,omitempty" json:"device"`
}

// EventController interacts with the event collection in the DB.
type EventController struct {
	collection *mongo.Collection
}

// GetAll retrieves all events.
func (ec *EventController) GetAll(ctx context.Context) ([]Event, error) {
	cursor, err := ec.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)

	for cursor.Next(ctx) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// GetByID retrieves a device given a MongoDB ID.
func (ec *EventController) GetByID(ctx context.Context, id primitive.ObjectID) (*Event, error) {
	var event Event
	eventResult := ec.collection.FindOne(ctx, bson.M{"_id": id})
	if eventResult.Err() != nil {
		if eventResult.Err() == mongo.ErrNoDocuments {
			return nil, errors.NewServerError("No events found", http.StatusNotFound)
		}
		return nil, eventResult.Err()
	}

	eventResult.Decode(&event)

	return &event, nil
}

// GetByDevice retrieves all events for a device
func (ec *EventController) GetByDevice(ctx context.Context, device primitive.ObjectID) ([]Event, error) {
	cursor, err := ec.collection.Find(ctx, bson.M{"device": device})
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)

	for cursor.Next(ctx) {
		var event Event
		if err := cursor.Decode(&event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// Create a new event in the db.
func (ec *EventController) Create(ctx context.Context, event *Event) error {
	_, err := ec.collection.InsertOne(ctx, *event)
	if err != nil {
		return err
	}

	return nil
}

// NewEventController creates a DeviceController with the correct collection.
func NewEventController() (*EventController, error) {
	db, err := db.Database()
	if err != nil {
		return nil, err
	}
	return &EventController{collection: db.Collection("events")}, nil
}
