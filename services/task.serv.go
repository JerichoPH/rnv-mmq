package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"rnv-mmq/models"
	"rnv-mmq/wrongs"
	"strings"
)

type (
	// TaskService 任务服务
	TaskService struct{ BaseService }

	// TaskQueryForm 任务查询表单
	TaskQueryForm struct {
		CreatedAtBetween  string `json:"created_at_between,omitempty"`
		createdAtOriginal string
		createdAtFinished string
		UpdatedAtBetween  string `json:"updated_at_between,omitempty"`
		updatedAtOriginal string
		updatedAtFinished string
		DeletedAtBetween  string `json:"deleted_at_between,omitempty"`
		deletedAtOriginal string
		deletedAtFinished string
		Uuid              string   `json:"uuid,omitempty"`
		Uuids             []string `json:"uuids,omitempty"`
		Sort              uint64   `json:"sort,omitempty"`
		Name              string   `json:"name,omitempty"`
		NameLike          string   `json:"name_like,omitempty"`
		Target            string   `json:"target,omitempty"`
		Targets           []string `json:"targets,omitempty"`
		Description       string   `json:"description,omitempty"`
		StatusCode        string   `json:"status_code,omitempty"`
		StatusCodes       []string `json:"status_codes,omitempty"`
		BusinessType      string   `json:"business_type,omitempty"`
		ContentFileUuid   string   `json:"content_file_uuid,omitempty"`
		ContentFileUuids  []string `json:"content_file_uuids"`
		Preloads          []string `json:"preloads"`
	}
	// TaskLogService 任务日志服务
	TaskLogService struct{ BaseService }
)

// ShouldBind 表单绑定
func (receiver TaskQueryForm) ShouldBind(ctx *gin.Context) TaskQueryForm {
	var (
		err                             error
		createdAt, updatedAt, deletedAt = make([]string, 2), make([]string, 2), make([]string, 2)
	)
	if err = ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate("表单验证失败：%s", err.Error())
	}

	if receiver.CreatedAtBetween != "" {
		createdAt = strings.Split("~", receiver.CreatedAtBetween)
		receiver.createdAtOriginal, receiver.createdAtFinished = createdAt[0], createdAt[1]
	}
	if receiver.UpdatedAtBetween != "" {
		updatedAt = strings.Split("~", receiver.UpdatedAtBetween)
		receiver.updatedAtOriginal, receiver.updatedAtFinished = updatedAt[0], updatedAt[1]
	}
	if receiver.DeletedAtBetween != "" {
		deletedAt = strings.Split("~", receiver.DeletedAtBetween)
		receiver.deletedAtOriginal, receiver.deletedAtFinished = deletedAt[0], deletedAt[1]
	}

	return receiver
}

// NewTaskService 构造函数
func NewTaskService(baseService BaseService) *TaskService {
	return &TaskService{BaseService: baseService}
}

// GetListByQuery 根据Query获取列表
func (receiver TaskService) GetListByQuery() *gorm.DB {
	return (receiver.Model).
		SetWheresEqual("name", "target", "content_file_uuid").
		SetWheresFuzzy(map[string]string{ // 模糊查询字段
			"description": "description like ?",
		}).
		SetWheresDateBetween("created_at", "updated_at").
		SetWheresExtraExists(map[string]func([]string, *gorm.DB) *gorm.DB{
			"targets[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("target in ?", values)
			},
			"content_file_uuids[]": func(strings []string, db *gorm.DB) *gorm.DB {
				return db.Where("content_file_uuid in ?", strings)
			},
		}).
		SetCtx(receiver.Ctx).
		GetDbUseQuery("").
		Table(fmt.Sprintf("%s as t", new(models.TaskModel).TableName()))
}

// GetListByPostParam 根据Post参数获取列表
func (receiver TaskService) GetListByPostParam() *gorm.DB {
	model := receiver.Model
	db := model.GetDb("")
	form := TaskQueryForm{}.ShouldBind(receiver.Ctx)

	if form.createdAtOriginal != "" && form.createdAtFinished != "" {
		db = db.Where("created_at between ? and ?", form.createdAtOriginal, form.createdAtFinished)
	}

	if form.updatedAtOriginal != "" && form.updatedAtFinished != "" {
		db = db.Where("updated_at between ? and ?", form.updatedAtOriginal, form.updatedAtFinished)
	}

	if form.deletedAtOriginal != "" && form.deletedAtFinished != "" {
		db = db.Where("deleted_at between ? and ?", form.deletedAtOriginal, form.deletedAtFinished)
	}

	if form.Uuid != "" {
		db = db.Where("uuid = ?", form.Uuid)
	}

	if len(form.Uuids) > 0 {
		db = db.Where("uuid in ?", form.Uuids)
	}

	if form.Name != "" {
		db = db.Where("name = ?", form.Name)
	}

	if form.NameLike != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", form.Name))
	}

	if form.Description != "" {
		db = db.Where("description like ?", fmt.Sprintf("%%%s%%", form.Description))
	}

	if form.StatusCode != "" {
		db = db.Where("status_code = ?", form.StatusCode)
	}

	if len(form.StatusCodes) > 0 {
		db = db.Where("status_code in ?", form.StatusCodes)
	}

	if form.BusinessType != "" {
		db = db.Where("business_type = ?", form.BusinessType)
	}

	if form.ContentFileUuid != "" {
		db = db.Where("content_file_uuid = ?", form.ContentFileUuid)
	}

	if len(form.ContentFileUuids) > 0 {
		db = db.Where("content_file_uuid in ?", form.ContentFileUuids)
	}

	if len(form.Preloads) > 0 {
		for _, preload := range form.Preloads {
			db = db.Preload(preload)
		}
	}

	return db
}

// NewTaskLogService 构造函数
func NewTaskLogService(baseService BaseService) *TaskLogService {
	return &TaskLogService{BaseService: baseService}
}

// GetListByQuery 根据Query获取列表
func (receiver TaskLogService) GetListByQuery() *gorm.DB {
	return (receiver.Model).
		SetWheresEqual("task_uuid").
		SetWheresFuzzy(map[string]string{ // 模糊查询字段
			"task_uuid": "task_uuid like ?",
		}).
		SetWheresDateBetween("created_at", "updated_at").
		SetWheresExtraExists(map[string]func([]string, *gorm.DB) *gorm.DB{
			"task_uuids[]": func(values []string, db *gorm.DB) *gorm.DB {
				return db.Where("task_uuid in ?", values)
			},
		}).
		SetCtx(receiver.Ctx).
		GetDbUseQuery("").
		Table(fmt.Sprintf("%s as tl", new(models.TaskLogModel).TableName()))
}
