package main

import (
	"flag"
	"fmt"
)

func main() {
	job := flag.String("job", "", "job to run")
	flag.Parse()
	fmt.Println(*job)
	if *job == "" {
		flag.Usage()
	}
	sampleWorker(*job)
}

func sampleWorker(key string) {
	fmt.Println("listen for key and act accordingly")
}
