package contorl

import (
	"encoding/hex"
	"fmt"
	"gateway/device"
	"gateway/httpServer/middleware"
	"gateway/httpServer/model"
	"gateway/setting"
	"gateway/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ApiSystemReboot(context *gin.Context) {

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "接收重启系统命令成功",
		Data:    "",
	})

	time.AfterFunc(3*time.Second, func() {
		setting.SystemReboot()
	})
}

func ApiSystemRebootService(context *gin.Context) {

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "接收重启服务命令成功",
		Data:    "",
	})

	time.AfterFunc(3*time.Second, func() {
		os.Exit(1)
	})
}

func GetDeviceOnline() {

	//更新设备在线率
	deviceTotalCnt := 0
	deviceOnlineCnt := 0
	for _, v := range device.CollectInterfaceMap.Coll {
		deviceTotalCnt += v.DeviceNodeCnt
		deviceOnlineCnt += v.DeviceNodeOnlineCnt
	}
	if (deviceOnlineCnt == 0) || (deviceTotalCnt == 0) {
		setting.SystemState.DeviceOnline = "0"
	} else {
		setting.SystemState.DeviceOnline = fmt.Sprintf("%2.1f", float32(deviceOnlineCnt*100.0/deviceTotalCnt))
	}
}

func GetDevicePacketLoss() {

	//更新设备丢包率
	deviceCommTotalCnt := 0
	deviceCommLossCnt := 0
	for _, v := range device.CollectInterfaceMap.Coll {
		for _, d := range v.DeviceNodeMap {
			deviceCommTotalCnt += d.CommTotalCnt
			deviceCommLossCnt += d.CommTotalCnt - d.CommSuccessCnt
		}
	}
	if (deviceCommLossCnt == 0) || (deviceCommTotalCnt == 0) {
		setting.SystemState.DevicePacketLoss = "0"
	} else {
		setting.SystemState.DevicePacketLoss = fmt.Sprintf("%2.1f", float32(deviceCommLossCnt*100.0/deviceCommTotalCnt))
	}
}

func ApiGetSystemStatus(context *gin.Context) {

	setting.GetCPUState(100 * time.Millisecond)
	setting.GetMemState()
	setting.GetDiskState()
	setting.GetRunTime()

	GetDeviceOnline()
	GetDevicePacketLoss()
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    setting.SystemState,
	})
}

func ApiGetPermissions(context *gin.Context) {
	// 通过http header中的token解析来认证
	token := context.Request.Header.Get("token")
	if token == "" {
		context.JSON(http.StatusOK, gin.H{
			"code":    "1",
			"message": "请求未携带token，无权限访问",
			"data":    "",
		})
		context.Abort()
		return
	}
	middleware.LoginResult.Token = token

	cValue, exists := context.Get("claims")
	if exists == false {
		context.JSON(http.StatusOK, gin.H{
			"code":    "1",
			"message": "用户名获取失败",
			"data":    "",
		})
		context.Abort()
		return
	}
	claims := cValue.(*middleware.CustomClaims)
	middleware.LoginResult.Name = claims.Name
	for _, v := range setting.PolicyWeb {
		if v.Role == claims.Name {
			middleware.LoginResult.Permissions = v.Policy
		}
	}
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    middleware.LoginResult,
	})
}

func ApiSystemGetLoginPassword(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		setting.GetAccountParam(),
	})
}

func ApiModifyLoginPassword(context *gin.Context) {

	type RolePwdTemplate struct {
		Role   string `json:"role"`
		OldPwd string `json:"oldPwd"`
		NewPwd string `json:"newPwd"`
	}

	account := RolePwdTemplate{}

	err := context.ShouldBindJSON(&account)
	if err != nil {
		setting.ZAPS.Errorf("修改登录账号json格式化错误 %v", err)

		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改登录账号json格式化错误",
			Data:    "",
		})
		return
	}

	err = setting.ModifyAccountParam(account.Role, account.OldPwd, account.NewPwd)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改登录账号成功",
		Data:    "",
	})
}

// 定义登陆逻辑
// model.LoginReq中定义了登陆的请求体(name,passwd)
func ApiLogin(c *gin.Context) {
	var loginReq middleware.LoginReq

	if c.BindJSON(&loginReq) == nil {
		// 登陆逻辑校验(查库，验证用户是否存在以及登陆信息是否正确)
		isPass, user, err := middleware.LoginCheck(loginReq)
		// 验证通过后为该次请求生成token
		if isPass {
			middleware.GenerateToken(c, user)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    "1",
				"message": "验证失败" + err.Error(),
				"data":    "",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    "1",
			"message": "用户数据解析失败",
			"data":    "",
		})
		return
	}
}

