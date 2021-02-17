package cmd

//var sensorCmd = &cobra.Command{
//	Use:   "sensor",
//	Short: "sensor debug",
//	Long: `sensor debug`,
//	Args: cobra.MinimumNArgs(2),
//	Run: func(cmd *cobra.Command, args []string) {
//		if err := config.Init();err != nil {
//			return
//		}
//		urls := utils.SensorCmdProcess(args)
//		fmt.Println(urls)
//		for _,url := range urls {
//			utils.Open(url)
//		}
//	},
//}
//
//
//
//
//func init() {
//	rootCmd.AddCommand(sensorCmd)
//}
