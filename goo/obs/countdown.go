package obs

import (
	"fmt"
	"sync"
	"time"

	"github.com/andreykaipov/goobs/api/requests/inputs"
)

var (
	countdownLock      sync.Mutex
	countdownTimestamp time.Time
	countdownTitle     string
)

func SetCountdown(title string, timestamp time.Time) {
	countdownLock.Lock()
	defer countdownLock.Unlock()

	countdownTitle = title
	countdownTimestamp = timestamp
}

func countdownRunner() {
	for {
		countdownLoop()
		time.Sleep(time.Second)
	}
}

func countdownLoop() {
	countdownLock.Lock()
	defer countdownLock.Unlock()

	if countdownTitle == "" {
		return
	}

	until := time.Until(countdownTimestamp)
	if until < 0 {
		countdownTitle = ""
		updateCountdownElement("", "")
		return
	}

	s := int(until.Seconds()) % 60
	tStr := fmt.Sprintf("%ds", s)

	m := int(until.Minutes())
	if m > 0 {
		tStr = fmt.Sprintf("%dm ", m) + tStr
	}

	updateCountdownElement(countdownTitle, tStr)
}

func updateCountdownElement(title, dur string) {
	scene, err := getScene()
	if err != nil {
		fmt.Println(err)
	}
	if !(scene == SCENE_LIVE || scene == SCENE_COMPLETE) {
		return
	}

	_, err = c.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName: "Countdown Title",
		InputSettings: map[string]interface{}{
			"text": fmt.Sprintf("%23s", title),
		},
	})
	if err != nil {
		fmt.Printf("error setting obs countdown title: %v\n", err)
	}

	_, err = c.Inputs.SetInputSettings(&inputs.SetInputSettingsParams{
		InputName: "Countdown Timer",
		InputSettings: map[string]interface{}{
			"text": fmt.Sprintf("%10s", dur),
		},
	})
	if err != nil {
		fmt.Printf("error setting obs countdown timer: %v\n", err)
	}
}
