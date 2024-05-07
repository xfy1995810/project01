package menu

import (
	"dcss/global"
	"dcss/global/param_check"
	"dcss/models"
	"dcss/pkg/resp"
	"dcss/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddMenu
// @Tags     菜单相关
// @Summary  添加菜单
// @Produce  application/json
// @Param    data  body      models.ReqAddMenu true "菜单名称 图标 路由名称 路由地址 文件路径 是否隐藏 菜单排序 开启菜单缓存"
// @Success  200   {object}  resp.Result{}  "添加菜单"
// @Security MenuKeyAuth
// @Router   /system/menu [post]
func AddMenu(ctx *gin.Context) {
	//获取参数
	var reqInfo models.ReqAddMenu
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorf("reqInfo bind err: %v", err)
		resp.Fail(ctx, nil, err.Error())
		return
	}
	global.LOG.Debugf("AddMenu,reqInfo: %v\n", utils.GetJson(reqInfo))

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	//  判断菜单是否存在
	var menu models.SysMenu
	// 获取数据库连接池
	db := global.DB

	err = db.Model(&models.SysMenu{}).Where("name = ?", reqInfo.Name).First(&menu).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Errorln("判断菜单是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断菜单是否存在报错")
		return
	}

	if menu.ID != 0 {
		resp.Fail(ctx, nil, "菜单名已存在")
		return
	}

	menu = models.SysMenu{
		ParentId:  reqInfo.ParentId,
		Name:      reqInfo.Name,
		Path:      reqInfo.Path,
		Component: reqInfo.Component,
		Meta: models.Meta{
			Icon:      reqInfo.Icon,
			ShowLink:  reqInfo.ShowLink,
			Rank:      reqInfo.Rank,
			KeepAlive: reqInfo.KeepAlive,
			Title:     reqInfo.Title,
		},
		Common: models.Common{
			CreatedID: ctx.GetUint("userID"),
		},
	}
	global.LOG.Debugln("新增菜单信息: ", utils.GetJson(menu))

	//	创建菜单
	err = db.Create(&menu).Error
	if err != nil {
		global.LOG.Errorln("创建菜单失败,err:", err)
		resp.Fail(ctx, nil, "创建菜单失败")
		return
	}

	resp.Response(ctx, 201, nil, "创建菜单成功")
}

// GetMenuList
// @Tags     菜单相关
// @Summary  获取所有的菜单
// @Produce  application/json
// @Success  200   {object}  resp.Result{data=models.RespGetMenuList}  "获取所有的菜单"
// @Security MenuKeyAuth
// @Router   /system/menus [get]
func GetMenuList(ctx *gin.Context) {
	menus, err := GetMenuTree(true)
	if err != nil {
		global.LOG.Errorln("获取菜单列表失败, err:", err)
		resp.Fail(ctx, nil, "获取菜单列表失败")
		return
	}

	resp.Success(ctx, models.RespGetMenuList{
		Menus: menus,
	}, "获取菜单列表成功")
}

// GetMenuByID
// @Tags     菜单相关
// @Summary  通过id来获取菜单信息
// @Produce  application/json
// @Param    id  path int true "MenuID"
// @Success  200   {object}  resp.Result{data=models.RespGetMenuByID}  "通过id来获取菜单信息"
// @Security MenuKeyAuth
// @Router   /system/menu/{id} [get]
func GetMenuByID(ctx *gin.Context) {
	//	获取菜单Id
	idStr := ctx.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	menu, err := GetMenuTreeByMenuId(uint(id))
	if err != nil {
		global.LOG.Errorln("获取菜单信息失败, err:", err)
		resp.Fail(ctx, nil, "获取菜单信息失败")
		return
	}

	resp.Success(ctx, menu, "获取菜单信息成功")
}

