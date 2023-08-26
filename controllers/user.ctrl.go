package controllers

import (
	"rnv-mmq/models"
	"rnv-mmq/services"
	"rnv-mmq/tools"
	"rnv-mmq/types"
	"rnv-mmq/wrongs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	// UserController 用户控制器
	UserController struct{}
	// UserStoreForm 用户表单
	UserStoreForm struct{}
)

// NewUserController 构造函数
func NewUserController() *UserController {
	return &UserController{}
}

// ShouldBind 表单绑定
func (receiver UserStoreForm) ShouldBind(ctx *gin.Context) UserStoreForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}

	return receiver
}

// Store 新建
func (UserController) Store(ctx *gin.Context) {
	var (
		ret *gorm.DB
		// repeat models.UserModel
	)

	// 新建
	user := &models.UserModel{}
	if ret = models.NewGorm().SetModel(models.UserModel{}).
		GetDb("").
		Create(&user); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Created(types.MapStringToAny{"user": user}).ToGinResponse())
}

// Delete 删除
func (UserController) Delete(ctx *gin.Context) {
	var (
		ret  *gorm.DB
		user models.UserModel
	)

	// 查询
	ret = models.NewGorm().SetModel(models.UserModel{}).
		SetWheres(types.MapStringToAny{"uuid": ctx.Param("uuid")}).
		GetDb("").
		First(&user)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	// 删除
	if ret := models.NewGorm().SetModel(models.UserModel{}).GetDb("").Delete(&user); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Deleted().ToGinResponse())
}

// Update 编辑
func (UserController) Update(ctx *gin.Context) {
	var (
		ret  *gorm.DB
		user models.UserModel
		// repeat  models.UserModel
	)

	// 表单
	// form := new(accountStoreForm).ShouldBind(ctx)

	// 查询
	ret = models.NewGorm().SetModel(models.UserModel{}).
		SetWheres(types.MapStringToAny{"uuid": ctx.Param("uuid")}).
		GetDb("").
		First(&user)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	// 编辑
	if ret = models.NewGorm().SetModel(models.UserModel{}).
		GetDb("").
		Where("uuid = ?", ctx.Param("uuid")).
		Save(&user); ret.Error != nil {
		wrongs.ThrowForbidden(ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Updated(types.MapStringToAny{"user": user}).ToGinResponse())
}

// Detail 详情
func (UserController) Detail(ctx *gin.Context) {
	var (
		ret  *gorm.DB
		user models.UserModel
	)
	ret = models.NewGorm().SetModel(models.UserModel{}).
		SetWheres(types.MapStringToAny{"uuid": ctx.Param("uuid")}).
		GetDb("").
		First(&user)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	ctx.JSON(tools.NewCorrectWithGinContext("", ctx).Datum(types.MapStringToAny{"user": user}).ToGinResponse())
}

func (UserController) listByQuery(ctx *gin.Context) *gorm.DB {
	return services.NewAccountService(services.BaseService{Model: models.NewGorm().SetModel(models.UserModel{}), Ctx: ctx}).GetListByQuery()
}

// List 列表
func (receiver UserController) List(ctx *gin.Context) {
	var users []models.UserModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForPager(
				receiver.listByQuery(ctx),
				func(db *gorm.DB) types.MapStringToAny {
					db.Find(&users)
					return types.MapStringToAny{"users": users}
				},
			).ToGinResponse(),
	)
}

// ListJdt jquery-dataTable分页列表
func (receiver UserController) ListJdt(ctx *gin.Context) {
	var users []models.UserModel

	ctx.JSON(
		tools.NewCorrectWithGinContext("", ctx).
			DataForJqueryDataTable(
				receiver.listByQuery(ctx),
				func(db *gorm.DB) types.MapStringToAny {
					db.Find(&users)
					return types.MapStringToAny{"users": users}
				},
			).ToGinResponse(),
	)
}
