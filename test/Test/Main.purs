module Test.Main where

import Prelude

import Effect (Effect)
import Effect.Aff (launchAff_)
import Test.Node.EventEmitter as EventEmitter
import Test.Spec.Reporter (consoleReporter)
import Test.Spec.Runner (runSpecPure)

main :: Effect Unit
main = launchAff_ $ runSpecPure [ consoleReporter ] do
  EventEmitter.spec
