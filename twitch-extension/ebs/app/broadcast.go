package app

/*
func (e *extensionSession) regularRobotStatusUpdate() {
	c := events.Subscribe()
	for {
		select {
		case <-e.exitCh:
			l.Println("exiting regular robot status update loop")
			return
		case sr := <-c:
			if sr == nil {
				e.updateRobotStatus(nil)
				continue
			}
			e.updateRobotStatus(&robotStatus{
				Status: sr.Status.String(),
			})
		}
	}
}

// ManualTriggerBroadcast sends early, used after a significant update to
// improve responsiveness
func (e *extensionSession) ManualTriggerBroadcast() {
	if e.cleanUpDone {
		return
	}
	e.triggerBroadcast <- struct{}{}
}

func (e *extensionSession) updateRobotStatus(data *robotStatus) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.broadcastDataCache.RobotStatus = data
}

// broadcasts the BroadcastData cache once per second
func (e *extensionSession) regularBroadcast() {
	// get marshaled data, protected by lock
	d := func() ([]byte, error) {
		e.lock.Lock()
		defer e.lock.Unlock()

		jsonData, err := json.Marshal(e.broadcastDataCache)
		if err != nil {
			return nil, err
		}
		return jsonData, nil
	}

	send := func() {
		data, err := d()
		if err != nil {
			l.Printf("failed to marshal broadcast data: %v\n", err)
			return
		}
		err = e.broadcastData(data)
		if err != nil {
			l.Printf("failed to send broadcast data: %v\n", err)
			return
		}
	}

	next := time.After(0)
	for {
		select {
		case <-e.exitCh:
			l.Println("exiting regularBroadcast loop")
			return
		case <-e.triggerBroadcast:
			next = time.After(time.Second * 2)
			send()
		case <-next:
			next = time.After(time.Millisecond * 1100)
			send()
		}
	}
}

*/
