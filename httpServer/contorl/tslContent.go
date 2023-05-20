package contorl

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gateway/device"
	"gateway/httpServer/model"
	"gateway/setting"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func ApiGetDeviceTSLContents(context *gin.Context) {

	type TSLLuaPropertyIntUintParamTemplate struct {
		Min         string `json:"min"`         //最小
		Max         string `json:"max"`         //最大
		MinMaxAlarm bool   `json:"minMaxAlarm"` //范围报警
		Step        string `json:"step"`        //步长
		StepAlarm   bool   `json:"stepAlarm"`   //阶跃报警
		Unit        string `json:"unit"`        //单位
	}

	type TSLLuaPropertyDoubleParamTemplate struct {
		Min         string `json:"min"`         //最小
		Max         string `json:"max"`         //最大
		MinMaxAlarm bool   `json:"minMaxAlarm"` //范围报警
		Step        string `json:"step"`        //步长
		StepAlarm   bool   `json:"stepAlarm"`   //阶跃报警
		Decimals    int    `json:"decimals"`    //小数位数
		Unit        string `json:"unit"`        //单位
	}

	type TSLLuaPropertyStringParamTemplate struct {
		DataLength      string `json:"dataLength,omitempty"`      //字符串长度
		DataLengthAlarm bool   `json:"dataLengthAlarm,omitempty"` //字符长度报警
	}

	type TSLLuaPropertyTemplate struct {
		Name       string      `json:"name"`       //属性名称，只可以是字母和数字的组合
		Explain    string      `json:"explain"`    //属性解释
		AccessMode int         `json:"accessMode"` //读写属性
		Type       int         `json:"type"`       //类型 uint32 int32 double string
		Params     interface{} `json:"params"`
	}

	type TSLLuaModelTemplate struct {
		Properties []TSLLuaPropertyTemplate       `json:"properties"` //属性
		Services   []device.TSLLuaServiceTemplate `json:"services"`   //服务
	}

	tslInfo := TSLLuaModelTemplate{}

	tslName := context.Query("name")
	for _, v := range device.TSLLuaMap {
		if v.Name == tslName {
			//tslInfo.Services = v.Services
			tslInfo.Properties = make([]TSLLuaPropertyTemplate, 0)
			property := TSLLuaPropertyTemplate{}
			for _, p := range v.Properties {
				property.Name = p.Name
				property.Explain = p.Label
				property.AccessMode = p.AccessMode
				property.Type = p.Type
				switch p.Type {
				case device.PropertyTypeUInt32:
					fallthrough
				case device.PropertyTypeInt32:
					{
						intUintPropertyParam := TSLLuaPropertyIntUintParamTemplate{
							Min:         p.Params.Min,
							Max:         p.Params.Max,
							MinMaxAlarm: p.Params.MinMaxAlarm,
							Step:        p.Params.Step,
							StepAlarm:   p.Params.StepAlarm,
							Unit:        p.Unit,
						}
						property.Params = intUintPropertyParam
					}
				case device.PropertyTypeDouble:
					{
						doublePropertyParam := TSLLuaPropertyDoubleParamTemplate{
							Min:         p.Params.Min,
							Max:         p.Params.Max,
							MinMaxAlarm: p.Params.MinMaxAlarm,
							Step:        p.Params.Step,
							StepAlarm:   p.Params.StepAlarm,
							Unit:        p.Unit,
							Decimals:    p.Decimals,
						}
						property.Params = doublePropertyParam
					}
				case device.PropertyTypeString:
					{
						stringPropertyParam := TSLLuaPropertyStringParamTemplate{
							DataLength:      p.Params.DataLength,
							DataLengthAlarm: p.Params.DataLengthAlarm,
						}
						property.Params = stringPropertyParam
					}
				}
				tslInfo.Properties = append(tslInfo.Properties, property)
			}

			context.JSON(http.StatusOK, model.ResponseData{
				Code:    "0",
				Message: "获取物模型模型内容成功",
				Data:    tslInfo,
			})
			return
		}
	}
	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "1",
		Message: "物模型模型名称不存在",
		Data:    "",
	})
}

