package events

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/asol-protos/go/topics_firmware"
	"github.com/gkstretton/dark/services/goo/ebsinterface"
	"github.com/gkstretton/dark/services/goo/email"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/types"
	"github.com/gkstretton/dark/services/goo/util/protoyaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var l = log.New(os.Stdout, "[events] ", log.Flags())
var latest_state_report *machinepb.StateReport
var srLock sync.Mutex

var subs = []chan *machinepb.StateReport{}

var subsLock sync.Mutex

func Start(sm *session.SessionManager, ebsApi ebsinterface.EbsApi) {
	mqtt.Subscribe(topics_firmware.TOPIC_STATE_REPORT_RAW, func(topic string, payload []byte) {
		l.Println("GOT STATE REPORT")
		t := time.Now().UnixMicro()

		sr := &machinepb.StateReport{}
		err := proto.Unmarshal(payload, sr)
		if err != nil {
			l.Printf("error unmarshalling state report: %v\n", err)
			return
		}
		sr.TimestampUnixMicros = uint64(t)
		sr.TimestampReadable = time.UnixMicro(t).
			Format("2006-01-02 15:04:05.000000")

		session, _ := sm.GetLatestSession()
		if session != nil {
			sr.Paused = session.Paused
			sr.LatestDslrFileNumber = filesystem.GetLatestDslrFileNumber(uint64(session.Id))
		}

		srLock.Lock()
		defer srLock.Unlock()

		latest_state_report = sr
		publishStateReport(sr)

		if ebsApi != nil {
			ebsApi.UpdateState(func(state *types.GooState) {
				state.X = sr.MovementDetails.TargetXUnit
				state.Y = sr.MovementDetails.TargetYUnit
				state.Status = types.GooStatusUnknown

				if sr.CollectionRequest == nil {
					state.CollectionState = nil
				} else {
					state.CollectionState = &types.CollectionState{
						VialNumber: int(sr.CollectionRequest.VialNumber),
						VolumeUl:   sr.CollectionRequest.VolumeUl,
						Completed:  sr.CollectionRequest.Completed,
					}
				}

				if sr.PipetteState == nil {
					state.DispenseState = nil
				} else {
					state.DispenseState = &types.DispenseState{
						VialNumber:        int(sr.PipetteState.VialHeld),
						VolumeRemainingUl: sr.PipetteState.VolumeTargetUl,
						Completed:         sr.PipetteState.Spent,
					}
				}
			})
		}

		// only save state reports for ongoing sessions
		if session != nil && !session.Complete {
			saveSessionStateReport(session, sr)
		}
	})

	mqtt.Subscribe(topics_backend.TOPIC_MARK_FAILED_DISPENSE, func(topic string, payload []byte) {
		go func() {
			srLock.Lock()
			defer srLock.Unlock()

			if latest_state_report.PipetteState.DispenseRequestNumber < 1 {
				l.Println("cannot mark dispense with number < 1, it means nothing's dispensed yet...")
				return
			}

			session, _ := sm.GetLatestSession()
			if session != nil && !session.Complete {
				appendFailedDispense(
					uint64(session.Id),
					latest_state_report.StartupCounter,
					uint64(latest_state_report.PipetteState.DispenseRequestNumber),
				)
			}
		}()
	})

	mqtt.Subscribe(topics_backend.TOPIC_MARK_DELAYED_DISPENSE, func(topic string, payload []byte) {
		go func() {
			srLock.Lock()
			defer srLock.Unlock()

			if latest_state_report.PipetteState.DispenseRequestNumber < 1 {
				l.Println("cannot mark dispense with number < 1, it means nothing's dispensed yet...")
				return
			}

			// could change to payload in future
			delayMs := uint64(1000)

			session, _ := sm.GetLatestSession()
			if session != nil && !session.Complete {
				appendDelayedDispense(
					uint64(session.Id),
					latest_state_report.StartupCounter,
					uint64(latest_state_report.PipetteState.DispenseRequestNumber),
					delayMs,
				)
			}
		}()
	})

	mqtt.Subscribe(topics_firmware.TOPIC_LOGS_CRIT, func(topic string, payload []byte) {
		go func() {
			critErr := string(payload)
			err := email.SendEmail(&machinepb.Email{
				Subject:   "Crit. f/w error: " + critErr,
				Body:      critErr,
				Recipient: machinepb.EmailRecipient_EMAIL_RECIPIENT_MAINTENANCE,
			})
			if err != nil {
				l.Printf("error sending email for critial firmware error (%s): %s\n", critErr, err)
				return
			}
			l.Printf("Emailed maintainer about critial error: %s\n", critErr)
		}()
	})

	RequestStateReport()

	listenForEbsConnect(ebsApi)
}

func listenForEbsConnect(ebsApi ebsinterface.EbsApi) {
	if ebsApi == nil {
		return
	}

	c := ebsApi.SubscribeMessages()
	defer ebsApi.UnsubscribeMessages(c)

	for {
		msg := <-c
		if msg.Type == types.EbsConnectedEvent {
			RequestStateReport()
		}
	}
}

func Subscribe() chan *machinepb.StateReport {
	subsLock.Lock()
	srLock.Lock()
	defer subsLock.Unlock()
	defer srLock.Unlock()

	c := make(chan *machinepb.StateReport, 10)
	subs = append(subs, c)

	// send latest state report to channel
	c <- latest_state_report

	return c
}

