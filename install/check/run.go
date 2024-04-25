package check

import (
	"fmt"
	"go-gateway/install/tool"
)

var (
	ForestGatewayPath	string = tool.ForestGatewayPath
	CmdRun			string = "cd %s && %s run main.go run -c %s/conf/dev/ -p control"
)

func RunGateway() error{
	boolRunGateway, err := tool.Confirm("run gatekeeper?", 2)
	if err != nil{
		return err
	}
	CmdRun := fmt.Sprintf(tool.CmdRun + "&&" + CmdRun, ForestGatewayPath, GoPath, ForestGatewayPath)
	if boolRunGateway {
		tool.LogInfo.Println(CmdRun)
		err := tool.RunCmd(CmdRun)
		if err != nil{
			return err
		}
	}
	return nil
}




