package user

import (
	"dcss/global"
	"dcss/global/param_check"
	"dcss/handler/api/v1/system/captcha"
	"dcss/models"
	"dcss/pkg/resp"
	"dcss/pkg/utils"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddUser
// @Tags     用户相关
// @Summary  添加用户账号
// @Produce  application/json
// @Param    data  body      models.ReqAddUser true "用户名 中文名 密码 电话 邮箱 状态"
// @Success  200   {object}  resp.Result{}  "添加用户账号"
// @Security ApiKeyAuth
// @Router   /auth/user [post]
func AddUser(ctx *gin.Context) {
	//获取参数
	var reqInfo models.ReqAddUser
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorf("reqInfo bind err: %v", err)
		resp.Fail(ctx, nil, err.Error())
		return
	}

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	//  判断用户是否存在
	var user models.SysUser
	// 获取数据库连接池
	db := global.DB

	db.Model(&models.SysUser{}).Where("username = ?", reqInfo.Username).First(&user)
	if user.ID != 0 {
		resp.Fail(ctx, nil, "用户名已存在")
		return
	}

	//	创建用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		resp.Response(ctx, 500, nil, "加密错误")
		return
	}

	reqInfo.Password = string(hashedPassword)
	err = db.Create(&models.SysUser{
		Username:    reqInfo.Username,
		Password:    reqInfo.Password,
		ChineseName: reqInfo.ChineseName,
		Phone:       reqInfo.Phone,
		Email:       reqInfo.Email,
		State:       reqInfo.State,
	}).Error

	if err != nil {
		global.LOG.Errorln("创建用户失败,err:", err)
		resp.Fail(ctx, nil, "创建用户失败")
		return
	}

	resp.Response(ctx, 201, nil, "注册成功")
}

// Login
// @Tags     用户相关
// @Summary  用户登录
// @Produce  application/json
// @Param    data  body      models.ReqLogin true "用户名 密码"
// @Success  200   {object}  resp.Result{data=models.RespLogin}  "用户登录"
// @Router   /login [post]
func Login(ctx *gin.Context) {
	// 获取数据库连接池
	db := global.DB
	// 获取参数
	var reqInfo models.ReqLogin
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		resp.Fail(ctx, nil, err.Error())
		return
	}

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	// 判断验证码是否有效
	if _, ok := global.CaptchaCache.Get(reqInfo.CaptchaID); !ok {
		resp.Fail(ctx, nil, "验证码已过期，请输入验证码后重新登录")
		return
	}

	verify := captcha.Store.Verify(reqInfo.CaptchaID, reqInfo.Captcha, true)
	if !verify {
		resp.Fail(ctx, nil, "验证码错误，请重新登录")
		return
	}

	// 判断用户是否存在
	var user models.SysUser
	err = db.Where("username = ?", reqInfo.Username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(ctx, nil, "用户名或密码错误")
		return
	}

	// 判断用户是否被禁用
	if !user.State {
		resp.Fail(ctx, nil, "用户已禁用")
		return
	}
	//	判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqInfo.Password)); err != nil {
		resp.Fail(ctx, nil, "用户名或密码错误")
		return
	}

	// 发放token
	token, err := utils.ReleaseToken(&user)

	if err != nil {
		resp.Fail(ctx, nil, "系统错误")
		return
	}
	// 返回结果
	resp.Success(ctx, models.RespLogin{
		ID:          user.ID,
		AccessToken: "Bearer " + token,
		Username:    user.Username,
		ChineseName: user.ChineseName,
	}, "登录成功")
}

