package device

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

const (
	TSLAccessModeRead int = iota
	TSLAccessModeWrite
	TSLAccessModeReadWrite
)

const (
	PropertyTypeUInt32 int = iota
	PropertyTypeInt32
	PropertyTypeDouble
	PropertyTypeString
)

const (
	TSLModelTypeLua int = iota
	TSLModelTypeS7
	TSLModelTypeModbus
	TSLModelTypeDLT6452007
)

type TSLPropertiesTemplate struct {
	Index      int                        `json:"index"`
	Name       string                     `json:"name"`       //属性名称，只可以是字母和数字的组合
	Label      string                     `json:"label"`      //属性解释
	AccessMode int                        `json:"accessMode"` //读写属性
	Type       int                        `json:"type"`       //类型 uint32 int32 double string
	Value      []TSLPropertyValueTemplate `json:"value"`
	Identity   string                     `json:"identity"` //唯一标识
}

type TSLPropertyValueTemplate struct {
	Index     int         `json:"index"`
	Value     interface{} `json:"value"`   //变量值，不可以是字符串
	Explain   interface{} `json:"explain"` //变量值解释，必须是字符串
	TimeStamp time.Time   `json:"timeStamp"`
}

type TSLEventTemplate struct {
	Topic string
	TSL   string
}

type TSLModelInterface interface {
	GetTSLModelIndex() int
	GetTSLModelName() string
	GetTSLModelLabel() string
	GetTSLModelType() int
	GetTSLModelParam() interface{}

	TSLModelPropertiesGet() interface{}
	TSLModelPropertiesAdd(property json.RawMessage) error
	TSLModelPropertiesDelete(propertiesName []string) error
	TSLModelPropertiesModify(property json.RawMessage) error
}

var TSLModels = make(map[string]TSLModelInterface)

func TSLModelsInit() {

	TSLLuaInit()
	for _, v := range TSLLuaMap {
		TSLModels[v.Name] = v
	}

	TSLModelS7Init()
	for _, v := range TSLModelS7Map {
		TSLModels[v.Name] = v
	}

	TSLModbusInit()
	TSLDLT6452007Init()
}

func AddTSLMode(name string, label string, modelType int) error {

	switch modelType {
	case TSLModelTypeLua:
		{
			_, ok := TSLModels[name]
			if ok {
				return errors.New("模型名称已经存在")
			}
			plugin := TSLLuaPluginTemplate{}
			luaModel := NewTSLLua(name, label, modelType, plugin)
			TSLLuaMap[name] = luaModel
			TSLModels[name] = luaModel
			WriteTSLLuaParamToJson()
		}
	case TSLModelTypeS7:
		{
			_, ok := TSLModels[name]
			if ok {
				return errors.New("模型名称已经存在")
			}
			s7Model := NewTSLModelS7(name, label, modelType)
			TSLModelS7Map[name] = s7Model
			TSLModels[name] = s7Model
			WriteTSLModelS7ParamToJson()
		}
	case TSLModelTypeModbus:
		{
			modbusModel := NewTSLModbus(name, label, modelType)
			TSLModbusMap[name] = modbusModel
			//TSLModels[name] = modbusModel
			WriteTSLModbusParamToJson()
		}
	case TSLModelTypeDLT6452007:
		{
			d07Model := NewTSLDLT6452007(name, label, modelType)
			TSLDLT6452007Map[name] = d07Model
			WriteTSLDLT6452007ParamToJson()
		}
	}
	return nil

}

func DeleteTSLMode(name string, modelType int) error {

	switch modelType {
	case TSLModelTypeLua:
		_, ok := TSLModels[name]
		if !ok {
			return errors.New("模型名称不存在")
		}
		err := DeleteTSLModelLua(name)
		if err != nil {
			return err
		}
	case TSLModelTypeS7:
		_, ok := TSLModels[name]
		if !ok {
			return errors.New("模型名称不存在")
		}
		err := DeleteTSLModelS7(name)
		if err != nil {
			return err
		}
	case TSLModelTypeModbus:
		err := DeleteTSLModbus(name)
		if err != nil {
			return err
		}
	case TSLModelTypeDLT6452007:
		err := DeleteTSLDLT6452007(name)
		if err != nil {
			return err
		}
	}

	delete(TSLModels, name)

	return nil
}

func ModifyTSLMode(name string, label string, modelType int, param json.RawMessage) error {

	err := errors.New("")
	switch modelType {
	case TSLModelTypeLua:
		_, ok := TSLModels[name]
		if !ok {
			errors.New("模型名称不存在")
		}
		err = ModifyTSLModelLua(name, label, param)
	case TSLModelTypeS7:
		err = ModifyTSLModelLua(name, label, param)
		_, ok := TSLModels[name]
		if !ok {
			errors.New("模型名称不存在")
		}
		err = ModifyTSLModelS7(name, label)
	case TSLModelTypeModbus:
		err = ModifyTSLModbus(name, label)
	case TSLModelTypeDLT6452007:
		err = ModifyTSLDLT6452007(name, label)
	}

	return err
}
