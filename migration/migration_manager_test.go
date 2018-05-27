package migration_test

import (
	"testing"

	"github.com/fossapps/starter/migration"
	"github.com/fossapps/starter/mock"

	"github.com/golang/mock/gomock"
)

func TestSeedCallsFirstTime(t *testing.T) {
	dbCtrl := gomock.NewController(t)
	mockCtrl := gomock.NewController(t)
	migrationCtrl := gomock.NewController(t)
	mockMigrationManager := mock.NewMockMigrationManager(migrationCtrl)
	mockDatabase := mock.NewMockDB(dbCtrl)
	mockMigration := mock.NewMockMigration(mockCtrl)
	// setup migration behavior
	mockMigration.EXPECT().GetKey().MinTimes(1).Return("sample_key")
	mockMigration.EXPECT().GetDescription().MinTimes(1).Return("description")
	mockMigration.EXPECT().Apply(gomock.Any())
	// setup mockMigrationManager behavior
	mockMigrationManager.EXPECT().IsApplied("sample_key").Return(false, nil)
	mockMigrationManager.EXPECT().MarkApplied("sample_key", "description")
	// setup database behavior
	mockDatabase.EXPECT().Migrations().MinTimes(1).Return(mockMigrationManager)
	migration.Apply(mockMigration, mockDatabase)
}

func TestSeedDoesNotExecuteDuplicates(t *testing.T) {
	dbCtrl := gomock.NewController(t)
	mockCtrl := gomock.NewController(t)
	migrationCtrl := gomock.NewController(t)
	mockMigrationManager := mock.NewMockMigrationManager(migrationCtrl)
	mockDatabase := mock.NewMockDB(dbCtrl)
	mockMigration := mock.NewMockMigration(mockCtrl)
	// setup migration behavior
	mockMigration.EXPECT().GetKey().MinTimes(1).Return("sample_key")
	mockMigration.EXPECT().Apply(gomock.Any()).Times(0)
	// setup mockMigrationManager behavior
	mockMigrationManager.EXPECT().IsApplied("sample_key").Return(true, nil)
	mockMigrationManager.EXPECT().MarkApplied("sample_key", "description")
	// setup database behavior
	mockDatabase.EXPECT().Migrations().Times(1).Return(mockMigrationManager)
	migration.Apply(mockMigration, mockDatabase)
}
