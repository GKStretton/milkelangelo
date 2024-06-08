package events

import "github.com/gkstretton/asol-protos/go/machinepb"

func ConditionWaiter(cond func(*machinepb.StateReport) bool) chan *machinepb.StateReport {
	c := Subscribe()
	defer Unsubscribe(c)

	filterChan := make(chan *machinepb.StateReport)
	go func() {
		for {
			r := <-c
			if cond(r) {
				filterChan <- r
				close(filterChan)
				return
			}
		}
	}()
	return filterChan
}
