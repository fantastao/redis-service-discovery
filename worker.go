package main

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Worker struct {
	Name    string
	Address string
	conn    redis.Conn
}

func NewWorker(name string, addr string) *Worker {
	return &Worker{Name: name, Address: addr}
}

func (w *Worker) Connect() {
	c, err := NewConn()
	if err != nil {
		log.Fatal(err)
	}
	w.conn = c
}

func (w *Worker) Close() {
	w.conn.Close()
}

func (w *Worker) HeartBeat() {
	w.Connect()
	defer w.Close()

	key := workerKeyPrefix + w.Name
	_, err := w.conn.Do("HMSET", key, "addr", w.Address, "name", w.Name)
	if err != nil {
		log.Println(err)
	}

	for {
		_, err := w.conn.Do("Expire", key, 5)
		if err != nil {
			break
		}
		time.Sleep(time.Second * 3)
	}

}
