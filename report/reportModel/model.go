package reportModel

import (
	"encoding/json"
	"errors"
	"gateway/setting"
	"gateway/utils"
	"time"
)

type ReportModelTemplate struct {
	Index      int                                     `json:"-"`
	Name       string                                  `json:"name"`       //名称，只可以是字母和数字的组合
	Label      string                                  `json:"label"`      //名称标签
	Code       string                                  `json:"code"`       //编码
	Properties map[string]*ReportModelPropertyTemplate `json:"properties"` //属性
}

type ReportModelPropertyTemplate struct {
	Index      int                              `json:"index"`
	Name       string                           `json:"name"` //上报模型名称，只可以是字母和数字的组合
	UploadName string                           `json:"uploadName"`
	Label      string                           `json:"label"`    //上报模型标签
	Type       int                              `json:"type"`     //类型 uint32 int32 double string
	Decimals   int                              `json:"decimals"` //小数位数
	Unit       string                           `json:"unit"`     //单位
	Params     ReportModelPropertyParamTemplate `json:"params"`
}

type ReportModelPropertyParamTemplate struct {
	Min             string `json:"min"`             //最小
	Max             string `json:"max"`             //最大
	MinMaxAlarm     bool   `json:"minMaxAlarm"`     //范围报警
	Step            string `json:"step"`            //步长
	StepAlarm       bool   `json:"stepAlarm"`       //阶跃报警
	DataLength      string `json:"dataLength"`      //字符串长度
	DataLengthAlarm bool   `json:"dataLengthAlarm"` //字符长度报警
}

var writeTimer *time.Timer
var ReportModels = make(map[string]*ReportModelTemplate)

func ReportModelInit() {
	_ = ReadReportModelParamFromJson()
	writeTimer = time.AfterFunc(time.Second, func() {
		_ = WriteReportModelParamToJson()
	})
	writeTimer.Stop()
}

func WriteReportModelParamToJson() error {
	utils.DirIsExist("./selfpara")

	data, err := json.Marshal(ReportModels)
	if err != nil {
		setting.ZAPS.Errorf("上报模型配置json格式化失败 %v", err)
		return errors.New("上报模型配置json格式化失败")
	}

	err = utils.FileWrite("./selfpara/reportModel.json", data)
	if err != nil {
		setting.ZAPS.Errorf("上报模型配置json文件写入失败 %v", err)
		return errors.New("上报模型配置json文件写入失败")
	}
	setting.ZAPS.Debug("上报模型配置json文件写入成功")

	return nil
}

func ReadReportModelParamFromJson() error {

	data, err := utils.FileRead("./selfpara/reportModel.json")
	if err != nil {
		setting.ZAPS.Debugf("上报模型配置json文件读取失败 %v", err)
		return errors.New("上报模型配置json文件读取失败")
	}

	err = json.Unmarshal(data, &ReportModels)
	if err != nil {
		setting.ZAPS.Errorf("上报模型配置json格式化失败 %v", err)
		return errors.New("上报模型配置json格式化失败")
	}

	setting.ZAPS.Info("上报模型配置json文件读取成功")
	return nil
}

func NewReportModel(name string, label string, code string) *ReportModelTemplate {

	return &ReportModelTemplate{
		Index:      len(ReportModels),
		Name:       name,
		Label:      label,
		Code:       code,
		Properties: make(map[string]*ReportModelPropertyTemplate),
	}
}

func AddReportModel(name string, label string, code string) error {

	_, ok := ReportModels[name]
	if ok {
		setting.ZAPS.Error("上报模型名称已经存在")
		return errors.New("上报模型名称已经存在")
	}

	model := NewReportModel(name, label, code)
	ReportModels[name] = model
	setting.ZAPS.Info("上报模型添加成功")
	_ = WriteReportModelParamToJson()
	return nil
}

func DeleteReportModel(name string) error {

	_, ok := ReportModels[name]
	if !ok {
		setting.ZAPS.Error("上报模型名称不存在")
		return errors.New("上报模型名称不存在")
	}

	delete(ReportModels, name)

	setting.ZAPS.Info("上报模型删除成功")
	_ = WriteReportModelParamToJson()
	return nil
}

func ModifyReportModel(name string, label string, code string) error {

	model, ok := ReportModels[name]
	if !ok {
		setting.ZAPS.Error("上报模型名称不存在")
		return errors.New("上报模型名称不存在")
	}

	model.Label = label
	model.Code = code
	setting.ZAPS.Info("上报模型修改成功")
	_ = WriteReportModelParamToJson()
	return nil
}

func NewReportModelProperty() *ReportModelPropertyTemplate {
	return &ReportModelPropertyTemplate{}
}

func AddReportModelProperty(modelName string, property *ReportModelPropertyTemplate) error {
	model, ok := ReportModels[modelName]
	if !ok {
		setting.ZAPS.Error("上报模型名称不存在")
		return errors.New("上报模型名称不存在")
	}

	_, ok = model.Properties[property.Name]
	if ok {
		setting.ZAPS.Error("上报模型属性名称已经存在")
		return errors.New("上报模型属性名称已经存在")
	}

	property.Index = len(model.Properties)
	model.Properties[property.Name] = property

	setting.ZAPS.Info("上报模型属性添加成功")
	writeTimer.Reset(time.Second)
	return nil
}

func ModifyReportModelProperty(modelName string, property *ReportModelPropertyTemplate) error {
	model, ok := ReportModels[modelName]
	if !ok {
		setting.ZAPS.Error("上报模型名称不存在")
		return errors.New("上报模型名称不存在")
	}

	_, ok = model.Properties[property.Name]
	if !ok {
		setting.ZAPS.Error("上报模型属性名称不存在")
		return errors.New("上报模型属性名称不存在")
	}

	model.Properties[property.Name] = property

	setting.ZAPS.Info("上报模型属性修改成功")
	_ = WriteReportModelParamToJson()
	return nil
}

func DeleteReportModelProperties(modelName string, propertyNames []string) error {
	model, ok := ReportModels[modelName]
	if !ok {
		setting.ZAPS.Error("上报模型名称不存在")
		return errors.New("上报模型名称不存在")
	}

	for _, v := range propertyNames {
		_, ok = model.Properties[v]
		if ok {
			delete(model.Properties, v)
		}
	}

	setting.ZAPS.Info("上报模型属性删除成功")
	_ = WriteReportModelParamToJson()
	return nil
}
