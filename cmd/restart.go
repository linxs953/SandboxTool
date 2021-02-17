package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"runtimeops/utils"
	_ "runtimeops/utils"
)

var restartPuid bool
var restartIt bool
var restartVersion string

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Sandbox Pod",
	Long:  "Restart  Sandbox Pod",
	Run: func(cmd *cobra.Command, args []string) {
		//if len(args) < 1 {
		//	fmt.Println("lack of args .\nUse restart -h for help")
		//	return
		//}
		var puid string
		var version string
		//var storageid  = ""
		//if len(args) > 2 {
		//	storageid = args[1]
		//}
		itFlagV, err := cmd.Flags().GetBool("it")
		if err != nil {
			log.Print("get [it] param error")
			return
		}
		versionFlagV, err := cmd.Flags().GetString("version")
		if err != nil {
			log.Print("get [version] param error")
			return
		}
		if versionFlagV == "" {
			version = "v2"
		} else {
			version = "v" + versionFlagV
		}
		puidFlagV, err := cmd.Flags().GetBool("puid")
		if err != nil {
			log.Print("get [puid] param error")
			return
		}
		switch puidFlagV {
		case true:
			if len(args) == 0 {
				puid = utils.DefaultPuid
			} else {
				puid = args[0]
			}
			utils.RestartByPuid(puid, version, itFlagV)
		case false:
			if len(args) == 0 {
				log.Print("unionid is empty")
				fmt.Println("Use restart -h for help")
				return
			}
			utils.RestartByUnionID(args[0], version, itFlagV)
		}
	},
}

func restartHelpMessage() string {
	helpMsg := "e.g :\n" +
		"  rtctl restart [unionid] --it=true \n" +
		"  rtctl restart [puid] --it=true --puid=true"
	helpMsg = fmt.Sprintf("%s\n\n%s\n%s\n", restartCmd.Long, restartCmd.UsageString(), helpMsg)
	return helpMsg
}

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.Flags().BoolVarP(&restartPuid, "puid", "p", false, "Use --puid=true to restart sandbox by puid")
	restartCmd.Flags().BoolVarP(&restartIt, "it", "", false, "Use --it=true to entrance sandbox when restarted")
	restartCmd.Flags().StringVarP(&restartVersion, "version", "v", "", "Use --version=1 to restart runtime v1 sandbox")
	restartCmd.SetHelpTemplate(restartHelpMessage())
}
