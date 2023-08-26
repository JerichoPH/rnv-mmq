package apiRoute

import (
	"github.com/gin-gonic/gin"
	"rnv-mmq/controllers"
	"rnv-mmq/middlewares"
)

// UserRouter 用户路由
type UserRouter struct{}

// NewUserRouter 构造函数
func NewUserRouter() UserRouter {
	return UserRouter{}
}

// Load 加载路由
func (UserRouter) Load(engine *gin.Engine) {
	r := engine.Group(
		"api",
		middlewares.CheckAuthorization(),
		// middlewares.CheckPermission(),
	)
	{
		// 新建
		r.POST("account/store", func(ctx *gin.Context) { controllers.NewUserController().Store(ctx) })

		// 删除
		r.DELETE("account/:uuid", func(ctx *gin.Context) { controllers.NewUserController().Delete(ctx) })

		// 编辑
		r.PUT("account/:uuid/update", func(ctx *gin.Context) { controllers.NewUserController().Update(ctx) })

		// 详情
		r.GET("account/:uuid", func(ctx *gin.Context) { controllers.NewUserController().Detail(ctx) })

		// 列表
		r.GET("account", func(ctx *gin.Context) { controllers.NewUserController().List(ctx) })

		// jquery-dataTable分页列表
		r.GET("account.jdt", func(ctx *gin.Context) { controllers.NewUserController().ListJdt(ctx) })
	}
}
