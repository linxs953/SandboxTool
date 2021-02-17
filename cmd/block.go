package cmd

//var blockCmd = &cobra.Command{
//	Use:   "block",
//	Short: "Block Sandbox User ",
//	Long:  "Block Sandbox User",
//	Run: func(cmd *cobra.Command, args []string) {
//		if len(args) < 2 {
//			fmt.Println("lack of args .\nUse block -h for help")
//			return
//		}
//
//		//if err := config.Init(); err != nil {
//		//	log.Error().Err(err).Msg("Init Config error")
//		//	return
//		//}
//		parmStr, err := cmd.Flags().GetString("puid")
//		if err != nil {
//			return
//		}
//		switch parmStr {
//		case "true":
//			utils.BlockByPuid(args)
//		default:
//			utils.BlockByUnionID(args)
//		}
//	},
//}
//
//func blockHelpMessage() string {
//	helpMsg :=  "e.g: \n" +
//		"  runtime block [unionid] [resource_group]\n" +
//		"  runtime block [puid] [resource_group] --puid=true"
//	helpMsg = fmt.Sprintf("%s\n%s\n%s",blockCmd.Long,blockCmd.UsageString(),helpMsg)
//	return helpMsg
//}
//
//func init() {
//	rootCmd.AddCommand(blockCmd)
//	blockCmd.Flags().StringVarP(&name, "puid", "p", "", "use `--puid=true` to get block sandbox by puid")
//	blockCmd.SetHelpTemplate(blockHelpMessage())
//}