func ApiSystemCPUUseList(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		setting.CPUDataStream.DataPoint,
	})
}

func ApiSystemMemoryUseList(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		setting.MemoryDataStream.DataPoint,
	})
}

func ApiSystemDiskUseList(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		setting.DiskDataStream.DataPoint,
	})
}

func ApiSystemDeviceOnlineList(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		setting.DeviceOnlineDataStream.DataPoint,
	})
}

func ApiSystemDevicePacketLossList(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		"0",
		"",
		setting.DevicePacketLossDataStream.DataPoint,
	})
}

func ApiSystemSetSystemRTC(context *gin.Context) {
	rRTC := &struct {
		SystemRTC string `json:"systemRTC"`
	}{}
	err := context.ShouldBindJSON(rRTC)
	if err != nil {
		fmt.Println("rRTC json unMarshall err,", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "json unMarshall err",
			Data:    "",
		})
		return
	}
	setting.SystemSetRTC(rRTC.SystemRTC)
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "",
		Data:    "",
	})
}

func ApiSystemSetSN(context *gin.Context) {
	rSN := struct {
		SN   string `json:"sn"`
		Name string `json:"name"`
	}{}

	err := context.BindJSON(&rSN)
	if err != nil {
		setting.ZAPS.Errorf("设置SN参数JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "设置SN参数JSON格式化错误",
			Data:    "",
		})
		return
	}

	setting.SetProduct(rSN.SN, rSN.Name)

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "设置SN成功",
		Data:    "",
	})
}

func ApiSystemGetSN(context *gin.Context) {

	sn := struct {
		SN   string `json:"sn"`
		Name string `json:"name"`
	}{
		SN:   setting.GetSN(),
		Name: setting.GetName(),
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取SN成功",
		Data:    sn,
	})
}

func ApiSystemSetNTPParam(context *gin.Context) {

	rNTPHostAddr := setting.NTPHostAddrTemplate{}

	err := context.ShouldBindJSON(&rNTPHostAddr)
	if err != nil {
		setting.ZAPS.Errorf("设置NTP参数失败,%v", err)

		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "设置NTP参数失败:" + err.Error(),
			Data:    "",
		})
		return
	}

	setting.NTPHostAddr = rNTPHostAddr
	setting.WriteNTPHostAddrToJson()
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "设置NTP参数成功",
		Data:    "",
	})
}

func ApiSystemSetNTPCmd(context *gin.Context) {
	err := setting.NTPGetTime()
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "校时成功",
		Data:    "",
	})
}

func ApiSystemGetNTPParam(context *gin.Context) {
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取NTP参数成功",
		Data:    setting.NTPHostAddr,
	})
}

func ApiBackupFiles(context *gin.Context) {

	backup := &struct {
		Names []string `json:"names"`
	}{}
	err := context.ShouldBindJSON(backup)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "JSON格式化错误",
			Data:    "",
		})
		return
	}

	err, name := setting.BackupFiles(backup.Names)
	if err == nil {
		//返回文件流
		context.Writer.Header().Add("Content-Disposition",
			fmt.Sprintf("attachment;filename=%s", filepath.Base(name)))
		context.File(name) //返回文件路径，自动调用http.ServeFile方法

	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
	}
}

func ApiBackupFilesToRemote(context *gin.Context) {
	err, _ := setting.BackupFilesToRemote()
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "备份到远端服务器错误," + err.Error(),
			Data:    "",
		})

	} else {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "0",
			Message: "备份到远端服务器成功",
			Data:    "",
		})
	}
}

func ApiRecoverFiles(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取文件错误," + err.Error(),
			Data:    "",
		})
		return
	}

	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + file.Filename

	//保存文件到服务器本地
	err = context.SaveUploadedFile(file, fileName)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存文件错误," + err.Error(),
			Data:    "",
		})
		return
	}

	//恢复
	err = setting.RecoverFiles(file.Filename)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: err.Error(),
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "上传恢复文件成功",
		Data:    "",
	})
}

func ApiRecoverFileInfoFromRemote(context *gin.Context) {

	sn := context.Query("sn")

	err, file := setting.GetRecoverFileListFromRemote(sn)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取备份文件列表失败",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取备份文件列表成功",
		Data:    file,
	})

}

