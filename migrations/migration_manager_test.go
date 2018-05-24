package migrations_test

import (
	"testing"

	"github.com/fossapps/starter/migrations"
	"github.com/fossapps/starter/mocks"

	"github.com/golang/mock/gomock"
)

func TestSeedCallsFirstTime(t *testing.T) {
	dbCtrl := gomock.NewController(t)
	mockCtrl := gomock.NewController(t)
	migrationCtrl := gomock.NewController(t)
	mockMigrationManager := mocks.NewMockMigrationManager(migrationCtrl)
	mockDatabase := mocks.NewMockDB(dbCtrl)
	mockMigration := mocks.NewMockMigration(mockCtrl)
	// setup migration behavior
	mockMigration.EXPECT().GetKey().MinTimes(1).Return("sample_key")
	mockMigration.EXPECT().GetDescription().MinTimes(1).Return("description")
	mockMigration.EXPECT().Apply(gomock.Any())
	// setup mockMigrationManager behavior
	mockMigrationManager.EXPECT().ShouldRun("sample_key").Return(true)
	mockMigrationManager.EXPECT().MarkApplied("sample_key", "description")
	// setup database behavior
	mockDatabase.EXPECT().Migrations().MinTimes(1).Return(mockMigrationManager)
	migrations.Apply(mockMigration, mockDatabase)
}

func TestSeedDoesNotExecuteDuplicates(t *testing.T) {
	dbCtrl := gomock.NewController(t)
	mockCtrl := gomock.NewController(t)
	migrationCtrl := gomock.NewController(t)
	mockMigrationManager := mocks.NewMockMigrationManager(migrationCtrl)
	mockDatabase := mocks.NewMockDB(dbCtrl)
	mockMigration := mocks.NewMockMigration(mockCtrl)
	// setup migration behavior
	mockMigration.EXPECT().GetKey().MinTimes(1).Return("sample_key")
	mockMigration.EXPECT().Apply(gomock.Any()).Times(0)
	// setup mockMigrationManager behavior
	mockMigrationManager.EXPECT().ShouldRun("sample_key").Return(false)
	mockMigrationManager.EXPECT().MarkApplied("sample_key", "description")
	// setup database behavior
	mockDatabase.EXPECT().Migrations().Times(1).Return(mockMigrationManager)
	migrations.Apply(mockMigration, mockDatabase)
}
