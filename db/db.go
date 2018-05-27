package db

import (
	"github.com/fossapps/starter/config"
	"github.com/globalsign/mgo"
)

// DatabaseLayer is implementation of DB interface
type DatabaseLayer struct {
	session *mgo.Session
}

// DB implementation is what newsletter depends on
// in case some other db is preferred, all one would have to do is to simply implement the db interface
type DB interface {
	Clone() DB
	Copy() DB
	Close()
	Users() UserManager
	RefreshTokens() RefreshTokenManager
	Permissions() PermissionManager
	Migrations() MigrationManager
	Devices() DeviceManager
}

// Copy works just like New, but preserves the exact authentication
// information from the original session.
func (db DatabaseLayer) Copy() DB {
	return GetDbImplementation(db.session.Copy())
}

// Clone works just like Copy, but also reuses the same socket as the original
// session, in case it had already reserved one due to its consistency
// guarantees.  This behavior ensures that writes performed in the old session
// are necessarily observed when using the new session, as long as it was a
// strong or monotonic session.  That said, it also means that long operations
// may cause other goroutines using the original session to wait.
func (db DatabaseLayer) Clone() DB {
	return GetDbImplementation(db.session.Clone())
}

// Close terminates the session. It's a runtime error to use a session
// after it has been closed.
func (db DatabaseLayer) Close() {
	db.session.Close()
}

// Users returns UserManager
func (db DatabaseLayer) Users() UserManager {
	return GetUserManager(db.session.DB(config.GetMongoConfig().DbName))
}

// RefreshTokens returns RefreshTokenManager
func (db DatabaseLayer) RefreshTokens() RefreshTokenManager {
	return GetRefreshTokenManager(db.session.DB(config.GetMongoConfig().DbName))
}

// Permissions returns PermissionManager
func (db DatabaseLayer) Permissions() PermissionManager {
	return GetPermissionManager(db.session.DB(config.GetMongoConfig().DbName))
}

// Migrations returns MigrationManager
func (db DatabaseLayer) Migrations() MigrationManager {
	return GetMigrationManager(db.session.DB(config.GetMongoConfig().DbName))
}

// Devices returns implementation of DeviceManager
func (db DatabaseLayer) Devices() DeviceManager {
	return GetDeviceManager(db.session.DB(config.GetMongoConfig().DbName))
}

// GetDbImplementation Returns database implementation of DB interface
func GetDbImplementation(session *mgo.Session) DB {
	return &DatabaseLayer{
		session: session,
	}
}
