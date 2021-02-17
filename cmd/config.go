package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtimeops/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set config env",
	Long:  "Set config env.",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) == 0 {
			fmt.Println("lack of args. \nUse `rtctl config list` for help")
			return
		}
		switch args[0] {
		case "list":
			fmt.Println(listHelpMessage())
		case "set":
			if len(args) < 3 {
				fmt.Println("lack of args. \n Use `rtctl config set [key] [value]`")
				return
			}
			//if err = config.Init(); err != nil {
			//	return
			//}
			keys := []string{args[1]}
			values := []string{args[2]}
			if err = config.SetEnv(keys, values); err != nil {
				return
			}
		default:
			return
		}
	},
}

func listHelpMessage() string {
	listHelpMsg := "You can use follow command to set config item\n " +
		"(required)rtctl config set BACKENDCICD string\n " +
		"(required)rtctl config set SECRETKEY string\n" +
		"Another config item\n" +
		"TOKENBYUNIONIDAPI ------- string (通过unionid获取token的API)\n" +
		"TOKENBYOPENAPI    ------- string (通过openid获取token的API)\n" +
		"PROD              ------- string (生产环境URL)\n" +
		"TEST              ------- string (测试环境URL)\n" +
		"RECYCLEAPI        ------- string (回收沙盒API)\n" +
		"SANDBOXGROUP      ------- []map[string]string (沙盒类型分组，包含sandboxid， sandboxtype(formal),protect)\n" +
		"BLOCKAPI          ------- string (封禁沙盒API)\n" +
		"ENTRTAPI          ------- string (沙盒entry API)\n" +
		"Or: You can create .env file (include all config item above) under  $HOME/ and add source $HOME/.env in $HOME/.zshrc"
	return listHelpMsg
}

func configHelpMessage() string {
	helpMsg := "e.g. \n  `rtctl config list` for supported config item\n" +
		"  `rtctl config set [key] [value]` to set config env\n"
	helpMsg = fmt.Sprintf("%s\n\n%s\n%s", configCmd.Long, configCmd.UsageString(), helpMsg)
	return helpMsg
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.SetHelpTemplate(configHelpMessage())
}
