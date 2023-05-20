package setting

import (
	"encoding/json"
	"errors"
	"gateway/utils"
)

type AccountParamTemplate struct {
	Role     string `json:"role"`
	Password string `json:"password"`
}

type PolicyGoTemplate struct {
	Sub string `json:"sub"`
	Obj string `json:"obj"`
	Act string `json:"act"`
}

type PolicyTemplate struct {
	Role     string                `json:"role"`
	Password string                `json:"password"`
	Policy   []PermissionsTemplate `json:"policy"`
}

type PermissionsTemplate struct {
	Name     string                       `json:"name"`
	Path     string                       `json:"path"`
	Meta     MetaTemplate                 `json:"meta"`
	Children []PermissionChildrenTemplate `json:"children"`
}

type MetaTemplate struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}

type PermissionChildrenTemplate struct {
	Name     string             `json:"name"`
	Path     string             `json:"path"`
	Meta     MetaTemplate       `json:"meta"`
	Children []ChildrenTemplate `json:"children"`
}

type ChildrenTemplate struct {
	Name string       `json:"name"`
	Path string       `json:"path"`
	Meta MetaTemplate `json:"meta"`
}

var PolicyGo = make([]PolicyGoTemplate, 0)
var PolicyWeb = make([]PolicyTemplate, 0)

var AccountParam = make(map[string]AccountParamTemplate)

func AddAccountParam(role, password string) error {
	account := AccountParamTemplate{
		Role:     role,
		Password: password,
	}

	_, ok := AccountParam[role]
	if !ok {
		AccountParam[role] = account
		return nil
	}
	return errors.New("用户名和密码已经存在")
}

func ModifyAccountParam(role, oldPwd, newPwd string) error {
	account := AccountParamTemplate{
		Role:     role,
		Password: newPwd,
	}

	roleParam, ok := AccountParam[role]
	if !ok {
		ZAPS.Error(errors.New("用户名和密码不存在"))
		return errors.New("用户名和密码不存在")
	}

	if roleParam.Password != oldPwd {
		ZAPS.Error(errors.New("旧密码错误"))
		return errors.New("旧密码错误")
	}

	AccountParam[role] = account

	index := -1
	for k, v := range PolicyWeb {
		if v.Role == role {
			index = k
			PolicyWeb[k].Password = newPwd
		}
	}
	if index != -1 {
		WritePolicyWebToJson()
	}

	return nil
}

func GetAccountParam() []AccountParamTemplate {

	account := make([]AccountParamTemplate, 0)
	for _, v := range AccountParam {
		account = append(account, v)
	}
	return account
}

func ReadPolicyGoFromJson() bool {

	data, err := utils.FileRead("./config/policyGo.json")
	if err != nil {
		ZAPS.Debugf("权限配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &PolicyGo)
	if err != nil {
		ZAPS.Errorf("权限配置json文件格式化失败")
		return false
	}
	ZAPS.Debug("读取Go权限配置json文件成功")

	return true
}

func WritePolicyGoToJson() {

	utils.DirIsExist("./config")
	sJson, _ := json.Marshal(PolicyGo)
	err := utils.FileWrite("./config/policyGo.json", sJson)
	if err != nil {
		ZAPS.Errorf("权限配置json文件写入失败")
		return
	}
	ZAPS.Debugf("权限配置json文件写入成功")
}

func ReadPolicyWebFromJson() bool {

	data, err := utils.FileRead("./config/policyWeb.json")
	if err != nil {
		ZAPS.Debugf("权限配置json文件读取失败 %v", err)
		return false
	}
	err = json.Unmarshal(data, &PolicyWeb)
	if err != nil {
		ZAPS.Errorf("权限配置json文件格式化失败")
		return false
	}
	ZAPS.Debug("读取Web权限配置json文件成功")

	return true
}

func WritePolicyWebToJson() {

	utils.DirIsExist("./config")
	sJson, _ := json.Marshal(PolicyWeb)
	err := utils.FileWrite("./config/policyWeb.json", sJson)
	if err != nil {
		ZAPS.Errorf("权限配置json文件写入失败")
		return
	}
	ZAPS.Debugf("权限配置json文件写入成功")
}
