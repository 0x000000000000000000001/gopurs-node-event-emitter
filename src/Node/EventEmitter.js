import EventEmitter from "node:events";

const newImpl = function () {
  return new EventEmitter();
}
export { newImpl };

// addEventListener - not implemented; alias to `on`
export const gopursUnsafeEmitFn1 = (emitter) => (eventName) => () => emitter.emit(eventName);
export const gopursUnsafeEmitFn2 = (emitter) => (eventName) => (arg1) => () => emitter.emit(eventName, arg1);
export const gopursUnsafeEmitFn3 = (emitter) => (eventName) => (arg1) => (arg2) => () => emitter.emit(eventName, arg1, arg2);
export const gopursUnsafeEmitFn4 = (emitter) => (eventName) => (arg1) => (arg2) => (arg3) => () => emitter.emit(eventName, arg1, arg2, arg3);
export const gopursUnsafeEmitFn5 = (emitter) => (eventName) => (arg1) => (arg2) => (arg3) => (arg4) => () => emitter.emit(eventName, arg1, arg2, arg3, arg4);
export const gopursUnsafeEmitFn6 = (emitter) => (eventName) => (arg1) => (arg2) => (arg3) => (arg4) => (arg5) => () => emitter.emit(eventName, arg1, arg2, arg3, arg4, arg5);
export const gopursUnsafeEmitFn7 = (emitter) => (eventName) => (arg1) => (arg2) => (arg3) => (arg4) => (arg5) => (arg6) => () => emitter.emit(eventName, arg1, arg2, arg3, arg4, arg5, arg6);
export const gopursUnsafeEmitFn8 = (emitter) => (eventName) => (arg1) => (arg2) => (arg3) => (arg4) => (arg5) => (arg6) => (arg7) => () => emitter.emit(eventName, arg1, arg2, arg3, arg4, arg5, arg6, arg7);
export const gopursUnsafeEmitFn9 = (emitter) => (eventName) => (arg1) => (arg2) => (arg3) => (arg4) => (arg5) => (arg6) => (arg7) => (arg8) => () => emitter.emit(eventName, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8);
export const eventNamesImpl = (emitter) => emitter.eventNames();
export const symbolOrStr = (left, right, sym) => typeof sym == "symbol" ? left(sym) : right(sym);
export const getMaxListenersImpl = (emitter) => emitter.getMaxListeners();
export const listenerCountImpl = (emitter, eventName) => emitter.listenerCount(eventName);
// listeners - not implemented; returned functions cannot be used in type-safe way.
export const gopursUnsafeOff = (emitter) => (eventName) => (cb) => () => emitter.off(eventName, cb);
export const gopursUnsafeOn = (emitter) => (eventName) => (cb) => () => emitter.on(eventName, cb);
export const gopursUnsafeOnce = (emitter) => (eventName) => (cb) => () => emitter.once(eventName, cb);
export const gopursUnsafePrependListener = (emitter) => (eventName) => (cb) => () => emitter.prependListener(eventName, cb);
export const gopursUnsafePrependOnceListener = (emitter) => (eventName) => (cb) => () => emitter.prependOnceListener(eventName, cb);
// removeAllListeners - not implemented; bad practice
// removeEventListener - not implemented; alias to `off`
export const setMaxListenersImpl = (emitter, max) => emitter.setMaxListeners(max);
// rawListeners - not implemented; returned functions cannot be used in type-safe way.

