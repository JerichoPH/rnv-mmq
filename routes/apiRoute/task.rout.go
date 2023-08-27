package apiRoute

import (
	"rnv-mmq/controllers"
	"rnv-mmq/middlewares"

	"github.com/gin-gonic/gin"
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
		taskRoute.POST("store", controllers.NewTaskController().Store)
		// 删除
		taskRoute.POST("/:uuid/delete", controllers.NewTaskController().Delete)
		// 编辑
		// taskRoute.PUT("/:uuid", controllers.NewTaskController().Update)
		// 详情
		taskRoute.POST("/:uuid/detail", controllers.NewTaskController().Detail)
		// 列表
		taskRoute.POST("list", controllers.NewTaskController().List)
		// jquery-dataTable数据列表
		taskRoute.POST("list.jdt", controllers.NewTaskController().ListJdt)
		// 标记执行
		taskRoute.POST("/:uuid/process", controllers.NewTaskController().Process)
		// 标记完成
		taskRoute.POST("/:uuid/finish", controllers.NewTaskController().Finish)
		// 标记失败
		taskRoute.POST("/:uuid/fail", controllers.NewTaskController().Fail)
		// 标记取消
		taskRoute.POST("/:uuid/cancel", controllers.NewTaskController().Cancel)
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
