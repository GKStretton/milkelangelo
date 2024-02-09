package scheduler

import "sync"

type AutomationLock struct {
	isRunning bool
	lock      sync.Mutex
}

func (a *AutomationLock) Set(isRunning bool) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.isRunning = isRunning
}

func (a *AutomationLock) Get() bool {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.isRunning
}
