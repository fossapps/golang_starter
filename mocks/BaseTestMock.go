package mocks

import (
	"github.com/cyberhck/captain"
	"time"
	"fmt"
)
type ExampleMockJob struct {
}
func (ExampleMockJob) LockProvider() captain.LockProvider {
	return nil
}

func (ExampleMockJob) Job() captain.Worker {
	return func(channels captain.CommChan) {
		fmt.Println("running mock job")
	}
}

func (ExampleMockJob) ResultProcessor() captain.ResultProcessor {
	return nil
}

func (ExampleMockJob) RuntimeProcessor() captain.RuntimeProcessor {
	return func(tick time.Time, message string, startTime time.Time) {
	}
}

func (ExampleMockJob) ShouldRun(key string) bool {
	return key == "key"
}
