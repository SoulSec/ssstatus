package ssstatus

import (
	"encoding/json"

	"github.com/SoulSec/ssstatus/validate"
)

type Config struct {
	Servers  Servers   `json:"servers"`
	Settings *Settings `json:"settings"`
}

func NewConfig(jsonData []byte) *Config {
	config := &Config{}
	err := json.Unmarshal(jsonData, config)
	if err != nil {
		panic("Error parsing json configuration data")
	}
	if err := validate.ValidateAll(config); err != nil {
		panic(err)
	}
	return config
}
