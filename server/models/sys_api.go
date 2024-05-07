package models

type SysApi struct {
	Common
	Name     string `gorm:"type:varchar(100);comment:'名称'" json:"name"`
	Method   string `gorm:"type:varchar(20);comment:'请求方式'" json:"method"`
	Path     string `gorm:"type:varchar(100);comment:'访问路径'" json:"path"`
	Category string `gorm:"type:varchar(50);comment:'所属类别'" json:"category"`
}

func (s *SysApi) TableName() string {
	return "sys_apis"
}

// ReqAddApi 新增接口 请求结构体
type ReqAddApi struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Method   string `json:"method" validate:"required,oneof=GET POST DELETE PUT"`
	Path     string `json:"path" validate:"required,min=1,max=100"`
	Category string `json:"category" validate:"required,min=1,max=50"`
}

// ReqEditApiByID 修改接口 请求结构体
type ReqEditApiByID struct {
	Name     string `json:"name" validate:"min=1,max=100"`
	Method   string `json:"method" validate:"required,oneof=GET POST DELETE PUT"`
	Path     string `json:"path" validate:"min=1,max=100"`
	Category string `json:"category" validate:"min=1,max=50"`
}

// RespGetApiByID 获取接口信息 响应结构体
type RespGetApiByID = SysApi

// RespGetApiList 获取接口列表 响应结构体
type RespGetApiList struct {
	Apis  []RespGetApiByID `json:"apis"`
	Total int64            `json:"total"`
}

// RespGetAllApiList 获取所有接口列表 响应结构体
type RespGetAllApiList struct {
	Apis []struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Method   string `json:"method"`
		Path     string `json:"path"`
		Category string `json:"category"`
	} `json:"apis"`
	Categories []string `json:"categories"`
	Total      int64    `json:"total"`
}

type CasbinInfo struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
