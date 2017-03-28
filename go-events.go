/*
 * https://co-sche.mit-license.org/
 */

package events

import (
	"reflect"
)

type Listener func(...interface{})

type EventEmitter struct {
	events map[string][]Listener
}

func NewEventEmitter() *EventEmitter {
	ee := new(EventEmitter)
	ee.events = make(map[string][]Listener)
	return ee
}

func (ee *EventEmitter) Emit(eventName string, args ...interface{}) bool {
	for _, listener := range ee.events[eventName] {
		listener(args...)
	}
	return true
}

func (ee *EventEmitter) On(eventName string, listener Listener) *EventEmitter {
	ee.events[eventName] = append(ee.events[eventName], listener)
	ee.Emit("newListener", eventName, listener)
	return ee
}

func (ee *EventEmitter) AddListener(eventName string, listener Listener) *EventEmitter {
	return ee.On(eventName, listener)
}

func (ee *EventEmitter) PrependListener(eventName string, listener Listener) *EventEmitter {
	ee.events[eventName] = append([]Listener{listener}, ee.events[eventName][0:]...)
	ee.Emit("newListener", eventName, listener)
	return ee
}

func (ee *EventEmitter) RemoveListener(eventName string, listener Listener) *EventEmitter {
	ee.events[eventName] = omit(ee.events[eventName], listener)
	if ee.ListenerCount(eventName) == 0 {
		delete(ee.events, eventName)
	}
	ee.Emit("removeListener", eventName, listener)
	return ee
}

func (ee *EventEmitter) RemoveAllListeners(eventName ...string) *EventEmitter {
	if len(eventName) == 0 {
		ee.events = make(map[string][]Listener)
	} else {
		delete(ee.events, eventName[0])
	}
	return ee
}

func (ee *EventEmitter) Listeners(eventName string) []Listener {
	return ee.events[eventName]
}

func (ee *EventEmitter) ListenerCount(eventName string) int {
	return len(ee.Listeners(eventName))
}

func (ee *EventEmitter) EventNames() []string {
	eventNames := make([]string, len(ee.events))
	i := 0
	for k, _ := range ee.events {
		eventNames[i] = k
		i++
	}
	return eventNames
}

/*
func (ee *EventEmitter) Once(eventName string, listener Listener) *EventEmitter {
	return ee
}

func (ee *EventEmitter) PrependOnceListener(eventName string, listener Listener) *EventEmitter {
	return ee
}
*/

/*
func (ee *EventEmitter) GetMaxListeners() uint {
	return 0
}

func (ee *EventEmitter) SetMaxListeners(n uint) *EventEmitter {
	return ee
}
*/

func omit(fs []Listener, f Listener) []Listener {
	for i := len(fs) - 1; i >= 0; i-- {
		if eq(f, fs[i]) {
			fs = del(fs, i)
		}
	}
	return fs
}

func del(fs[]Listener, i int) []Listener {
	j := len(fs) - 1
	copy(fs[i:], fs[i+1:])
	fs[j] = nil
	fs = fs[:j]
	return fs
}

func eq(a Listener, b Listener) bool {
	return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
}
