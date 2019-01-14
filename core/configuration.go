package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration structure
type Configuration struct {
	Security SecurityConfiguration
}

// SecurityConfiguration security configuration struct
type SecurityConfiguration struct {
	Secret string
}

// GetConfiguration returns a single instance of configuration
func GetConfiguration() (*Configuration, error) {

	evar := os.Getenv("ORI_ENV")
	scnf := "development"
	if evar != "" {
		scnf = evar
	}

	fil, err := os.Open(fmt.Sprintf("settings.%s.json", scnf))
	if err != nil {
		return nil, err
	}

	cnf := &Configuration{}
	if err := json.NewDecoder(fil).Decode(cnf); err != nil {
		return nil, err
	}

	if err := fil.Close(); err != nil {
		return nil, err
	}

	return cnf, nil
}
