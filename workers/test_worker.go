package workers

import (
	"time"

	"github.com/cyberhck/captain"
)

// TestWorker a sample worker
type TestWorker struct {}

// LockProvider no lock provider implementation for now
func (job TestWorker) LockProvider() captain.LockProvider {
	return nil
}

// Job actual task to be performed by worker
func (job TestWorker) Job() captain.Worker {
	return func(channels captain.CommChan) {
		time.Sleep(2 * time.Second)
	}
}

// ResultProcessor act on result
func (job TestWorker) ResultProcessor() captain.ResultProcessor {
	return nil
}

// RuntimeProcessor monitors worker every tick
func (job TestWorker) RuntimeProcessor() captain.RuntimeProcessor {
	return func(tick time.Time, message string, startTime time.Time) {
	}
}

// ShouldRun returns weather or not this worker should run given a particular key
func (job TestWorker) ShouldRun(key string) bool {
	return key == "test_worker"
}
