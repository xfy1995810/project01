package models

type SampleInfos struct {
	Common
	Name   string `gorm:"column:name;type:varchar(255);comment:'样品登记'" json:"name"`
	Num    string `gorm:"column:num;type:varchar(255);comment:'样品编号'" json:"num"`
	Batch  string `gorm:"column:batch;type:text;comment:'样品批号'"   json:"batch"`
	Remark string `gorm:"column:remark;type:varchar(255);comment:'备注'" json:"remark"`
}

// RespGetSampleByID 获取样品信息 响应结构体
type RespGetSampleByID = SampleInfos

// ReqAddSample 新增样品信息 请求结构体
type ReqAddSample struct {
	Name   string `json:"name" validate:"required,min=2,max=50"`
	Num    string `json:"num" validate:"required,min=2,max=50"`
	Batch  string `json:"batch" validate:"required,min=2,max=50"`
	Remark string `json:"remark"`
}

// ReqAddSampleByID 新增样品信息 请求结构体
type ReqAddSampleByID = ReqAddSample

// ReqEditSampleByID 编辑样品信息 请求结构体
type ReqEditSampleByID = ReqAddSample

type RespGetSampleList struct {
	Samples []SampleInfos `json:"Samples"`
	Total   int64         `json:"total"`
}