func ApiImportDeviceTSLContents(context *gin.Context) {

	// 获取物模型名称
	tslName := context.PostForm("name")

	tslType := -1
	tslLua, ok := device.TSLLuaMap[tslName]
	if ok {
		tslType = 0
	}
	tslS7, ok := device.TSLModelS7Map[tslName]
	if ok {
		tslType = 1
	}

	if tslType == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("物模型模版待导入CSV文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版待导入CSV文件不存在",
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

	result := setting.LoadCsvCfg(fileName, 1, 2, 1) //标题在第2行，从第3行取数据，第2列取数据
	if result == nil {
		setting.ZAPS.Errorf("加载CSV文件错误 %v", err)
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "加载CSV文件错误",
			Data:    "",
		})
		return
	}

	for _, record := range result.Records {
		if record.GetString("ContentType") == "property" {
			if tslType == 0 {
				property := device.TSLLuaPropertyTemplate{
					Name:       record.GetString("Name"),
					Label:      record.GetString("Label"),
					AccessMode: record.GetInt("AccessMode"),
					Type:       record.GetInt("Type"),
					Decimals:   record.GetInt("Decimals"),
					Unit:       record.GetString("Unit"),
					//Params:     propertyParam,
				}
				propertyJson, err := json.Marshal(property)
				if err != nil {
					setting.ZAPS.Errorf("添加属性JSON格式化错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "添加属性JSON格式化错误",
						Data:    "",
					})
					return
				}
				err = tslLua.TSLModelPropertiesAdd(propertyJson)
				if err != nil {
					setting.ZAPS.Errorf("从CSV文件中添加属性错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "从CSV文件中添加属性错误",
						Data:    "",
					})
					return
				}
			} else if tslType == 1 {
				propertyParam := device.TSLModelS7PropertyParamTemplate{
					DBNumber:  record.GetString("DBNumber"),
					DataType:  record.GetInt("DataType"),
					StartAddr: record.GetString("StartAddr"),
				}
				property := device.TSLModelS7PropertyTemplate{
					Name:       record.GetString("Name"),
					Label:      record.GetString("Label"),
					AccessMode: record.GetInt("AccessMode"),
					Type:       record.GetInt("Type"),
					Decimals:   record.GetInt("Decimals"),
					Unit:       record.GetString("Unit"),
					Params:     propertyParam,
				}
				propertyJson, err := json.Marshal(property)
				if err != nil {
					setting.ZAPS.Errorf("添加属性JSON格式化错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "添加属性JSON格式化错误",
						Data:    "",
					})
					return
				}
				_ = tslS7.TSLModelPropertiesAdd(propertyJson)
			}
		} else if record.GetString("ContentType") == "service" {
			if tslType == 0 {
				service := device.TSLLuaServiceTemplate{
					Name:   record.GetString("Name"),
					Label:  record.GetString("Label"),
					Params: make(map[string]interface{}),
				}
				callType := record.GetString("CallType")
				if callType == "synchronous" {
					service.CallType = 0
				} else if callType == "asynchronous" {
					service.CallType = 1
				}

				_, err = tslLua.TSLLuaServicesAdd(service)
				if err != nil {
					setting.ZAPS.Errorf("从CSV文件中添加服务错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "从CSV文件中添加服务错误",
						Data:    "",
					})
					return
				}
			}
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "物模型模版导入CSV成功",
		Data:    "",
	})
	return
}

