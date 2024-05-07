package models

type SysRole struct {
	Common
	Name   string    `gorm:"type:varchar(20);not null;unique" json:"name"`
	Remark string    `gorm:"type:varchar(100);comment:'备注'" json:"remark"`
	State  bool      `gorm:"type:bool;not null;default:false;comment:'1启用，0禁用'" json:"state"`
	Users  []SysUser `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"users"`
	Menus  []SysMenu `gorm:"many2many:menu_roles;" json:"menus"` // 角色菜单多对多关系
}

func (s *SysRole) TableName() string {
	return "sys_roles"
}

// ReqAddRole 新增角色 请求结构体
type ReqAddRole struct {
	Name    string `json:"name" validate:"required,min=1,max=20"`
	Remark  string `json:"remark" validate:"min=0,max=100"`
	State   bool   `json:"state" validate:"boolean"`
	UserIds []uint `json:"user_ids"`
}

// ReqEditRoleByID 修改角色 请求结构体
type ReqEditRoleByID struct {
	Name    string `json:"name" validate:"required,min=1,max=20"`
	Remark  string `json:"remark" validate:"min=0,max=100"`
	State   bool   `json:"state" validate:"boolean"`
	UserIds []uint `json:"user_ids"`
}

// ReqEditRoleStateByID 修改角色状态 请求结构体
type ReqEditRoleStateByID struct {
	State bool `json:"state" validate:"boolean"`
}

// ReqSetRoleApiAuth 通过角色id来绑定Api权限 请求结构体
type ReqSetRoleApiAuth struct {
	ApiInfos []string `json:"api_infos" validate:"required"`
}

//roleId: number;
//apiIds: Array<number>;

// ReqSetRoleMenuAuth 通过角色id来绑定Api权限 请求结构体
type ReqSetRoleMenuAuth struct {
	MenuInfos []uint `json:"menu_infos" validate:"required"`
}

// RespGetRoleByID 获取角色信息 响应结构体
// type RespGetRoleByID = SysRole
type RespGetRoleByID struct {
	SysRole
	Users []CommonSysUser `gorm:"many2many:user_roles;joinForeignKey:SysRoleID;joinReferences:SysUserID;" json:"users"`
}

// RespGetRoleList 获取角色列表 响应结构体
type RespGetRoleList struct {
	Roles []RespGetRoleByID `json:"roles"`
	Total int64             `json:"total"`
}

type CommonSysRole struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
	State  bool   `json:"state"`
}

func (c *CommonSysRole) TableName() string {
	return "sys_roles"
}
