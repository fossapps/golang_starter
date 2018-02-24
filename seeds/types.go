package seeds

import (
	"crazy_nl_backend/helpers"
	"time"
)

type ISeeder interface {
	GetKey() string
	GetDescription() string
	Seed(db helpers.IDatabase)
	Remove(db helpers.IDatabase)
}
type SeedInfo struct {
	Key string `json:"key"`
	Description string `json:"description"`
	AppliedAt time.Time `json:"applied_at"`
}
