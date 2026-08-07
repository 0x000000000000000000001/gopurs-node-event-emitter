-- | ## Handling events emitted by an `EventEmitter`
-- |
-- | One can add callbacks to an `EventEmitter` on two major axes:
-- | - whether listener is added to the **end** (i.e. `on`) or **start** (i.e. `prependListener`) of the array
-- | - whether a listener is automatically removed after the first event (i.e. `once` or `prependOnceListener`).
-- |
-- | This module provides functions for each of the above 4 callback-adding functions
-- | If `<fn>` is either `on`, `once`, `prependListener`, or `prependOnceListener`, then this module exposes
-- | 1. `<fn>` - returns a callback that removes the listener
-- | 2. `<fn>_` - does not return a callback that can remove the listener
-- |
-- | ## Defining events emitted by an `EventEmitter`
-- |
-- | Below, we'll provide an example for how to define an event handler for a type. Let's assume the following:
-- | - There is a type `Foo` that exends `EventEmitter`
-- | - `Foo` values can handle "bar" events
-- | - a "bar" event takes the following callback: `EffectFn2 (Nullable Error) String Unit`
-- | - the `String` value is always either "red", "green", or "blue"
-- |
-- | Then we would write
-- | ```
-- | data Color 
-- |   = Red 
-- |   | Green 
-- |   | Blue
-- |
-- | -- Note: see docs on `EventHandle` 
-- | -- for the below naming convention justification 
-- | -- of suffixing an event name with `H`.
-- | barH 
-- |   :: EventHandle 
-- |        Foo 
-- |        (Maybe Error -> Color -> Effect Unit) 
-- |        (EffectFn1 (Nullable Error) String Unit)
-- | barH = EventHandle "bar" $ \psCb -> 
-- |   mkEffectFn2 \nullableError str ->
-- |     psCb (toMaybe nullableError) case str of
-- |       "red" -> Red
-- |       "green" -> Green
-- |       "blue" -> Blue
-- |       _ -> 
-- |         unsafeCrashWith $ 
-- |           "Impossible String value for event 'bar': " <> show str
-- | ```
-- |
-- | ## Emitting events via an `EventEmitter`
-- |
-- | Unfortunately, there isn't a good way to emit events safely in PureScript. If one wants to emit an event
-- | in PureScript code that will be consumed by PureScript code, there are better abstractions to use than `EventEmitter`.
-- | If one wants to emit an event in PureScript code that will be consumed by JavaScript code, then
-- | the `unsafeEmitFn` function can be used to call n-ary functions. However, this is very unsafe. See its docs for more context.
module Node.EventEmitter
  ( EventEmitter
  , new
  , SymbolOrStr
  , eventNames
  , getMaxListeners
  , listenerCount
  , setMaxListeners
  , setUnlimitedListeners
  , unsafeEmitFn1
  , unsafeEmitFn2
  , unsafeEmitFn3
  , EventHandle(..)
  , newListenerH
  , removeListenerH
  , on
  , on_
  , once
  , once_
  , prependListener
  , prependListener_
  , prependOnceListener
  , prependOnceListener_
  ) where

import Prelude

import Data.Either (Either(..))
import Data.Function.Uncurried (Fn3, runFn3)
import Effect (Effect)
import Effect.Uncurried (EffectFn1, EffectFn2, EffectFn3, EffectFn4, EffectFn5, EffectFn6, EffectFn7, EffectFn8, EffectFn9, EffectFn10, mkEffectFn1, mkEffectFn2, mkEffectFn3, mkEffectFn4, mkEffectFn5, mkEffectFn6, mkEffectFn7, mkEffectFn8, mkEffectFn9, mkEffectFn10, runEffectFn1, runEffectFn2, runEffectFn3, runEffectFn4)
import Node.Symbol (JsSymbol)
import Unsafe.Coerce (unsafeCoerce)

foreign import data EventEmitter :: Type

-- | Create a new event emitter
foreign import newImpl :: Effect EventEmitter

new :: Effect EventEmitter
new = newImpl

foreign import data SymbolOrStr :: Type

