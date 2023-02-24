package obs

import (
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/filters"
	"github.com/gkstretton/dark/services/goo/config"
	"github.com/gkstretton/dark/services/goo/keyvalue"
	"gopkg.in/yaml.v3"
)

// // I can't find a call for v5 websockets in this client api
// func setSessionNumber(number int) {
// 	if c == nil {
// 		fmt.Println("cannot set obs session number because client is nil")
// 		return
// 	}
// 	_, err := c.Sources.SetTextFreetype2Properties(&sources.SetTextFreetype2PropertiesParams{
// 		Source: "Session Number",
// 		Text:   fmt.Sprintf("%d", number),
// 	})
// 	if err != nil {
// 		fmt.Printf("error setting obs session number: %v\n", err)
// 		return
// 	}
// }

func setCropConfig() {
	if c == nil {
		fmt.Println("cannot set obs crop config because client is nil")
		return
	}

	err := setSourceCrop("TopCam", config.CC_TOP_CAM)
	if err != nil {
		fmt.Printf("failed to set OBS TopCam crop settings: %v\n", err)
	}

	err = setSourceCrop("FrontCam", config.CC_FRONT_CAM)
	if err != nil {
		fmt.Printf("failed to set OBS FrontCam crop settings: %v\n", err)
	}
}

func setSourceCrop(sourceName, cropConfigKey string) error {
	cc := &config.CropConfig{}
	b := keyvalue.Get(cropConfigKey)
	err := yaml.Unmarshal(b, cc)
	if err != nil {
		return fmt.Errorf("failed to unmarshall crop config: %v", err)
	}

	settings := map[string]interface{}{
		"left":   cc.LeftRel,
		"top":    cc.TopRel,
		"right":  cc.RightRel,
		"bottom": cc.BottomRel,
	}
	fmt.Printf("setting %s in obs to settings %+v\n", sourceName, settings)

	resp, err := c.Filters.SetSourceFilterSettings(&filters.SetSourceFilterSettingsParams{
		SourceName:     sourceName,
		FilterName:     "Crop",
		FilterSettings: settings,
	})
	if err != nil {
		return fmt.Errorf("failed to set source filter settings: %v", err)
	}
	fmt.Println(resp)

	return nil
}
