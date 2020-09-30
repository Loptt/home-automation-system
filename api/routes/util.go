package routes

import (
	"context"

	"github.com/Loptt/home-automation-system/api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setUpdateOn(ctx context.Context, deviceID primitive.ObjectID) error {
	dc, err := models.NewDeviceController()
	if err != nil {
		return err
	}

	uc, err := models.NewUserController()
	if err != nil {
		return err
	}

	cc, err := models.NewConfigurationController()
	if err != nil {
		return err
	}

	device, err := dc.GetByID(ctx, deviceID)
	if err != nil {
		return err
	}

	user, err := uc.GetByID(ctx, device.User)
	if err != nil {
		return err
	}

	configuration, err := cc.GetByUser(ctx, user.ID)
	if err != nil {
		return err
	}

	configuration.Update = true

	if err := cc.Update(ctx, configuration.ID, configuration); err != nil {
		return err
	}

	return nil
}
