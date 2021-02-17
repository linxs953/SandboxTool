package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"runtimeops/utils"
)

var startPuid bool
var version string
var startIt bool

const puidLength = 36

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a sandbox",
	Long:  "Start a sandbox",
	Run: func(cmd *cobra.Command, args []string) {
		//if len(args) < 1 {
		//	fmt.Println("lack of args .\nUse start -h for help")
		//	return
		//}
		var storageId string
		var puid string
		//if err := config.Init(); err != nil {
		//	log.Error().Err(err).Msg("Init config error")
		//	return
		//}
		puidFlag, err := cmd.Flags().GetBool("puid")
		if err != nil {
			return
		}
		itFlag, err := cmd.Flags().GetBool("it")
		if err != nil {
			return
		}
		versionFlag, err := cmd.Flags().GetString("version")
		if err != nil {
			return
		}
		if versionFlag == "" {
			versionFlag = "v2"
		} else {
			versionFlag = "v" + versionFlag
		}
		switch puidFlag {
		case true:
			if len(args) == 0 {
				puid = utils.DefaultPuid
			} else {
				puid = args[0]
				if len(puid) < puidLength {
					log.Print("puid format is not invalid")
					return
				}
			}
			if len(args) < 2 {
				storageId = ""
			} else {
				storageId = args[1]
			}
			utils.StartSandboxByPuid(puid, storageId, versionFlag, itFlag)
		case false:
			if len(args) == 0 {
				log.Print("unionid is empty")
				fmt.Println("Use start -h for help")
				return
			}
			utils.StartSandboxByUnionID(args[0], versionFlag, itFlag)
		}
	},
}

func startHelpMessage() string {
	helpMsg := "e.g : \n" +
		"  rtctl start [puid] [storage_id] --puid=true --it=true --version=1 (已携带puid情况下storage_id可为空) \n" +
		"  rtctl start [unionid] --version=1 --it=true"
	helpMsg = fmt.Sprintf("%s\n\n%s\n%s", startCmd.Long, startCmd.UsageString(), helpMsg)
	return helpMsg
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&startPuid, "puid", "p", false, "Use --puid=true to get sandbox by puid")
	startCmd.Flags().StringVarP(&version, "version", "v", "", "Use --version=1 to get runtime v1 sandbox")
	startCmd.Flags().BoolVarP(&startIt, "it", "", false, "Use --it=true to entrance sandbox when it started")
	startCmd.SetHelpTemplate(startHelpMessage())
}
