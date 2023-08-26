package middlewares

import (
	"fmt"
	"reflect"
	"rnv-mmq/models"
	"rnv-mmq/tools"
	"rnv-mmq/types"
	"rnv-mmq/wrongs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckAuthorization 检查Jwt是否合法
func CheckAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取令牌
		split := strings.Split(tools.GetJwtFromHeader(ctx), " ")
		if len(split) != 2 {
			wrongs.ThrowUnAuth("令牌格式错误")
		}
		tokenType := split[0]
		token := split[1]

		var (
			user models.UserModel
			ret  *gorm.DB
		)
		if token == "" {
			wrongs.ThrowUnAuth("令牌不存在")
		} else {
			switch tokenType {
			case "JWT":
				claims, err := tools.ParseJwt(token)

				// 判断令牌是否有效
				if err != nil {
					wrongs.ThrowUnAuth("令牌解析失败")
				} else if time.Now().Unix() > claims.ExpiresAt {
					wrongs.ThrowUnAuth("令牌过期")
				}

				// 判断用户是否存在
				if reflect.DeepEqual(claims, tools.Claims{}) {
					wrongs.ThrowUnAuth("令牌解析失败：用户不存在")
				}

				// 获取用户信息
				ret = models.NewGorm().SetModel(models.UserModel{}).GetDb("").Where("uuid", claims.Uuid).First(&user)
				wrongs.ThrowWhenIsEmpty(ret, fmt.Sprintf("令牌指向用户(JWT) %s %v ", token, claims))
			case "AU":
				ret = models.NewGorm().SetModel(models.UserModel{}).SetWheres(types.MapStringToAny{"open_id": token}).GetDb("").First(&user)
				wrongs.ThrowWhenIsEmpty(ret, fmt.Sprintf("令牌指向用户(AU) %s", token))
			default:
				wrongs.ThrowForbidden("权鉴认证方式不支持")
			}

			ctx.Set(types.CURRENT_USER, user) // 设置用户信息
		}

		ctx.Next()
	}
}