func ApiRecoverFileFromRemote(context *gin.Context) {

	backupIdString := context.Query("backupId")
	sn := context.Query("sn")

	id, _ := strconv.Atoi(backupIdString)
	err := setting.GetRecoverFileFromRemote(sn, id)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取备份文件失败",
			Data:    "",
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取备份文件成功",
		Data:    "",
	})

}

func ApiSystemUpdate(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "",
			Data:    "",
		})
		return
	}

	fileName := "./" + file.Filename

	//保存文件到服务器本地
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "Save File Error",
			Data:    "",
		})
		return
	}

	//升级文件解析
	//status := setting.Update(file.Filename)
	//if status == true {
	//	context.JSON(http.StatusOK, model.ResponseData{
	//		Code:    "0",
	//		Message: "",
	//		Data:    "",
	//	})
	//	setting.SystemReboot()
	//}

	setting.SystemReboot()
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "上传升级文件成功",
		Data:    "",
	})

}

func ApiSystemSendPing(context *gin.Context) {
	rPing := &struct {
		IP string `json:"ip"`
	}{}
	err := context.ShouldBindJSON(rPing)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "JSON格式化错误",
			Data:    "",
		})
		return
	}

	err, str := setting.SendPing(rPing.IP)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "发送ping命令失败",
			Data:    str,
		})
		return
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "发送ping命令成功",
		Data:    str,
	})
}

func ApiSendDirectDataToCollInterface(context *gin.Context) {

	type AckDataTemplate struct {
		Date string `json:"date"`
		Data string `json:"data"`
		Type int    `json:"type"`
	}

	data := struct {
		CollInterfaceName string `json:"collInterfaceName"`
		DirectData        string `json:"directData"`
		CheckSum          int    `json:"checkSum"`
	}{}

	ackDatas := make([]AckDataTemplate, 0)
	err := context.BindJSON(&data)
	if err != nil {
		setting.ZAPS.Errorf("通信调试助手报文JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "通信调试助手报文JSON格式化错误",
			Data:    ackDatas,
		})
		return
	}

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[data.CollInterfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "通信调试助手 采集接口不存在",
			Data:    ackDatas,
		})
		return
	}

	//去掉字符串中的空格
	data.DirectData = strings.ReplaceAll(data.DirectData, " ", "")
	if len(data.DirectData)%2 != 0 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "16进制数据格式不正确，需要偶数个数字",
			Data:    ackDatas,
		})
		return
	}

	reqData, _ := hex.DecodeString(data.DirectData)
	if data.CheckSum == 1 {
		crc16 := setting.CRC16(reqData)
		reqData = append(reqData, byte(crc16%256))
		reqData = append(reqData, byte(crc16/256))
	}
	setting.ZAPS.Debugf("reqData %X", reqData)
	ackData := AckDataTemplate{
		Date: time.Now().Format("2006-01-02 15:04:05.999"),
		Data: strings.ToUpper(fmt.Sprintf("%X ", reqData)),
		Type: 1,
	}
	ackDatas = append(ackDatas, ackData)
	req := device.CommunicationDirectDataReqTemplate{
		CollInterfaceName: data.CollInterfaceName,
		Data:              reqData,
	}
	ack := coll.CommQueueManage.CommunicationManageAddDirectData(req)
	ackData.Date = time.Now().Format("2006-01-02 15:04:05.999")
	ackData.Type = 0
	if ack.Status == true {
		ackData.Data = strings.ToUpper(fmt.Sprintf("%X ", ack.Data))
	} else {
		ackData.Data = ""
	}
	ackDatas = append(ackDatas, ackData)

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "通信成功",
		Data:    ackDatas,
	})
}

func ApiSystemExportAuth(context *gin.Context) {

	//exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	name := setting.GenerateAuthFile()
	if name == "" {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "生成授权文件失败",
			Data:    "",
		})
		return
	}

	fileName := name
	//setting.ZAPS.Debugf("fileName %s", fileName)

	defer os.Remove(fileName)
	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", filepath.Base(fileName)))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

}

func ApiSystemImportAuth(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取授权文件失败",
			Data:    "",
		})
		return
	}

	fileName := "./" + file.Filename

	//保存文件到服务器本地
	err = utils.FileCreate("/" + file.Filename)
	if err != nil {
		setting.ZAPS.Errorf("创建文件错误 %v", err)
	}
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存授权文件失败",
			Data:    "",
		})
		return
	}

	//defer os.Remove(fileName)

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "导入授权文件成功",
		Data:    "",
	})
}
