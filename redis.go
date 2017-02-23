package main

import (
	"github.com/garyburd/redigo/redis"
)

var RedisAddr = ":6379"

func NewConn() (redis.Conn, error) {
	return redis.Dial("tcp", RedisAddr)
}
