package middleware

import (
	"errors"
	"fmt"
	"gateway/setting"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 定义一个jwt对象
type JWT struct {
	// 声明签名信息
	SigningKey []byte
}

// 自定义有效载荷(这里采用自定义的Name和Email作为有效载荷的一部分)
type CustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	// StandardClaims结构体实现了Claims接口(Valid()函数)
	jwt.StandardClaims
}

// 构造用户表
type User struct {
	Id        int32  `gorm:"AUTO_INCREMENT"`
	Name      string `json:"username"`
	Pwd       string `json:"password"`
	Phone     int64  `gorm:"DEFAULT:0"`
	Email     string `gorm:"type:varchar(20);unique_index;"`
	CreatedAt *time.Time
	UpdateTAt *time.Time
}

// LoginReq请求参数
type LoginReq struct {
	Name string `json:"username"`
	Pwd  string `json:"password"`
}

// 登陆结果
type LoginResultTemplate struct {
	Token       string                        `json:"token"`
	Name        string                        `json:"username"`
	Permissions []setting.PermissionsTemplate `json:"permissions"`
}

var (
	TokenExpired error = errors.New("Token is expired")

	LoginResult LoginResultTemplate
)

// 初始化jwt对象
func NewJWT() *JWT {
	return &JWT{
		[]byte("wdtACU300"),
	}
}

// 调用jwt-go库生成token
// 指定编码的算法为jwt.SigningMethodHS256
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#Token
	// 返回一个token的结构体指针
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// token解码
func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	// https://gowalker.org/github.com/dgrijalva/jwt-go#ParseWithClaims
	// 输入用户自定义的Claims结构体对象,token,以及自定义函数来解析token字符串为jwt的Token结构体指针
	// Keyfunc是匿名函数类型: type Keyfunc func(*Token) (interface{}, error)
	// func ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		// https://gowalker.org/github.com/dgrijalva/jwt-go#ValidationError
		// jwt.ValidationError 是一个无效token的错误结构
		if ve, ok := err.(*jwt.ValidationError); ok {
			// ValidationErrorMalformed是一个uint常量，表示token不可用
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token不可用")
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token过期")
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("无效的token")
			} else {
				return nil, fmt.Errorf("token不可用")
			}

		}
	}

	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token无效")

}

// 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "1",
				"message": "请求未携带token，无权限访问",
				"data":    "",
			})
			c.Abort()
			return
		}
		// 初始化一个JWT对象实例，并根据结构体方法来解析token
		j := NewJWT()
		// 解析token中包含的相关信息(有效载荷)
		claims, err := j.ParserToken(token)
		if err != nil {
			// token过期
			if err.Error() == "token不可用" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    "1",
					"message": "token不可用",
					"data":    "",
				})
				c.Abort()
				return
			} else if err.Error() == "token过期" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    "1",
					"message": "token授权已过期，请重新申请授权",
					"data":    "",
				})
				c.Abort()
				return
			}
			// 其他错误
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "1",
				"message": err.Error(),
				"data":    "",
			})
			c.Abort()
			return
		}

		// 将解析后的有效载荷claims重新写入gin.Context引用对象中
		c.Set("claims", claims)
	}
}

// LoginCheck验证
func LoginCheck(login LoginReq) (bool, User, error) {
	userData := User{}
	userExist := false

	for _, v := range setting.PolicyWeb {
		if v.Role == login.Name && v.Password == login.Pwd {
			userExist = true
			userData.Name = login.Name
			userData.Email = ""
		}
	}

	if !userExist {
		return userExist, userData, fmt.Errorf("登陆信息有误")
	}
	return userExist, userData, nil
}

// token生成器
// md 为上面定义好的middleware中间件
func GenerateToken(c *gin.Context, user User) {
	// 构造SignKey: 签名和解签名需要使用一个值
	j := NewJWT()

	// 构造用户claims信息(负荷)
	claims := CustomClaims{
		user.Name,
		user.Email,
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),        // 签名生效时间
			ExpiresAt: time.Now().Unix() + 3600, // 签名过期时间
			Issuer:    "iot",                    // 签名颁发者
		},
	}

	// 根据claims生成token对象
	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    "1",
			"message": err.Error(),
			"data":    "",
		})
		return
	}

	for _, v := range setting.PolicyWeb {
		if v.Role == claims.Name {
			LoginResult.Permissions = v.Policy
		}
	}

	data := LoginResultTemplate{
		Name:        user.Name,
		Token:       token,
		Permissions: LoginResult.Permissions,
	}
	LoginResult = data
	c.JSON(http.StatusOK, gin.H{
		"code":    "0",
		"message": "登录成功",
		"data":    data,
	})
	return

}
