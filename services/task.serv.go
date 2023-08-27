package services

import (
	"fmt"
	"gorm.io/gorm"
	"rnv-mmq/models"
)

type (
	// TaskService 任务服务
	TaskService struct{ BaseService }
	// TaskLogService 任务日志服务
	TaskLogService struct{ BaseService }
)

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
