package contorl

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gateway/controllers"
	"gateway/device"
	"gateway/httpServer/model"
	"gateway/setting"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func ApiAddNode(context *gin.Context) {

	nodeInfo := struct {
		InterfaceName string `json:"interfaceName"`
		Name          string `json:"name"`
		Label         string `json:"label"`
		Addr          string `json:"addr"`
		TSL           string `json:"tsl"`
		DeviceType    string `json:"deviceType"`
	}{}

	emController := controllers.NewEMController()
	emController.AddEmDevice(context)

	err := context.ShouldBindBodyWith(&nodeInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加设备JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加设备JSON格式化错误",
			Data:    "",
		})
		return
	}

	nodeCnt := 0
	for _, v := range device.CollectInterfaceMap.Coll {
		nodeCnt += v.DeviceNodeCnt
	}

	setting.ZAPS.Debugf("设备数量 %v", nodeCnt)

	for _, v := range device.CollectInterfaceMap.Coll {
		for _, d := range v.DeviceNodeMap {
			if nodeInfo.Name == d.Name {
				context.JSON(http.StatusOK, model.ResponseData{
					Code:    "1",
					Message: "设备名称已经存在",
					Data:    "",
				})
				return
			}
		}
	}

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[nodeInfo.InterfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口名称不存在",
			Data:    "",
		})
		return
	}

	_ = coll.AddDeviceNode(nodeInfo.Name, nodeInfo.TSL, nodeInfo.Addr, nodeInfo.Label, nodeInfo.DeviceType)
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "增加设备成功",
		Data:    "",
	})
	return

}

func ApiAddNodesFromXlsx(context *gin.Context) {

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("待导入xlsx文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "待导入xlsx文件不存在",
			Data:    "",
		})
		return
	}

	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + file.Filename

	//保存文件到服务器本地
	err = utils.FileCreate(fileName)
	if err != nil {
		setting.ZAPS.Errorf("创建xlsx文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "创建xlsx文件错误",
			Data:    "",
		})
		return
	}
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		setting.ZAPS.Errorf("保存CSV文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
			Data:    "",
		})
		return
	}

	defer os.Remove(fileName)

	err, cells := setting.ReadExcel(fileName) //标题在第2行，从第3行取数据，第2列取数据
	if err != nil {
		setting.ZAPS.Errorf("加载xlsx文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "加载xlsx文件错误",
			Data:    "",
		})
		return
	}

	setting.ZAPS.Debugf("cells %v", cells)
	for _, cell := range cells {
		if len(cell) < 6 {
			continue
		}
		nodeInfo := struct {
			CollInterface string `json:"collInterface"`
			Name          string `json:"name"`
			Label         string `json:"label"`
			Addr          string `json:"addr"`
			TSL           string `json:"tsl"`
			DeviceType    string `json:"deviceType"`
		}{
			CollInterface: cell[0],
			Name:          cell[1],
			Label:         cell[2],
			Addr:          cell[3],
			TSL:           cell[4],
			DeviceType:    cell[5],
		}

		// 获取采集接口名称
		device.CollectInterfaceMap.Lock.Lock()
		coll, ok := device.CollectInterfaceMap.Coll[nodeInfo.CollInterface]
		device.CollectInterfaceMap.Lock.Unlock()
		if !ok {
			continue
		}

		//判断设备名称是否已经存在
		_, ok = coll.DeviceNodeMap[nodeInfo.Name]
		if ok {
			continue
		}

		_ = coll.AddDeviceNode(nodeInfo.Name, nodeInfo.TSL, nodeInfo.Addr, nodeInfo.Label, nodeInfo.DeviceType)
		// 导入cmd写入sqlite
		emController := controllers.NewEMController()
		emController.AddEmDeviceFromXlsx(nodeInfo.Name, nodeInfo.TSL, nodeInfo.Addr, nodeInfo.Label, nodeInfo.CollInterface)

	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "从xlsx中批量添加设备成功",
		Data:    "",
	})
	return

}

