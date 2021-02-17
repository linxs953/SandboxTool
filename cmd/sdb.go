package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	"os"
	"runtimeops/utils"
)

var wss string
var endpoint string
var token string
var origin string
var keepAlive bool

const helpMsgWithoutH = `sdb (Sandbox Debug Bridge) allows you to interactive with specific sandbox by websocket url.
Use sdb -h for help`

var sdbCmd = &cobra.Command{
	Use:   "sdb",
	Short: "Sandbox Debug Bridge",
	Long: `Sandbox Debug Bridge`,
	Run: func(cmd *cobra.Command, args []string) {
		para, _ := utils.ParseCliParameter(os.Args)
		if len(para) == 0 {
			fmt.Println(helpMsgWithoutH)
			return
		}
		utils.Connect(para, cmd)
	},
}


func sdbHelpMessage() string {
	helpMsg := "e.g. \n" +
								"  1. sdb --url='wss_url' --alive=true\n" +
								"  2. sdb --ep='endpoint_url' -t='sandbox_token'\n"
	helpMsg = fmt.Sprintf("%s\n\n%s\n%s\n",sdbCmd.Long,sdbCmd.UsageString(),helpMsg)
return helpMsg
}


func init() {
	rootCmd.AddCommand(sdbCmd)
	sdbCmd.Flags().StringVarP(&wss, "url", "u", "", "websocket url")
	sdbCmd.Flags().StringVarP(&endpoint, "ep", "", "", "sandbox endpoint")
	sdbCmd.Flags().StringVarP(&token, "token", "t", "", "sandbox token")
	sdbCmd.Flags().StringVarP(&origin, "origin", "s", "", "websocket origin")
	sdbCmd.Flags().BoolVarP(&keepAlive,"alive","a",false,"connection keep-alive")
	sdbCmd.SetHelpTemplate(sdbHelpMessage())
}
