package controllers

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"rnv-mmq/models"
	"rnv-mmq/tools"
	"rnv-mmq/types"
	"rnv-mmq/wrongs"
)

type (
	// AuthorizationController 权鉴控制器
	AuthorizationController struct{}
	// authorizationRegisterForm 注册表单
	authorizationRegisterForm struct {
		Username             string `json:"username"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
		Nickname             string `json:"nickname"`
	}
)

// NewAuthorizationController 构造函数
func NewAuthorizationController() *AuthorizationController {
	return &AuthorizationController{}
}

// ShouldBind 绑定表单
//
//	@receiver ins
//	@param ctx
//	@return authorizationRegisterForm
func (receiver authorizationRegisterForm) ShouldBind(ctx *gin.Context) authorizationRegisterForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if receiver.Username == "" {
		wrongs.ThrowValidate("账号必填")
	}
	if receiver.Password == "" {
		wrongs.ThrowValidate("密码必填")
	}
	if receiver.Nickname == "" {
		wrongs.ThrowValidate("昵称必填")
	}
	if len(receiver.Password) < 6 || len(receiver.Password) > 18 {
		wrongs.ThrowValidate("密码不可小于6位或大于18位")
	}
	if receiver.Password != receiver.PasswordConfirmation {
		wrongs.ThrowValidate("两次密码输入不一致")
	}

	return receiver
}

// AuthorizationLoginForm 登录表单
type AuthorizationLoginForm struct {
	Username string `form:"account" json:"account" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// ShouldBind 绑定表单
func (receiver AuthorizationLoginForm) ShouldBind(ctx *gin.Context) AuthorizationLoginForm {
	if err := ctx.ShouldBind(&receiver); err != nil {
		wrongs.ThrowValidate(err.Error())
	}
	if receiver.Username == "" {
		wrongs.ThrowValidate("账号必填")
	}
	if receiver.Password == "" {
		wrongs.ThrowValidate("密码必填")
	}
	if len(receiver.Password) < 6 || len(receiver.Password) > 18 {
		wrongs.ThrowValidate("密码不可小于6位或大于18位")
	}

	return receiver
}

// Register 注册
func (AuthorizationController) Register(ctx *gin.Context) {
	// 表单验证
	form := (&authorizationRegisterForm{}).ShouldBind(ctx)

	// 检查重复项（用户名）
	var repeat models.UserModel
	var ret *gorm.DB
	ret = (&models.GormModel{}).
		SetWheres(types.MapStringToAny{"username": form.Username}).
		GetDb("").
		First(&repeat)
	wrongs.ThrowWhenIsRepeat(ret, "用户名")
	ret = (&models.GormModel{}).
		SetWheres(types.MapStringToAny{"nickname": form.Nickname}).
		GetDb("").
		First(&repeat)
	wrongs.ThrowWhenIsRepeat(ret, "昵称")

	// 密码加密
	bytes, _ := bcrypt.GenerateFromPassword([]byte(form.Password), 14)

	// 保存新用户
	account := &models.UserModel{
		GormModel: models.GormModel{Uuid: uuid.NewV4().String()},
		Username:  form.Username,
		Password:  string(bytes),
		Nickname:  form.Nickname,
	}
	if ret = models.NewGorm().SetModel(models.UserModel{}).
		SetOmits(clause.Associations).
		GetDb("").
		Create(&account); ret.Error != nil {
		wrongs.ThrowForbidden("创建失败：" + ret.Error.Error())
	}

	ctx.JSON(tools.NewCorrectWithGinContext("注册成功", ctx).Created(types.MapStringToAny{"account": account}).ToGinResponse())
}

// Login 登录
func (AuthorizationController) Login(ctx *gin.Context) {
	// 表单验证
	form := (&AuthorizationLoginForm{}).ShouldBind(ctx)

	var (
		user models.UserModel
		ret  *gorm.DB
	)
	// 获取用户
	ret = models.NewGorm().
		SetModel(models.UserModel{}).
		SetWheres(types.MapStringToAny{"user": form.Username}).
		GetDb("").
		First(&user)
	wrongs.ThrowWhenIsEmpty(ret, "用户")

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		wrongs.ThrowUnAuth("账号或密码错误")
	}

	// 生成Jwt
	if token, err := tools.GenerateJwt(
		user.Uuid,
		user.Nickname,
	); err != nil {
		// 生成jwt错误
		wrongs.ThrowForbidden(err.Error())
	} else {
		ctx.JSON(tools.NewCorrectWithGinContext("登陆成功", ctx).Datum(types.MapStringToAny{
			"token": token,
			"user": types.MapStringToAny{
				"username": user.Username,
				"uuid":     user.Uuid,
			},
		}).ToGinResponse())
	}
}

// GetMenus 获取当前用户菜单
// func (AuthorizationController) GetMenus(ctx *gin.Context) {
//	var ret *gorm.GetDb
//	if accountUuid, exists := ctx.Get(tools.AccountOpenIdFieldName); !exists {
//		wrongs.ThrowUnLogin("用户未登录")
//	} else {
//		// 获取当前用户信息
//		var account models.UserModel
//		ret = models.NewGorm().SetModel(models.UserModel{}).
//			SetWheres(types.MapStringToAny{"uuid": accountUuid}).
//			SetPreloads("RbacRoles", "RbacRoles.Menus").
//			GetDb("",nil).
//			FindOneUseQuery(&account)
//		if !wrongs.ThrowWhenIsEmpty(ret, "") {
//			wrongs.ThrowUnLogin("当前令牌指向用户不存在")
//		}
//
//		var menus []models.MenuModel
//		models.NewGorm().SetModel(models.MenuModel{}).
//			GetDb("",nil).
//			Joins("join pivot_rbac_role_and_menus prram on menus.uuid = prram.menu_uuid").
//			Joins("join rbac_roles r on prram.rbac_role_uuid = r.uuid").
//			Joins("join pivot_rbac_role_and_accounts prraa on r.uuid = prraa.rbac_role_uuid").
//			Joins("join accounts a on prraa.account_uuid = a.uuid").
//			Where("a.uuid = ?", account.GormModel.Uuid).
//			Where("menus.deleted_at is null").
//			Where("menus.parent_uuid = ''").
//			Order("menus.sort asc").
//			Order("menus.id asc").
//			Preload("Subs").
//			Find(&menus)
//
//		ctx.JSON(tools.CorrectInit("", ctx).Datum(types.MapStringToAny{"menus": menus}))
//	}
// }
