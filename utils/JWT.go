package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

type MyClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// NotNeedToken 这里填入不需要token验证的路径
var NotNeedToken = make(map[string]struct{}) //这里用空结构体加map实现了一个集合（空结构体不占内存）

func init() {
	//将不需要验证token的路径添加到集合中
	NotNeedToken["/douyin/user/register/"] = struct{}{}
	NotNeedToken["/douyin/user/login/"] = struct{}{}
	NotNeedToken["/douyin/feed/"] = struct{}{}
}

// TokenExpireDuration 设置过期时间
const TokenExpireDuration = time.Hour * 24 * 180

// Secret 密钥,自行设定
var Secret = []byte("DouYin HelloWorld")

// GenToken 传入用户ID生成token，有效期半年
func GenToken(USerID int64) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		USerID, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间(半年)
			Issuer:    "ADuiDuiDui",                               // 签发人（啊对对队）
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(Secret)
}

// JwtVerify 验证token（这个是作为一个中间件被添加到全局路由中了，所有的路由都会执行这个路由）
func JwtVerify(c *gin.Context) {
	//过滤是否验证token
	currentRouter := c.Request.RequestURI //获取当前路由 "
	index := strings.Index(currentRouter, "?")
	if index != -1 {
		currentRouter = currentRouter[0:index] //去掉query参数
	}
	_, ok := NotNeedToken[currentRouter]
	if ok { //不需要验证token
		log.Println(currentRouter + "：当前路径不需要token验证")
		return
	}
	token := c.Query("token")
	if token == "" {
		token = c.PostForm("token")
	}
	if token == "" {
		panic("未携带token")
	}
	claims, err := ParseToken(token)
	if err != nil {
		panic("invalid token")
	}
	c.Set("UserID", claims.UserID) //后续可以使用c.get(UserID)获取到用户ID
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	// 校验token
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
