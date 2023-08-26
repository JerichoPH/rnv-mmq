package tools

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"rnv-mmq/wrongs"
	"time"
)

var jwtSecret = []byte("jwt_key") // 加密密钥

// Claims Jwt 表单
type Claims struct {
	Uuid     string `json:"open_id"`
	Nickname string `json:"nickname"`
	jwt.StandardClaims
}

// GenerateJwt 生成Jwt
func GenerateJwt(
	uuid string,
	nickname string,
) (string, error) {
	// 设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(168 * time.Hour)

	claims := Claims{
		Uuid:     uuid,
		Nickname: nickname,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			Issuer: "rnv-mmq",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseJwt 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseJwt(token string) (*Claims, error) {
	// 用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// GetJwtFromHeader 从header中获取jwt
func GetJwtFromHeader(ctx *gin.Context) string {
	tokens := ctx.Request.Header["Authorization"]

	if len(tokens) == 0 {
		wrongs.ThrowUnAuth("令牌不存在")
	}
	return tokens[0]
}
