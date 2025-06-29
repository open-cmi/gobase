package pubsub

import (
	evbus "github.com/asaskevich/EventBus"
)

var bus evbus.Bus

func Publish(topic string, args ...interface{}) {
	bus.Publish(topic, args...)
}

func RegisterSubscribe(topic string, fn interface{}) error {
	bus.Subscribe(topic, fn)
	return nil
}

func init() {
	bus = evbus.New()
}
