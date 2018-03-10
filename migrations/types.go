package migrations

import (
	"crazy_nl_backend/helpers"
	"time"
)

type IMigration interface {
	GetKey() string
	GetDescription() string
	Apply(db helpers.IDatabase)
	Remove(db helpers.IDatabase)
}
type MigrationInfo struct {
	Key string `json:"key"`
	Description string `json:"description"`
	AppliedAt time.Time `json:"applied_at"`
}
