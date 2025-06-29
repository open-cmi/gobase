package ticker

import (
	"errors"
	"fmt"

	"github.com/open-cmi/gobase/essential/logger"
	"github.com/open-cmi/gobase/initial"
	"github.com/robfig/cron/v3"
)

var initialized bool

type Ticker struct {
	Name string
	Spec string
	Func func(name string, data interface{})
	Data interface{}
}

var tickers map[string]Ticker = make(map[string]Ticker)
var cronMap map[string]*cron.Cron = make(map[string]*cron.Cron)

func Register(name string, spec string, f func(string, interface{}), data interface{}) error {
	_, found := tickers[name]
	if found {
		errMsg := fmt.Sprintf("ticker %s registered failed", name)
		return errors.New(errMsg)
	}
	tickers[name] = Ticker{
		Name: name,
		Spec: spec,
		Func: f,
		Data: data,
	}
	if initialized {
		// 如果已经初始化了，此时需要立即创建
		ins := cron.New(cron.WithSeconds())
		ins.AddFunc(spec, func() {
			logger.Debugf("start to run timer %s\n", name)
			f(name, data)
		})
		go ins.Run()
		cronMap[name] = ins
	}
	return nil
}

func Init() error {
	for i := range tickers {
		t := tickers[i]
		ins := cron.New(cron.WithSeconds())
		_, err := ins.AddFunc(t.Spec, func() {
			logger.Debugf("start to run timer %s\n", t.Name)
			t.Func(t.Name, t.Data)
		})
		if err != nil {
			return err
		}
		go ins.Start()
		cronMap[t.Name] = ins
	}

	initialized = true
	return nil
}

func Close() {
	for _, c := range cronMap {
		c.Stop()
	}
}

func Remove(name string) error {
	ins, ok := cronMap[name]
	if !ok {
		return errors.New("cron task is not existing")
	}
	ins.Stop()
	delete(cronMap, name)
	delete(tickers, name)
	return nil
}

func init() {
	initial.Register("ticker", initial.PhaseFinal, Init)
}
