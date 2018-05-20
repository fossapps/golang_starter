package db

import (
	"starter/config"

	"github.com/globalsign/mgo"
)

// DatabaseLayer is implementation of Db interface
type DatabaseLayer struct {
	session *mgo.Session
}

// Db implementation is what newsletter depends on
// in case some other db is preferred, all one would have to do is to simply implement the db interface
type Db interface {
	Clone() Db
	Copy() Db
	Close()
	Users() IUserManager
	RefreshTokens() IRefreshTokenManager
	Permissions() IPermissionManager
	Migrations() IMigrationManager
	Devices() IDeviceManager
}

// Copy works just like New, but preserves the exact authentication
// information from the original session.
func (db DatabaseLayer) Copy() Db {
	return GetDbImplementation(db.session.Copy())
}

// Clone works just like Copy, but also reuses the same socket as the original
// session, in case it had already reserved one due to its consistency
// guarantees.  This behavior ensures that writes performed in the old session
// are necessarily observed when using the new session, as long as it was a
// strong or monotonic session.  That said, it also means that long operations
// may cause other goroutines using the original session to wait.
func (db DatabaseLayer) Clone() Db {
	return GetDbImplementation(db.session.Clone())
}

// Close terminates the session. It's a runtime error to use a session
// after it has been closed.
func (db DatabaseLayer) Close() {
	db.session.Close()
}

// Users returns UserManager
func (db DatabaseLayer) Users() IUserManager {
	return GetUserManager(db.session.DB(config.GetMongoConfig().DbName))
}

// RefreshTokens returns RefreshTokenManager
func (db DatabaseLayer) RefreshTokens() IRefreshTokenManager {
	return GetRefreshTokenManager(db.session.DB(config.GetMongoConfig().DbName))
}

// Permissions returns PermissionManager
func (db DatabaseLayer) Permissions() IPermissionManager {
	return GetPermissionManager(db.session.DB(config.GetMongoConfig().DbName))
}

// Migrations returns MigrationManager
func (db DatabaseLayer) Migrations() IMigrationManager {
	return GetMigrationManager(db.session.DB(config.GetMongoConfig().DbName))
}

// Devices returns implementation of IDeviceManager
func (db DatabaseLayer) Devices() IDeviceManager {
	return GetDeviceManager(db.session.DB(config.GetMongoConfig().DbName))
}

// GetDbImplementation Returns database implementation of Db interface
func GetDbImplementation(session *mgo.Session) Db {
	return &DatabaseLayer{
		session: session,
	}
}
