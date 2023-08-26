package apiRoute

import (
	"github.com/gin-gonic/gin"
	"rnv-mmq/controllers"
	"rnv-mmq/middlewares"
)

// TaskRouter 路由
type TaskRouter struct{}

// NewTaskRouter 构造函数
func NewTaskRouter() TaskRouter { return TaskRouter{} }

// Load 加载路由
func (TaskRouter) Load(engine *gin.Engine) {
	taskRoute := engine.Group(
		"api/task",
		middlewares.CheckAuthorization(),
		// middlewares.CheckPermission(),
	)
	{
		// 新建
		taskRoute.POST("", controllers.NewTaskController().Store)
		// 删除
		taskRoute.DELETE("/:uuid", controllers.NewTaskController().Delete)
		// 编辑
		taskRoute.PUT("/:uuid", controllers.NewTaskController().Update)
		// 详情
		taskRoute.GET("/:uuid", controllers.NewTaskController().Detail)
		// 列表
		taskRoute.GET("", controllers.NewTaskController().List)
		// jquery-dataTable数据列表
		taskRoute.GET(".jdt", controllers.NewTaskController().ListJdt)
	}

	taskLogRoute := engine.Group(
		"api/task-log",
		middlewares.CheckAuthorization(),
		// middlewares.CheckPermission(),
	)
	{
		taskLogRoute.GET("", func(context *gin.Context) {
			context.JSON(200, gin.H{"message": "task-log"})
		})
	}
}
