package models

import (
	"gorm.io/gorm"
	"time"
)

type SysUser struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Username    string         `gorm:"type:varchar(20);not null;unique;" json:"username"`
	ChineseName string         `gorm:"type:varchar(20);not null;" json:"chinese_name"`
	Password    string         `gorm:"type:varchar(64);not null;" json:"-"`
	Roles       []SysRole      `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"roles"`
	Phone       string         `gorm:"type:varchar(11);not null;" json:"phone"`
	Email       string         `gorm:"type:varchar(64);not null;" json:"email"`
	//  State 1 启用，2 禁用
	State   bool `gorm:"type:bool;not null;default:false;" json:"state"`
	IsSuper bool `gorm:"type:bool;not null;default:false;" json:"is_super"`
}

type CommonSysUser struct {
	ID          uint   `json:"id"`
	Username    string `json:"username" example:"用户名"`
	ChineseName string `json:"chinese_name" example:"中文名"`
}

func (s *SysUser) TableName() string {
	return "sys_users"
}

func (c *CommonSysUser) TableName() string {
	return "sys_users"
}

// ReqAddUser 新增用户 请求结构体
type ReqAddUser struct {
	Username    string `json:"username" example:"用户名" validate:"required,min=2,max=20"`
	ChineseName string `json:"chinese_name" example:"中文名" validate:"required,min=2,max=20"`
	Password    string `json:"password" example:"密码" validate:"required,min=6,max=60"`
	Phone       string `json:"phone" example:"电话号" validate:"required,checkMobile"`
	Email       string `json:"email" example:"邮箱" validate:"required,email"`
	State       bool   `json:"state" validate:"boolean"`
}

// ReqLogin 用户登录 请求结构体
type ReqLogin struct {
	Username  string `json:"username" example:"用户名" validate:"required,min=2,max=20"`
	Password  string `json:"password" example:"密码" validate:"required,min=6,max=60"`
	CaptchaID string `json:"captcha_id" validate:"required"`
	Captcha   string `json:"captcha" example:"623123" validate:"required"`
}

// ReqEditUserByID 修改用户 请求结构体
type ReqEditUserByID struct {
	ChineseName string `json:"chinese_name" example:"中文名,min=2,max=20"`
	Phone       string `json:"phone" example:"电话号" validate:"required,checkMobile"`
	Email       string `json:"email" example:"邮箱" validate:"required,email"`
	State       bool   `json:"state" validate:"boolean"`
}

// ReqEditUserStateByID 修改用户状态 请求结构体
type ReqEditUserStateByID struct {
	State bool `json:"state" validate:"boolean"`
}

// ReqEditUserPasswordByID 修改用户状态 请求结构体
type ReqEditUserPasswordByID struct {
	Password    string `json:"password" example:"当前密码" validate:"required,min=6,max=60"`
	NewPassword string `json:"new_password" example:"新密码" validate:"required,min=6,max=60"`
}

// RespLogin 用户登录 响应结构体
type RespLogin struct {
	ID          uint   `json:"id"`
	AccessToken string `json:"accessToken" example:"访问token"`
	Username    string `json:"username" example:"用户名"`
	ChineseName string `json:"chinese_name" example:"中文名"`
}

// RespGetUserByID 获取用户信息 响应结构体
type RespGetUserByID struct {
	ID          uint            `json:"id"`
	CreatedAt   time.Time       `json:"created_at" example:"创建时间"`
	UpdatedAt   time.Time       `json:"updated_at" example:"更新时间"`
	Username    string          `json:"username" example:"用户名"`
	ChineseName string          `json:"chinese_name" example:"中文名"`
	Phone       string          `json:"phone" example:"手机号"`
	Email       string          `json:"email" example:"邮箱"`
	State       bool            `json:"state"`
	Roles       []CommonSysRole `gorm:"many2many:user_roles;joinForeignKey:SysUserID;joinReferences:SysRoleID;" json:"roles"`
}

// RespGetUserList 获取用户列表 响应结构体
type RespGetUserList struct {
	Users []RespGetUserByID `json:"users"`
	Total int64             `json:"total"`
}

// RespGetAllUserList 获取所有用户列表 响应结构体
type RespGetAllUserList struct {
	Users []struct {
		ID          uint   `json:"id"`
		Username    string `json:"username" example:"用户名"`
		ChineseName string `json:"chinese_name" example:"中文名"`
	} `json:"users"`
	Total int64 `json:"total"`
}
