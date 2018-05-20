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

// IDeviceManager interface to satisfied for managing devices
type IDeviceManager interface {
	Register(token string) error
	Exists(token string) bool
	FindByToken(token string) *Device
}

// DeviceManager implementation of IDeviceManager
type DeviceManager struct {
	db *mgo.Database
}

// Register saves a new device to database
func (deviceManager DeviceManager) Register(token string) error {
	if deviceManager.Exists(token) {
		return errors.New("device already exists")
	}
	return deviceManager.db.C("devices").Insert(Device{
		Token: token,
	})
}

// Exists checks if a device already exists
func (deviceManager DeviceManager) Exists(token string) bool {
	return deviceManager.FindByToken(token) != nil
}

// FindByToken returns a device for a given token
func (deviceManager DeviceManager) FindByToken(token string) *Device {
	var device Device
	deviceManager.db.C("devices").Find(bson.M{
		"token": token,
	}).One(&device)
	if device.Token == "" {
		return nil
	}
	return &device
}

// GetDeviceManager returns an implementation of IDeviceManager
func GetDeviceManager(db *mgo.Database) IDeviceManager {
	return DeviceManager{
		db: db,
	}
}