foreign import eventNamesImpl :: EventEmitter -> Array SymbolOrStr

eventNames :: EventEmitter -> Array (Either JsSymbol String)
eventNames ee = map (\x -> runFn3 symbolOrStr Left Right x) $ eventNamesImpl ee

foreign import symbolOrStr
  :: Fn3
       (forall a. JsSymbol -> Either JsSymbol a)
       (forall b. String -> Either b String)
       SymbolOrStr
       (Either JsSymbol String)

foreign import getMaxListenersImpl :: EffectFn1 EventEmitter Int

-- | By default, an event emitter can only have a maximum of 10 listeners
-- | for a given event.
getMaxListeners :: EventEmitter -> Effect Int
getMaxListeners = runEffectFn1 getMaxListenersImpl

foreign import listenerCountImpl :: EffectFn2 EventEmitter String Int

listenerCount :: EventEmitter -> String -> Effect Int
listenerCount emitter eventName = runEffectFn2 listenerCountImpl emitter eventName

foreign import setMaxListenersImpl :: EffectFn2 EventEmitter Int Unit

setMaxListeners :: Int -> EventEmitter -> Effect Unit
setMaxListeners max emitter = runEffectFn2 setMaxListenersImpl emitter max

setUnlimitedListeners :: EventEmitter -> Effect Unit
setUnlimitedListeners = setMaxListeners 0

-- | THIS IS UNSAFE! REALLY UNSAFE!
-- | Gets the `emit` function for a particular `EventEmitter`, so that one can call n-ary functions.
-- |
-- | Given `http2session.goaway([code[, lastStreamID[, opaqueData]]])` as an example...
-- | - https://nodejs.org/dist/latest-v18.x/docs/api/http2.html#event-goaway
-- | - https://nodejs.org/dist/latest-v18.x/docs/api/http2.html#http2sessiongoawaycode-laststreamid-opaquedata
-- |
-- | We can then write a single function that handles all four cases:
-- | ```
-- | goAway
-- |   :: Http2Session
-- |   -> Maybe Code
-- |   -> Maybe LastStreamId
-- |   -> Maybe OpaqueData
-- |   -> Effect Unit
-- | goAway h2s = case _, _, _ of
-- |   Just c, Just id, Just d ->
-- |     runEffectFn4 (unsafeEmitFn h2s :: EffectFn4 String Code LastStreamId OpaqueData Unit) "goaway" c id d
-- |   Just c, Just id, Nothing ->
-- |     -- If you're feeling lucky, omit the type annotations completely
-- |     runEffectFn3 (unsafeEmitFn h2s) "goaway" c id
-- |   Just c, Nothing, Nothing ->
-- |     runEffectFn2 (unsafeEmitFn h2s :: EffectFn2 String Code LastStreamId Unit) "goaway" c
-- |   _, _, _ ->
-- |     runEffectFn1 (unsafeEmitFn h2s :: EffectFn1 String Unit) "goaway"
-- | ```
-- | 
-- | Synchronously calls each of the listeners registered for the event named `eventName`, 
-- | in the order they were registered, passing the supplied arguments to each.
-- | Returns `true` if the event had listeners, `false` otherwise.


