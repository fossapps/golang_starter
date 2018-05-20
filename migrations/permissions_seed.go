package migrations

import (
	"starter/db"

	"github.com/globalsign/mgo"
)

// PermissionSeeds seed for adding initial list of permissions
type PermissionSeeds struct{}

// Permissions struct to hold permission information
type Permissions struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

// GetKey returns key for permission seeds
func (PermissionSeeds) GetKey() string {
	return "INITIAL_PERMISSION_SEEDS"
}

// GetDescription returns description for permission seeds
func (PermissionSeeds) GetDescription() string {
	return "Create basic permission systems"
}

// Apply adds permission to database
func (PermissionSeeds) Apply(dbLayer db.Db) {
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
			Key:         "User.List",
			Description: "Permission to list all users",
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
			Key:         "Permission.List",
			Description: "Permission to list all the permission there exists",
		},
		{
			Key:         "sudo",
			Description: "Special Permission, this includes all permissions",
		},
	}
	for _, permission := range permissions {
		dbLayer.Permissions().Create(permission.Key, permission.Description)
		// todo fires one query at a time, optimize for bulk inserts
	}
}

// Remove remove un does the migration
func (PermissionSeeds) Remove(db *mgo.Database) {

}