// EditMenuByID
// @Tags     菜单相关
// @Summary  通过id来修改菜单信息
// @Produce  application/json
// @Param    id  path int true "MenuID"
// @Param    data  body      models.ReqEditMenuByID true "菜单名称 图标 路由名称 路由地址 文件路径 是否隐藏 菜单排序 开启菜单缓存"
// @Success  200   {object}  resp.Result{}  "通过id来修改菜单信息"
// @Security MenuKeyAuth
// @Router   /system/menu/{id} [put]
func EditMenuByID(ctx *gin.Context) {
	//获取接口参数
	id := ctx.Param("id")

	var reqInfo models.ReqEditMenuByID
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorf("reqInfo bind err: %v", err)
		resp.Fail(ctx, nil, err.Error())
		return
	}
	global.LOG.Debugf("EditMenu,reqInfo: %v\n", utils.GetJson(reqInfo))

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	// 判断菜单是否存在
	var menu models.SysMenu
	err = global.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "菜单不存在")
			return
		}

		global.LOG.Errorln("判断菜单是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断菜单是否存在报错")
		return
	}

	// 判断修改的菜单名是否存在
	var count int64
	err = global.DB.Model(&models.SysMenu{}).Where("name = ?", reqInfo.Name).Not("id = ?", id).Count(&count).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Errorln("判断菜单是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断菜单是否存在报错")
		return
	}

	if count > 0 {
		resp.Fail(ctx, nil, "菜单名已存在")
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.SysMenu{}).Where("id = ?", menu.ID).
			Updates(map[string]interface{}{
				"name":       reqInfo.Name,
				"path":       reqInfo.Path,
				"component":  reqInfo.Component,
				"icon":       reqInfo.Icon,
				"show_link":  reqInfo.ShowLink,
				"rank":       reqInfo.Rank,
				"keep_alive": reqInfo.KeepAlive,
				"title":      reqInfo.Title,
				"UpdatedID":  ctx.GetUint("userID"),
			}).Error; err != nil {

			return err
		}

		return nil
	})

	if err != nil {
		global.LOG.Errorln("更新菜单失败, err:", err)
		resp.Fail(ctx, nil, "更新菜单失败")
		return
	}

	resp.Success(ctx, nil, "更新菜单信息成功")
}

// DeleteMenuByID
// @Tags     菜单相关
// @Summary  通过id来删除菜单
// @Produce  application/json
// @Param    id  path int true "MenuID"
// @Success  200   {object}  resp.Result{}  "通过id来删除菜单"
// @Security MenuKeyAuth
// @Router   /system/menu/{id} [delete]
func DeleteMenuByID(ctx *gin.Context) {
	//获取接口参数
	id := ctx.Param("id")

	// 判断菜单是否存在
	var menu models.SysMenu
	err := global.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "菜单不存在")
			return
		}

		global.LOG.Errorln("判断菜单是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断菜单是否存在报错")
		return
	}
	idSlice, err := GetSelfAndChildrenIdList(&menu)
	if err != nil {
		global.LOG.Errorln("获取相关菜单ID列表失败, err:", err)
		resp.Fail(ctx, nil, "获取相关菜单ID列表失败")
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Unscoped().Delete(&models.SysMenu{}, &idSlice).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.LOG.Errorln("菜单删除失败, err:", err)
		resp.Fail(ctx, nil, "菜单删除失败")
		return
	}

	resp.Success(ctx, nil, "删除菜单信息成功")
}

// GetDynamicMenuList
// @Tags     菜单相关
// @Summary  获取当前用户的动态菜单
// @Produce  application/json
// @Success  200   {object}  resp.Result{data=models.SysMenu}  "获取当前用户的动态菜单"
// @Security MenuKeyAuth
// @Router   /system/dynamic_menus [get]
func GetDynamicMenuList(ctx *gin.Context) {
	u, exist := ctx.Get("user")
	if !exist {
		global.LOG.Errorln("获取动态菜单失败，ctx获取user失败")
		resp.Fail(ctx, nil, "获取动态菜单失败")
		return
	}
	user, ok := u.(*models.SysUser)
	if !ok {
		global.LOG.Errorln("获取动态菜单失败，类型断言user失败")
		resp.Fail(ctx, nil, "获取动态菜单失败")
		return
	}

	var menus []models.SysMenu
	var err error

	if user.IsSuper {
		menus, err = GetMenuTree(false)
		if err != nil {
			global.LOG.Errorln("获取菜单列表失败, err:", err)
			resp.Fail(ctx, nil, "获取菜单列表失败")
			return
		}

		resp.Success(ctx, menus, "获取菜单列表成功")
		return
	}

	var userRoles = user.Roles
	var haveSubMenus []models.SysMenu

	err = global.DB.Model(&userRoles).Association("Menus").Find(&haveSubMenus)
	if err != nil {
		global.LOG.Errorln("获取角色菜单列表失败, err:", err)
		resp.Fail(ctx, nil, "获取角色菜单列表失败")
		return
	}
	global.LOG.Debugln("haveSubMenus: ", utils.GetJson(haveSubMenus))

	menus, err = GetMenuTreeBySubMenuList(haveSubMenus)
	if err != nil {
		global.LOG.Errorln("获取菜单列表失败, err:", err)
		resp.Fail(ctx, nil, "获取菜单列表失败")
		return
	}
	global.LOG.Debugln("haveMenus: ", utils.GetJson(menus))

	resp.Success(ctx, menus, "获取菜单列表成功")
}
