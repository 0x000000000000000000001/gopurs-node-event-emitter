package Node_EventEmitter

import (
	"gopurs/output/gopurs_runtime"
)

type listener struct {
	cb   gopurs_runtime.Value
	once bool
}

type EventEmitter struct {
	listeners    map[string][]listener
	maxListeners int64
}

func NewImpl(_ interface{}) interface{} {
	return &EventEmitter{
		listeners:    make(map[string][]listener),
		maxListeners: 10,
	}
}

func EventNamesImpl(emitter gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
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
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	return e.maxListeners
}

func ListenerCountImpl(emitter gopurs_runtime.Value, eventName string, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	return int64(len(e.listeners[eventName]))
}

func SetMaxListenersImpl(emitter gopurs_runtime.Value, max int64, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	e.maxListeners = max
	return nil
}

func GopursUnsafeEmitFn1(emitter gopurs_runtime.Value, eventName string, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.Apply(l.cb, gopurs_runtime.Box[any](nil))
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn2(emitter gopurs_runtime.Value, eventName string, arg1 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.Apply(l.cb, arg1)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn3(emitter gopurs_runtime.Value, eventName string, arg1, arg2 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp2(l.cb, arg1, arg2)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeOn(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	e.listeners[eventName] = append(e.listeners[eventName], listener{cb: cb, once: false})
	return nil
}

func GopursUnsafeOff(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
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
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	e.listeners[eventName] = append(e.listeners[eventName], listener{cb: cb, once: true})
	return nil
}

func GopursUnsafePrependListener(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	e.listeners[eventName] = append([]listener{{cb: cb, once: false}}, e.listeners[eventName]...)
	return nil
}

func GopursUnsafePrependOnceListener(emitter gopurs_runtime.Value, eventName string, cb gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	e.listeners[eventName] = append([]listener{{cb: cb, once: true}}, e.listeners[eventName]...)
	return nil
}

func GopursUnsafeEmitFn4(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp3(l.cb, arg1, arg2, arg3)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn5(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp4(l.cb, arg1, arg2, arg3, arg4)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn6(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp5(l.cb, arg1, arg2, arg3, arg4, arg5)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn7(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5, arg6 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp6(l.cb, arg1, arg2, arg3, arg4, arg5, arg6)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn8(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5, arg6, arg7 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp7(l.cb, arg1, arg2, arg3, arg4, arg5, arg6, arg7)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}

func GopursUnsafeEmitFn9(emitter gopurs_runtime.Value, eventName string, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 gopurs_runtime.Value, _ interface{}) interface{} {
	e := gopurs_runtime.Unbox[*EventEmitter](emitter)
	list := e.listeners[eventName]
	if len(list) == 0 {
		return false
	}
	var next []listener
	for _, l := range list {
		gopurs_runtime.UncurriedApp8(l.cb, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8)
		if !l.once {
			next = append(next, l)
		}
	}
	e.listeners[eventName] = next
	return true
}
