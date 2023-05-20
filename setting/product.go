package setting

import (
	"encoding/json"
	"gateway/utils"
)

type ProductParamTemplate struct {
	SN      string `json:"SN"`
	HardVer string `json:"hardVer,omitempty"`
	Name    string `json:"name"`
}

var ProductParam ProductParamTemplate

func ReadProductParamFromJson() error {
	data, err := utils.FileRead("./product.json")
	if err != nil {
		ZAPS.Debugf("打开生产配置json文件失败 %v", err)
		return err
	}
	err = json.Unmarshal(data, &ProductParam)
	if err != nil {
		ZAPS.Errorf("生产配置json文件格式化失败 %v", err)
		return err
	}
	ZAPS.Debugf("打开生产配置json文件成功")

	SystemState.SN = ProductParam.SN
	SystemState.Name = ProductParam.Name

	return nil
}

func WriteProductParamToJson() error {
	sJson, _ := json.Marshal(ProductParam)
	err := utils.FileWrite("./product.json", sJson)
	if err != nil {
		ZAPS.Errorf("写入生产配置参数json文件 %s %v", "失败", err)
		return err
	}
	ZAPS.Infof("写入生产配置参数json文件 %s", "成功")

	return nil
}

func SetSN(sn string) {

	ProductParam.SN = sn
	SystemState.SN = sn
	_ = WriteProductParamToJson()
}

func GetSN() string {

	return ProductParam.SN
}

func GetName() string {

	return ProductParam.Name
}

func SetProduct(sn string, name string) {

	ProductParam.SN = sn
	SystemState.SN = sn

	ProductParam.Name = name
	SystemState.Name = name
	_ = WriteProductParamToJson()
}

func SetHardVer(hardVer string) {

	ProductParam.HardVer = hardVer
	_ = WriteProductParamToJson()
}

func GetHardVer() string {

	return ProductParam.HardVer
}
