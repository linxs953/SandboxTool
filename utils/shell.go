package utils

import (
	"github.com/rs/zerolog/log"
	"os/exec"
)

func Open(url string) {
	openCmd := exec.Command("open","-a","Google Chrome", url)
	err := openCmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Open URL error")
		return
	}
	log.Print("Already Open in Browser")
}
