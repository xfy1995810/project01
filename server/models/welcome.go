package models

// RespWelcome 获取主页相关数据
type RespWelcome struct {
	UserCount int64 `json:"user_count"`
}

type TaskStatusCount struct {
	Status int64 `json:"status"`
	Count  int64 `json:"count"`
}
