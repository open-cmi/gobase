package events

import (
	"github.com/open-cmi/gobase/initial"
	"github.com/open-cmi/gobase/pkg/evchan"
)

var echan *evchan.EventChan

func Register(event string, handler func(ev string, data interface{})) error {
	return echan.RegisterEvent(event, handler)
}

func Notify(event string, data interface{}) {
	echan.NotifyEvent(event, data)
}

func Init() error {
	echan.Run()
	return nil
}

func init() {
	echan = evchan.NewEventChan()
	initial.Register("chan-event", initial.PhaseDefault, Init)
}
