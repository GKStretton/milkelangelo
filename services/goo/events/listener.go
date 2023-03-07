package events

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/util/protoyaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var subs = []chan *machinepb.StateReport{}

func Run(sm *session.SessionManager) {
	mqtt.Subscribe(config.TOPIC_STATE_REPORT_RAW, func(topic string, payload []byte) {
		t := time.Now().UnixMicro()

		sr := &machinepb.StateReport{}
		err := proto.Unmarshal(payload, sr)
		if err != nil {
			fmt.Printf("error unmarshalling state report: %v\n", err)
			return
		}
		sr.TimestampUnixMicros = uint64(t)
		// fmt.Printf("%+v\n", sr)

		// Abort unless session is active
		session, _ := sm.GetLatestSession()
		if session == nil || session.Complete || session.Paused {
			return
		}

		saveSessionStateReport(session, sr)
		publishStateReport(sr)
	})
}

func Subscribe() chan *machinepb.StateReport {
	c := make(chan *machinepb.StateReport, 10)
	subs = append(subs, c)
	return c
}

// publish to internal channels and to broker
func publishStateReport(sr *machinepb.StateReport) {
	// internal
	for _, c := range subs {
		select {
		case c <- sr:
		default:
		}
	}

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
	err = mqtt.Publish(config.TOPIC_STATE_REPORT_JSON, string(b))
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

	result := strings.Replace(string(output), "StateReports:\n", "", 1)

	p := filesystem.GetStateReportPath(uint64(s.Id))

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file for state report storage: %v\n", err)
	}
	defer f.Close()
	f.Write([]byte(result))
}
