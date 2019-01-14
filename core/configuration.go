package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseConfiguration returns a single instance of configuration
func ParseConfiguration(configuration interface{}) error {

	evar := os.Getenv("ORI_ENV")
	scnf := "development"
	if evar != "" {
		scnf = evar
	}

	fil, err := os.Open(fmt.Sprintf("settings.%s.json", scnf))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(fil).Decode(configuration); err != nil {
		return err
	}

	if err := fil.Close(); err != nil {
		return err
	}

	return nil
}