func ApiImportDeviceTSLContentsFromXlsx(context *gin.Context) {

	// 获取物模型名称
	tslName := context.PostForm("name")

	tslType := -1
	tslLua, ok := device.TSLLuaMap[tslName]
	if ok {
		tslType = 0
	}
	tslS7, ok := device.TSLModelS7Map[tslName]
	if ok {
		tslType = 1
	}

	if tslType == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版不存在",
			Data:    "",
		})
		return
	}

	// 获取文件头
	file, err := context.FormFile("fileName")
	if err != nil {
		setting.ZAPS.Errorf("物模型模版待导入xlsx文件不存在")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版待导入xlsx文件不存在",
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
		setting.ZAPS.Errorf("保存xlsx文件错误 %v", err)
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

	for _, cell := range cells {
		if cell[1] == "property" {
			if tslType == 0 {
				if len(cell) < 16 {
					continue
				}
				property := device.TSLLuaPropertyTemplate{
					Name:       setting.GetString(cell[2]),
					Label:      setting.GetString(cell[3]),
					AccessMode: setting.GetInt(cell[4]),
					Type:       setting.GetInt(cell[5]),
					Decimals:   setting.GetInt(cell[6]),
					Unit:       setting.GetString(cell[7]),
				}
				propertyJson, err := json.Marshal(property)
				if err != nil {
					setting.ZAPS.Errorf("添加属性JSON格式化错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "添加属性JSON格式化错误",
						Data:    "",
					})
					return
				}
				err = tslLua.TSLModelPropertiesAdd(propertyJson)
				if err != nil {
					setting.ZAPS.Errorf("从xlsx文件中添加属性错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "从xlsx文件中添加属性错误",
						Data:    "",
					})
					return
				}
			} else if tslType == 1 {
				if len(cell) < 11 {
					continue
				}
				propertyParam := device.TSLModelS7PropertyParamTemplate{
					DBNumber:  setting.GetString(cell[8]),
					DataType:  setting.GetInt(cell[9]),
					StartAddr: setting.GetString(cell[10]),
				}
				property := device.TSLModelS7PropertyTemplate{
					Name:       setting.GetString(cell[2]),
					Label:      setting.GetString(cell[3]),
					AccessMode: setting.GetInt(cell[4]),
					Type:       setting.GetInt(cell[5]),
					Decimals:   setting.GetInt(cell[6]),
					Unit:       setting.GetString(cell[7]),
					Params:     propertyParam,
				}
				propertyJson, err := json.Marshal(property)
				if err != nil {
					setting.ZAPS.Errorf("添加属性JSON格式化错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "添加属性JSON格式化错误",
						Data:    "",
					})
					return
				}
				_ = tslS7.TSLModelPropertiesAdd(propertyJson)
			}
		}
		//else if record.GetString("ContentType") == "service" {
		//	if tslType == 0 {
		//		service := device.TSLLuaServiceTemplate{
		//			Name:   record.GetString("Name"),
		//			Label:  record.GetString("Label"),
		//			Params: make(map[string]interface{}),
		//		}
		//		callType := record.GetString("CallType")
		//		if callType == "synchronous" {
		//			service.CallType = 0
		//		} else if callType == "asynchronous" {
		//			service.CallType = 1
		//		}
		//
		//		_, err = tslLua.TSLLuaServicesAdd(service)
		//		if err != nil {
		//			setting.ZAPS.Errorf("从CSV文件中添加服务错误 %v", err)
		//			context.JSON(http.StatusOK, model.ResponseData{
		//				Code:    "1",
		//				Message: "从CSV文件中添加服务错误",
		//				Data:    "",
		//			})
		//			return
		//		}
		//	}
		//}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "采集物模型导入xlsx成功",
		Data:    "",
	})
	return
}

