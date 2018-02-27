package main

import (
	"flag"
	"fmt"

	"crazy_nl_backend/cmd/cron/workers"
)

func main() {
	job := flag.String("job", "", "job to run")
	flag.Parse()
	fmt.Println(*job)
	if *job == "" {
		flag.Usage()
	}
	importDevices(*job)
}

func importDevices(key string) {
	job := workers.ImportDevices{}
	workers.Run(key, job)
}
