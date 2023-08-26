package controllers

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"rnv-mmq/models"
	"rnv-mmq/services"
	"rnv-mmq/tools"
	"rnv-mmq/types"
	"rnv-mmq/wrongs"
)

type (
	TaskController struct{}
	// TaskStoreForm 任务表单（新建）
	TaskStoreForm struct {
		Name        string   `json:"name"`
		Targets     []string `json:"targets"`
		Description string   `json:"description"`
	}
	// TaskUpdateForm 任务表单（修改）
	TaskUpdateForm struct {
		Name        string `json:"name"`
		Target      string `json:"target"`
		Description string `json:"description"`
	}
)

// ShouldBind 绑定表单
func (receiver TaskStoreForm) ShouldBind(ctx *gin.Context) TaskStoreForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}

	if receiver.Name == "" {
		wrongs.ThrowValidate("任务名称必填")
	}
	if len(receiver.Targets) == 0 {
		wrongs.ThrowValidate("任务目标必填")
	}

	return receiver
}

// ShouldBind 绑定表单
func (receiver TaskUpdateForm) ShouldBind(ctx *gin.Context) TaskUpdateForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}

	if receiver.Name == "" {
		wrongs.ThrowValidate("任务名称必填")
	}
	if receiver.Target == "" {
		wrongs.ThrowValidate("任务目标必填")
	}

	return receiver
}

// NewTaskController 构造函数
func NewTaskController() *TaskController {
	return &TaskController{}
}

// Store 新建
func (TaskController) Store(ctx *gin.Context) {
	var (
		ret   *gorm.DB
		tasks []*models.TaskModel
	)

	// 表单
	form := (&TaskStoreForm{}).ShouldBind(ctx)

	wrongs.ThrowWhenIsRepeat(ret, "任务名称")

	// 新建
	tasks = make([]*models.TaskModel, len(form.Targets))
	for idx, target := range form.Targets {
		tasks[idx] = &models.TaskModel{
			GormModel:   models.GormModel{Uuid: uuid.NewV4().String()},
			Name:        form.Name,
			Target:      target,
			Description: form.Description,
			StatusCode:  "ORIGINAL",
		}
	}
	if ret = models.NewGorm().
		SetModel(models.TaskModel{}).
		GetDb("").
		Create(&tasks); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(types.MapStringToAny{"tasks": tasks}).ToGinResponse())
}

// Delete 删除
func (TaskController) Delete(ctx *gin.Context) {
	var (
		ret  *gorm.DB
		task models.TaskModel
	)

	// 查询
	ret = models.NewGorm().
		SetModel(models.TaskModel{}).
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 删除
	if ret := models.NewGorm().
		SetModel(models.TaskModel{}).
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Delete(&task); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (TaskController) Update(ctx *gin.Context) {
	var (
		ret  *gorm.DB
		task models.TaskModel
	)

	// 表单
	form := TaskUpdateForm{}.ShouldBind(ctx)

	// 查询
	ret = models.NewGorm().
		SetModel(models.TaskModel{}).
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 编辑
	task.Name = form.Name
	task.Target = form.Target
	task.Description = form.Description
	if ret = models.NewGorm().
		SetModel(models.TaskModel{}).
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&task); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(types.MapStringToAny{"task": task}).ToGinResponse())
}

// Detail 详情
func (TaskController) Detail(ctx *gin.Context) {
	var (
		ret  *gorm.DB
		task models.TaskModel
	)
	ret = models.NewGorm().
		SetModel(models.TaskModel{}).
		SetCtx(ctx).
		GetDbUseQuery("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(types.MapStringToAny{"task": task}).ToGinResponse())
}

func (TaskController) listUseQuery(ctx *gin.Context) *gorm.DB {
	return services.NewTaskService(services.BaseService{Model: models.NewGorm().SetModel(models.TaskModel{}), Ctx: ctx}).GetListByQuery()
}

// List 列表
func (receiver TaskController) List(ctx *gin.Context) {
	var tasks []*models.TaskModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				receiver.listUseQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&tasks)
					return types.MapStringToAny{"tasks": tasks}
				},
			).
			ToGinResponse(),
	)
}

// ListJdt jquery-dataTable后端分页数据
func (receiver TaskController) ListJdt(ctx *gin.Context) {
	var tasks []*models.TaskModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				receiver.listUseQuery(ctx),
				func(db *gorm.DB) map[string]interface{} {
					db.Find(&tasks)
					return types.MapStringToAny{"tasks": tasks}
				},
			).
			ToGinResponse(),
	)
}