func ApiGetDeviceTSLContentsFromPluginXlsx(context *gin.Context) {

	// 获取物模型名称
	tslName := context.Query("name")

	tslType := -1
	tslLua, ok := device.TSLLuaMap[tslName]
	if ok {
		tslType = 0
	}
	tslS7, ok := device.TSLModelS7Map[tslName]
	if ok {
		tslType = 1
	}

	if tslType == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版不存在",
			Data:    "",
		})
		return
	}

	fileName := "./plugin/" + tslName + "/" + tslName + ".xlsx"

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

	for _, cell := range cells {
		if cell[1] == "property" {
			if tslType == 0 {
				if len(cell) < 16 {
					continue
				}
				property := device.TSLLuaPropertyTemplate{
					Name:       setting.GetString(cell[2]),
					Label:      setting.GetString(cell[3]),
					AccessMode: setting.GetInt(cell[4]),
					Type:       setting.GetInt(cell[5]),
					Decimals:   setting.GetInt(cell[6]),
					Unit:       setting.GetString(cell[7]),
				}
				propertyJson, err := json.Marshal(property)
				if err != nil {
					setting.ZAPS.Errorf("添加属性JSON格式化错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "添加属性JSON格式化错误",
						Data:    "",
					})
					return
				}
				err = tslLua.TSLModelPropertiesAdd(propertyJson)
				if err != nil {
					setting.ZAPS.Errorf("从xlsx文件中添加属性错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "从xlsx文件中添加属性错误",
						Data:    "",
					})
					return
				}
			} else if tslType == 1 {
				if len(cell) < 11 {
					continue
				}
				propertyParam := device.TSLModelS7PropertyParamTemplate{
					DBNumber:  setting.GetString(cell[2]),
					DataType:  setting.GetInt(cell[3]),
					StartAddr: setting.GetString(cell[4]),
				}
				property := device.TSLModelS7PropertyTemplate{
					Name:       setting.GetString(cell[2]),
					Label:      setting.GetString(cell[3]),
					AccessMode: setting.GetInt(cell[4]),
					Type:       setting.GetInt(cell[5]),
					Decimals:   setting.GetInt(cell[6]),
					Unit:       setting.GetString(cell[7]),
					Params:     propertyParam,
				}
				propertyJson, err := json.Marshal(property)
				if err != nil {
					setting.ZAPS.Errorf("添加属性JSON格式化错误 %v", err)
					context.JSON(http.StatusOK, model.ResponseData{
						Code:    "1",
						Message: "添加属性JSON格式化错误",
						Data:    "",
					})
					return
				}
				_ = tslS7.TSLModelPropertiesAdd(propertyJson)
			}
		}
	}

	context.JSON(http.StatusOK, model.ResponseData{
		Code:    "0",
		Message: "从plugin中获取物模型内容成功",
		Data:    "",
	})
	return
}