func ApiAddNodesFromCSV(context *gin.Context) {

	// 获取采集接口名称
	collName := context.PostForm("name")
	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[collName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		setting.ZAPS.Errorf("采集接口名称不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口名称不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("待导入CSV文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "待导入CSV文件不存在",
			Data:    "",
		})
		return
	}

	fileName := "./tmp/" + file.Filename

	utils.DirIsExist("./tmp")
	//保存文件到服务器本地
	err = utils.FileCreate(fileName)
	if err != nil {
		setting.ZAPS.Errorf("创建CSV文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "创建CSV文件错误",
			Data:    "",
		})
		return
	}
	if err := context.SaveUploadedFile(file, fileName); err != nil {
		setting.ZAPS.Errorf("保存CSV文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存CSV文件错误",
			Data:    "",
		})
		return
	}

	defer os.Remove(fileName)

	result := setting.LoadCsvCfg(fileName, 1, 2, 0) //标题在第2行，从第3行取数据，第2列取数据
	if result == nil {
		setting.ZAPS.Errorf("加载CSV文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "加载CSV文件错误",
			Data:    "",
		})
		return
	}

	setting.ZAPS.Debugf("records %v", result.Records)
	for _, record := range result.Records {
		nodeInfo := struct {
			Name       string `json:"name"`
			Label      string `json:"label"`
			Addr       string `json:"addr"`
			TSL        string `json:"tsl"`
			DeviceType string `json:"deviceType"`
		}{
			Name:  record.GetString("name"),
			Label: record.GetString("label"),
			Addr:  record.GetString("addr"),
			TSL:   record.GetString("tsl"),
		}

		//判断设备名称是否已经存在
		_, ok := coll.DeviceNodeMap[nodeInfo.Name]
		if ok {
			continue
		}

		_ = coll.AddDeviceNode(nodeInfo.Name, nodeInfo.TSL, nodeInfo.Addr, nodeInfo.Label, nodeInfo.DeviceType)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "物模型模版导入CSV成功",
		Data:    "",
	})
	return

}

func ApiExportNodesToCSV(context *gin.Context) {

	collName := context.Query("name")

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[collName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		setting.ZAPS.Errorf("采集接口名称不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口名称不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileName := exeCurDir + "/selfpara/" + collName + ".xlsx"

	excelRecords := [][]string{
		{"设备名称", "设备标签", "通信地址", "采集模型", "设备类型"},
		{"name", "label", "addr", "tsl", "deviceType"},
	}
	for _, v := range coll.DeviceNodeMap {
		record := make([]string, 0)
		record = append(record, v.Name)
		record = append(record, v.Label)
		record = append(record, v.Addr)
		record = append(record, v.TSL)
		record = append(record, v.DeviceType)
		excelRecords = append(excelRecords, record)
	}

	fs, err := os.Create(fileName)
	if err != nil {
		setting.ZAPS.Errorf("creat tsl.csv err,%v", err)
		return
	}

	defer os.Remove(fileName)
	defer fs.Close()
	// 写入UTF-8 BOM
	//_, err = fs.WriteString("\xEF\xBB\xBF")

	//创建一个新的写入文件流
	csvFile := csv.NewWriter(fs)
	csvRecords := [][]string{
		{"设备名称", "设备标签", "通信地址", "采集模型", "设备类型"},
		{"name", "label", "addr", "tsl", "deviceType"},
	}

	for _, v := range coll.DeviceNodeMap {
		record := make([]string, 0)
		record = append(record, v.Name)
		record = append(record, v.Label)
		record = append(record, v.Addr)
		record = append(record, v.TSL)
		record = append(record, v.DeviceType)
		csvRecords = append(csvRecords, record)
	}

	err = csvFile.WriteAll(csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("保存CSV文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存CSV文件错误",
			Data:    "",
		})
		return
	}
	csvFile.Flush()

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(fileName))))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return
}

