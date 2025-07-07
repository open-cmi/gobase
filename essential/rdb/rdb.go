package rdb

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/essential/logger"
)

// Config cache config
type Config struct {
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Password string `json:"password,omitempty"`
}

type Client struct {
	Index int
	Mutex sync.Mutex
	Conn  *redis.Client
}

var gConf Config

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

func GetConf() *Config {
	return &gConf
}

func Init() error {
	addr := net.JoinHostPort(gConf.Host, strconv.Itoa(gConf.Port))
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: gConf.Password,
		DB:       0,
	})
	defer c.Close()
	pong, err := c.Ping(context.TODO()).Result()
	if err != nil {
		return err
	}
	logger.Debugf("redis ping pong: %s\n", pong)
	return nil
}

// Parse db init
func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}

	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	gConf.Host = "127.0.0.1"
	gConf.Port = 25431
	config.RegisterConfig("rdb", Parse, Save)
}
