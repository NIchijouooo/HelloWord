package middleware

import (
	"net/http"
	"gateway/setting"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

var Enforcer *casbin.Enforcer

func PrivilegeInit() error {
	var err error

	Enforcer, err = casbin.NewEnforcer("./config/model.conf", "./config/policy.csv")
	if err != nil {
		return err
	}

	setting.ReadPolicyGoFromJson()
	err = AddPolicies()
	if err != nil {
		return err
	}

	setting.ReadPolicyWebFromJson()
	for _, v := range setting.PolicyWeb {
		_ = setting.AddAccountParam(v.Role, v.Password)
	}
	setting.ZAPS.Debug("权限配置json文件格式化成功")
	return err
}

func AddPolicies() error {
	for _, v := range setting.PolicyGo {
		_, _ = Enforcer.AddPolicy(v.Sub, v.Obj, v.Act)
	}

	return nil
}

func Privilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		cValue, exists := c.Get("claims")
		if exists == false {
			c.JSON(http.StatusOK, gin.H{
				"code":    "1",
				"message": "token获取失败",
				"data":    "",
			})
			c.Abort()
			return
		}
		claims := cValue.(*CustomClaims)
		name := claims.Name
		path := c.Request.URL.Path
		method := c.Request.Method
		setting.ZAPS.Infof("http请求:role[%v] api[%v] method[%v]", name, path, method)
		isPass, err := Enforcer.Enforce(name, path, method)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    "1",
				"message": "权限检查错误",
				"data":    "",
			})
			c.Abort()
			return
		}
		if isPass == false {
			c.JSON(http.StatusOK, gin.H{
				"code":    "1",
				"message": "无权限访问",
				"data":    "",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
