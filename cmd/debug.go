package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"runtimeops/utils"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug academy-system user",
	Long:  "Debug academy-system user",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("lack of args .\nUse debug -h for help")
			return
		}
		var reqURL string
		//var err error
		//err = config.Init()
		//if err != nil {
		//	log.Error().Err(err).Msg("Init Config error")
		//	return
		//}
		unionid := args[1]
		env := args[0]
		if unionid == "" || env == "" {
			log.Print("unionid is nil")
			return
		}
		tokenString := utils.GetTokenByUnionid(unionid)
		if tokenString == "" {
			log.Print("get token nil")
			return
		}
		backendCICD := utils.GetStringEnv("BACKENDCICD", "")
		if backendCICD == "" {
			log.Print("Get env nil")
			return
		}
		if env == "local" {
			reqURL = utils.GetDebugURL(env, backendCICD)
		} else {
			reqURL = utils.GetDebugURL(env, "")
		}
		debugURL := fmt.Sprintf("%s?token=%s", reqURL, tokenString)
		utils.Open(debugURL)
	},
}

func debugUsage() string {
	helpMsg := "e.g: \n" +
		"  rtctl debug [env] [unionid] (env can be prod or test)\n"
	helpMsg = fmt.Sprintf("%s\n\n%s\n%s", debugCmd.Long, debugCmd.UsageString(), helpMsg)
	return helpMsg
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.SetHelpTemplate(debugUsage())
}
