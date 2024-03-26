package contentscheduler

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"github.com/gkstretton/dark/services/goo/socialmedia"
)

// control access to contentplan
var lock sync.Mutex

type manager struct {
	s *socialmedia.SocialManager
}

func getManager() *manager {
	return &manager{
		s: socialmedia.NewSocialManager(),
	}
}

// poll every minute, process session content plan to upload whenever one's due
func Start(sm *session.SessionManager) {
	m := getManager()

	mqtt.Subscribe(topics_backend.TRIGGER_UPLOAD_FROM_CONTENT_PLAN, func(topic string, payload []byte) {
		pl := string(payload)
		sessionNum, err := strconv.Atoi(pl)
		if err != nil {
			fmt.Printf("failed to get sessionnumber from payload '%s'\n", pl)
			return
		}
		err = m.processSession(uint64(sessionNum))
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	next := time.After(time.Second)
	for {
		<-next
		next = time.After(time.Minute * 10)
		if !keyvalue.GetBool("ENABLE_CONTENT_SCHEDULER_LOOP") {
			fmt.Println("content scheduler loop not enabled")
			next = time.After(time.Hour)
			continue
		}

		err := m.processLatestProductionSession(sm)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func Test(sm *session.SessionManager) {
	m := getManager()
	fmt.Println(m.processLatestProductionSession(sm))
}

func (m *manager) processLatestProductionSession(sm *session.SessionManager) error {
	session, err := sm.GetLatestProdutionSession()
	if err != nil {
		fmt.Printf("error getting prod session in contentscheduler: %v\n", err)
	}

	return m.processSession(uint64(session.Id))
}

func (m *manager) processSession(sessionNumber uint64) error {
	path := filesystem.GetContentPlanPath(sessionNumber)
	return m.processContentPlan(path, sessionNumber)
}
