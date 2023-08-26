package apiRoute

import (
	"github.com/gin-gonic/gin"
	"rnv-mmq/middlewares"
)

// TaskRouter 路由
type TaskRouter struct{}

// NewTaskRouter 构造函数
func NewTaskRouter() TaskRouter {
	return TaskRouter{}
}

// Load 加载路由
func (TaskRouter) Load(engine *gin.Engine) {
	r := engine.Group(
		"api/task",
		middlewares.CheckAuthorization(),
		// middlewares.CheckPermission(),
	)
	{
		// 新建
		r.POST("", controllers.New().Store)

		// 删除
		r.DELETE("/:uuid", controllers.New().Delete)

		// 编辑
		r.PUT("/:uuid", controllers.New().Update)

		// 详情
		r.GET("/:uuid", controllers.New().Detail)

		// 列表
		r.GET("", controllers.New().List)

		// jquery-dataTable数据列表
		r.GET(".jdt", controllers.New().ListJdt)
	}
}
