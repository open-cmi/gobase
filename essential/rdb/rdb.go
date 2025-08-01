package rdb

import (
	"context"
	"net"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/gobase/essential/logger"
)

type Client struct {
	Index int
	Mutex sync.Mutex
	Conn  *redis.Client
}

var gClientPoolMutex sync.Mutex
var gClientPool map[int]*Client = make(map[int]*Client)

// GetClient get client
func GetClient(dbIndex int) *Client {
	if dbIndex < 0 || dbIndex > 15 {
		logger.Errorf("dbIndex should be between 0 and 15\n")
		return nil
	}
	c, ok := gClientPool[dbIndex]
	if ok {
		return c
	}
	addr := net.JoinHostPort(gConf.Host, strconv.Itoa(gConf.Port))
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: gConf.Password,
		DB:       dbIndex,
	})
	if cli == nil {
		return nil
	}
	gClientPoolMutex.Lock()
	defer gClientPoolMutex.Unlock()
	gClientPool[dbIndex] = &Client{
		Conn: cli,
	}
	return gClientPool[dbIndex]
}

func (c *Client) Reconnect() error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Conn.Close()

	addr := net.JoinHostPort(gConf.Host, strconv.Itoa(gConf.Port))
	cli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: gConf.Password,
		DB:       c.Index,
	})
	_, err := cli.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	c.Conn = cli
	return nil
}
