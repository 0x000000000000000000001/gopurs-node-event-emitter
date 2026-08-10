package Node_EventEmitter

import (
	"gopurs/output/gopurs_runtime"
	"sync"
)

type listener struct {
	cb   gopurs_runtime.Value
	once bool
}

type EventEmitter struct {
	listeners    map[string][]listener
	mu           sync.RWMutex
	maxListeners int64
	Any          any
}

func (e *EventEmitter) GetEventEmitter() *EventEmitter {
	return e
}

func NewImpl(_ interface{}) interface{} {
	return &EventEmitter{
		listeners:    make(map[string][]listener),
		maxListeners: 10,
	}
}

func EventNamesImpl(emitter gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.RLock()
	defer e.mu.RUnlock()
	names := make([]gopurs_runtime.Value, 0, len(e.listeners))
	for k := range e.listeners {
		names = append(names, gopurs_runtime.Str(k))
	}
	return names
}

func SymbolOrStr(left, right, sym gopurs_runtime.Value) interface{} {
	return gopurs_runtime.Apply(right, sym)
}

func GetMaxListenersImpl(emitter gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.maxListeners
}

func ListenerCountImpl(emitter gopurs_runtime.Value, eventName string, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.RLock()
	defer e.mu.RUnlock()
	return int64(len(e.listeners[eventName]))
}

func SetMaxListenersImpl(emitter gopurs_runtime.Value, max int64, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.maxListeners = max
	return nil
}

func GopursUnsafeEmitFn1(emitter gopurs_runtime.Value, eventName string, arg1 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.Apply(l.cb, gopurs_runtime.Value{})
	}
	return true
}

func GopursUnsafeEmitFn2(emitter gopurs_runtime.Value, eventName string, arg1 gopurs_runtime.Value, arg2 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.Apply(l.cb, arg1)
	}
	return true
}

func GopursUnsafeEmitFn3(emitter gopurs_runtime.Value, eventName string, arg1, arg2 gopurs_runtime.Value, arg3 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp2(l.cb, arg1, arg2)
	}
	return true
}

func GopursUnsafeOn(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[eventName] = append(e.listeners[eventName], listener{cb: cb, once: false})
	return nil
}

func GopursUnsafeOff(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	defer e.mu.Unlock()
	list := e.listeners[eventName]
	for i, l := range list {
		if l.cb.UnsafePtr == cb.UnsafePtr {
			e.listeners[eventName] = append(list[:i], list[i+1:]...)
			break
		}
	}
	return nil
}

func GopursUnsafeOnce(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[eventName] = append(e.listeners[eventName], listener{cb: cb, once: true})
	return nil
}

func GopursUnsafePrependListener(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[eventName] = append([]listener{{cb: cb, once: false}}, e.listeners[eventName]...)
	return nil
}

func GopursUnsafePrependOnceListener(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners[eventName] = append([]listener{{cb: cb, once: true}}, e.listeners[eventName]...)
	return nil
}

func GopursUnsafeEmitFn4(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3 gopurs_runtime.Value, arg4 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp3(l.cb, arg1, arg2, arg3)
	}
	return true
}

func GopursUnsafeEmitFn5(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4 gopurs_runtime.Value, arg5 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp4(l.cb, arg1, arg2, arg3, arg4)
	}
	return true
}

func GopursUnsafeEmitFn6(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5 gopurs_runtime.Value, arg6 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp5(l.cb, arg1, arg2, arg3, arg4, arg5)
	}
	return true
}

func GopursUnsafeEmitFn7(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5, arg6 gopurs_runtime.Value, arg7 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp6(l.cb, arg1, arg2, arg3, arg4, arg5, arg6)
	}
	return true
}

func GopursUnsafeEmitFn8(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5, arg6, arg7 gopurs_runtime.Value, arg8 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp7(l.cb, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
	}
	return true
}

func GopursUnsafeEmitFn9(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 gopurs_runtime.Value, arg9 interface{}) interface{} {
	e := ExtractEventEmitter(emitter)
	e.mu.Lock()
	list := e.listeners[eventName]
	if len(list) == 0 {
		e.mu.Unlock()
		return false
	}
	var next []listener
	for _, l := range list {
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	e.mu.Unlock()

	for _, l := range list {
		gopurs_runtime.UncurriedApp8(l.cb, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
	}
	return true
}

type EventEmitterLike interface {
	GetEventEmitter() *EventEmitter
}

func ExtractEventEmitter(emitter gopurs_runtime.Value) *EventEmitter {
	var val any = emitter
	for {
		if v, ok := val.(gopurs_runtime.Value); ok {
			if v.Type == gopurs_runtime.TypeAny {
				val = v.PtrVal()
			} else {
				break
			}
		} else {
			break
		}
	}
	if e, ok := val.(*EventEmitter); ok {
		return e
	}
	if e, ok := val.(EventEmitterLike); ok {
		return e.GetEventEmitter()
	}
	panic("Cannot extract EventEmitter from value")
}
