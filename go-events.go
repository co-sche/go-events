/*
 * https://co-sche.mit-license.org/
 */

package events

type Event interface{}

type EventListener struct {
	callback func(...Event)
}

type EventEmitter struct {
	events map[string][]*EventListener
}

func NewEventListener(callback func(...Event)) *EventListener {
	return &EventListener{callback}
}

func NewEventEmitter() *EventEmitter {
	ee := new(EventEmitter)
	ee.events = make(map[string][]*EventListener)
	return ee
}

func (ee *EventEmitter) Emit(eventName string, args ...Event) bool {
	for _, listener := range ee.events[eventName] {
		listener.callback(args...)
	}
	return true
}

func (ee *EventEmitter) On(eventName string, listener *EventListener) *EventEmitter {
	ee.events[eventName] = append(ee.events[eventName], listener)
	ee.Emit("newListener", eventName, listener)
	return ee
}

func (ee *EventEmitter) AddListener(eventName string, listener *EventListener) *EventEmitter {
	return ee.On(eventName, listener)
}

func (ee *EventEmitter) PrependListener(eventName string, listener *EventListener) *EventEmitter {
	ee.events[eventName] = append([]*EventListener{listener}, ee.events[eventName][0:]...)
	ee.Emit("newListener", eventName, listener)
	return ee
}

func (ee *EventEmitter) RemoveListener(eventName string, listener *EventListener) *EventEmitter {
	ee.events[eventName] = omit(ee.events[eventName], listener)
	if ee.ListenerCount(eventName) == 0 {
		delete(ee.events, eventName)
	}
	ee.Emit("removeListener", eventName, listener)
	return ee
}

func (ee *EventEmitter) RemoveAllListeners(eventName ...string) *EventEmitter {
	if len(eventName) == 0 {
		ee.events = make(map[string][]*EventListener)
	} else {
		delete(ee.events, eventName[0])
	}
	return ee
}

func (ee *EventEmitter) Listeners(eventName string) []*EventListener {
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

func omit(fs []*EventListener, f *EventListener) []*EventListener {
	for i := len(fs) - 1; i >= 0; i-- {
		if f == fs[i] {
			fs = del(fs, i)
		}
	}
	return fs
}

func del(fs[]*EventListener, i int) []*EventListener {
	j := len(fs) - 1
	copy(fs[i:], fs[i+1:])
	fs[j] = nil
	fs = fs[:j]
	return fs
}