func ApiExportDeviceTSLContentsToCSV(context *gin.Context) {

	tslName := context.Query("name")

	tslType := -1
	tslLua, ok := device.TSLLuaMap[tslName]
	if ok {
		tslType = 0
	}
	tslS7, ok := device.TSLModelS7Map[tslName]
	if ok {
		tslType = 1
	}

	if tslType == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileName := exeCurDir + "/selfpara/" + tslName + ".csv"

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
	csvRecords := make([][]string, 0)
	csvFile := csv.NewWriter(fs)
	if tslType == 0 {
		csvRecords = [][]string{
			{"模型名称", "功能类型", "功能名称", "标识符", "读写类型", "数据类型", "小数位",
				"单位", "最小值", "最大值", "范围报警", "步长", "步长报警", "字符串长度", "字符串长度报警", "服务调用方式"},
			{"TSLName", "ContentType", "Name", "Label", "AccessMode", "Type", "Decimals",
				"Unit", "Min", "Max", "MinMaxAlarm", "Step", "StepAlarm", "DataLength", "DataLengthAlarm", "CallType"},
		}

		for _, v := range tslLua.Properties {
			record := make([]string, 0)
			record = append(record, tslLua.Name)
			record = append(record, "property")
			record = append(record, v.Name)
			record = append(record, v.Label)
			record = append(record, strconv.Itoa(v.AccessMode))
			record = append(record, strconv.Itoa(v.Type))
			record = append(record, fmt.Sprintf("%d", v.Decimals))
			record = append(record, v.Unit)
			record = append(record, v.Params.Min)
			record = append(record, v.Params.Max)
			str := ""
			if v.Params.MinMaxAlarm == true {
				str = "true"
			} else {
				str = "false"
			}
			record = append(record, str)
			record = append(record, v.Params.Step)
			if v.Params.StepAlarm == true {
				str = "true"
			} else {
				str = "false"
			}
			record = append(record, str)
			record = append(record, v.Params.DataLength)
			if v.Params.DataLengthAlarm == true {
				str = "true"
			} else {
				str = "false"
			}
			record = append(record, str)
			record = append(record, "-")
			csvRecords = append(csvRecords, record)
		}

		for _, v := range tslLua.Services {
			record := make([]string, 0)
			record = append(record, tslLua.Name)
			record = append(record, "service")
			record = append(record, v.Name)
			record = append(record, v.Label)
			//空
			for i := 0; i < 11; i++ {
				record = append(record, "-")
			}
			if v.CallType == 0 {
				record = append(record, "synchronous")
			} else {
				record = append(record, "asynchronous")
			}
			csvRecords = append(csvRecords, record)
		}
	} else if tslType == 1 {
		csvRecords = [][]string{
			{"模型名称", "功能类型", "功能名称", "标识符", "读写类型", "数据类型", "小数位",
				"单位", "数据块", "数据类型", "数据地址"},
			{"TSLName", "ContentType", "Name", "Label", "AccessMode", "Type", "Decimals",
				"Unit", "DBNumber", "DataType", "StartAddr"},
		}

		for _, v := range tslS7.Properties {
			record := make([]string, 0)
			record = append(record, tslS7.Name)
			record = append(record, "property")
			record = append(record, v.Name)
			record = append(record, v.Label)
			record = append(record, strconv.Itoa(v.AccessMode))
			record = append(record, strconv.Itoa(v.Type))
			record = append(record, strconv.Itoa(v.Decimals))
			record = append(record, v.Unit)
			record = append(record, v.Params.DBNumber)
			record = append(record, fmt.Sprintf("%d", v.Params.DataType))
			record = append(record, v.Params.StartAddr)

			csvRecords = append(csvRecords, record)
		}
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

func ApiExportDeviceTSLContentsToXlsx(context *gin.Context) {

	tslName := context.Query("name")

	tslType := -1
	tslLua, ok := device.TSLLuaMap[tslName]
	if ok {
		tslType = 0
	}
	tslS7, ok := device.TSLModelS7Map[tslName]
	if ok {
		tslType = 1
	}

	if tslType == -1 {
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "物模型模版不存在",
			Data:    "",
		})
		return
	}

	//创建文件
	utils.DirIsExist("./tmp")
	fileName := "./tmp/" + tslName + ".xlsx"

	//创建一个新的写入文件流
	csvRecords := make([][]string, 0)
	if tslType == 0 {
		csvRecords = [][]string{
			{"模型名称", "功能类型", "功能名称", "标识符", "读写类型", "数据类型", "小数位",
				"单位", "最小值", "最大值", "范围报警", "步长", "步长报警", "字符串长度", "字符串长度报警", "服务调用方式"},
			{"TSLName", "ContentType", "Name", "Label", "AccessMode", "Type", "Decimals",
				"Unit", "Min", "Max", "MinMaxAlarm", "Step", "StepAlarm", "DataLength", "DataLengthAlarm", "CallType"},
		}

		for _, v := range tslLua.Properties {
			record := make([]string, 0)
			record = append(record, tslLua.Name)
			record = append(record, "property")
			record = append(record, v.Name)
			record = append(record, v.Label)
			record = append(record, strconv.Itoa(v.AccessMode))
			record = append(record, strconv.Itoa(v.Type))
			record = append(record, fmt.Sprintf("%d", v.Decimals))
			record = append(record, v.Unit)
			record = append(record, v.Params.Min)
			record = append(record, v.Params.Max)
			str := ""
			if v.Params.MinMaxAlarm == true {
				str = "true"
			} else {
				str = "false"
			}
			record = append(record, str)
			record = append(record, v.Params.Step)
			if v.Params.StepAlarm == true {
				str = "true"
			} else {
				str = "false"
			}
			record = append(record, str)
			record = append(record, v.Params.DataLength)
			if v.Params.DataLengthAlarm == true {
				str = "true"
			} else {
				str = "false"
			}
			record = append(record, str)
			record = append(record, "-")
			csvRecords = append(csvRecords, record)
		}

		for _, v := range tslLua.Services {
			record := make([]string, 0)
			record = append(record, tslLua.Name)
			record = append(record, "service")
			record = append(record, v.Name)
			record = append(record, v.Label)
			//空
			for i := 0; i < 11; i++ {
				record = append(record, "-")
			}
			if v.CallType == 0 {
				record = append(record, "synchronous")
			} else {
				record = append(record, "asynchronous")
			}
			csvRecords = append(csvRecords, record)
		}
	} else if tslType == 1 {
		csvRecords = [][]string{
			{"模型名称", "功能类型", "功能名称", "标识符", "读写类型", "数据类型", "小数位",
				"单位", "数据块", "数据类型", "数据地址"},
			{"TSLName", "ContentType", "Name", "Label", "AccessMode", "Type", "Decimals",
				"Unit", "DBNumber", "DataType", "StartAddr"},
		}

		for _, v := range tslS7.Properties {
			record := make([]string, 0)
			record = append(record, tslS7.Name)
			record = append(record, "property")
			record = append(record, v.Name)
			record = append(record, v.Label)
			record = append(record, strconv.Itoa(v.AccessMode))
			record = append(record, strconv.Itoa(v.Type))
			record = append(record, strconv.Itoa(v.Decimals))
			record = append(record, v.Unit)
			record = append(record, v.Params.DBNumber)
			record = append(record, fmt.Sprintf("%d", v.Params.DataType))
			record = append(record, v.Params.StartAddr)

			csvRecords = append(csvRecords, record)
		}
	}

	err := setting.WriteExcel(fileName, csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("保存xlsx文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "保存xlsx文件错误",
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

func ApiExportDeviceTSLContentsTemplate(context *gin.Context) {

	//创建文件
	fileName := "./tmp/TSLModelTemplate.csv"

	fs, err := os.Create(fileName)
	if err != nil {
		setting.ZAPS.Errorf("创建物模型模版CSV文件错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "创建物模型模版CSV文件错误",
			Data:    "",
		})
		return
	}
	defer os.Remove(fileName)
	defer fs.Close()
	// 写入UTF-8 BOM
	//_, err = fs.WriteString("\xEF\xBB\xBF")

	//创建一个新的写入文件流
	csvFile := csv.NewWriter(fs)
	csvRecords := [][]string{
		{"模型名称", "功能类型", "功能名称", "标识符", "读写类型", "数据类型", "小数位",
			"单位", "最小值", "最大值", "范围报警", "步长", "步长报警", "字符串长度", "字符串长度报警", "服务调用方式"},
		{"TSLName", "ContentType", "Name", "Explain", "AccessMode", "Type", "Decimals",
			"Unit", "Min", "Max", "MinMaxAlarm", "Step", "StepAlarm", "DataLength", "DataLengthAlarm", "CallType"},
	}

	err = csvFile.WriteAll(csvRecords)
	if err != nil {
		setting.ZAPS.Errorf("向物模型模版CSV文件写入记录错误")
		context.JSON(http.StatusOK, model.ResponseData{
			Code:    "1",
			Message: "向物模型模版CSV文件写入记录错误",
			Data:    "",
		})
		return
	}
	csvFile.Flush()

	//返回文件流
	context.Writer.Header().Add("Content-Disposition",
		fmt.Sprintf("attachment;filename=%s", filepath.Base(fileName)))
	context.File(fileName) //返回文件路径，自动调用http.ServeFile方法

	return
}
