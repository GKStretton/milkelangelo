package twitchapi

// consists of latest state report
type BroadcastData struct {
	RobotStatus *robotStatus
}

type robotStatus struct {
	Status string
}