func Unsubscribe(c chan *machinepb.StateReport) {
	subsLock.Lock()
	defer subsLock.Unlock()

	for i, sub := range subs {
		if sub == c {
			subs = append(subs[:i], subs[i+1:]...)
			close(c)
			break
		}
	}
}

// publish to internal channels and to broker
func publishStateReport(sr *machinepb.StateReport) {
	subsLock.Lock()
	// internal
	for _, c := range subs {
		select {
		case c <- sr:
		default:
		}
	}
	subsLock.Unlock()

	// broker
	m := protojson.MarshalOptions{
		Multiline:       true,
		UseProtoNames:   true,
		Indent:          "\t",
		EmitUnpopulated: true,
	}
	b, err := m.Marshal(sr)
	if err != nil {
		l.Printf("error marshalling state report: %v\n", err)
		return
	}
	go func() {
		err = mqtt.Publish(topics_backend.TOPIC_STATE_REPORT_JSON, string(b))
		if err != nil {
			l.Printf("error publishing json state report: %v\n", err)
			return
		}
	}()
}

func saveSessionStateReport(s *session.Session, sr *machinepb.StateReport) {
	list := &machinepb.StateReportList{
		StateReports: []*machinepb.StateReport{
			sr,
		},
	}
	output, err := protoyaml.Marshal(list)
	if err != nil {
		l.Printf("error marshalling state report to yaml: %v\n", err)
	}

	p := filesystem.GetStateReportPath(uint64(s.Id))

	var result = string(output)
	// remove the main key from additional reports
	if _, err := os.Stat(p); err == nil {
		result = strings.Replace(result, "StateReports:\n", "", 1)
	}

	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.Printf("error opening file for state report storage: %v\n", err)
	}
	defer f.Close()
	f.Write([]byte(result))
}

func appendFailedDispense(sessionId, startupCounter, dispenseNumber uint64) {
	p := filesystem.GetDispenseMetadataPath(sessionId)

	meta := &machinepb.DispenseMetadataMap{}

	data, err := os.ReadFile(p)
	if err != nil && !os.IsNotExist(err) {
		l.Printf("Error reading failed dispenses file: %v\n", err)
		return
	}

	if len(data) > 0 {
		err = protoyaml.Unmarshal(data, meta)
		if err != nil {
			l.Printf("Error unmarshalling failed dispenses: %v\n", err)
			return
		}
	}

	if meta.DispenseMetadata == nil {
		meta.DispenseMetadata = map[string]*machinepb.DispenseMetadata{}
	}

	key := fmt.Sprintf("%d_%d", startupCounter, dispenseNumber)

	if _, ok := meta.DispenseMetadata[key]; ok {
		meta.DispenseMetadata[key].FailedDispense = true
		meta.DispenseMetadata[key].DispenseDelayMs = 0
	} else {
		meta.DispenseMetadata[key] = &machinepb.DispenseMetadata{
			FailedDispense:  true,
			DispenseDelayMs: 0,
		}
	}

	data, err = protoyaml.Marshal(meta)
	if err != nil {
		l.Printf("Error marshalling failed dispenses: %v\n", err)
		return
	}

	err = os.WriteFile(p, data, 0644)
	if err != nil {
		l.Printf("Error writing failed dispenses file: %v\n", err)
		return
	}

	l.Printf("wrote failed dispense to file (session %d, startup %d, dispense %d)\n", sessionId, startupCounter, dispenseNumber)
}

func appendDelayedDispense(sessionId, startupCounter, dispenseNumber, delayMs uint64) {
	p := filesystem.GetDispenseMetadataPath(sessionId)

	meta := &machinepb.DispenseMetadataMap{}

	data, err := os.ReadFile(p)
	if err != nil && !os.IsNotExist(err) {
		l.Printf("Error reading failed dispenses file: %v\n", err)
		return
	}

	if len(data) > 0 {
		err = protoyaml.Unmarshal(data, meta)
		if err != nil {
			l.Printf("Error unmarshalling failed dispenses: %v\n", err)
			return
		}
	}

	if meta.DispenseMetadata == nil {
		meta.DispenseMetadata = map[string]*machinepb.DispenseMetadata{}
	}

	key := fmt.Sprintf("%d_%d", startupCounter, dispenseNumber)

	if _, ok := meta.DispenseMetadata[key]; ok {
		meta.DispenseMetadata[key].FailedDispense = false
		meta.DispenseMetadata[key].DispenseDelayMs = delayMs
	} else {
		meta.DispenseMetadata[key] = &machinepb.DispenseMetadata{
			FailedDispense:  false,
			DispenseDelayMs: delayMs,
		}
	}

	data, err = protoyaml.Marshal(meta)
	if err != nil {
		l.Printf("Error marshalling delayed dispenses: %v\n", err)
		return
	}

	err = os.WriteFile(p, data, 0644)
	if err != nil {
		l.Printf("Error writing delayed dispenses file: %v\n", err)
		return
	}

	l.Printf("wrote delayed dispense to file (session %d, startup %d, dispense %d)\n", sessionId, startupCounter, dispenseNumber)
}

func GetLatestStateReportCopy() *machinepb.StateReport {
	srLock.Lock()
	defer srLock.Unlock()

	return proto.Clone(latest_state_report).(*machinepb.StateReport)
}

func RequestStateReport() {
	err := mqtt.Publish(topics_firmware.TOPIC_STATE_REPORT_REQUEST, "")
	if err != nil {
		l.Printf("failed to request state report from firmware: %v\n", err)
	}
}
