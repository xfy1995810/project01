package db

import (
	"dcss/global"
	"dcss/models"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// initUserData 初始化超级管理员账号
func initUserData() error {
	var user models.SysUser
	err := global.DB.First(&user, "username = ?", "xiaotangyuan").Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Infoln("管理员账号不存在，新建中...")
		password := "xty0731"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash richmail user password failed, err: %s", err)
		}

		err = global.DB.Create(&models.SysUser{
			ID:          1,
			Username:    "xiaotangyuan",
			Password:    string(hashedPassword),
			ChineseName: "小汤圆",
			Phone:       "8888888",
			Email:       "vip888888@xty.cn",
			State:       true,
			IsSuper:     true,
		}).Error

		return err

	}
	return nil
}

// initApiData 初始化api数据
func initApiData() error {
	apis := []*models.SysApi{
		// 主页
		{
			Method:   "GET",
			Path:     "/welcome/",
			Category: "主页",
			Name:     "欢迎页",
		},
		// 用户
		{
			Method:   "GET",
			Path:     "/auth/users/",
			Category: "用户",
			Name:     "获取用户列表",
		},
		{
			Method:   "POST",
			Path:     "/auth/user/",
			Category: "用户",
			Name:     "添加用户",
		},
		{
			Method:   "GET",
			Path:     "/auth/user/:id/",
			Category: "用户",
			Name:     "获取指定id用户",
		},
		{
			Method:   "PUT",
			Path:     "/auth/user/:id/",
			Category: "用户",
			Name:     "修改用户",
		},
		{
			Method:   "DELETE",
			Path:     "/auth/user/:id/",
			Category: "用户",
			Name:     "删除用户",
		},
		{
			Method:   "PUT",
			Path:     "/auth/user/:id/state/",
			Category: "用户",
			Name:     "修改用户状态",
		},
		{
			Method:   "PUT",
			Path:     "/auth/user/:id/password/",
			Category: "用户",
			Name:     "修改用户密码",
		},
		{
			Method:   "GET",
			Path:     "/auth/all_users/",
			Category: "角色",
			Name:     "获取所有用户(供角色选取用户使用)",
		},
		// 角色
		{
			Method:   "GET",
			Path:     "/auth/roles/",
			Category: "角色",
			Name:     "获取角色列表",
		},
		{
			Method:   "GET",
			Path:     "/auth/role/:id/",
			Category: "角色",
			Name:     "获取指定id角色信息",
		},
		{
			Method:   "POST",
			Path:     "/auth/role/",
			Category: "角色",
			Name:     "创建角色",
		},
		{
			Method:   "PUT",
			Path:     "/auth/role/:id/",
			Category: "角色",
			Name:     "修改角色",
		},
		{
			Method:   "DELETE",
			Path:     "/auth/role/:id/",
			Category: "角色",
			Name:     "删除角色",
		},
		{
			Method:   "PUT",
			Path:     "/auth/role/:id/state/",
			Category: "角色",
			Name:     "修改角色状态",
		},
		{
			Method:   "POST",
			Path:     "/auth/role/:id/api_auth/",
			Category: "角色",
			Name:     "设置角色跟Api的权限关系",
		},
		{
			Method:   "GET",
			Path:     "/auth/role/:id/api_auth/",
			Category: "角色",
			Name:     "获取角色跟Api的权限关系",
		},
		{
			Method:   "POST",
			Path:     "/auth/role/:id/menu_auth/",
			Category: "角色",
			Name:     "设置角色跟菜单的权限关系",
		},
		{
			Method:   "GET",
			Path:     "/auth/role/:id/menu_auth/",
			Category: "角色",
			Name:     "获取角色跟菜单的权限关系",
		},
		{
			Method:   "GET",
			Path:     "/auth/all_apis/",
			Category: "角色",
			Name:     "获取所有的api(供角色赋权使用)",
		},
		// api
		{
			Method:   "GET",
			Path:     "/system/apis/",
			Category: "api",
			Name:     "获取api列表",
		},
		{
			Method:   "GET",
			Path:     "/system/api/:id/",
			Category: "api",
			Name:     "获取指定id的api列表",
		},
		{
			Method:   "POST",
			Path:     "/system/api/",
			Category: "api",
			Name:     "创建api",
		},
		{
			Method:   "PUT",
			Path:     "/system/api/:id/",
			Category: "api",
			Name:     "修改api",
		},
		{
			Method:   "DELETE",
			Path:     "/system/api/:id/",
			Category: "api",
			Name:     "删除api",
		},
		// 菜单
		{
			Method:   "GET",
			Path:     "/system/menus/",
			Category: "菜单",
			Name:     "获取菜单列表",
		},
		{
			Method:   "GET",
			Path:     "/system/menu/:id/",
			Category: "菜单",
			Name:     "获取指定id的菜单列表",
		},
		{
			Method:   "POST",
			Path:     "/system/menu/",
			Category: "菜单",
			Name:     "创建菜单",
		},
		{
			Method:   "PUT",
			Path:     "/system/menu/:id/",
			Category: "菜单",
			Name:     "修改菜单",
		},
		{
			Method:   "DELETE",
			Path:     "/system/menu/:id/",
			Category: "菜单",
			Name:     "删除菜单",
		},
		//样品
		{
			Method:   "GET",
			Path:     "/sample/sample/:id/",
			Category: "样品",
			Name:     "获取指定id的样品信息列表",
		},
		{
			Method:   "POST",
			Path:     "/sample/sample/",
			Category: "样品",
			Name:     "创建样品信息",
		},
		{
			Method:   "PUT",
			Path:     "/sample/sample/:id/",
			Category: "样品",
			Name:     "修改样品信息",
		},
		{
			Method:   "DELETE",
			Path:     "/sample/sample/:id/",
			Category: "样品",
			Name:     "删除样品信息",
		},
	}

	var api models.SysApi
	err := global.DB.First(&api, "path = ?", "/welcome/").Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Infoln("Api数据不存在，新建中...")
		for i := range apis {
			apis[i].CreatedID = 1
			if !strings.HasSuffix(apis[i].Path, "/") {
				apis[i].Path += "/"
			}
		}

		return global.DB.Create(apis).Error
	}
	return err
}

