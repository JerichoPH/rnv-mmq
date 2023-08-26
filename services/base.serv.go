package services

import (
	"rnv-mmq/models"

	"github.com/gin-gonic/gin"
)

type (
	// BaseService 基础服务
	BaseService struct {
		Model      *models.GormModel
		Ctx        *gin.Context
		DbConnName string
	}
)