func ApiExportNodesToXlsx(context *gin.Context) {

	param := struct {
		Names []string `json:"names"`
	}{
		Names: make([]string, 0),
	}

	err := context.BindJSON(&param)
	if err != nil {
		setting.ZAPS.Errorf("导出采集接口设备到xlsx文件请求参数JSON格式化错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "导出采集接口设备到xlsx文件请求参数JSON格式化错误",
			Data:    "",
		})
		return
	}

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/collInterfaceDevices.xlsx"

	cells := [][]string{
		{"采集接口名称", "设备名称", "设备标签", "通信地址", "采集模型名称", "设备类型"},
		{"collInterface", "name", "label", "addr", "tsl", "deviceType"},
	}

	for _, v := range param.Names {
		device.CollectInterfaceMap.Lock.Lock()
		coll, ok := device.CollectInterfaceMap.Coll[v]
		device.CollectInterfaceMap.Lock.Unlock()
		if !ok {
			continue
		}

		for _, d := range coll.DeviceNodeMap {
			record := make([]string, 0)
			record = append(record, coll.CollInterfaceName)
			record = append(record, d.Name)
			record = append(record, d.Label)
			record = append(record, d.Addr)
			record = append(record, d.TSL)
			record = append(record, d.DeviceType)
			cells = append(cells, record)
		}
	}

	err = setting.WriteExcel(fileName, cells)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "写Excel文件错误",
			Data:    "",
		})
		return
	}

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", url.QueryEscape(filepath.Base(fileName))))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return
}

func ApiModifyNode(context *gin.Context) {

	nodeInfo := struct {
		InterfaceName string `json:"interfaceName"`
		Name          string `json:"name"`
		Label         string `json:"label"`
		Addr          string `json:"addr"`
		TSL           string `json:"tsl"`
		DeviceType    string `json:"deviceType"`
	}{}

	emController := controllers.NewEMController()
	emController.UpdateEmDevice(context)

	err := context.ShouldBindBodyWith(&nodeInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("增加设备JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "增加设备JSON格式化错误",
			Data:    "",
		})
		return
	}

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[nodeInfo.InterfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口名称不存在",
			Data:    "",
		})
		return
	}

	err = coll.ModifyDeviceNode(nodeInfo.Name, nodeInfo.TSL, nodeInfo.Addr, nodeInfo.Label, nodeInfo.DeviceType)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "修改设备失败",
			Data:    "",
		})
		return
	}
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "修改设备成功",
		Data:    "",
	})
}

