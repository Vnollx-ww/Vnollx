package middlerware

import (
	"Vnollx/kitex_gen/base"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

var Mysecret = []byte("vnollxvnollx")

type MyClaims struct {
	UserId int64 `json:"UserId"`
	jwt.RegisteredClaims
}

func NewToken(userid int64) (string, error) {
	c := MyClaims{
		userid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("JWT_ISSUER"),
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(Mysecret)
}
func ParseToken(Token string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(Token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Mysecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效token")
}
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 定义一个结构体，用来解析请求体中的 token
		var requestBody struct {
			Token string // 假设token在请求体中是 "token" 字段
		}
		// 解析请求体
		if err := ctx.ShouldBind(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, base.BaseResponse{
				Code: http.StatusBadRequest,
				Msg:  "请求体格式有误",
			})
			ctx.Abort()
			return
		}
		// 获取token
		token := requestBody.Token
		// 如果token为空，返回401 Unauthorized
		if token == "" {
			ctx.JSON(http.StatusBadRequest, base.BaseResponse{
				Code: http.StatusBadRequest,
				Msg:  "请求体中token为空",
			})
			ctx.Abort()
			return
		}
		// 解析token
		mc, err := ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, base.BaseResponse{
				Code: http.StatusBadRequest,
				Msg:  "无效的token",
			})
			ctx.Abort()
			return
		}
		// 将用户名保存到上下文中
		ctx.Set("uid", mc.UserId)
		// 继续处理请求
		ctx.Next()
	}
}
