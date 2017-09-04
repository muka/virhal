package emitter

import (
	"github.com/chuckpreslar/emission"
)

var instance *emission.Emitter

//Event wrap a message
type Event interface {
	GetRaw() interface{}
}

//GetEmitter return the event emitter instance
func GetEmitter() *emission.Emitter {
	if instance == nil {
		instance = emission.NewEmitter()
	}
	return instance
}

//On register for events
func On(event string, fn func(ev Event)) {
	instance.On(event, fn)
}

//Off unregister for events
func Off(event string, fn func(ev Event)) {
	instance.Off(event, fn)
}
