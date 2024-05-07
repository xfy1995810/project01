package role

import (
	"dcss/global"
	"dcss/global/casbin"
	"dcss/global/param_check"
	"dcss/models"
	"dcss/pkg/resp"
	"dcss/pkg/utils"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddRole
// @Tags     角色相关
// @Summary  添加角色
// @Produce  application/json
// @Param    data  body      models.ReqAddRole true "角色名 备注 状态 []{用户}"
// @Success  200   {object}  resp.Result{}  "添加角色"
// @Security ApiKeyAuth
// @Router   /auth/role [post]
func AddRole(ctx *gin.Context) {
	//获取参数
	var reqInfo models.ReqAddRole
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorf("reqInfo bind err: %v", err)
		resp.Fail(ctx, nil, err.Error())
		return
	}
	global.LOG.Debugf("AddRole,reqInfo: %v\n", reqInfo)

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	//  判断角色是否存在
	var role models.SysRole
	// 获取数据库连接池
	db := global.DB

	err = db.Model(&models.SysRole{}).Where("name = ?", reqInfo.Name).First(&role).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Errorln("判断角色是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断角色是否存在报错")
		return
	}

	if role.ID != 0 {
		resp.Fail(ctx, nil, "角色名已存在")
		return
	}

	var users []models.SysUser
	if len(reqInfo.UserIds) > 0 {
		err = db.Select("id").Find(&users, reqInfo.UserIds).Error
		if err != nil {
			resp.Fail(ctx, nil, "角色获取用户失败")
			return
		}
	}

	role = models.SysRole{
		Name:   reqInfo.Name,
		Remark: reqInfo.Remark,
		State:  reqInfo.State,
		Users:  users,
		Common: models.Common{
			CreatedID: ctx.GetUint("userID"),
		},
	}
	global.LOG.WithFields(logrus.Fields{
		"data": utils.GetJson(&role),
	}).Debugln("新增角色信息")

	//	创建角色
	err = db.Create(&role).Error
	if err != nil {
		global.LOG.Errorln("创建角色失败,err:", err)
		resp.Fail(ctx, nil, "创建角色失败")
		return
	}

	resp.Response(ctx, 201, nil, "创建角色成功")
}

// GetRoleList
// @Tags     角色相关
// @Summary  获取所有的角色
// @Produce  application/json
// @Param    query  query  string     false  "过滤关键字"
// @Param    pagenum   query   int     false  "页码"
// @Param    pagesize   query    int     false  "返回数量"
// @Success  200   {object}  resp.Result{data=models.RespGetRoleList}  "获取所有的角色"
// @Security ApiKeyAuth
// @Router   /auth/roles [get]
func GetRoleList(ctx *gin.Context) {
	db := global.DB.Model(&models.SysRole{})
	// 从请求的接口中获取相关数据
	queryStr := ctx.DefaultQuery("query", "")
	pagenum, _ := strconv.Atoi(ctx.DefaultQuery("pagenum", "1"))
	pagenum = pagenum - 1
	pagesize, _ := strconv.Atoi(ctx.DefaultQuery("pagesize", "10"))
	// 定义数据库接收的role结构体变量
	var roleList []models.RespGetRoleByID
	var err error
	// 角色总数量
	var total int64

	if pagesize > 100 || pagesize < 0 {
		pagesize = 20
	}

	if len(queryStr) != 0 {
		queryStr = "%" + queryStr + "%"
		db = db.Where("name LIKE ?", queryStr).
			Or("remark LIKE ?", queryStr)
	}

	err = db.Count(&total).Error
	if err != nil {
		global.LOG.Errorln("获取角色列表总数失败, err:", err)
		resp.Fail(ctx, nil, "获取角色列表总数失败")
		return
	}

	err = db.Scopes(models.HandleCommonUserPreload).Preload("Users").
		Offset(pagenum * pagesize).
		Limit(pagesize).
		Order("id").
		Omit("deleted_at").
		Find(&roleList).Error
	if err != nil {
		global.LOG.Errorln("获取角色列表失败, err:", err)
		resp.Fail(ctx, nil, "获取角色列表失败")
		return
	}

	resp.Success(ctx, models.RespGetRoleList{
		Roles: roleList,
		Total: total,
	}, "获取角色列表成功")
}

// GetRoleByID
// @Tags     角色相关
// @Summary  通过id来获取角色信息
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Success  200   {object}  resp.Result{data=models.RespGetRoleByID}  "通过id来获取角色信息"
// @Security ApiKeyAuth
// @Router   /auth/role/{id} [get]
func GetRoleByID(ctx *gin.Context) {
	//	获取角色id
	id := ctx.Param("id")
	// 获取db
	db := global.DB.Model(&models.SysRole{})
	var role models.RespGetRoleByID
	// 判断角色是否存在
	err := db.Scopes(models.HandleCommonUserPreload).Preload("Users").
		Where("id = ?", id).
		Omit("deleted_at").First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "角色不存在")
			return
		}
		global.LOG.Errorln("获取角色信息失败, err:", err)
		resp.Fail(ctx, nil, "获取角色信息失败")
		return
	}

	resp.Success(ctx, role, "获取角色信息成功")
}

