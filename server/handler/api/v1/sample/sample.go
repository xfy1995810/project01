package sample

import (
	"dcss/global"
	"dcss/global/param_check"
	"dcss/models"
	"dcss/pkg/resp"
	"dcss/pkg/utils"

	//"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetSampleObjList(ctx *gin.Context) {
	db := global.DB.Model(&models.SampleInfos{})
	// 定义数据库接收的资产结构体变量
	var sampleInfoList []models.SampleInfos
	var err error
	// 样品信息总数量
	var total int64

	err = db.Count(&total).Error
	if err != nil {
		global.LOG.Errorln("获取样品信息列表总数失败, err:", err)
		resp.Fail(ctx, nil, "获取样品信息列表总数失败")
		return
	}

	err = db.Scopes(models.HandleCommonUserPreload).
		Order("-id").
		Omit("deleted_at").
		Find(&sampleInfoList).Error
	if err != nil {
		global.LOG.Errorln("获取样品信息列表失败, err:", err)
		resp.Fail(ctx, nil, "获取样品信息列表失败")
		return
	}

	resp.Success(ctx, models.RespGetSampleList{
		Samples: sampleInfoList,
		Total:   total,
	}, "获取样品信息列表成功")
}

func AddSampleObj(ctx *gin.Context) {
	var reqInfo models.ReqAddSample
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorf("reqInfo bind err: %v", err)
		resp.Fail(ctx, nil, err.Error())
		return
	}

	//  判断样品信息是否存在
	var tmpSampleInfos models.SampleInfos
	// 获取数据库连接池
	db := global.DB
	err = db.Model(&models.SampleInfos{}).Where("name = ?", reqInfo.Name).First(&tmpSampleInfos).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Errorln("判断添加样品信息任务是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断添加样品信息任务是否存在报错")
		return
	}

	if tmpSampleInfos.ID != 0 {
		resp.Fail(ctx, nil, "样品登记已存在")
		return
	}

	SampleInfos := models.SampleInfos{
		Name:   reqInfo.Name,
		Num:    reqInfo.Num,
		Batch:  reqInfo.Batch,
		Remark: reqInfo.Remark,
		Common: models.Common{
			CreatedID: ctx.GetUint("userID"),
		},
	}

	global.LOG.WithFields(logrus.Fields{
		"data": utils.GetJson(&SampleInfos),
	}).Debugln("新增样品信息任务")

	//      创建样品信息任务
	err = db.Create(&SampleInfos).Error
	if err != nil {
		global.LOG.Errorln("创建样品信息任务失败,err:", err)
		resp.Fail(ctx, nil, "创建样品信息任务失败")
		return
	}

	resp.Response(ctx, 200, nil, "创建样品信息任务成功")
}

func DeleteSampleObjByID(ctx *gin.Context) {
	// 样品信息是否存在
	var SampleInfos models.SampleInfos
	err := checkFileInfoIsExist(ctx, nil, &SampleInfos)
	if err != nil {
		return
	}

	err = global.DB.Unscoped().Delete(&SampleInfos).Error
	if err != nil {
		resp.Fail(ctx, nil, fmt.Sprintf("样品信息删除失败, err: %s", err))
		return
	}

	resp.Success(ctx, nil, "删除样品信息信息成功")
}

func checkFileInfoIsExist(ctx *gin.Context, reqInfo any, asset *models.SampleInfos) error {
	//获取接口参数
	id, _ := strconv.Atoi(ctx.Param("id"))

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

	err := db.Model(&models.SampleInfos{}).Where("id = ?", id).First(asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "样品信息项id不存在")
			return err
		}

		global.LOG.Errorln("文件配置任务是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断文件配置任务是否存在报错")
		return err
	}
	return nil
}

func EditSampleObjByID(ctx *gin.Context) {
	var reqInfo models.ReqAddSample
	var SampleInfos models.SampleInfos
	err := checkFileInfoIsExist(ctx, &reqInfo, &SampleInfos)
	if err != nil {
		return
	}
	var checkNameCount int64
	err = global.DB.Model(&models.SampleInfos{}).Where("name = ?", reqInfo.Name).Not("id = ?", SampleInfos.ID).Count(&checkNameCount).Error
	if err != nil {
		global.LOG.Errorln("样品信息获取失败, err:", err)
		resp.Fail(ctx, nil, "样品信息获取失败")
		return
	}
	if checkNameCount > 0 {
		resp.Fail(ctx, nil, "样品信息任务名重复，修改失败")
		return
	}

	if err := global.DB.Model(&models.SampleInfos{}).Where("id = ?", SampleInfos.ID).
		Updates(map[string]interface{}{
			"name":   reqInfo.Name,
			"num":    reqInfo.Num,
			"batch":  reqInfo.Batch,
			"remark": reqInfo.Remark,
		}).Error; err != nil {

		global.LOG.Errorln(fmt.Sprintf("更新样品信息任务错误,err: %s", err))
		resp.Fail(ctx, nil, "更新样品信息任务失败")
		return
	}

	resp.Success(ctx, nil, "更新样品信息任务信息成功")

}

// func WriteFileContext(FileName,FileContext string )(string ,error) {
//     workDir, err := os.Getwd()
//     if err != nil {
//         return "Error getting current directory:",err
//     }
//     // 拼接当前工作目录和子目录'FileInfo_setting_dirs/'
//     LocalDirs := path.Join(workDir, "FileInfo_setting_dirs/")

//     if _, err := os.Stat(LocalDirs); os.IsNotExist(err) {
// 		err := os.MkdirAll(LocalDirs, 0755)
// 		if err != nil {
// 			panic("create FileInfo_setting_dirs failed, err: " + err.Error())
// 		}
// 	}

//     LocalFilePath := path.Join(LocalDirs, FileName)

//     // 打开文件
//     file, err0 := os.OpenFile(LocalFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0744)
//     if err0 != nil {
//         return "",err0
//     }
//     defer file.Close()

//     // 将命令写入脚本文件
//     context := []byte(FileContext)
//     errw := ioutil.WriteFile(LocalFilePath, context, 0744)
//     if errw != nil {
//       fmt.Println("无法写入内容到本地文件路径:", errw)
//       return "无法写入内容到本地文件路径",errw
//     }else{
//       return LocalFilePath,nil
//     }
// }
