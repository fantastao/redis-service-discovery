package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup

	master := NewMaster("test/*")
	wg.Add(1)
	go func() {
		defer wg.Done()
		master.WatchWorkers()
	}()

	worker := NewWorker("worker", "localhost:1210")
	wg.Add(1)
	go func() {
		defer wg.Done()
		worker.HeartBeat()
	}()

	wg.Wait()
}
