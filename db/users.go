package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

// IUserManager deals with anything related to User Persistence
type IUserManager interface {
	FindByEmail(email string) *User
	FindById(id string) *User
	Create(user User) error
}

// User is representation of a user
type User struct {
	ID          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Permissions []string      `json:"permissions"`
}

// UserLayer is implementation of IUserManager interface
type UserLayer struct {
	db *mgo.Database
}

// FindByEmail given email, returns a User
func (dbLayer UserLayer) FindByEmail(email string) *User {
	var user User
	dbLayer.db.C("users").Find(bson.M{
		"email": email,
	}).One(&user)
	if user.Email == "" {
		return nil
	}
	return &user
}

// FindById given id of user, returns a User
func (dbLayer UserLayer) FindById(id string) *User {
	user := new(User)
	dbLayer.db.C("users").Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(&user)
	if user.ID == "" {
		return nil
	}
	return user
}

// Create a user, returns error if there's one
func (dbLayer UserLayer) Create(user User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	return dbLayer.db.C("users").Insert(user)
}

// GetUserManager returns implementation of IUserManager
func GetUserManager(db *mgo.Database) IUserManager {
	return &UserLayer{
		db: db,
	}
}
