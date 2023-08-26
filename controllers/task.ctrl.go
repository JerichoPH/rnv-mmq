package controllers

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"rnv-mmq/models"
	"rnv-mmq/services"
	"rnv-mmq/tools"
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

	// 新建
	tasks = make([]*models.TaskModel, len(form.Targets))
	for idx, target := range form.Targets {
		tasks[idx] = &models.TaskModel{
			GormModel:   models.GormModel{Uuid: uuid.NewV4().String()},
			Name:        form.Name,
			Target:      target,
			Description: form.Description,
			StatusCode:  models.TaskModelStatusCodeOriginal,
		}
	}
	if ret = models.NewGorm().
		SetModel(models.TaskModel{}).
		GetDb("").
		Create(&tasks); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(map[string]any{"tasks": tasks}).ToGinResponse())
}

// Delete 删除
func (TaskController) Delete(ctx *gin.Context) {
	var (
		ret        *gorm.DB
		task       models.TaskModel
		canIDelete = false
		reason     string
	)

	// 查询
	ret = models.NewTaskModelGorm().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 判断任务是否可删除
	if canIDelete, reason = task.CanIDelete(); !canIDelete {
		wrongs.ThrowForbidden(reason)
	}

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

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(map[string]any{"task": task}).ToGinResponse())
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

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(map[string]any{"task": task}).ToGinResponse())
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
				func(db *gorm.DB) map[string]any {
					db.Find(&tasks)
					return map[string]any{"tasks": tasks}
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
				func(db *gorm.DB) map[string]any {
					db.Find(&tasks)
					return map[string]any{"tasks": tasks}
				},
			).
			ToGinResponse(),
	)
}

// PostProcess 标记执行
func (TaskController) PostProcess(ctx *gin.Context) {
	var (
		task        *models.TaskModel
		ret         *gorm.DB
		canIProcess = false
		reason      string
	)

	ret = models.NewTaskModelGorm().GetDb("").Where("uuid = ?", ctx.Param("uuid")).First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 判断任务是否可以【标记执行】
	if canIProcess, reason = task.CanIProcess(); !canIProcess {
		wrongs.ThrowForbidden(reason)
	}

	// 执行任务
	task.StatusCode = models.TaskModelStatusCodeProcessing
	if ret = models.NewTaskModelGorm().GetDb("").Where("uuid = ?", task.Uuid).Save(&task); ret.Error != nil {
		wrongs.ThrowForbidden("【标记执行】失败：%s", ret.Error.Error())
	}

	// 写入执行日志
	if ret = models.NewTaskLogModelGorm().GetDb("").Create(&models.TaskLogModel{
		GormModel: models.GormModel{Uuid: uuid.NewV4().String()},
		Name:      "执行",
		TaskUuid:  task.Uuid,
	}); ret.Error != nil {
		wrongs.ThrowForbidden("写入【标记执行】日志失败：%s", ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Blank().ToGinResponse())
}

// PostFinish 标记完成
func (TaskController) PostFinish(ctx *gin.Context) {
	var (
		task       *models.TaskModel
		ret        *gorm.DB
		canIFinish = false
		reason     string
	)

	ret = models.NewTaskModelGorm().GetDb("").Where("uuid = ?", ctx.Param("uuid")).First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 判断任务是否可以【标记完成】
	if canIFinish, reason = task.CanIFinish(); !canIFinish {
		wrongs.ThrowForbidden(reason)
	}

	// 执行任务
	task.StatusCode = models.TaskModelStatusCodeFinished
	if ret = models.NewTaskModelGorm().GetDb("").Where("uuid = ?", task.Uuid).Save(&task); ret.Error != nil {
		wrongs.ThrowForbidden("【标记完成】失败：%s", ret.Error.Error())
	}

	// 写入执行日志
	if ret = models.NewTaskLogModelGorm().GetDb("").Create(&models.TaskLogModel{
		GormModel: models.GormModel{Uuid: uuid.NewV4().String()},
		Name:      "完成",
		TaskUuid:  task.Uuid,
	}); ret.Error != nil {
		wrongs.ThrowForbidden("写入【标记完成】日志失败：%s", ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Blank().ToGinResponse())
}

// PostFail 标记失败
func (TaskController) PostFail(ctx *gin.Context) {
	var (
		task     *models.TaskModel
		ret      *gorm.DB
		canIFail = false
		reason   string
	)

	ret = models.
		NewTaskModelGorm().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 判断任务是否可以【标记失败】
	if canIFail, reason = task.CanIFail(); !canIFail {
		wrongs.ThrowForbidden(reason)
	}

	// 取消任务
	task.StatusCode = models.TaskModelStatusCodeFailed
	if ret = models.NewTaskModelGorm().GetDb("").Where("uuid = ?", task.Uuid).Save(&task); ret.Error != nil {
		wrongs.ThrowForbidden("【标记失败】失败：%s", ret.Error.Error())
	}

	// 写入执行日志
	if ret = models.NewTaskLogModelGorm().GetDb("").Create(&models.TaskLogModel{
		GormModel: models.GormModel{Uuid: uuid.NewV4().String()},
		Name:      "失败",
		TaskUuid:  task.Uuid,
	}); ret.Error != nil {
		wrongs.ThrowForbidden("写入【标记失败】日志失败：%s", ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Blank().ToGinResponse())
}

// PostCancel 标记取消
func (TaskController) PostCancel(ctx *gin.Context) {
	var (
		task       *models.TaskModel
		ret        *gorm.DB
		canICancel = false
		reason     string
	)

	ret = models.
		NewTaskModelGorm().
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		First(&task)
	wrongs.ThrowWhenIsEmpty(ret, "任务")

	// 判断任务是否可以【标记取消】
	if canICancel, reason = task.CanICancel(); !canICancel {
		wrongs.ThrowForbidden(reason)
	}

	// 取消任务
	task.StatusCode = models.TaskModelStatusCodeCancel
	if ret = models.NewTaskModelGorm().GetDb("").Where("uuid = ?", task.Uuid).Save(&task); ret.Error != nil {
		wrongs.ThrowForbidden("【标记取消】失败：%s", ret.Error.Error())
	}

	// 写入执行日志
	if ret = models.NewTaskLogModelGorm().GetDb("").Create(&models.TaskLogModel{
		GormModel: models.GormModel{Uuid: uuid.NewV4().String()},
		Name:      "取消",
		TaskUuid:  task.Uuid,
	}); ret.Error != nil {
		wrongs.ThrowForbidden("写入【标记取消】日志失败：%s", ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Blank().ToGinResponse())
}
