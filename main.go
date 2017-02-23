package main

func main() {
	master := NewMaster("test/*")
	go master.WatchWorkers()
	worker := NewWorker("worker", "localhost:1210")
	worker.HeartBeat()
}
