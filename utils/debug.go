package utils

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func GetDebugURL(env string, backendURL string) string {
	if env == "prod" {
		return GetStringEnv("PROD", "")
	} else if env == "test" {
		return GetStringEnv("TEST", "")
	} else if env == "local" {
		local := GetStringEnv("LOCAL", "")
		if local == "" {
			return ""
		}
		url := fmt.Sprintf(local, backendURL)
		return url
	} else {
		log.Printf("Can not support env %s", env)
		return ""
	}
}
