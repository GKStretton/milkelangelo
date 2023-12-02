// vialprofiles handles saving a snapshot of the system vial profile config
// when a session starts.
package vialprofiles

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gkstretton/asol-protos/go/machinepb"
	"github.com/gkstretton/asol-protos/go/topics_backend"
	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"github.com/gkstretton/dark/services/goo/mqtt"
	"github.com/gkstretton/dark/services/goo/session"
	"google.golang.org/protobuf/encoding/protojson"
)

func Start(sm *session.SessionManager) {
	// debug topic, session number as payload
	mqtt.Subscribe("asol/testing/save-profile-snapshot", func(topic string, payload []byte) {
		id, err := strconv.Atoi(string(payload))
		if err != nil {
			fmt.Printf("cannot convert id '%s' to int for profile snapshot: %v\n", payload, err)
		}
		saveSnapshot(uint64(id))
	})

	go func() {
		c := sm.SubscribeToEvents()

		for {
			e := <-c
			if e.Type == session.SESSION_STARTED {
				// save system vial configuration snapshot for this session
				saveSnapshot(uint64(e.SessionID))
			}
		}
	}()
}

func saveSnapshot(sessionId uint64) {
	snapshot := GetSystemVialConfigurationSnapshot()

	path := filesystem.GetVialProfileSnapshotPath(sessionId)
	err := filesystem.WriteProtoYaml(path, snapshot)
	if err != nil {
		fmt.Printf("failed to save vial profile snapshot at start of session: %v\n", err)
	}
}

// returns the map of vial position -> profile
// probably for saving at the start of a session
func GetSystemVialConfigurationSnapshot() *machinepb.SystemVialConfigurationSnapshot {
	all := getAllProfiles()
	positionConfig := getSystemProfileConfiguration()

	snapshot := &machinepb.SystemVialConfigurationSnapshot{
		Profiles: map[uint64]*machinepb.VialProfile{},
	}

	for pos, profileId := range positionConfig.Vials {
		snapshot.Profiles[pos] = all.Profiles[profileId]
	}

	return snapshot
}

func getAllProfiles() *machinepb.VialProfileCollection {
	allProfiles := &machinepb.VialProfileCollection{}
	raw := keyvalue.Get(topics_backend.KV_KEY_ALL_VIAL_PROFILES)
	err := protojson.Unmarshal(raw, allProfiles)
	if err != nil {
		fmt.Printf("failed to unmarshal vial profiles from kv: %v", err)
		return &machinepb.VialProfileCollection{}
	}

	return allProfiles
}

func getSystemProfileConfiguration() *machinepb.SystemVialConfiguration {
	raw := keyvalue.Get(topics_backend.KV_KEY_SYSTEM_VIAL_PROFILES)
	systemConf := &machinepb.SystemVialConfiguration{}
	err := protojson.Unmarshal(raw, systemConf)
	if err != nil {
		fmt.Printf("failed to unmarshal system vial conf from kv: %v", err)
		return &machinepb.SystemVialConfiguration{}
	}

	return systemConf
}

// for a given robot vial position, get the vial profile
// does a lookup against the system configuration
// returns nil if can't find
func GetSystemVialProfile(vialPosition int) *machinepb.VialProfile {
	snapshot := GetSystemVialConfigurationSnapshot()
	return snapshot.Profiles[uint64(vialPosition)]
}

func GetVialOptionsAndMap() ([]string, map[uint64]string) {
	vialConfig := GetSystemVialConfigurationSnapshot()

	options := []string{}
	vialPosToName := map[uint64]string{}
	for posNo, profile := range vialConfig.GetProfiles() {
		vialPosToName[posNo] = profile.Name
		options = append(options, strings.ToLower(profile.Name))
	}
	return options, vialPosToName
}
