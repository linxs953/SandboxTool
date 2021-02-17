package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"runtimeops/utils"
	"strings"
)

var recyclePuid bool
var recycleEndpoint string
var recycleToken string
var recycleCmd = &cobra.Command{
	Use:   "recycle",
	Short: "Recycle Sandbox",
	Long:  "Recycle Sandbox",
	Run: func(cmd *cobra.Command, args []string) {
		var success bool
		epFlagV, err := utils.GetFlagValueStringP(cmd, "ep")
		if err != nil {
			return
		}
		tokenFlagV, err := utils.GetFlagValueStringP(cmd, "token")
		if err != nil {
			return
		}
		puidFlagV, err := utils.GetFlagValueBool(cmd, "puid")
		if err != nil {
			return
		}
		if tokenFlagV == nil || epFlagV == nil {
			return
		}

		if *tokenFlagV != "" && *epFlagV != "" {
			success = utils.RecycleSandbox(utils.Sandbox{}, *epFlagV, *tokenFlagV)
			if success {
				parts := strings.Split(*epFlagV, "/")
				if len(parts) < 5 {
					return
				}
				log.Printf("Recycle sandbox [%s/%s] successfully", parts[3], parts[5])
			}
		} else {
			if (*tokenFlagV == "" || *epFlagV == "") && !puidFlagV {
				log.Print("token and endpoint need appear at the same time")
				return
			}
			if len(args) < 1 {
				fmt.Println("Use recycle -h for help")
				return
			}
			switch puidFlagV {
			case true:
				token := utils.GenerateToken(args[0], "")
				if token == "" {
					return
				}
				sb := utils.GetV2Sandbox(token, false)
				success = utils.RecycleSandbox(sb)
				if success {
					ep := sb.Data.Endpoint
					parts := strings.Split(ep, "/")
					if len(parts) < 5 {
						return
					}
					p, r := parts[3], parts[4]
					log.Printf("Recycle sandbox [%s/%s] Successfully", p, r)
					return
				}
			case false:
				log.Print("lack of flag --puid")
				return
			}
		}
	},
}

func recycleHelpMessage() string {
	helpMsg := "e.g : \n" +
		"  rtctl recycle [puid]  --puid=true\n" +
		"  rtctl recycle --ep='sandbox_endpoint' --token='sandbox_token'\n\n"
	helpMsg = fmt.Sprintf("%s\n\n%s\n%s\n", recycleCmd.Long, recycleCmd.UsageString(), helpMsg)
	return helpMsg
}

func init() {
	rootCmd.AddCommand(recycleCmd)
	recycleCmd.Flags().BoolVarP(&recyclePuid, "puid", "p", false, "use `--puid=true` to recycle sandbox by puid")
	recycleCmd.Flags().StringVarP(&recycleEndpoint, "ep", "", "", "Use --ep to set sandbox which will be recycled endpoint")
	recycleCmd.Flags().StringVarP(&recycleToken, "token", "t", "", "Use --token to set sandbox which will be recycled token")
	recycleCmd.SetHelpTemplate(recycleHelpMessage())
}
