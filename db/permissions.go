package db

import (
	"errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Permission struct, key is the unique permission key, description is human readable description of permission
type Permission struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

// PermissionManager deals with Permission persistence
type PermissionManager interface {
	Create(key string, description string) error
	Exists(key string) bool
	List() ([]Permission, error)
}

// permissionManager struct which implements PermissionManager
type permissionManager struct {
	db *mgo.Database
}

// Create a permission given key and description
func (pLayer permissionManager) Create(key string, description string) error {
	if pLayer.Exists(key) {
		return errors.New("permission already exists")
	}
	return pLayer.db.C("permissions").Insert(Permission{
		Description: description,
		Key:         key,
	})
}

// Exists returns weather or not a permission already exists
func (pLayer permissionManager) Exists(key string) bool {
	var perm Permission
	pLayer.db.C("permissions").Find(bson.M{
		"key": key,
	}).One(&perm)
	return perm.Key == key
}

// List all permissions available
func (pLayer permissionManager) List() ([]Permission, error) {
	var permission []Permission
	err := pLayer.db.C("permissions").Find(nil).All(&permission)
	return permission, err
}

// GetPermissionManager returns an implementation of PermissionManager
func GetPermissionManager(db *mgo.Database) PermissionManager {
	return permissionManager{
		db: db,
	}
}
