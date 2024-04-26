package lib

import (
	"errors"
	funcLib "go-gateway/lib/func"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var confPath string
var panelType string

var cmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run a forest gateway application",
	Long:  `Run a forest gateway application by parameter`,
	Args: func(cmd *cobra.Command, args []string) error {
		panelType = cmd.Flag("panel_type").Value.String()
		confPath = cmd.Flag("conf_path").Value.String()
		if ok, _ := funcLib.PathExists(confPath); !ok {
			return errors.New("conf_path is not a real dir")
		}
		if !funcLib.IsInArrayString(panelType, []string{"proxy", "control", "both"}) {
			return errors.New("panel_type error，choose one from both/proxy/control")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func CmdExecute() error {
	cmdRun.Flags().StringVarP(&panelType, "panel_type", "p", "", "set panel type like 'both/proxy/control'")
	cmdRun.Flags().StringVarP(&confPath, "conf_path", "c", "", "set configuration path like './conf/dev/'")
	cmdRun.MarkFlagRequired("panel_type")
	cmdRun.MarkFlagRequired("conf_path")
	
	var rootCmd = &cobra.Command{
		Use:               "",
		Short:             "Gatekeeper command manager",
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true, DisableNoDescFlag: true, DisableDescriptions: true},
	}
	rootCmd.AddCommand(cmdRun)
	gin.SetMode(gin.ReleaseMode)
	return rootCmd.Execute()
}

func GetCmdConfPath() string {
	return confPath
}

func GetCmdPanelType() string {
	return panelType
}

func SetCmdConfPath(path string) {
	confPath = path
}

func SetCmdPanelType(paneltype string) {
	panelType = paneltype
}
