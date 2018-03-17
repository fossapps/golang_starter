package migrations

import (
	"github.com/globalsign/mgo"
)

type PermissionSeeds struct{}

type Permissions struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

func (PermissionSeeds) GetKey() string {
	return "INITIAL_PERMISSION_SEEDS"
}

func (PermissionSeeds) GetDescription() string {
	return "Create basic permission systems"
}

func (PermissionSeeds) Apply(db *mgo.Database) {
	permissions := []Permissions{
		{
			Key:         "User.Create",
			Description: "Permission to create a new user",
		},
		{
			Key:         "User.Edit",
			Description: "Permission to edit user",
		},
		{
			Key:         "User.Delete",
			Description: "Permission to remove a user",
		},
		{
			Key:         "Metric.View",
			Description: "Permission to view Metric data",
		},
		{
			Key:         "Notification.Create",
			Description: "Permission to create a new Notification",
		},
		{
			Key:         "Notification.View",
			Description: "Permission to create a View existing Notification",
		},
		{
			Key:         "Notification.Delete",
			Description: "Permission to remove an existing notification",
		},
		{
			Key:         "sudo",
			Description: "Special Permission, this includes all permissions",
		},
	}
	for _, permission := range permissions {
		db.C("permissions").Insert(permission)
		// todo fires one query at a time, optimize for bulk inserts
	}
}

func (PermissionSeeds) Remove(db *mgo.Database) {

}
