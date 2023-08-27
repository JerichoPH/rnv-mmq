package models

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"rnv-mmq/tools"
	"rnv-mmq/wrongs"
)

type (
	// FileModel 文件模型
	FileModel struct {
		GormModel
		Filename          string      `gorm:"varchar(128);not null;default:'';comment:文件保存名;" json:"filename,omitempty"`
		OriginalFilename  string      `gorm:"type:varchar(128);not null;default:'';comment:原始文件名称" json:"original_filename,omitempty"`
		OriginalExtension string      `gorm:"type:varchar(128);not null;default:'';comment:原始文件扩展名" json:"original_extension,omitempty"`
		Size              uint64      `gorm:"type:bigint(20);not null;default:0;comment:文件大小" json:"size,omitempty"`
		FileType          FileType    `gorm:"type:varchar(128);not null;default:'';comment:文件类型" json:"file_type,omitempty"`
		PrefixPath        string      `gorm:"type:varchar(128);not null;default:'';comment:文件前缀路径" json:"prefix_path,omitempty"`
		Tasks             []TaskModel `gorm:"foreignKey:content_file_uuid;references:uuid;comment:相关任务（请求）" json:"tasks,omitempty"`
		Content           string      `gorm:"-" json:"content,omitempty"`
		Url               string      `gorm:"-" json:"url,omitempty"`
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
func (*FileModel) TableName() string {
	return "files"
}

// StoreOne 保存单个文件
func (*FileModel) StoreOne(file *FileModel, content string) *FileModel {
	var (
		ret      *gorm.DB
		filename string
	)

	if ret = NewFileModelGorm().
		GetDb("").
		Create(&file); ret.Error != nil {
		wrongs.ThrowForbidden("保存文件失败：%s", ret.Error.Error())
	}

	filename = fmt.Sprintf("%s%s", file.Uuid, file.OriginalExtension)
	file.Filename = filename
	if ret = NewFileModelGorm().
		GetDb("").
		Where("uuid = ?", file.Uuid).
		Save(&file); ret.Error != nil {
		wrongs.ThrowForbidden("保存文件失败：%s", ret.Error.Error())
	}

	tools.NewFileSystem(tools.JoinWithoutEmpty([]string{"./static", file.PrefixPath, file.Filename}, "/")).WriteString(content)

	return file
}

// IsExist 查看文件是否存在
func (receiver *FileModel) IsExist() bool {
	if _, err := os.Stat(tools.JoinWithoutEmpty([]string{"./static", receiver.PrefixPath, receiver.Filename}, "/")); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		wrongs.ThrowForbidden(err.Error())
		return false
	}
	return true
}

// AfterFind 自动读取文件内容
func (receiver *FileModel) AfterFind(tx *gorm.DB) (err error) {
	if receiver.IsExist() {
		switch receiver.FileType {
		case FileTypeJson:
			receiver.Content = tools.NewFileSystem(
				tools.JoinWithoutEmpty([]string{"./static", receiver.PrefixPath, receiver.Filename}, "/"),
			).ReadString()
		case FileTypeText:
			receiver.Content = tools.NewFileSystem(
				tools.JoinWithoutEmpty([]string{"./static", receiver.PrefixPath, receiver.Filename}, "/"),
			).ReadString()
		}
	}

	return nil
}
