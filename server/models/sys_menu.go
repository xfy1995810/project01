package models

type SysMenu struct {
	Common
	ParentId  uint      `json:"parentId"`
	Children  []SysMenu `gorm:"foreignkey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"children"`
	Name      string    `gorm:"type:varchar(64);comment:'路由名称'" json:"name"`
	Path      string    `gorm:"type:varchar(256);comment:'路由地址'" json:"path"`
	Component string    `gorm:"type:varchar(256);comment:'文件路径'" json:"component"`
	Meta      Meta      `gorm:"embedded;comment:附加属性" json:"meta"`
	Roles     []SysRole `gorm:"many2many:menu_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"roles"`
}

type Meta struct {
	KeepAlive bool   `gorm:"type:bool;default:false;comment:'开启菜单缓存'" json:"keepAlive"`
	Title     string `gorm:"type:varchar(64);comment:'菜单名称'" json:"title"`
	Icon      string `gorm:"type:varchar(64);comment:'图标'" json:"icon"`
	ShowLink  bool   `gorm:"type:bool;default:false;comment:'是否隐藏'" json:"showLink"`
	Rank      int    `gorm:"type:int;default:100;comment:'菜单排序'" json:"rank,omitempty"`
}

func (s *SysMenu) TableName() string {
	return "sys_menus"
}

// ReqAddMenu 新增菜单 请求结构体
type ReqAddMenu struct {
	ParentId  uint   `json:"parentId"`
	Title     string `json:"title" validate:"required,min=1,max=64"`
	Icon      string `json:"icon" validate:"max=64" `
	Name      string `json:"name" validate:"required,min=1,max=64"`
	Path      string `json:"path" validate:"required,min=1,max=64"`
	Component string `json:"component" validate:"max=64"`
	ShowLink  bool   `json:"showLink" validate:"boolean"`
	Rank      int    `json:"rank" validate:"number,gt=0"`
	KeepAlive bool   `json:"keepAlive" validate:"boolean"`
}

// ReqEditMenuByID 修改菜单 请求结构体
type ReqEditMenuByID struct {
	Title     string `json:"title" validate:"required,min=1,max=64"`
	Icon      string `json:"icon" validate:"max=64" `
	Name      string `json:"name" validate:"required,min=1,max=64"`
	Path      string `json:"path" validate:"required,min=1,max=64"`
	Component string `json:"component" validate:"max=64"`
	ShowLink  bool   `json:"showLink" validate:"boolean"`
	Rank      int    `json:"rank" validate:"number,gt=0"`
	KeepAlive bool   `json:"keepAlive" validate:"boolean"`
}

// RespGetMenuByID 获取菜单信息 响应结构体
type RespGetMenuByID = SysMenu

// RespGetMenuList 获取菜单列表 响应结构体
type RespGetMenuList struct {
	Menus []RespGetMenuByID `json:"menus"`
	//Total int64             `json:"total"`
}
