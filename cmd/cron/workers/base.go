package workers

import "github.com/cyberhck/captain"

type ICronJob interface {
	Job() captain.Worker
	LockProvider() captain.LockProvider
	ResultProcessor() captain.ResultProcessor
	RuntimeProcessor() captain.RuntimeProcessor
	ShouldRun(key string) bool
}

func Run(key string, task ICronJob) {
	if !task.ShouldRun(key) {
		return
	}
	job := captain.CreateJob()
	job.WithRuntimeProcessor(task.RuntimeProcessor())
	job.WithResultProcessor(task.ResultProcessor())
	job.WithLockProvider(task.LockProvider())
	job.SetWorker(task.Job())
	job.Run()
}
