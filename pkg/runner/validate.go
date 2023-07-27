package runner

import (
	"fmt"
	"github.com/iahfdoa/m3u8Find/pkg/utils"
	"github.com/projectdiscovery/fileutil"
)

const (
	errorDefaultNotInArray = `不在默认参数中：%v -> %v`
)

func (userOptions *UserOptions) validateOptions() error {

	// 保证参数合规
	if userOptions.Limit <= 0 {
		return fmt.Errorf("limit 不能小于等于0 %v (default %v)", userOptions.Limit, Limit)
	}

	if userOptions.Timeout <= 0 {
		return fmt.Errorf("连接延迟时间不能小于等于0 %v (default %v)", userOptions.Timeout, Timeout)
	}
	if userOptions.ReconnectionNum <= 0 {
		return fmt.Errorf("重连次数不能小于等于0 %v (default %v)", userOptions.ReconnectionNum, ReconnectionNum)
	}

	// 检查output是否为文件
	if !fileutil.FolderExists(userOptions.Output) {
		if fileutil.FileExists(userOptions.Output) {
			return fmt.Errorf("output %s 文件已经存在", userOptions.Output)
		}
	}
	// 设置 不在指定数组
	if userOptions.PrimaryTag != "" && !utils.StrInArray(userOptions.PrimaryTag, PrimaryTags) {
		return fmt.Errorf(errorDefaultNotInArray, userOptions.PrimaryTag, PrimaryTags)
	}

	if userOptions.GroupRole != "" && !utils.StrInArray(userOptions.GroupRole, GroupRoles) {
		return fmt.Errorf(errorDefaultNotInArray, userOptions.GroupRole, GroupRoles)
	}
	return nil

}
