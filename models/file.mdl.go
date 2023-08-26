package models

import (
	"fmt"
	"os"
	"rnv-mmq/wrongs"
)

// FileModel 文件模型
type FileModel struct {
	GormModel
	Name              string      `gorm:"type:varchar(128);not null;default:'';comment:文件名称" json:"name,omitempty"`
	OriginalFilename  string      `gorm:"type:varchar(128);not null;default:'';comment:原始文件名称" json:"original_filename,omitempty"`
	OriginalExtension string      `gorm:"type:varchar(128);not null;default:'';comment:原始文件扩展名" json:"original_extension,omitempty"`
	Size              uint64      `gorm:"type:bigint(20);not null;default:0;comment:文件大小" json:"size,omitempty"`
	PrefixPath        string      `gorm:"type:varchar(128);not null;default:'';comment:文件前缀路径" json:"prefix_path,omitempty"`
	Tasks             []TaskModel `gorm:"foreignKey:file_uuid;references:uuid;comment:相关任务" json:"tasks,omitempty"`
}

// TableName 文件表名称
func (FileModel) TableName() string {
	return "files"
}

// IsExist 查看文件是否存在
func (receiver FileModel) IsExist() bool {
	if _, err := os.Stat(fmt.Sprintf("%s/%s", receiver.PrefixPath, receiver.Name)); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		wrongs.ThrowForbidden(err.Error())
		return false
	}
	return true
}