-- | Packs all the type information we need to call `on`/`once`/`prependListener`/`prependOnceListener`
-- | with the correct callback function type.
-- |
-- | **Naming convention**: 
-- | If the name of an event is `foo`, 
-- | the corresponding PureScript `EventHandle` value should be called `fooH`.
-- | The `H` suffix is what prevent name conflicts in two situations:
-- | 1. similarly-named methods (e.g. the `"close"` event and the `close` method)
-- | 2. PureScript keywords (e.g. the `"data"` event)
-- |
-- | If an event, `foo`, can have two different kinds of callbacks, (e.g. See `Node.Stream`'s `data` event),
-- | one of two things should happen:
-- | 1. a suffix should follow the `H` to distinguish between the two (e.g. `dataHString`/`dataHBuffer`)
-- | 2. a prime character (i.e. `'`) should follow the `H` to distinguish between the two (e.g. `dataH`/`dataH'`)
data EventHandle :: Type -> Type -> Type -> Type
data EventHandle emitterType pureScriptCallback javaScriptCallback =
  EventHandle String (pureScriptCallback -> javaScriptCallback)

type role EventHandle representational representational representational

newListenerH :: EventHandle EventEmitter (Either JsSymbol String -> Effect Unit) (EffectFn1 SymbolOrStr Unit)
newListenerH = EventHandle "newListener" $ \cb -> mkEffectFn1 \jsSymbol ->
  cb $ runFn3 symbolOrStr Left Right jsSymbol

removeListenerH :: EventHandle EventEmitter (Either JsSymbol String -> Effect Unit) (EffectFn1 SymbolOrStr Unit)
removeListenerH = EventHandle "removeListener" $ \cb -> mkEffectFn1 \jsSymbol ->
  cb $ runFn3 symbolOrStr Left Right jsSymbol

-- | Adds the listener to the **end** of the `listeners` array.
-- | Returns a callback that will remove the listener from the event emitter's `listeners` array.
-- | If the listener removal callback isn't needed, use `on_`.
-- |
-- | Intended usage:
-- | ```
-- | removeLoggerCallback <- eventEmitter # on errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- |   log $ "This listener will now be removed."
-- | -- sometime later...
-- | removeLoggerCallback
-- | ```
on :: forall event emitter a. EventHandle emitter a event -> a -> emitter -> Effect (Effect Unit)
on (EventHandle eventName toJsCb) psCb eventEmitter = do
  let jsCb = toJsCb psCb
  runEffectFn3 unsafeOn (unsafeCoerce eventEmitter) eventName jsCb
  pure $ runEffectFn3 unsafeOff (unsafeCoerce eventEmitter) eventName jsCb

-- | Adds the callback to the **end** of the `listeners` array and provides no way to remove the listener in the future.
-- | If you need a callback to remove the listener in the future, use `on`.
-- | Intended usage:
-- | ```
-- | eventEmitter # on_ errorHandle  \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- | ```
on_
  :: forall emitter psCb jsCb
   . EventHandle emitter psCb jsCb
  -> psCb
  -> emitter
  -> Effect Unit
on_ (EventHandle eventName toJsCb) psCb eventEmitter =
  runEffectFn3 unsafeOn (unsafeCoerce eventEmitter) eventName $ toJsCb psCb

-- | Adds the listener to the **end** of the `listeners` array. The listener will be removed after it is invoked once.
-- | Returns a callback that will remove the listener from the event emitter's listeners array.
-- | If the listener removal callback isn't needed, use `once_`.
-- |
-- | Intended usage:
-- | ```
-- | removeLoggerCallback <- eventEmitter # once errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- |   log $ "This listener will now be removed."
-- | -- sometime later...
-- | removeLoggerCallback
-- | ```
once :: forall event emitter a. EventHandle emitter a event -> a -> emitter -> Effect (Effect Unit)
once (EventHandle eventName toJsCb) psCb eventEmitter = do
  let jsCb = toJsCb psCb
  runEffectFn3 unsafeOnce (unsafeCoerce eventEmitter) eventName jsCb
  pure $ runEffectFn3 unsafeOff (unsafeCoerce eventEmitter) eventName jsCb

-- | Adds the listener to the **end** of the `listeners` array. The listener will be removed after it is invoked once.
-- | Returns a callback that will remove the listener from the event emitter's listeners array.
-- | If you need a callback to remove the listener in the future, use `once`.
-- |
-- | Intended usage:
-- | ```
-- | eventEmitter # once_ errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- | ```
once_
  :: forall emitter psCb jsCb
   . EventHandle emitter psCb jsCb
  -> psCb
  -> emitter
  -> Effect Unit
once_ (EventHandle eventName toJsCb) psCb eventEmitter =
  runEffectFn3 unsafeOnce (unsafeCoerce eventEmitter) eventName $ toJsCb psCb

-- | Adds the listener to the **start** of the `listeners` array.
-- | Returns a callback that will remove the listener from the event emitter's listeners array.
-- | If the listener removal callback isn't needed, use `prependListener_`.
-- |
-- | Intended usage:
-- | ```
-- | removeLoggerCallback <- eventEmitter # prependListener errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- |   log $ "This listener will now be removed."
-- | -- sometime later...
-- | removeLoggerCallback
-- | ```
prependListener :: forall event emitter a. EventHandle emitter a event -> a -> emitter -> Effect (Effect Unit)
prependListener (EventHandle eventName toJsCb) psCb eventEmitter = do
  let jsCb = toJsCb psCb
  runEffectFn3 unsafePrependListener (unsafeCoerce eventEmitter) eventName jsCb
  pure $ runEffectFn3 unsafeOff (unsafeCoerce eventEmitter) eventName jsCb

-- | Adds the listener to the **start** of the `listeners` array.
-- | Returns a callback that will remove the listener from the event emitter's listeners array.
-- | If the listener removal callback isn't needed, use `prependListener`.
-- |
-- | Intended usage:
-- | ```
-- | eventEmitter # prependListener_ errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- | ```
prependListener_
  :: forall emitter psCb jsCb
   . EventHandle emitter psCb jsCb
  -> psCb
  -> emitter
  -> Effect Unit
prependListener_ (EventHandle eventName toJsCb) psCb eventEmitter =
  runEffectFn3 unsafePrependListener (unsafeCoerce eventEmitter) eventName $ toJsCb psCb

-- | Adds the listener to the **start** of the `listeners` array. The listener will be removed after it is invoked once.
-- | Returns a callback that will remove the listener from the event emitter's listeners array.
-- | If the listener removal callback isn't needed, use `prependOnceListener_`.
-- |
-- | Intended usage:
-- | ```
-- | removeLoggerCallback <- eventEmitter # prependOnceListener errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- |   log $ "This listener will now be removed."
-- | -- sometime later...
-- | removeLoggerCallback
-- | ```
prependOnceListener :: forall event emitter a. EventHandle emitter a event -> a -> emitter -> Effect (Effect Unit)
prependOnceListener (EventHandle eventName toJsCb) psCb eventEmitter = do
  let jsCb = toJsCb psCb
  runEffectFn3 unsafePrependOnceListener (unsafeCoerce eventEmitter) eventName jsCb
  pure $ runEffectFn3 unsafeOff (unsafeCoerce eventEmitter) eventName jsCb

-- | Adds the listener to the **start** of the `listeners` array. The listener will be removed after it is invoked once.
-- | Returns a callback that will remove the listener from the event emitter's listeners array.
-- | If you need a callback to remove the listener in the future, use `prependOnceListener`.
-- |
-- | Intended usage:
-- | ```
-- | eventEmitter # prependOnceListener_ errorHandle \error -> do
-- |   log $ "Got error: " <> Exception.message error
-- | ```
prependOnceListener_
  :: forall emitter psCb jsCb
   . EventHandle emitter psCb jsCb
  -> psCb
  -> emitter
  -> Effect Unit
prependOnceListener_ (EventHandle eventName toJsCb) psCb eventEmitter =
  runEffectFn3 unsafePrependOnceListener (unsafeCoerce eventEmitter) eventName $ toJsCb psCb

foreign import gopursUnsafeOn :: forall f. EventEmitter -> String -> f -> Effect Unit
foreign import gopursUnsafeOff :: forall f. EventEmitter -> String -> f -> Effect Unit
foreign import gopursUnsafeOnce :: forall f. EventEmitter -> String -> f -> Effect Unit
foreign import gopursUnsafePrependListener :: forall f. EventEmitter -> String -> f -> Effect Unit
foreign import gopursUnsafePrependOnceListener :: forall f. EventEmitter -> String -> f -> Effect Unit

unsafeOn :: forall f. EffectFn3 EventEmitter String f Unit
unsafeOn = mkEffectFn3 \a b c -> gopursUnsafeOn a b c

unsafeOff :: forall f. EffectFn3 EventEmitter String f Unit
unsafeOff = mkEffectFn3 \a b c -> gopursUnsafeOff a b c

unsafeOnce :: forall f. EffectFn3 EventEmitter String f Unit
unsafeOnce = mkEffectFn3 \a b c -> gopursUnsafeOnce a b c

unsafePrependListener :: forall f. EffectFn3 EventEmitter String f Unit
unsafePrependListener = mkEffectFn3 \a b c -> gopursUnsafePrependListener a b c

unsafePrependOnceListener :: forall f. EffectFn3 EventEmitter String f Unit
unsafePrependOnceListener = mkEffectFn3 \a b c -> gopursUnsafePrependOnceListener a b c

foreign import gopursUnsafeEmitFn1 :: EventEmitter -> String -> Effect Boolean
foreign import gopursUnsafeEmitFn2 :: forall a. EventEmitter -> String -> a -> Effect Boolean
foreign import gopursUnsafeEmitFn3 :: forall a b. EventEmitter -> String -> a -> b -> Effect Boolean
foreign import gopursUnsafeEmitFn4 :: forall a b c. EventEmitter -> String -> a -> b -> c -> Effect Boolean
foreign import gopursUnsafeEmitFn5 :: forall a b c d. EventEmitter -> String -> a -> b -> c -> d -> Effect Boolean
foreign import gopursUnsafeEmitFn6 :: forall a b c d e. EventEmitter -> String -> a -> b -> c -> d -> e -> Effect Boolean
foreign import gopursUnsafeEmitFn7 :: forall a b c d e f. EventEmitter -> String -> a -> b -> c -> d -> e -> f -> Effect Boolean
foreign import gopursUnsafeEmitFn8 :: forall a b c d e f g. EventEmitter -> String -> a -> b -> c -> d -> e -> f -> g -> Effect Boolean
foreign import gopursUnsafeEmitFn9 :: forall a b c d e f g h. EventEmitter -> String -> a -> b -> c -> d -> e -> f -> g -> h -> Effect Boolean

unsafeEmitFn1 :: EffectFn2 EventEmitter String Boolean
unsafeEmitFn1 = mkEffectFn2 \a b -> gopursUnsafeEmitFn1 a b

unsafeEmitFn2 :: forall a. EffectFn3 EventEmitter String a Boolean
unsafeEmitFn2 = mkEffectFn3 \a b c -> gopursUnsafeEmitFn2 a b c

unsafeEmitFn3 :: forall a b. EffectFn4 EventEmitter String a b Boolean
unsafeEmitFn3 = mkEffectFn4 \a b c d -> gopursUnsafeEmitFn3 a b c d

unsafeEmitFn4 :: forall a b c. EffectFn5 EventEmitter String a b c Boolean
unsafeEmitFn4 = mkEffectFn5 \a b c d e -> gopursUnsafeEmitFn4 a b c d e

unsafeEmitFn5 :: forall a b c d. EffectFn6 EventEmitter String a b c d Boolean
unsafeEmitFn5 = mkEffectFn6 \a b c d e f -> gopursUnsafeEmitFn5 a b c d e f

unsafeEmitFn6 :: forall a b c d e. EffectFn7 EventEmitter String a b c d e Boolean
unsafeEmitFn6 = mkEffectFn7 \a b c d e f g -> gopursUnsafeEmitFn6 a b c d e f g

unsafeEmitFn7 :: forall a b c d e f. EffectFn8 EventEmitter String a b c d e f Boolean
unsafeEmitFn7 = mkEffectFn8 \a b c d e f g h -> gopursUnsafeEmitFn7 a b c d e f g h

unsafeEmitFn8 :: forall a b c d e f g. EffectFn9 EventEmitter String a b c d e f g Boolean
unsafeEmitFn8 = mkEffectFn9 \a b c d e f g h i -> gopursUnsafeEmitFn8 a b c d e f g h i

unsafeEmitFn9 :: forall a b c d e f g h. EffectFn10 EventEmitter String a b c d e f g h Boolean
unsafeEmitFn9 = mkEffectFn10 \a b c d e f g h i j -> gopursUnsafeEmitFn9 a b c d e f g h i j
