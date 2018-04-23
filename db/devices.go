package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"errors"
)

type Device struct {
	Token string
}

type IDeviceManager interface {
	Register(token string) error
	Exists(token string) bool
	FindByToken(token string) *Device
}

type DeviceManager struct {
	db *mgo.Database
}

func (deviceManager DeviceManager) Register(token string) error {
	if deviceManager.Exists(token) {
		return errors.New("device already exists")
	}
	return deviceManager.db.C("devices").Insert(Device{
		Token: token,
	})
}

func (deviceManager DeviceManager) Exists(token string) bool {
	return deviceManager.FindByToken(token) != nil
}

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

func GetDeviceManager(db *mgo.Database) IDeviceManager {
	return DeviceManager{
		db: db,
	}
}
