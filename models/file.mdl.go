package models

import (
	"gorm.io/gorm"
	"os"
	"rnv-mmq/tools"
	"rnv-mmq/wrongs"
)

type (
	// FileModel 文件模型
	FileModel struct {
		GormModel
		OriginalFilename  string      `gorm:"type:varchar(128);not null;default:'';comment:原始文件名称" json:"original_filename,omitempty"`
		OriginalExtension string      `gorm:"type:varchar(128);not null;default:'';comment:原始文件扩展名" json:"original_extension,omitempty"`
		Size              uint64      `gorm:"type:bigint(20);not null;default:0;comment:文件大小" json:"size,omitempty"`
		FileType          FileType    `gorm:"type:varchar(128);not null;default:'';comment:文件类型" json:"file_type,omitempty"`
		PrefixPath        string      `gorm:"type:varchar(128);not null;default:'';comment:文件前缀路径" json:"prefix_path,omitempty"`
		TasksForRequest   []TaskModel `gorm:"foreignKey:request_file_uuid;references:uuid;comment:相关任务（请求）" json:"tasks_for_request,omitempty"`
		Content           string      `gorm:"-"`
		Url               string      `gorm:"-"`
	}

	// FileType 文件类型
	FileType string
)

const (
	FileTypeText FileType = "TEXT"
	FileTypeJson FileType = "JSON"
	FileTypePdf  FileType = "PDF"
)

// NewFileModelGorm 新建Gorm模型
func NewFileModelGorm() *GormModel {
	return NewGorm().SetModel(FileModel{})
}

// TableName 文件表名称
func (FileModel) TableName() string {
	return "files"
}

// IsExist 查看文件是否存在
func (receiver FileModel) IsExist() bool {
	if _, err := os.Stat(tools.JoinWithoutEmpty([]string{"./static", receiver.PrefixPath, receiver.Uuid}, "/")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		wrongs.ThrowForbidden(err.Error())
		return false
	}
	return true
}

// AfterFind 自动读取文件内容
func (receiver FileModel) AfterFind(tx *gorm.DB) (err error) {
	if receiver.IsExist() {
		switch receiver.FileType {
		case FileTypeJson:
			receiver.Content = tools.NewFileSystem(
				tools.JoinWithoutEmpty([]string{"./static", receiver.PrefixPath, receiver.Uuid}, "/"),
			).ReadJson().(string)
		case FileTypeText:
			receiver.Content = tools.NewFileSystem(
				tools.JoinWithoutEmpty([]string{"./static", receiver.PrefixPath, receiver.Uuid}, "/"),
			).ReadString()
		}
	}

	return nil
}
