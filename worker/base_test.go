package worker_test

import (
	"testing"

	"github.com/cyberhck/captain"
	"github.com/fossapps/starter/mock"
	"github.com/fossapps/starter/worker"
	"github.com/golang/mock/gomock"
)

func TestRunDoesNotCallAnythingIfItShouldNotRun(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	job := mock.NewMockICronJob(mockCtrl)
	job.EXPECT().ShouldRun("key").Return(false)
	job.EXPECT().Job().Times(0)
	job.EXPECT().LockProvider().Times(0)
	job.EXPECT().ResultProcessor().Times(0)
	worker.Run("key", job)
}

func TestRunCallsJob(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	job := mock.NewMockICronJob(mockCtrl)
	job.EXPECT().ShouldRun("key").Return(true)
	job.EXPECT().Job().Times(1).Return(func(channel captain.CommChan) {})
	job.EXPECT().RuntimeProcessor().Times(1)
	job.EXPECT().ResultProcessor().Times(1)
	job.EXPECT().LockProvider().Times(1)
	worker.Run("key", job)
}

func ExampleRun() {
	mock := mock.ExampleMockJob{}
	worker.Run("key", mock)
	// Output:
	// running mock job
}
