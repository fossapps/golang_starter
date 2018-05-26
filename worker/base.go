package worker

import "github.com/cyberhck/captain"

// CronJob any cron must satisfy this permission to run
type CronJob interface {
	Job() captain.Worker
	LockProvider() captain.LockProvider
	ResultProcessor() captain.ResultProcessor
	RuntimeProcessor() captain.RuntimeProcessor
	ShouldRun(key string) bool
}

// Run a particular job
func Run(key string, task CronJob) {
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