// GetUserList
// @Tags     用户相关
// @Summary  获取用户列表
// @Produce  application/json
// @Param    query  query  string     false  "过滤关键字"
// @Param    pagenum   query   int     false  "页码"
// @Param    pagesize   query    int     false  "返回数量"
// @Success  200   {object}  resp.Result{data=models.RespGetUserList}  "获取用户列表"
// @Security ApiKeyAuth
// @Router   /auth/users [get]
func GetUserList(ctx *gin.Context) {
	db := global.DB.Model(&models.SysUser{})
	// 从请求的接口中获取相关数据
	queryStr := ctx.DefaultQuery("query", "")
	pagenum, _ := strconv.Atoi(ctx.DefaultQuery("pagenum", "1"))
	pagenum = pagenum - 1
	pagesize, _ := strconv.Atoi(ctx.DefaultQuery("pagesize", "10"))
	// 定义数据库接收的user结构体变量
	var userList []models.RespGetUserByID
	var err error
	// 用户总数量
	var total int64

	if pagesize > 100 || pagesize < 0 {
		pagesize = 20
	}

	if len(queryStr) != 0 {
		queryStr = "%" + queryStr + "%"
		db = db.Where("username LIKE ?", queryStr).
			Or("phone LIKE ?", queryStr).
			Or("email LIKE ?", queryStr)
	}

	err = db.Count(&total).Error
	if err != nil {
		global.LOG.Errorln("获取用户列表总数失败, err:", err)
		resp.Fail(ctx, nil, "获取用户列表总数失败")
		return
	}

	err = db.Preload("Roles").
		Offset(pagenum*pagesize).
		Limit(pagesize).
		Order("id").
		Omit("deleted_at", "password").
		Find(&userList).Error
	if err != nil {
		global.LOG.Errorln("获取用户列表失败, err:", err)
		resp.Fail(ctx, nil, "获取用户列表失败")
		return
	}

	resp.Success(ctx, models.RespGetUserList{
		Users: userList,
		Total: total,
	}, "获取用户列表成功")
}

// GetUserByID
// @Tags     用户相关
// @Summary  通过id来获取用户信息
// @Produce  application/json
// @Param    id  path int true "用户ID"
// @Success  200   {object}  resp.Result{data=models.RespGetUserByID}  "通过id来获取用户信息"
// @Security ApiKeyAuth
// @Router   /auth/user/{id} [get]
func GetUserByID(ctx *gin.Context) {
	//	获取用户id
	id := ctx.Param("id")
	// 获取db
	db := global.DB.Model(&models.SysUser{})
	var user models.RespGetUserByID
	// 判断用户是否存在
	err := db.Preload("Roles").Where("id = ?", id).Omit("deleted_at", "password").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(ctx, nil, "用户不存在")
		return
	}
	resp.Success(ctx, user, "获取用户信息成功")
}

// EditUserByID
// @Tags     用户相关
// @Summary  通过id来修改用户信息
// @Produce  application/json
// @Param    id  path int true "用户ID"
// @Param    data  body      models.ReqEditUserByID true "中文名 电话 邮箱 状态"
// @Success  200   {object}  resp.Result{}  "通过id来修改用户信息"
// @Security ApiKeyAuth
// @Router   /auth/user/{id} [put]
func EditUserByID(ctx *gin.Context) {
	//	获取用户id
	id := ctx.Param("id")
	//	获取接口信息
	var reqInfo models.ReqEditUserByID
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Debugln("获取reqInfo, err: ", err)
		resp.Fail(ctx, nil, "获取reqInfo失败")
		return
	}

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	// 获取全局db
	db := global.DB
	user := models.SysUser{}

	// 判断用户是否存在
	err = db.Model(&models.SysUser{}).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(ctx, nil, "用户不存在")
		return
	}

	// 不允许自己修改自己的状态
	ctxID := ctx.GetUint("userID")
	if user.ID == ctxID && user.State != reqInfo.State {
		resp.Fail(ctx, nil, "自己不能修改自己状态")
		return
	}

	err = db.Model(&models.SysUser{}).Where("id = ?", id).Updates(map[string]interface{}{
		"state":        reqInfo.State,
		"chinese_name": reqInfo.ChineseName,
		"phone":        reqInfo.Phone,
		"email":        reqInfo.Email,
	}).Error

	if err != nil {
		global.LOG.Errorln("更新用户失败, err: ", err)
		resp.Fail(ctx, nil, "更新用户失败")
		return
	}
	resp.Success(ctx, nil, "更新用户信息成功")
}

// DeleteUserByID
// @Tags     用户相关
// @Summary  通过id来删除用户
// @Produce  application/json
// @Param    id  path int true "用户ID"
// @Success  200   {object}  resp.Result{}  "通过id来删除用户"
// @Security ApiKeyAuth
// @Router   /auth/user/{id} [delete]
func DeleteUserByID(ctx *gin.Context) {
	//	获取接口参数
	id := ctx.Param("id")
	db := global.DB
	// 判断用户是否存在
	var user models.SysUser
	err := db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(ctx, nil, "用户不存在")
		return
	}
	// 不允许自己删除自己
	ctxID := ctx.GetUint("userID")
	if user.ID == ctxID {
		resp.Fail(ctx, nil, "自己不能删除自己")
		return
	}

	// 不允许删除初始化管理员
	if user.ID == 1 {
		resp.Fail(ctx, nil, "不能删除初始化管理员")
		return
	}

	err = db.Unscoped().Delete(&user).Error
	if err != nil {
		global.LOG.Errorln("用户删除失败, err: ", err)
		resp.Fail(ctx, nil, "用户删除失败")
		return
	}

	resp.Success(ctx, nil, "删除用户信息成功")
}

