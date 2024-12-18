package ebsinterface

import "github.com/gkstretton/dark/services/goo/types"

func (e *extensionSession) GetEbsState() *types.EbsStateReport {
	e.ebsStateLock.Lock()
	defer e.ebsStateLock.Unlock()

	return e.ebsState
}
