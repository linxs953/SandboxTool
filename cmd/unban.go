package cmd

//import (
//	"fmt"
//	"github.com/spf13/cobra"
//	"runtimeops/utils"
//)
//
//var unbanCmd = &cobra.Command{
//	Use:   "unban",
//	Short: "Unban User Sandbox",
//	Long:  "Unban User Sandbox",
//	Run: func(cmd *cobra.Command, args []string) {
//		if len(args) < 2 {
//			fmt.Println("lack of args .\nUse unban -h for help")
//			return
//		}
//		//if err := config.Init(); err != nil {
//		//	log.Error().Err(err).Msg("Init config error")
//		//	return
//		//}
//		parmStr, err := cmd.Flags().GetString("puid")
//		if err != nil {
//			return
//		}
//		switch parmStr {
//		case "true":
//			utils.UnbanByPuid(args)
//		default:
//			utils.UnbanByUnionid(args)
//		}
//	},
//}
//
//func unbanHelpMessage() string {
//	helpMsg := "e.g : \n " +
//		"  runtime unban [unionid] [resource_group]\n"
//	helpMsg = fmt.Sprintf("%s\n\n%s\n%s\n",unbanCmd.Long,unbanCmd.UsageString(),helpMsg)
//	return helpMsg
//}
//
//func init() {
//	rootCmd.AddCommand(unbanCmd)
//	unbanCmd.Flags().StringVarP(&name, "puid", "p", "", "use `--puid=true` to get sandbox by puid")
//	unbanCmd.SetHelpTemplate(unbanHelpMessage())
//}