// EditUserStateByID
// @Tags     用户相关
// @Summary  通过id来修改用户状态
// @Produce  application/json
// @Param    id  path int true "用户ID"
// @Param    data  body      models.ReqEditUserStateByID true "中文名 电话 邮箱 状态"
// @Success  200   {object}  resp.Result{}  "通过id来修改用户状态"
// @Security ApiKeyAuth
// @Router   /auth/user/{id}/state/ [put]
func EditUserStateByID(ctx *gin.Context) {
	//获取接口参数
	id := ctx.Param("id")
	var reqInfo models.ReqEditUserStateByID
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorln("获取reqInfo, err: ", err)
		resp.Fail(ctx, nil, "获取reqInfo失败")
		return
	}

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	db := global.DB

	// 判断用户是否存在
	var user models.SysUser
	err = db.Model(&models.SysUser{}).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(ctx, nil, "用户不存在")
		return
	}

	// 不允许自己修改自己的状态
	ctxID := ctx.GetUint("userID")
	if user.ID == ctxID {
		resp.Fail(ctx, nil, "自己不能修改自己状态")
		return
	}

	user.State = reqInfo.State
	err = db.Select("State").Save(&user).Error
	if err != nil {
		global.LOG.Errorln("更新用户状态失败, err: ", err)
		resp.Fail(ctx, nil, "更新用户状态失败")
		return
	}

	resp.Success(ctx, 200, "更新用户状态成功")
}

// EditUserPasswordByID 通过id来修改用户密码
// @Tags     用户相关
// @Summary  通过id来修改用户密码
// @Produce  application/json
// @Param    id  path int true "用户ID"
// @Param    data  body      models.ReqEditUserPasswordByID true "密码 确认密码 新密码"
// @Success  200   {object}  resp.Result{}  "通过id来修改用户密码"
// @Security ApiKeyAuth
// @Router   /auth/user/{id}/password/ [put]
func EditUserPasswordByID(ctx *gin.Context) {
	//获取接口参数
	id := ctx.Param("id")
	var reqInfo models.ReqEditUserPasswordByID
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		resp.Fail(ctx, nil, "获取reqInfo失败")
		return
	}

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	if reqInfo.Password == reqInfo.NewPassword {
		resp.Fail(ctx, nil, "新旧密码一样")
		return
	}

	db := global.DB
	// 判断用户是否存在
	var user models.SysUser
	err = db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp.Fail(ctx, nil, "用户不存在")
		return
	}

	//	判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqInfo.Password)); err != nil {
		resp.Fail(ctx, nil, "密码错误")
		return
	}

	// 生成新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqInfo.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		resp.Response(ctx, 500, nil, "加密错误")
		return
	}

	err = db.Model(&user).Update("password", string(hashedPassword)).Error
	if err != nil {
		global.LOG.Errorln("更新用户密码失败, err: ", err)
		resp.Fail(ctx, nil, "更新用户密码失败")
		return
	}

	resp.Success(ctx, nil, "更新用户密码成功")
}

// GetAllUserList
// @Tags     用户相关
// @Summary  获取所有的用户,提供给角色添加用户使用
// @Produce  application/json
// @Success  200   {object}  resp.Result{data=models.RespGetAllUserList}  "获取所有的用户"
// @Security ApiKeyAuth
// @Router   /auth/all_users [get]
func GetAllUserList(ctx *gin.Context) {
	db := global.DB.Model(&models.SysUser{})
	// 定义数据库接收的user结构体变量
	var respUserList models.RespGetAllUserList
	var err error

	err = db.Count(&respUserList.Total).Error
	if err != nil {
		global.LOG.Errorln("获取所有的用户列表总数失败, err:", err)
		resp.Fail(ctx, nil, "获取所有的用户列表总数失败")
		return
	}

	err = db.Order("id").
		Find(&respUserList.Users).Error
	if err != nil {
		global.LOG.Errorln("获取用户列表失败, err:", err)
		resp.Fail(ctx, nil, "获取用户列表失败")
		return
	}

	resp.Success(ctx, respUserList, "获取用户列表成功")
}