// EditRoleByID
// @Tags     角色相关
// @Summary  通过id来修改角色信息
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Param    data  body      models.ReqEditRoleByID true "名称 备注 状态 []{用户ID}"
// @Success  200   {object}  resp.Result{}  "通过id来修改角色信息"
// @Security ApiKeyAuth
// @Router   /auth/role/{id} [put]
func EditRoleByID(ctx *gin.Context) {
	//	获取接口信息
	var reqInfo models.ReqEditRoleByID

	// 判断角色是否存在
	var role models.SysRole
	err := checkRoleIsExist(ctx, &reqInfo, &role)
	if err != nil {
		return
	}

	var users []models.SysUser
	if len(reqInfo.UserIds) > 0 {
		err = global.DB.Select("id").Find(&users, reqInfo.UserIds).Error
		if err != nil {
			global.LOG.Errorln("角色获取用户失败, err:", err)
			resp.Fail(ctx, nil, "角色获取用户失败")
			return
		}
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&role).
			Association("Users").Replace(users); err != nil {
			return fmt.Errorf("清除角色的用户信息错误,err: %s", err)
		}

		if err := tx.Model(&models.SysRole{}).Where("id = ?", role.ID).
			Updates(map[string]interface{}{
				"id":        role.ID,
				"name":      reqInfo.Name,
				"remark":    reqInfo.Remark,
				"state":     reqInfo.State,
				"UpdatedID": ctx.GetUint("userID"),
			}).Error; err != nil {
			return fmt.Errorf("更新角色信息错误,err: %s", err)
		}

		return nil
	})

	if err != nil {
		global.LOG.Errorln(err)
		resp.Fail(ctx, nil, "更新角色失败")
		return
	}
	resp.Success(ctx, nil, "更新角色信息成功")
}

// DeleteRoleByID
// @Tags     角色相关
// @Summary  通过id来删除角色
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Success  200   {object}  resp.Result{}  "通过id来删除角色"
// @Security ApiKeyAuth
// @Router   /auth/role/{id} [delete]
func DeleteRoleByID(ctx *gin.Context) {
	// 判断角色是否存在
	var role models.SysRole
	err := checkRoleIsExist(ctx, nil, &role)
	if err != nil {
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&role).Association("Users").Clear()
		if err != nil {
			return fmt.Errorf("清除角色的用户信息错误,err: %s", err)
		}

		err = tx.Unscoped().Delete(&role).Error
		if err != nil {
			return fmt.Errorf("角色删除失败, err: %s", err)
		}

		err = casbin.RemoveFilteredPolicy(tx, role.ID)
		if err != nil {
			return fmt.Errorf("角色删除失败, err: %s", err)
		}

		return nil
	})

	if err != nil {
		global.LOG.Errorln(err)
		resp.Fail(ctx, nil, "角色删除失败")
		return
	}

	err = casbin.FreshCasbin()
	if err != nil {
		global.LOG.Errorln("重新加载角色Api失败, err:", err)
		resp.Fail(ctx, nil, "重新加载角色Api失败")
		return
	}

	resp.Success(ctx, nil, "删除角色信息成功")
}

// EditRoleStateByID
// @Tags     角色相关
// @Summary  通过id来修改角色状态
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Param    data  body      models.ReqEditRoleStateByID true "中文名 电话 邮箱 状态"
// @Success  200   {object}  resp.Result{}  "通过id来修改角色状态"
// @Security ApiKeyAuth
// @Router   /auth/role/{id}/state/ [put]
func EditRoleStateByID(ctx *gin.Context) {
	var reqInfo models.ReqEditRoleStateByID

	// 判断角色是否存在
	var role models.SysRole
	err := checkRoleIsExist(ctx, &reqInfo, &role)
	if err != nil {
		return
	}

	role.State = reqInfo.State
	err = global.DB.Select("State").Save(&role).Error
	if err != nil {
		global.LOG.Errorln("更新角色状态失败, err: ", err)
		resp.Fail(ctx, nil, "更新角色状态失败")
		return
	}

	resp.Success(ctx, 200, "更新角色状态成功")
}

