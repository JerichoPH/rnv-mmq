package models

type (
	// TaskModel 任务模型
	TaskModel struct {
		GormModel
		Name             string               `gorm:"type:varchar(128);not null;default:'';comment:任务名称;" json:"name,omitempty"`
		Target           string               `gorm:"type:varchar(128);not null;default:'';comment:任务目标;" json:"target,omitempty"`
		Description      string               `gorm:"type:varchar(128);not null;default:'';comment:任务说明;" json:"description,omitempty"`
		StatusCode       TaskModelStatusCodes `gorm:"type:enum('ORIGINAL','PROCESSING','FINISHED','FAILED','CANCEL');not null;default:'ORIGINAL';comment:任务状态码" json:"status_code,omitempty"`
		StatusText       string               `gorm:"->;type:varchar(128) as ((case status_code when 'ORIGINAL' then '未执行' when 'PROCESSING' then '执行中' when 'FINISHED' then '已完成' when 'FAILED' then '失败' when 'CANCEL' then '已取消' else '' end));comment:任务状态文本" json:"status_text,omitempty"`
		RequestFileUuid  string               `gorm:"index;type:char(36);not null;comment:所属请求文件uuid" json:"request_file_uuid,omitempty"`
		RequestFile      *FileModel           `gorm:"foreignKey:request_file_uuid;references:uuid;comment:所属文件;" json:"request_file,omitempty"`
		ResponseFileUuid string               `gorm:"index;type:char(36);not null;comment:所属响应文件uuid" json:"response_file_uuid,omitempty"`
		ResponseFile     *FileModel           `gorm:"foreignKey:response_file_uuid;references:uuid;comment:所属响应文件;" json:"response_file,omitempty"`
		TaskLogs         []*TaskLogModel      `gorm:"foreignKey:task_uuid;references:uuid;comment:相关任务日志;" json:"task_logs,omitempty"`
	}

	// TaskModelStatusCodes 任务状态代码
	TaskModelStatusCodes string

	// TaskLogModel 任务日志模型
	TaskLogModel struct {
		GormModel
		Name     string    `gorm:"type:varchar(128);not null;default:'';comment:任务日志名称;" json:"name,omitempty"`
		TaskUuid string    `gorm:"index;type:char(36);not null;comment:所属任务uuid;" json:"task_uuid,omitempty"`
		Task     TaskModel `gorm:"foreignKey:task_uuid;references:uuid;comment:所属任务;" json:"task,omitempty"`
	}
)

// NewTaskModelGorm 新建Gorm模型
func NewTaskModelGorm() *GormModel {
	return NewGorm().SetModel(TaskModel{})
}

const (
	TaskModelStatusCodeOriginal   TaskModelStatusCodes = "ORIGINAL"
	TaskModelStatusCodeProcessing TaskModelStatusCodes = "PROCESSING"
	TaskModelStatusCodeFinished   TaskModelStatusCodes = "FINISHED"
	TaskModelStatusCodeFailed     TaskModelStatusCodes = "FAILED"
	TaskModelStatusCodeCancel     TaskModelStatusCodes = "CANCEL"
)

// TableName 任务表名称
func (TaskModel) TableName() string {
	return "tasks"
}

// CanIProcess 判断任务是否可以【标记执行】
func (receiver TaskModel) CanIProcess() (bool, string) {
	var (
		canIProcess bool
		reason      string
	)

	switch TaskModelStatusCodes(receiver.StatusCode) {
	case TaskModelStatusCodeOriginal:
	case TaskModelStatusCodeFailed:
		canIProcess = true
	default:
		reason = "【未开始】、【失败】的任务才可以【标记开始】"
	}

	return canIProcess, reason
}

// CanIFinish 判断任务是否可以【标记完成】
func (receiver TaskModel) CanIFinish() (bool, string) {
	var (
		canIFinish bool
		reason     string
	)

	switch receiver.StatusCode {
	case TaskModelStatusCodeProcessing:
		canIFinish = true
	default:
		reason = "只有【执行中】任务可以【标记完成】"
	}

	return canIFinish, reason
}

// CanIFail 判断任务是否可以【标记失败】
func (receiver TaskModel) CanIFail() (bool, string) {
	var (
		canIFail = false
		reason   string
	)

	switch receiver.StatusCode {
	case TaskModelStatusCodeProcessing:
		canIFail = false
	default:
		reason = "只有【执行中】任务可以【标记失败】"
	}

	return canIFail, reason
}

// CanIDelete 判断任务是否可以【删除】
func (receiver TaskModel) CanIDelete() (bool, string) {
	var (
		canIDelete = false
		reason     string
	)

	switch TaskModelStatusCodes(receiver.StatusCode) {
	case TaskModelStatusCodeFailed:
		canIDelete = true
		reason = "只有【失败】的任务可以被【删除】"
	}

	return canIDelete, reason
}

// CanICancel 判断任务是否可以【标记取消】
func (receiver TaskModel) CanICancel() (bool, string) {
	var (
		canICancel = true
		reason     string
	)

	switch TaskModelStatusCodes(receiver.StatusCode) {
	case TaskModelStatusCodeProcessing:
	case TaskModelStatusCodeCancel:
		canICancel = false
		reason = "【执行中】、【已取消】的任务不能【标记取消】"
	}
	return canICancel, reason
}

// NewTaskLogModelGorm 新建Gorm模型
func NewTaskLogModelGorm() *GormModel {
	return NewGorm().SetModel(TaskLogModel{})
}

// TableName 任务日志表名称
func (TaskLogModel) TableName() string {
	return "task_logs"
}
