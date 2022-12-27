package keyvalue

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gkstretton/dark/services/goo/filesystem"
)

func setKeyValue(key string, value []byte) error {
	lock.Lock()
	defer lock.Unlock()

	// open file
	p := filepath.Join(filesystem.GetKeyValueStorePath(), key)
	err := os.WriteFile(p, value, 0666)
	if err != nil {
		return fmt.Errorf("failed to write value to key %s at %s: %v", key, p, err)
	}
	return nil
}

func getKeyValue(key string) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()

	// open file
	p := filepath.Join(filesystem.GetKeyValueStorePath(), key)
	value, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read value to key %s at %s: %v", key, p, err)
	}
	return value, nil
}
