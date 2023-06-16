package device

import (
	"encoding/json"
	"gateway/device/eventBus"
	"gateway/setting"
	"gateway/utils"
	"github.com/pkg/errors"
	"time"
)

// 物模型 Thing Specification Language
type TSLModelS7Template struct {
	Index      int                           `json:"index"`
	Name       string                        `json:"name"`       //名称，只可以是字母和数字的组合
	Label      string                        `json:"label"`      //名称标签
	Type       int                           `json:"type"`       //模型类型
	Properties []*TSLModelS7PropertyTemplate `json:"properties"` //属性
	Event      eventBus.Bus                  `json:"-"`          //事件队列
}

type TSLModelS7PropertyTemplate struct {
	Name       string                            `json:"name"`       //属性名称，只可以是字母和数字的组合
	Label      string                            `json:"label"`      //属性解释
	AccessMode int                               `json:"accessMode"` //读写属性
	Type       int                               `json:"type"`       //类型 uint32 int32 double string
	Params     TSLModelS7PropertyParamTemplate   `json:"params"`
	Decimals   int                               `json:"decimals"` //小数位数
	Unit       string                            `json:"unit"`     //单位
	Value      []TSLModelS7PropertyValueTemplate `json:"value"`
}

type TSLModelS7PropertyParamTemplate struct {
	DBNumber    string `json:"dbNumber"`  //数据块
	DataType    int    `json:"dataType"`  //数据类型
	StartAddr   string `json:"startAddr"` //起始地址
	IotDataType string `json:"IotDataType"`
}

type TSLModelS7PropertyValueTemplate struct {
	Index     int         `json:"index"`
	Value     interface{} `json:"value"`   //变量值，不可以是字符串
	Explain   interface{} `json:"explain"` //变量值解释，必须是字符串
	TimeStamp time.Time   `json:"timeStamp"`
}

type TSLModelS7ServiceTemplate struct {
	Name     string                 `json:"name"`     //服务名称
	Label    string                 `json:"explain"`  //服务名称说明
	CallType int                    `json:"callType"` //服务调用方式
	Params   map[string]interface{} `json:"params"`   //服务参数
}

var TSLModelS7Map = make(map[string]*TSLModelS7Template)
var s7WriteTimer *time.Timer

func TSLModelS7Init() {

	ReadTSLModelS7ParamFromJson()
	s7WriteTimer = time.AfterFunc(time.Second, func() {
		WriteTSLModelS7ParamToJson()
	})
	s7WriteTimer.Stop()
}

func ReadTSLModelS7ParamFromJson() bool {

	utils.DirIsExist("./selfpara")

	data, err := utils.FileRead("./selfpara/TSLModelS7Param.json")
	if err != nil {
		setting.ZAPS.Debugf("打开物模型配置json文件失败 %v", err)
		return false
	}

	err = json.Unmarshal(data, &TSLModelS7Map)
	if err != nil {
		setting.ZAPS.Errorf("物模型配置json文件格式化失败 %v", err)
		return false
	}
	for _, v := range TSLModelS7Map {
		v.Event = eventBus.NewBus()
	}
	setting.ZAPS.Debugf("打开物模型配置json文件成功")
	return true
}

func WriteTSLModelS7ParamToJson() {

	utils.DirIsExist("./selfpara")

	sJson, _ := json.Marshal(TSLModelS7Map)
	err := utils.FileWrite("./selfpara/TSLModelS7Param.json", sJson)
	if err != nil {
		setting.ZAPS.Errorf("物模型配置json文件写入失败")
		return
	}
	setting.ZAPS.Infof("物模型配置json文件写入成功")
}

func NewTSLModelS7(name string, label string, modelType int) *TSLModelS7Template {
	return &TSLModelS7Template{
		Index:      len(TSLModelS7Map),
		Name:       name,
		Label:      label,
		Type:       modelType,
		Properties: make([]*TSLModelS7PropertyTemplate, 0),
		Event:      eventBus.NewBus(),
	}
}

func (t *TSLModelS7Template) GetTSLModelIndex() int {
	return t.Index
}

