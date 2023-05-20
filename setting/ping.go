package setting

import (
	"bytes"
	"github.com/go-ping/ping"
	"os/exec"
)

type PINGCmdParamTemplate struct {
	DstIP string `json:"DstIP"` //目标服务器IP
	Count int    `json:"Count"` //ping次数
}

func PINGCmdInit() {
	pingCmd, err := ping.NewPinger("www.baidu.com")
	if err != nil {
		ZAPS.Errorf("PingCmd输入参数错误 %v", err)
		return
	}
	pingCmd.Count = 3
	pingCmd.Source = "192.168.1.24"

	err = pingCmd.Run()
	if err != nil {
		ZAPS.Errorf("PingCmd执行错误 %v", err)
	}
	stats := pingCmd.Statistics()
	ZAPS.Infof("PingCmd执行结果 %v", stats)
}

func SendPing(ip string) (error, string) {

	//cmd := exec.Command("ping", ip, "-c", "4", "-W", "5")
	cmd := exec.Command("ping", ip, "-c", "4")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		ZAPS.Debugf("ping发送失败 %v", err)
		return err, ""
	} else {
		ZAPS.Debug("ping发送成功")
	}
	ZAPS.Debugln(out.String())
	return nil, out.String()
}
