package main

import (
	"github.com/garyburd/redigo/redis"
)

var (
	RedisAddr       = ":6379"
	keyspacePrefix  = "__keyspace@0__:"
	workerKeyPrefix = "test/"
)

func NewConn() (redis.Conn, error) {
	return redis.Dial("tcp", RedisAddr)
}
