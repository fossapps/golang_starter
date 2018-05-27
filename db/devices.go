package db

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Device struct to represent a device
type Device struct {
	Token string
}

// DeviceManager interface to satisfied for managing devices
type DeviceManager interface {
	Register(token string) error
	Exists(token string) (bool, error)
	FindByToken(token string) (*Device, error)
}

// deviceManager implementation of DeviceManager
type deviceManager struct {
	db *mgo.Database
}

// Register saves a new device to database
func (deviceManager deviceManager) Register(token string) error {
	exists, err := deviceManager.Exists(token)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("device already exists")
	}
	return deviceManager.db.C("devices").Insert(Device{
		Token: token,
	})
}

// Exists checks if a device already exists
func (deviceManager deviceManager) Exists(token string) (bool, error) {
	device, err := deviceManager.FindByToken(token)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return device != nil, err
}

// FindByToken returns a device for a given token
func (deviceManager deviceManager) FindByToken(token string) (*Device, error) {
	var device Device
	err := deviceManager.db.C("devices").Find(bson.M{
		"token": token,
	}).One(&device)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if device.Token == "" {
		return nil, nil
	}
	return &device, nil
}

// GetDeviceManager returns an implementation of DeviceManager
func GetDeviceManager(db *mgo.Database) DeviceManager {
	return deviceManager{
		db: db,
	}
}
