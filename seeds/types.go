package seeds

import "crazy_nl_backend/helpers"

type ISeeder interface {
	GetKey() string
	GetDescription() string
	Seed(db helpers.IMongoClient)
	Remove()
}
