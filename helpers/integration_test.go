package helpers_test

import "crazy_nl_backend/helpers"

func DropDb(db helpers.IDatabase) {
	db.DropDatabase()
}