package Node_EventEmitter

import (
	"gopurs/output/gopurs_runtime"
	"reflect"
	"testing"
	"time"
)

type listenerAddedProbe struct {
	notify func(string)
}

func (p *listenerAddedProbe) OnListenerAdded(eventName string) {
	p.notify(eventName)
}

type listenerAdder func(gopurs_runtime.Value, string, gopurs_runtime.Value, interface{}) interface{}

func TestListenerAddedAfterUnlockPreservesListenerBehavior(t *testing.T) {
	for _, tc := range []struct {
		name    string
		add     listenerAdder
		once    bool
		prepend bool
	}{
		{"on", GopursUnsafeOn, false, false},
		{"once", GopursUnsafeOnce, true, false},
		{"prepend", GopursUnsafePrependListener, false, true},
		{"prependOnce", GopursUnsafePrependOnceListener, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			e := NewImpl(nil).(*EventEmitter)
			boxed := gopurs_runtime.Box(e)
			var calls []string
			callback := func(name string) gopurs_runtime.Value {
				return gopurs_runtime.Func(func(_ gopurs_runtime.Value) gopurs_runtime.Value {
					calls = append(calls, name)
					return gopurs_runtime.Value{}
				})
			}
			GopursUnsafeOn(boxed, "data", callback("initial"), nil)
			notifications := 0
			e.Any = &listenerAddedProbe{notify: func(eventName string) {
				notifications++
				if eventName != "data" {
					t.Errorf("notification event = %q, want data", eventName)
				}
				// Both readers and writers must be able to reenter the emitter.
				if count := ListenerCountImpl(boxed, eventName, nil); count != int64(2) {
					t.Errorf("new listener not visible to notification: count = %v", count)
				}
				SetMaxListenersImpl(boxed, 17, nil)
				GopursUnsafeEmitFn1(boxed, eventName, nil)
			}}
			done := make(chan struct{})
			go func() {
				defer close(done)
				tc.add(boxed, "data", callback("new"), nil)
			}()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("listener notification deadlocked on emitter reentry")
			}
			if notifications != 1 {
				t.Fatalf("notifications = %d, want 1", notifications)
			}
			first := []string{"initial", "new"}
			if tc.prepend {
				first = []string{"new", "initial"}
			}
			if !reflect.DeepEqual(calls, first) {
				t.Fatalf("first emit = %v, want %v", calls, first)
			}
			GopursUnsafeEmitFn1(boxed, "data", nil)
			want := append(append([]string{}, first...), first...)
			if tc.once {
				want = append(append([]string{}, first...), "initial")
			}
			if !reflect.DeepEqual(calls, want) {
				t.Errorf("two emits = %v, want %v", calls, want)
			}
			if got := GetMaxListenersImpl(boxed, nil); got != int64(17) {
				t.Errorf("reentrant mutation = %v, want 17", got)
			}
		})
	}
}

func TestListenerAddedOptionalForExistingEmitters(t *testing.T) {
	for _, payload := range []any{nil, struct{}{}} {
		e := NewImpl(nil).(*EventEmitter)
		e.Any = payload
		boxed := gopurs_runtime.Box(e)
		calls := 0
		cb := gopurs_runtime.Func(func(_ gopurs_runtime.Value) gopurs_runtime.Value {
			calls++
			return gopurs_runtime.Value{}
		})
		GopursUnsafeOn(boxed, "data", cb, nil)
		GopursUnsafeOnce(boxed, "data", cb, nil)
		GopursUnsafePrependListener(boxed, "data", cb, nil)
		GopursUnsafePrependOnceListener(boxed, "data", cb, nil)
		GopursUnsafeEmitFn1(boxed, "data", nil)
		GopursUnsafeEmitFn1(boxed, "data", nil)
		if calls != 6 {
			t.Errorf("payload %T: callbacks = %d, want 6", payload, calls)
		}
	}
}