func ApiModifyNodes(context *gin.Context) {

	type DeleteAck struct {
		Name   string
		Status bool
	}

	aParam := struct {
		Code    string      `json:"Code"`
		Message string      `json:"Message"`
		Data    []DeleteAck `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    make([]DeleteAck, 0),
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)
	log.Println(string(bodyBuf[:n]))

	nodeInfo := &struct {
		InterfaceName string   `json:"CollInterfaceName"`
		DTSL          string   `json:"TSL"`
		Name          []string `json:"Name"`
	}{
		InterfaceName: "",
		DTSL:          "",
		Name:          make([]string, 0),
	}

	err := json.Unmarshal(bodyBuf[:n], nodeInfo)
	if err != nil {
		fmt.Println("nodeInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"

		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for _, v := range nodeInfo.Name {
		ack := DeleteAck{
			Name:   v,
			Status: false,
		}
		aParam.Data = append(aParam.Data, ack)
	}

	for k, v := range nodeInfo.Name {
		device.CollectInterfaceMap.Lock.Lock()
		coll, ok := device.CollectInterfaceMap.Coll[nodeInfo.InterfaceName]
		device.CollectInterfaceMap.Lock.Unlock()
		if !ok {
			continue
		}

		for _, d := range coll.DeviceNodeMap {
			if d.Name == v {
				d.TSL = nodeInfo.DTSL
				device.WriteCollectInterfaceManageToJson()
				aParam.Data[k].Status = true
			}
		}
	}

	aParam.Code = "0"
	aParam.Message = ""
	sJson, _ := json.Marshal(aParam)
	context.String(http.StatusOK, string(sJson))
}

func ApiDeleteNode(context *gin.Context) {

	nodeInfo := struct {
		InterfaceName string   `json:"InterfaceName"`
		DNames        []string `json:"deviceNames"`
	}{}

	emController := controllers.NewEMController()
	emController.DeleteEmDevice(context)

	err := context.ShouldBindBodyWith(&nodeInfo, binding.JSON)
	if err != nil {
		setting.ZAPS.Error("删除设备JSON格式化错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "删除设备JSON格式化错误",
			Data:    "",
		})
		return
	}

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[nodeInfo.InterfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口名称不存在",
			Data:    "",
		})
		return
	}

	coll.DeleteDeviceNodes(nodeInfo.DNames)
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "删除设备成功",
		Data:    "",
	})
}

func ApiGetNode(context *gin.Context) {

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")

	aParam := &struct {
		Code    string
		Message string
		Data    device.DeviceNodeTemplate
	}{}

	coll, ok := device.CollectInterfaceMap.Coll[sName]
	if !ok {
		aParam.Code = "1"
		aParam.Message = "CollInterfaceName is no exist"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	for _, n := range coll.DeviceNodeMap {
		if n.Addr == sAddr {
			aParam.Code = "0"
			aParam.Message = ""
			aParam.Data = *coll.GetDeviceNode(sAddr)
			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))
			return
		}
	}
}

func ApiGetNodes(context *gin.Context) {

	params := struct {
		Names []string `json:"names"`
	}{}

	err := context.BindJSON(&params)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "获取获取指定采集接口JSON格式化错误",
			Data:    "",
		})
		return
	}

	type DeviceNodeTemplate struct {
		Index             int    `json:"index"`             //设备偏移量
		CollInterfaceName string `json:"collInterfaceName"` //设备采集接口名称
		CommInterfaceName string `json:"commInterfaceName"` //设备通信接口名称
		Name              string `json:"name"`              //设备名称
		Label             string `json:"label"`             //设备标签
		Addr              string `json:"addr"`              //设备地址
		TSL               string `json:"tsl"`               //设备物模型
		DeviceType        string `json:"deviceType"`        //设备类型
		LastCommRTC       string `json:"lastCommRTC"`       //最后一次通信时间戳
		CommTotalCnt      int    `json:"commTotalCnt"`      //通信总次数
		CommSuccessCnt    int    `json:"commSuccessCnt"`    //通信成功次数
		CurCommFailCnt    int    `json:"-"`                 //当前通信失败次数
		CommStatus        string `json:"commStatus"`        //通信状态
		CommMaxTime       int    `json:"commMaxTime"`       //通信最长用时
		CommMinTime       int    `json:"commMinTime"`       //通信最短用时
	}

	nodes := struct {
		Node           []DeviceNodeTemplate `json:"node"`
		CommTotalCnt   int                  `json:"commTotalCnt"`
		CommSuccessCnt int                  `json:"commSuccessCnt"`
	}{
		Node:           make([]DeviceNodeTemplate, 0),
		CommTotalCnt:   0,
		CommSuccessCnt: 0,
	}

	node := DeviceNodeTemplate{}
	index := 0
	for _, n := range params.Names {
		for _, v := range device.CollectInterfaceMap.Coll {
			if n != v.CollInterfaceName {
				continue
			}
			for _, d := range v.DeviceNodeMap {
				node.Index = index
				node.CollInterfaceName = v.CollInterfaceName
				node.CommInterfaceName = v.CommInterfaceName
				node.Name = d.Name
				node.Label = d.Label
				node.Addr = d.Addr
				node.TSL = d.TSL
				node.DeviceType = d.DeviceType
				node.LastCommRTC = d.LastCommRTC
				node.CommTotalCnt = d.CommTotalCnt
				node.CommSuccessCnt = d.CommSuccessCnt
				node.CurCommFailCnt = d.CurCommFailCnt
				node.CommStatus = d.CommStatus
				nodes.Node = append(nodes.Node, node)
				index++
			}
			nodes.CommTotalCnt += v.DeviceNodeCnt
			nodes.CommSuccessCnt += v.DeviceNodeOnlineCnt
		}
	}

	//排序，方便前端页面显示
	//sort.Slice(nodes, func(i, j int) bool {
	//	iName := nodes.Node[i].Index
	//	jName := nodes.Node[j].Index
	//	return iName < jName
	//})

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "获取指定采集接口成功",
		Data:    nodes,
	})
}

/*
*
从缓存中获取设备变量
*/
func ApiGetNodeVariableFromCache(context *gin.Context) {

	type VariableTemplate struct {
		Index     int         `json:"index"`     // 变量偏移量
		Name      string      `json:"name"`      // 变量名
		Label     string      `json:"lable"`     // 变量标签
		Value     interface{} `json:"value"`     // 变量值
		Explain   interface{} `json:"explain"`   // 变量值说明
		TimeStamp string      `json:"timestamp"` // 变量时间戳
		Type      string      `json:"type"`      // 变量类型
		Unit      string      `json:"unit"`      // 变量单位
	}

	interfaceName := context.Query("collInterfaceName")
	deviceName := context.Query("deviceName")

	//查找设备是否存在
	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[interfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口不存在",
			Data:    "",
		})
		return
	}

	node, ok := coll.DeviceNodeMap[deviceName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "设备名称接口不存在",
			Data:    "",
		})
		return
	}

	variables := make([]VariableTemplate, 0)
	for _, v := range node.Properties {
		variable := VariableTemplate{}
		variable.Name = v.Name
		variable.Label = v.Label
		variable.Unit = node.GetTSLModelUnit(v.Name)
		if v.Type == device.PropertyTypeUInt32 {
			variable.Type = "uint32"
		} else if v.Type == device.PropertyTypeInt32 {
			variable.Type = "int32"
		} else if v.Type == device.PropertyTypeDouble {
			variable.Type = "double"
		} else if v.Type == device.PropertyTypeString {
			variable.Type = "string"
		}
		// 取出切片中最后一个值
		if len(v.Value) > 0 {
			index := len(v.Value) - 1
			variable.Index = v.Value[index].Index
			variable.Value = v.Value[index].Value
			variable.Explain = v.Value[index].Explain
			variable.TimeStamp = v.Value[index].TimeStamp.Format("2006-01-02 15:04:05.000")
		} else {
			variable.Value = ""
			variable.Explain = ""
			variable.TimeStamp = ""
		}
		variables = append(variables, variable)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "读取设备缓存成功",
		Data:    variables,
	})
}

func ApiGetNodeHistoryVariableFromCache(context *gin.Context) {

	sName := context.Query("CollInterfaceName")
	sAddr := context.Query("Addr")
	sVariable := context.Query("VariableName")

	aParam := &struct {
		Code    string
		Message string
		Data    []device.TSLPropertyValueTemplate
	}{}

	coll, ok := device.CollectInterfaceMap.Coll[sName]
	if !ok {
		aParam.Code = "1"
		aParam.Message = "node is no exist"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))

		return
	}

	for _, d := range coll.DeviceNodeMap {
		if d.Addr == sAddr {
			aParam.Code = "0"
			aParam.Message = ""
			for _, v := range d.Properties {
				if v.Name == sVariable {
					aParam.Data = append(aParam.Data, v.Value...)
				}
			}

			sJson, _ := json.Marshal(aParam)
			context.String(http.StatusOK, string(sJson))

			return
		}
	}
}

func ApiGetNodeRealVariable(context *gin.Context) {

	type VariableTemplate struct {
		Index     int         `json:"index"`     // 变量偏移量
		Name      string      `json:"name"`      // 变量名
		Label     string      `json:"label"`     // 变量标签
		Value     interface{} `json:"value"`     // 变量值
		Explain   interface{} `json:"explain"`   //变量值解释
		TimeStamp string      `json:"timestamp"` // 变量时间戳
		Type      string      `json:"type"`      // 变量类型
		Unit      string      `json:"unit"`      // 变量单位
	}

	interfaceName := context.Query("collInterfaceName")
	deviceName := context.Query("deviceName")

	//查找设备是否存在
	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[interfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "采集接口不存在",
			Data:    "",
		})
		return
	}

	node, ok := coll.DeviceNodeMap[deviceName]
	if !ok {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "设备名称接口不存在",
			Data:    "",
		})
		return
	}

	//发送命令到响应的采集接口
	err, properties := coll.ReadDeviceRealVariable(deviceName)
	if err != nil {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "读取设备实时变量错误",
			Data:    "",
		})
		return
	}

	variables := make([]VariableTemplate, 0)
	for _, v := range properties {
		variable := VariableTemplate{}
		variable.Name = v.Name
		variable.Label = v.Label
		variable.Unit = node.GetTSLModelUnit(v.Name)
		if v.Type == device.PropertyTypeUInt32 {
			variable.Type = "uint32"
		} else if v.Type == device.PropertyTypeInt32 {
			variable.Type = "int32"
		} else if v.Type == device.PropertyTypeDouble {
			variable.Type = "double"
		} else if v.Type == device.PropertyTypeString {
			variable.Type = "string"
		}
		// 取出切片中最后一个值
		if len(v.Value) > 0 {
			index := len(v.Value) - 1
			variable.Index = v.Index
			variable.Value = v.Value[index].Value
			variable.Explain = v.Value[index].Explain
			variable.TimeStamp = v.Value[index].TimeStamp.Format("2006-01-02 15:04:05.000")
		} else {
			variable.Value = ""
			variable.Explain = ""
			variable.TimeStamp = ""
		}
		variables = append(variables, variable)
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "读取设备实时变量成功",
		Data:    variables,
	})
}

func ApiInvokeService(context *gin.Context) {

	aParam := struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
		Data    string `json:"Data"`
	}{
		Code:    "1",
		Message: "",
		Data:    "",
	}

	bodyBuf := make([]byte, 1024)
	n, _ := context.Request.Body.Read(bodyBuf)

	serviceInfo := struct {
		CollInterfaceName string                 `json:"collInterfaceName"`
		DeviceName        string                 `json:"deviceName"`
		ServiceName       string                 `json:"serviceName"`
		ServiceParam      map[string]interface{} `json:"serviceParam"`
	}{}
	err := json.Unmarshal(bodyBuf[:n], &serviceInfo)
	if err != nil {
		fmt.Println("serviceInfo json unMarshall err,", err)

		aParam.Code = "1"
		aParam.Message = "json unMarshall err"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	device.CollectInterfaceMap.Lock.Lock()
	coll, ok := device.CollectInterfaceMap.Coll[serviceInfo.CollInterfaceName]
	device.CollectInterfaceMap.Lock.Unlock()
	if !ok {
		aParam.Code = "1"
		aParam.Message = "device is not exist"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
		return
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = serviceInfo.CollInterfaceName
	cmd.DeviceName = serviceInfo.DeviceName
	cmd.FunName = serviceInfo.ServiceName
	paramStr, _ := json.Marshal(serviceInfo.ServiceParam)
	cmd.FunPara = string(paramStr)
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		aParam.Code = "0"
		aParam.Message = ""
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
	} else {
		aParam.Code = "1"
		aParam.Message = "device is not return"
		sJson, _ := json.Marshal(aParam)
		context.String(http.StatusOK, string(sJson))
	}
}

/*
*
从缓存中获取设备变量
*/
func ApiV2GetNodeVariableFromCache(context *gin.Context) {

	type VariableTemplate struct {
		Index     int         `json:"index"`     // 变量偏移量
		Name      string      `json:"name"`      // 变量名
		Label     string      `json:"lable"`     // 变量标签
		Value     interface{} `json:"value"`     // 变量值
		Explain   interface{} `json:"explain"`   // 变量值说明
		TimeStamp string      `json:"timestamp"` // 变量时间戳
		Type      string      `json:"type"`      // 变量类型
		Unit      string      `json:"unit"`      // 变量单位
	}

	nodeName := context.Query("name")

	node := &device.DeviceNodeTemplate{}
	collName := ""
	device.CollectInterfaceMap.Lock.Lock()
	for k, v := range device.CollectInterfaceMap.Coll {
		n, ok := v.DeviceNodeMap[nodeName]
		if ok {
			collName = k
			node = n
			break
		}
	}
	device.CollectInterfaceMap.Lock.Unlock()

	variables := make([]VariableTemplate, 0)
	if collName == "" {
		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "设备名称不存在",
			"Data":    variables,
		})
		return
	}

	index := 0
	variable := VariableTemplate{}
	for _, p := range node.Properties {
		variable.Name = p.Name
		variable.Label = p.Label
		variable.Unit = node.GetTSLModelUnit(p.Name)
		if p.Type == device.PropertyTypeUInt32 {
			variable.Type = "uint32"
		} else if p.Type == device.PropertyTypeInt32 {
			variable.Type = "int32"
		} else if p.Type == device.PropertyTypeDouble {
			variable.Type = "double"
		} else if p.Type == device.PropertyTypeString {
			variable.Type = "string"
		}
		// 取出切片中最后一个值
		if len(p.Value) > 0 {
			index = len(p.Value) - 1
			variable.Index = p.Value[index].Index
			variable.Value = p.Value[index].Value
			variable.Explain = p.Value[index].Explain
			variable.TimeStamp = p.Value[index].TimeStamp.Format("2006-01-02 15:04:05.000")
		} else {
			variable.Value = ""
			variable.Explain = ""
			variable.TimeStamp = ""
		}
		variables = append(variables, variable)
	}
	context.JSON(http.StatusOK, gin.H{
		"Code":    "0",
		"Message": "",
		"Data":    variables,
	})

}

func ApiV2GetNodeRealtimeVariable(context *gin.Context) {

	type VariableTemplate struct {
		Index     int         `json:"index"` // 变量偏移量
		Name      string      `json:"name"`  // 变量名
		Label     string      `json:"lable"` // 变量标签
		Value     interface{} `json:"value"` // 变量值
		Explain   interface{} `json:"explain"`
		TimeStamp string      `json:"timestamp"` // 变量时间戳
		Type      string      `json:"type"`      // 变量类型
		Unit      string      `json:"unit"`      // 变量单位
	}

	nodeName := context.Query("name")
	coll := &device.CollectInterfaceTemplate{}
	node := &device.DeviceNodeTemplate{}
	collName := ""
	device.CollectInterfaceMap.Lock.Lock()
	for k, v := range device.CollectInterfaceMap.Coll {
		n, ok := v.DeviceNodeMap[nodeName]
		if ok {
			coll = v
			collName = k
			node = n
			break
		}
	}
	device.CollectInterfaceMap.Lock.Unlock()

	variables := make([]VariableTemplate, 0)
	if collName == "" {
		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "设备名称不存在",
			"Data":    variables,
		})
		return
	}

	//发送命令到响应的采集接口
	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = coll.CollInterfaceName
	cmd.DeviceName = node.Name
	cmd.FunName = "GetRealVariables"
	cmd.FunPara = ""
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		index := 0
		variable := VariableTemplate{}
		for _, v := range node.Properties {
			//variable.Index = v.Index
			variable.Name = v.Name
			variable.Label = v.Label
			variable.Unit = node.GetTSLModelUnit(v.Name)
			// 取出切片中最后一个值
			if len(v.Value) > 0 {
				index = len(v.Value) - 1
				variable.Value = v.Value[index].Value
				variable.Explain = v.Value[index].Explain
				variable.TimeStamp = v.Value[index].TimeStamp.Format("2006-01-02 15:04:05.000")
			} else {
				variable.Value = ""
				variable.Explain = ""
				variable.TimeStamp = ""
			}
			//variable.Type = v.Type
			variables = append(variables, variable)
		}
		context.JSON(http.StatusOK, gin.H{
			"Code":    "0",
			"Message": "",
			"Data":    variables,
		})
		return
	} else {
		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "设备无相应",
			"Data":    variables,
		})
		return
	}
}

func ApiV2InvokeNodeService(context *gin.Context) {

	serviceInfo := struct {
		DeviceName   string
		ServiceName  string
		ServiceParam map[string]interface{}
	}{}

	err := context.BindJSON(&serviceInfo)
	if err != nil {
		setting.ZAPS.Errorf("调用设备服务JSON格式错误 %v", err)

		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "调用设备服务JSON格式错误",
			"Data":    "",
		})
		return
	}

	coll := &device.CollectInterfaceTemplate{}
	collName := ""
	device.CollectInterfaceMap.Lock.Lock()
	for k, v := range device.CollectInterfaceMap.Coll {
		_, ok := v.DeviceNodeMap[serviceInfo.DeviceName]
		if ok {
			coll = v
			collName = k
			break
		}
	}
	device.CollectInterfaceMap.Lock.Unlock()
	if collName == "" {
		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "设备名称不存在",
			"Data":    "",
		})
		return
	}

	cmd := device.CommunicationCmdTemplate{}
	cmd.CollInterfaceName = coll.CollInterfaceName
	cmd.DeviceName = serviceInfo.DeviceName
	cmd.FunName = serviceInfo.ServiceName
	paramStr, _ := json.Marshal(serviceInfo.ServiceParam)
	cmd.FunPara = string(paramStr)
	cmdRX := coll.CommQueueManage.CommunicationManageAddEmergency(cmd)
	if cmdRX.Status == true {
		context.JSON(http.StatusOK, gin.H{
			"Code":    "0",
			"Message": "",
			"Data":    "",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"Code":    "1",
			"Message": "设备无返回",
			"Data":    "",
		})
	}

}
