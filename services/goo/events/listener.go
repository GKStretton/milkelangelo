package events

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/util/protoyaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var latest_state_report *machinepb.StateReport

var subs = []chan *machinepb.StateReport{}
var lock sync.Mutex

func Run(sm *session.SessionManager) {
	mqtt.Subscribe(topics_firmware.TOPIC_STATE_REPORT_RAW, func(topic string, payload []byte) {
		t := time.Now().UnixMicro()

		sr := &machinepb.StateReport{}
		err := proto.Unmarshal(payload, sr)
		if err != nil {
			fmt.Printf("error unmarshalling state report: %v\n", err)
			return
		}
		sr.TimestampUnixMicros = uint64(t)
		sr.TimestampReadable = time.UnixMicro(t).
			Format("2006-03-02 15:04:05.000000")

		session, _ := sm.GetLatestSession()
		if session != nil {
			sr.Paused = session.Paused
		}

		latest_state_report = sr
		publishStateReport(sr)

		// Abort unless session is active or paused
		if session != nil && !session.Complete {
			saveSessionStateReport(session, sr)
		}
	})

	mqtt.Subscribe(topics_backend.TOPIC_MARK_FAILED_DISPENSE, func(topic string, payload []byte) {
		session, _ := sm.GetLatestSession()
		if session != nil && !session.Complete {
			appendFailedDispense(
				uint64(session.Id),
				latest_state_report.StartupCounter,
				uint64(latest_state_report.PipetteState.DispenseRequestNumber),
			)
		}
	})
}

func Subscribe() chan *machinepb.StateReport {
	lock.Lock()
	defer lock.Unlock()
	c := make(chan *machinepb.StateReport, 10)
	subs = append(subs, c)
	return c
}

// publish to internal channels and to broker
func publishStateReport(sr *machinepb.StateReport) {
	lock.Lock()
	// internal
	for _, c := range subs {
		select {
		case c <- sr:
		default:
		}
	}
	lock.Unlock()

	// broker
	m := protojson.MarshalOptions{
		Multiline:       true,
		UseProtoNames:   true,
		Indent:          "\t",
		EmitUnpopulated: true,
	}
	b, err := m.Marshal(sr)
	if err != nil {
		fmt.Printf("error marshalling state report: %v\n", err)
		return
	}
	err = mqtt.Publish(topics_backend.TOPIC_STATE_REPORT_JSON, string(b))
	if err != nil {
		fmt.Printf("error publishing json state report: %v\n", err)
		return
	}
}

func saveSessionStateReport(s *session.Session, sr *machinepb.StateReport) {
	list := &machinepb.StateReportList{
		StateReports: []*machinepb.StateReport{
			sr,
		},
	}
	output, err := protoyaml.Marshal(list)
	if err != nil {
		fmt.Printf("error marshalling state report to yaml: %v\n", err)
	}

	p := filesystem.GetStateReportPath(uint64(s.Id))

	var result = string(output)
	// remove the main key from additional reports
	if _, err := os.Stat(p); err == nil {
		result = strings.Replace(result, "StateReports:\n", "", 1)
	}

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file for state report storage: %v\n", err)
	}
	defer f.Close()
	f.Write([]byte(result))
}

func appendFailedDispense(sessionId, startupCounter, dispenseNumber uint64) {
	p := filesystem.GetFailedDispensesPath(sessionId)

	failedDispenses := &machinepb.DispenseMetadataMap{}

	data, err := os.ReadFile(p)
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error reading failed dispenses file: %v\n", err)
		return
	}

	if len(data) > 0 {
		err = protoyaml.Unmarshal(data, failedDispenses)
		if err != nil {
			fmt.Printf("Error unmarshalling failed dispenses: %v\n", err)
			return
		}
	}

	if failedDispenses.DispenseMetadata == nil {
		failedDispenses.DispenseMetadata = map[string]*machinepb.DispenseMetadata{}
	}

	key := fmt.Sprintf("%d_%d", startupCounter, dispenseNumber)

	if _, ok := failedDispenses.DispenseMetadata[key]; ok {
		failedDispenses.DispenseMetadata[key].FailedDispense = true
	} else {
		failedDispenses.DispenseMetadata[key] = &machinepb.DispenseMetadata{
			StartupCounter:        startupCounter,
			DispenseRequestNumber: dispenseNumber,
			FailedDispense:        true,
		}
	}

	data, err = protoyaml.Marshal(failedDispenses)
	if err != nil {
		fmt.Printf("Error marshalling failed dispenses: %v\n", err)
		return
	}

	err = os.WriteFile(p, data, 0644)
	if err != nil {
		fmt.Printf("Error writing failed dispenses file: %v\n", err)
		return
	}

	fmt.Printf("wrote failed dispense to file (session %d, startup %d, dispense %d)\n", sessionId, startupCounter, dispenseNumber)
}
