package ebsinterface

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gkstretton/dark/services/goo/types"
)

// UpdateState updates the state that is sent to the EBS
func (e *extensionSession) UpdateState(f func(state *types.GooState)) {
	e.gooStateLock.Lock()
	defer e.gooStateLock.Unlock()

	f(&e.gooState)

	select {
	case e.gooStateChan <- struct{}{}:
	default:
	}
}

func (e *extensionSession) stateSender() {
	for {
		<-e.gooStateChan
		e.sendState()
	}
}

// SendState sends the current state to the EBS
func (e *extensionSession) sendState() {
	l.Printf("sending state to EBS...")
	result, err := url.JoinPath(e.ebsAddress, "/update-state")
	if err != nil {
		l.Printf("error forming ebs update state url: %s", err)
		return
	}

	// Marshal state data
	stateData, err := json.Marshal(e.gooState)
	if err != nil {
		l.Printf("error marshalling state data: %s", err)
		return
	}

	// PUT to /update-state
	req, err := http.NewRequest(http.MethodPut, result, bytes.NewReader(stateData))
	if err != nil {
		l.Printf("error creating request: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+e.ebsToken)

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.Printf("error sending request: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Printf("unexpected status code: %d", resp.StatusCode)
	}

	l.Printf("state sent to EBS: %+v", e.gooState)
}
