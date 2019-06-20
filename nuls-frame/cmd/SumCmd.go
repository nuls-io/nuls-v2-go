package cmd

import (
	"github.com/nuls-io/go-nuls/nuls-frame/types"
)

func init()  {
	cmd := &SumCmd{
		&CmdInfo {
			Name : 		"sum",
			Version:	1.0,
		},
	}
	AddCmd(cmd)
}

type SumCmd BaseCmd

func (cmd *SumCmd) GetCmdInfo() *CmdInfo {
	return cmd.CmdInfo
}

func (cmd *SumCmd) Do(params *map[string]interface{}) (*types.Response, error) {
	//log.Println(params)
	replay := make(map[string]interface{})
	replay["msg"] = "this is Server replay of message [XXX]"
	response := types.NewResponse("", "1", "test sum", "0", &replay)

	return response, nil
}