// initMenuData 初始化菜单数据
func initMenuData() error {
	menus := []models.SysMenu{

		{ParentId: 0, Name: "auth", Path: "/auth", Meta: models.Meta{Title: "权限管理", Icon: "ep:hide", Rank: 7, ShowLink: true}},
		{ParentId: 1, Name: "user", Path: "/auth/user", Component: "/auth/user/index.vue", Meta: models.Meta{Title: "用户管理", ShowLink: true, Rank: 1}},
		{ParentId: 1, Name: "role", Path: "/auth/role", Component: "/auth/role/index.vue", Meta: models.Meta{Title: "角色管理", ShowLink: true, Rank: 2}},
		{ParentId: 0, Name: "system", Path: "/system", Meta: models.Meta{Title: "系统管理", Icon: "ep:setting", Rank: 8, ShowLink: true}},
		{ParentId: 4, Name: "api", Path: "/system/api", Component: "/system/api/index.vue", Meta: models.Meta{Title: "Api管理", ShowLink: true, Rank: 1}},
		{ParentId: 0, Name: "xiaoTY", Path: "/sample", Meta: models.Meta{Title: "样品管理", Icon: "ep:data-analysis", Rank: 5, ShowLink: true}},
		{ParentId: 6, Name: "info", Path: "/sample/info", Component: "/sample/info/index.vue", Meta: models.Meta{Title: "信息管理", ShowLink: true, Rank: 1}},
		{ParentId: 6, Name: "data", Path: "/sample/data", Component: "/sample/data/index.vue", Meta: models.Meta{Title: "数据管理", ShowLink: true, Rank: 2}},
		{ParentId: 4, Name: "menus", Path: "/system/menu", Component: "/system/menus/index.vue", Meta: models.Meta{Title: "菜单管理", ShowLink: true, Rank: 2}},
		{ParentId: 4, Name: "setting", Path: "/system/setting", Component: "/system/setting/index.vue", Meta: models.Meta{Title: "通用设置", ShowLink: true, Rank: 3}},
	}

	var menu models.SysMenu
	err := global.DB.First(&menu, "name = ?", "user").Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Infoln("菜单数据不存在，新建中...")
		for i := range menus {
			menus[i].CreatedID = 1
		}

		return global.DB.Create(menus).Error
	}
	return err
}
