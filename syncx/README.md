# `syncx` and `atomx` packages

## Deterministic Race Detection

The main problem with running tests with race detector is that your test become flaky: sometimes it can find race,
sometimes it cannot, sometimes it find another one. You can also face with another problem: when running with too few
GOMAXPROCS, races can be stayed undetected. To bring consistency and determinism, and force random switches between
parallel goroutines, these two packages use fault injection when are built with `race` flag, and use standard atomic
package without the flag. To test your algorithm, use `synctest.Stress` function. 

## More primitives

The `syncx` library appears with `Semaphore` and `Barrier` primitives. `Semaphore` is very useful when you want to
limit amount of parallel processes running in critical section (for example, in cgo calls to prevent spawning too
many threads), and `Barrier` is useful for tests (use it for synchronizing starts of parallel processes).
