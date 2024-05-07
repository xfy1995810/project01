package menu

import (
	"dcss/global"
	"dcss/models"
	"errors"
)

// GetMenuTree 获取动态菜单树
func GetMenuTree(needSubRank bool) (menus []models.SysMenu, err error) {
	treeMap, err := getMenuTreeMap(needSubRank)
	if err != nil {
		return menus, err
	}

	menus = treeMap[0]
	for i := 0; i < len(menus); i++ {
		getChildrenList(&menus[i], treeMap)
	}
	return menus, err
	/*
		menus := []any{
			gin.H{
				"path": "/system",
				"meta": gin.H{
					"title": "动态菜单",
					"icon":  "ep:eleme",
					"rank":  5,
				},
				"children": []any{
					gin.H{
						"path":      "/system/api",
						"component": "/system/api/index",
						"name":      "api",
						"meta": gin.H{
							"title": "api",
						},
					},
					gin.H{
						"path":      "/system/menu",
						"component": "/system/menu/index",
						"name":      "menu",
						"meta": gin.H{
							"title": "菜单",
						},
					},
				},
			}
	*/
}

// GetMenuTreeByMenuId 获取指定菜单id的动态菜单树
func GetMenuTreeByMenuId(id uint) (menu models.SysMenu, err error) {
	allMenus, err := getAllMenus(true)
	if err != nil {
		return models.SysMenu{}, err
	}

	treeMap := make(map[uint][]models.SysMenu, 0)
	for _, v := range allMenus {
		if v.ID == id {
			menu = v
		}
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}

	menus, ok := treeMap[id]
	if ok {
		for i := 0; i < len(menus); i++ {
			getChildrenList(&menus[i], treeMap)
		}
		menu.Children = menus
	}

	if menu.ID == 0 {
		return menu, errors.New("指定ID菜单不存在")
	}

	return menu, err
}

// GetMenuTreeBySubMenuList 获取动态菜单树
func GetMenuTreeBySubMenuList(subMenus []models.SysMenu) (menus []models.SysMenu, err error) {
	allMenus, err := getAllMenus(false)
	if err != nil {
		return nil, err
	}

	allMenusMap := make(map[uint]models.SysMenu)
	treeMap := make(map[uint]models.SysMenu)
	for _, v := range allMenus {
		allMenusMap[v.ID] = v
	}

	treeMap = getParentMenuList(subMenus, allMenusMap, treeMap)
	rootMenu, ok := treeMap[0]
	if !ok {
		return menus, errors.New("获取动态菜单树失败")
	}

	return rootMenu.Children, nil
}

// getMenuTreeMap 获取路由总树map
func getMenuTreeMap(needSubRank bool) (treeMap map[uint][]models.SysMenu, err error) {
	allMenus, err := getAllMenus(needSubRank)
	if err != nil {
		return nil, err
	}
	treeMap = make(map[uint][]models.SysMenu)
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

// getChildrenList 获取子菜单
func getChildrenList(menu *models.SysMenu, treeMap map[uint][]models.SysMenu) {
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		getChildrenList(&menu.Children[i], treeMap)
	}
}

// getParentMenuList 根据子菜单列表获取父菜单列表
func getParentMenuList(subMenus []models.SysMenu, allMenusMap, treeMap map[uint]models.SysMenu) map[uint]models.SysMenu {
	var menus []models.SysMenu

	for _, subMenu := range subMenus {
		if subMenu.ID == 0 {
			continue
		}

		// 二级菜单的rank不能有值,有值 前端有bug
		if subMenu.ParentId != 0 {
			subMenu.Meta.Rank = 0
		}

		tmpSubMenu, ok := treeMap[subMenu.ID]
		if ok {
			subMenu = tmpSubMenu
		}

		parentMenu, ok := treeMap[subMenu.ParentId]
		if !ok {
			parentMenu = allMenusMap[subMenu.ParentId]
		}
		menus = append(menus, parentMenu)

		// 如果主菜单里面有对应id的子菜单，则进行替换，不追加
		menuChildrenIsExist := false
		for i, c := range parentMenu.Children {
			if c.ID == subMenu.ID {
				menuChildrenIsExist = true
				parentMenu.Children[i] = subMenu
			}
		}
		if !menuChildrenIsExist {
			parentMenu.Children = append(parentMenu.Children, subMenu)
		}

		treeMap[subMenu.ParentId] = parentMenu
	}

	if len(menus) > 0 {
		treeMap = getParentMenuList(menus, allMenusMap, treeMap)
	}

	return treeMap
}

// getChildrenIdList 获取自身和子菜单的id列表
func getChildrenIdList(treeMap map[uint][]models.SysMenu, menuId uint, idSlice []uint) []uint {
	idSlice = append(idSlice, menuId)
	for i := 0; i < len(treeMap[menuId]); i++ {
		idSlice = getChildrenIdList(treeMap, treeMap[menuId][i].ID, idSlice)
	}
	return idSlice
}

func GetSelfAndChildrenIdList(menu *models.SysMenu) (idSlice []uint, err error) {
	treeMap, err := getMenuTreeMap(true)
	if err != nil {
		return idSlice, err
	}

	idSlice = getChildrenIdList(treeMap, menu.ID, idSlice)
	return idSlice, err
}

// getAllMenus 从数据库中获取所有的菜单信息
func getAllMenus(needRank bool) ([]models.SysMenu, error) {
	var allMenus []models.SysMenu
	err := global.DB.Order("rank").Scopes(models.HandleCommonUserPreload).Find(&allMenus).Error
	if err != nil {
		return nil, err
	}
	// 二级菜单的rank不能有值,有值 前端有bug
	if !needRank {
		for i, v := range allMenus {
			if v.ParentId != 0 {
				allMenus[i].Meta.Rank = 0
			}
		}
	}

	return allMenus, nil
}