func (t *TSLModelS7Template) GetTSLModelName() string {
	return t.Name
}

func (t *TSLModelS7Template) GetTSLModelLabel() string {
	return t.Label
}

func (t *TSLModelS7Template) GetTSLModelType() int {
	return t.Type
}

func (t *TSLModelS7Template) GetTSLModelParam() interface{} {
	return ""
}

func AddTSLModelS7(name string, label string, modelType int) error {

	_, ok := TSLModelS7Map[name]
	if !ok {
		TSLModelS7Map[name] = NewTSLModelS7(name, label, modelType)
		WriteTSLModelS7ParamToJson()
		return nil
	}

	return errors.New("物模型模版已经存在")
}

func DeleteTSLModelS7(name string) error {

	for _, v := range CollectInterfaceMap.Coll {
		for _, d := range v.DeviceNodeMap {
			if d.TSL == name {
				return errors.New("物模型已被使用，不可以删除")
			}
		}
	}

	_, ok := TSLModelsName[name]
	if !ok {
		return errors.New("物模型不存在")
	}

	delete(TSLModelS7Map, name)
	delete(TSLModelsName, name)
	WriteTSLModelS7ParamToJson()

	return nil
}

func ModifyTSLModelS7(name string, label string) error {

	tsl, ok := TSLModelS7Map[name]
	if !ok {
		return errors.New("物模型模版不存在")
	}

	tsl.Label = label
	tsl.Type = TSLModelTypeS7
	setting.ZAPS.Infof("物模型[%s]Plugin[%s]插件发生变化", name)
	eventMsg := TSLEventTemplate{
		TSL:   name,
		Topic: "modify",
	}
	_ = tsl.Event.Publish("modify", eventMsg)
	WriteTSLModelS7ParamToJson()

	return nil
}

func (d *TSLModelS7Template) TSLModelPropertiesGet() interface{} {
	return d.Properties
}

func (d *TSLModelS7Template) TSLModelPropertiesAdd(property json.RawMessage) error {

	s7Property := TSLModelS7PropertyTemplate{}
	err := json.Unmarshal(property, &s7Property)
	if err != nil {
		return err
	}

	index := -1
	for k, v := range d.Properties {
		if v.Name == s7Property.Name {
			d.Properties[k] = &s7Property
			index = k
		}
	}
	if index == -1 {
		d.Properties = append(d.Properties, &s7Property)
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		_ = d.Event.Publish("modify", eventMsg)
	} else {
		return errors.New("属性名称已经存在")
	}

	s7WriteTimer.Reset(time.Second)

	return nil
}

func (d *TSLModelS7Template) TSLModelPropertiesDelete(propertiesName []string) error {

	cnt := 0
	for _, p := range propertiesName {
		for k := 0; k < len(d.Properties); k++ {
			if d.Properties[k].Name == p {
				setting.ZAPS.Debugf("Properties index %v", k)
				cnt = 1
				d.Properties = append(d.Properties[:k], d.Properties[k+1:]...)
			}
		}
	}
	if cnt != 0 {
		WriteTSLModelS7ParamToJson()
		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}
	return errors.New("property is not exist")

}

func (d *TSLModelS7Template) TSLModelPropertiesModify(property json.RawMessage) error {

	s7Property := TSLModelS7PropertyTemplate{}
	err := json.Unmarshal(property, &s7Property)
	if err != nil {
		return err
	}

	properties := make([]*TSLModelS7PropertyTemplate, 0)
	for _, v := range d.Properties {
		properties = append(properties, v)
	}
	index := -1
	for k, v := range properties {
		if v.Name == s7Property.Name {
			properties[k] = &s7Property
			index = k
		}
	}
	if index != -1 {
		d.Properties = properties
		WriteTSLModelS7ParamToJson()

		eventMsg := TSLEventTemplate{
			TSL:   d.Name,
			Topic: "modify",
		}
		d.Event.Publish("modify", eventMsg)
		return nil
	}

	return errors.New("属性名称不存在")
}
