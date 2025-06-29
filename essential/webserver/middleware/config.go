package middleware

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/open-cmi/gobase/essential/config"
	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/essential/rdb"
	"github.com/open-cmi/gobase/initial"
	"github.com/open-cmi/gobase/pkg/memstore"
	"github.com/open-cmi/gobase/pkg/redistore"
)

type Config struct {
	Store  string `json:"store"`
	MaxAge int    `json:"max_age"`
}

var gConf Config

func Init() error {
	var err error

	switch gConf.Store {
	case "memory":
		memoryStore = memstore.NewMemStore([]byte("memorystore"),
			[]byte("enckey12341234567890123456789012"))
		memoryStore.MaxAge(gConf.MaxAge)
	case "redis":
		rdbConf := rdb.GetConf()
		host := fmt.Sprintf("%s:%d", rdbConf.Host, rdbConf.Port)
		redisStore, err = redistore.NewRediStoreWithDB(100, "tcp", host, rdbConf.Password, "2")
		if err != nil {
			logger.Error("redis store new failed\n")
			return err
		}

		redisStore.SetKeyPrefix("koa-sess-")
		redisStore.SetSerializer(redistore.JSONSerializer{})
		redisStore.SetMaxAge(gConf.MaxAge)
	default:
		return errors.New("middleware store type not supported")
	}

	return nil
}

func Parse(raw json.RawMessage) error {
	err := json.Unmarshal(raw, &gConf)
	if err != nil {
		return err
	}
	if gConf.MaxAge == 0 {
		gConf.MaxAge = 3600
	}
	return nil
}

func Save() json.RawMessage {
	raw, _ := json.Marshal(&gConf)
	return raw
}

func init() {
	// default config
	gConf.Store = "memory"

	config.RegisterConfig("middleware", Parse, Save)
	initial.Register("middleware", initial.PhaseDefault, Init)
}
