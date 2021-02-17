package main

import (
	"github.com/rs/zerolog/log"
	"runtimeops/cmd"
	"runtimeops/config"
)

func init() {
	if err := config.Init(); err != nil {
		log.Error().Err(err).Msg("Init config error")
		return
	}
}

func main() {
	cmd.Execute()
}