// SetRoleApiAuth
// @Tags     角色相关
// @Summary  通过角色id来绑定Api权限
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Param    data  body      models.ReqSetRoleApiAuth true "ApiIds[]int"
// @Success  200   {object}  resp.Result{}  "通过角色id来绑定Api权限"
// @Security ApiKeyAuth
// @Router   /auth/role/{id}/api_auth [post]
func SetRoleApiAuth(ctx *gin.Context) {
	//获取参数
	// 判断角色是否存在
	var role models.SysRole
	var reqInfo models.ReqSetRoleApiAuth
	err := checkRoleIsExist(ctx, &reqInfo, &role)
	if err != nil {
		return
	}
	var sysApiList []models.CasbinInfo

	for _, apiInfo := range reqInfo.ApiInfos {
		apiInfoList := strings.Split(apiInfo, "|")
		sysApiList = append(sysApiList, models.CasbinInfo{
			Method: apiInfoList[0],
			Path:   apiInfoList[1],
		})
	}

	err = casbin.UpdateCasbins(role.ID, sysApiList)
	if err != nil {
		global.LOG.Errorln("角色绑定Api失败, err: ", err)
		resp.Fail(ctx, nil, "角色绑定Api失败")
		return
	}

	resp.Success(ctx, nil, "绑定角色Api权限成功")
}

// GetRoleApiAuthByID
// @Tags     角色相关
// @Summary  根据角色id来获取api权限
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Success  200   {object}  resp.Result{data=[]models.CasbinInfo}  "根据角色id来获取api权限"
// @Security ApiKeyAuth
// @Router   /auth/roles/{id}/api_auth [get]
func GetRoleApiAuthByID(ctx *gin.Context) {
	//获取接口参数
	id := ctx.Param("id")

	resp.Success(ctx, casbin.GetPolicyPathByAuthorityId(id), "获取角色Api权限成功")
}

// SetRoleMenuAuth
// @Tags     角色相关
// @Summary  通过角色id来绑定菜单权限
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Param    data  body      models.ReqSetRoleApiAuth true "MenuIds[]int"
// @Success  200   {object}  resp.Result{}  "通过角色id来绑定菜单权限"
// @Security ApiKeyAuth
// @Router   /auth/role/{id}/menu_auth [post]
func SetRoleMenuAuth(ctx *gin.Context) {
	//获取参数
	// 判断角色是否存在
	var role models.SysRole
	var reqInfo models.ReqSetRoleMenuAuth
	err := checkRoleIsExist(ctx, &reqInfo, &role)
	if err != nil {
		return
	}

	var menus []models.SysMenu
	if len(reqInfo.MenuInfos) > 0 {
		err = global.DB.Select("id").Find(&menus, reqInfo.MenuInfos).Error
		if err != nil {
			global.LOG.Errorln("角色获取菜单失败, err:", err)
			resp.Fail(ctx, nil, "角色获取菜单失败")
			return
		}
	}

	err = global.DB.Model(&role).Association("Menus").Replace(menus)
	if err != nil {
		global.LOG.Errorln("绑定角色菜单信息错误, err:", err)
		resp.Fail(ctx, nil, "绑定角色菜单信息错误")
		return
	}

	resp.Success(ctx, nil, "绑定角色菜单信息权限成功")
}

// GetRoleMenuAuthByID
// @Tags     角色相关
// @Summary  根据角色id来获取菜单权限
// @Produce  application/json
// @Param    id  path int true "角色ID"
// @Success  200   {object}  resp.Result{data=[]uint}  "根据角色id来获取菜单权限"
// @Security ApiKeyAuth
// @Router   /auth/roles/{id}/menu_auth [get]
func GetRoleMenuAuthByID(ctx *gin.Context) {
	//获取接口参数
	id := ctx.Param("id")
	// 判断角色是否存在
	var role models.SysRole
	err := global.DB.Model(&models.SysRole{}).Preload("Menus").
		Where("id = ?", id).First(&role).Error
	//Select("Menus", "ID").Where("id = ?", id).First(role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "角色不存在")
			return
		}

		global.LOG.Errorln("判断角色是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断角色是否存在报错")
		return
	}

	var menuIds = make([]uint, 0)
	for _, menu := range role.Menus {
		menuIds = append(menuIds, menu.ID)
	}

	resp.Success(ctx, menuIds, "获取角色Api权限成功")
}

func checkRoleIsExist(ctx *gin.Context, reqInfo any, role *models.SysRole) error {
	//获取接口参数
	id := ctx.Param("id")

	if reqInfo != nil {
		err := ctx.ShouldBindJSON(reqInfo)
		if err != nil {
			global.LOG.Errorln("获取reqInfo, err: ", err)
			resp.Fail(ctx, nil, "获取reqInfo失败")
			return err
		}
		global.LOG.Debugln(utils.GetJson(reqInfo))

		// 校验
		err = param_check.VerifyParam(reqInfo)
		if err != nil {
			resp.Fail(ctx, err.Error(), "参数校验失败")
			return err
		}
	}

	db := global.DB

	err := db.Model(&models.SysRole{}).Where("id = ?", id).First(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "角色不存在")
			return err
		}

		global.LOG.Errorln("判断角色是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断角色是否存在报错")
		return err
	}
	return nil
}
