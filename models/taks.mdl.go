package models

type (
	// TaskModel 任务模型
	TaskModel struct {
		GormModel
		Name        string         `gorm:"type:varchar(128);not null;default:'';comment:任务名称" json:"name,omitempty"`
		Target      string         `gorm:"type:varchar(128);not null;default:'';comment:任务目标;" json:"target,omitempty"`
		Description string         `gorm:"type:varchar(128);not null;default:'';comment:任务说明;" json:"description,omitempty"`
		StatusCode  string         `gorm:"type:enum('ORIGINAL','PROCESSING','FINISHED','FAILED');not null;default:'pending';comment:任务状态码" json:"status_code,omitempty"`
		StatusText  string         `gorm:"->;type:varchar(128) as ((case status_code when 'ORIGINAL' then '未执行' when 'PROCESSING' then '执行中' when 'FINISHED' then '已完成' when 'FAILED' then '失败' else '' end));not null;default:'';comment:任务状态文本" json:"status_text,omitempty"`
		FileUuid    string         `gorm:"index;type:char(36);not null;comment:所属文件uuid" json:"file_uuid,omitempty"`
		File        *FileModel     `gorm:"foreignKey:file_uuid;references:uuid;comment:所属文件" json:"file,omitempty"`
		TaskLogs    []TaskLogModel `gorm:"foreignKey:task_uuid:references:uuid;comment:相关日志;" json:"task_logs,omitempty"`
	}
	// TaskLogModel 任务日志模型
	TaskLogModel struct {
		GormModel
		TaskUuid string    `gorm:"index;type:char(36);not null;comment:所属任务uuid" json:"task_uuid,omitempty"`
		Task     TaskModel `gorm:"foreignKey:task_uuid:references:uuid;comment:所属任务;" json:"task,omitempty"`
	}
)

// TableName 任务表名称
func (TaskModel) TableName() string {
	return "tasks"
}

// TableName 任务日志表名称
func (TaskLogModel) TableName() string {
	return "task_logs"
}
