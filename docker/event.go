package docker

import (
	"github.com/chuckpreslar/emission"
	"github.com/docker/docker/api/types/events"
)

var emitter *emission.Emitter

//Event wrap a message
type Event struct {
	Message events.Message
}

//GetEmitter return the event emitter instance
func GetEmitter() *emission.Emitter {
	if emitter == nil {
		emitter = emission.NewEmitter()
	}
	return emitter
}

//On register for events
func On(event string, fn func(ev Event)) {
	emitter.On(event, fn)
}

//Off unregister for events
func Off(event string, fn func(ev Event)) {
	emitter.Off(event, fn)
}
