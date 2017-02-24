package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type Master struct {
	members         map[string]*Member
	conn            redis.Conn
	keyspacePattern string
}

type Member struct {
	Name string
	Addr string
}

func NewMaster(kp string) *Master {
	c, err := NewConn()
	if err != nil {
		log.Fatal(err)
	}

	master := &Master{
		members:         make(map[string]*Member),
		conn:            c,
		keyspacePattern: kp,
	}
	return master
}

func (m *Master) AddMember(channel string) {
	key := strings.TrimPrefix(channel, keyspacePrefix)

	// create new conn for m.conn is subscribing
	c, err := NewConn()
	if err != nil {
		log.Println(err)
	}
	r, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		log.Println(err)
	}

	name := r["name"]
	addr := r["addr"]
	m.members[name] = &Member{Name: name, Addr: addr}

	log.Println(m.members)
}

func (m *Master) RmMember(channel string) {
	key := strings.TrimPrefix(channel, keyspacePrefix)
	name := strings.TrimPrefix(key, workerKeyPrefix)
	delete(m.members, name)

	log.Println(m.members)
}

func (m *Master) WatchWorkers() {
	channel := keyspacePrefix + m.keyspacePattern
	psc := redis.PubSubConn{Conn: m.conn}
	psc.PSubscribe(channel)

	for {
		switch n := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("Message: %s %s\n", n.Channel, n.Data)
		case redis.PMessage:
			fmt.Printf("PMessage: %s %s %s\n", n.Pattern, n.Channel, n.Data)
			switch string(n.Data) {
			case "hset":
				m.AddMember(n.Channel)
			case "expired":
				m.RmMember(n.Channel)
			case "del":
				m.RmMember(n.Channel)
			}
		case redis.Subscription:
			fmt.Printf("Subscription: %s %s %d\n", n.Kind, n.Channel, n.Count)
			if n.Count == 0 {
				return
			}
		case error:
			fmt.Printf("error: %v\n", n)
			return
		}
	}
}
