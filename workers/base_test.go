package workers_test

import (
	"testing"
	"github.com/golang/mock/gomock"
	"crazy_nl_backend/mocks"
	"crazy_nl_backend/workers"
	"github.com/cyberhck/captain"
)

func TestRunDoesNotCallAnythingIfItShouldNotRun(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	job := mocks.NewMockICronJob(mockCtrl)
	job.EXPECT().ShouldRun("key").Return(false)
	job.EXPECT().Job().Times(0)
	job.EXPECT().LockProvider().Times(0)
	job.EXPECT().ResultProcessor().Times(0)
	workers.Run("key", job)
}

func TestRunCallsJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	job := mocks.NewMockICronJob(mockCtrl)
	job.EXPECT().ShouldRun("key").Return(true)
	job.EXPECT().Job().Times(1).Return(func(channel captain.CommChan) {})
	job.EXPECT().RuntimeProcessor().Times(1)
	job.EXPECT().ResultProcessor().Times(1)
	job.EXPECT().LockProvider().Times(1)
	workers.Run("key", job)
}

func ExampleRun() {
	mock := mocks.ExampleMockJob{}
	workers.Run("key", mock)
	// Output:
	// running mock job
}
