package db

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

// UserManager deals with anything related to User Persistence
type UserManager interface {
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Create(user User) error
	List() ([]User, error)
	Update(id string, user User) error
}

// User is representation of a user
type User struct {
	ID          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Email       string        `json:"email"`
	Password    string        `json:"password"`
	Permissions []string      `json:"permissions"`
}

// UserLayer is implementation of UserManager interface
type UserLayer struct {
	db *mgo.Database
}

// FindByEmail given email, returns a User
func (dbLayer UserLayer) FindByEmail(email string) (*User, error) {
	var user User
	err := dbLayer.db.C("users").Find(bson.M{
		"email": email,
	}).One(&user)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if user.Email == "" {
		return nil, nil
	}
	return &user, nil
}

// FindByID given id of user, returns a User
func (dbLayer UserLayer) FindByID(id string) (*User, error) {
	user := new(User)
	err := dbLayer.db.C("users").Find(bson.M{
		"_id": bson.ObjectIdHex(id),
	}).One(&user)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, nil
	}
	return user, nil
}

// Create a user, returns error if there's one
func (dbLayer UserLayer) Create(user User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hash)
	return dbLayer.db.C("users").Insert(user)
}

// List users
func (dbLayer UserLayer) List() ([]User, error) {
	var users []User
	err := dbLayer.db.C("users").Find(nil).All(&users)
	return users, err
}

// Update a user by id
func (dbLayer UserLayer) Update(id string, user User) error {
	// todo if there's password, it needs to be crypted as well.
	return dbLayer.db.C("users").Update(bson.M{
		"_id": bson.ObjectIdHex(id),
	}, user)
}

// GetUserManager returns implementation of UserManager
func GetUserManager(db *mgo.Database) UserManager {
	return &UserLayer{
		db: db,
	}
}
