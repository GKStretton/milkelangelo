package livecapture

import (
	"fmt"
	"os"

	"github.com/gkstretton/dark/services/goo/filesystem"
	"github.com/gkstretton/dark/services/goo/keyvalue"
)

func SaveCropConfig(ccKey string, contentPath string) error {
	// e.g. 1.mp4.yml
	ymlPath := contentPath + ".yml"
	config := keyvalue.Get(ccKey)
	if config == nil {
		return fmt.Errorf("cannot saveCropConfig of %s for '%s' because key not found", ccKey, contentPath)
	}
	err := os.WriteFile(ymlPath, config, 0666)
	if err != nil {
		return fmt.Errorf("failed to write cropConfig of %s to '%s': %v", ccKey, ymlPath, err)
	}
	filesystem.SetPerms(ymlPath)
	return nil
}
