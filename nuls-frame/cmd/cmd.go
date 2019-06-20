package cmd

import (
	"github.com/nuls-io/go-nuls/nuls-frame/types"
)

var cmds = make([]Cmd, 0)

type Cmd interface {
	GetCmdInfo() *CmdInfo
	Do(params *map[string]interface{}) (*types.Response, error)
}

type CmdInfo struct {
	Name		string
	Version		float32
}

type BaseCmd struct {
	CmdInfo *CmdInfo
}

func (cmd *BaseCmd) GetCmdInfo() *CmdInfo {
	return cmd.CmdInfo
}

func (cmd *BaseCmd) Do(params *map[string]interface{}) (*types.Response, error) {
	response := &types.Response{
		ResponseComment: 				"This is an unimplemented cmd",
	}
	return response, nil
}

// 获取对应的Cmd，当一个接口对应多个版本时，自动返回最新的版本
func GetCmd(cmdName string) *Cmd {
	var cmd *Cmd
	for _, c := range cmds {
		if c.GetCmdInfo() == nil || c.GetCmdInfo().Name != cmdName {
			continue
		}
		if cmd == nil ||  c.GetCmdInfo().Version > (*cmd).GetCmdInfo().Version {
			cmd = &c
		}
	}
	return cmd
}

// 根据名称和最小版本号获取Cmd，返回的Cmd版本需要和指定的最小版本的大版本号一致
func GetCmdByVersion(cmdName string, minVersion float32) *Cmd {
	var cmd *Cmd
	for _, c := range cmds {
		if c.GetCmdInfo() == nil || c.GetCmdInfo().Name != cmdName {
			continue
		}
		// 大版本不一样，跳过
		// Big version is different, skip
		if int(minVersion) != int(c.GetCmdInfo().Version) {
			continue
		}
		if cmd == nil {
			cmd = &c
			continue
		}
		if c.GetCmdInfo().Version > (*cmd).GetCmdInfo().Version {
			cmd = &c
		}
	}
	return cmd
}

func AddCmd(cmd Cmd) {
	cmds = append(cmds, cmd)
}