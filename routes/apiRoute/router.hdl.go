package apiRoute

import (
	"github.com/gin-gonic/gin"
)

// RouterHandle 分组路由
type RouterHandle struct{}

// Register 组册路由
func (RouterHandle) Register(engine *gin.Engine) {
	NewAuthorizationRouter().Load(engine) // 权鉴
	NewUserRouter().Load(engine)          // 用户
}
