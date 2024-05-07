package api

import (
	"dcss/global"
	"dcss/global/casbin"
	"dcss/global/param_check"
	"dcss/models"
	"dcss/pkg/resp"
	"dcss/pkg/utils"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddApi
// @Tags     Api相关
// @Summary  添加Api
// @Produce  application/json
// @Param    data  body      models.ReqAddApi true "方法 路径 目录 备注"
// @Success  200   {object}  resp.Result{}  "添加Api"
// @Security ApiKeyAuth
// @Router   /system/api [post]
func AddApi(ctx *gin.Context) {
	//获取参数
	var reqInfo models.ReqAddApi
	err := ctx.ShouldBindJSON(&reqInfo)
	if err != nil {
		global.LOG.Errorf("reqInfo bind err: %v", err)
		resp.Fail(ctx, nil, err.Error())
		return
	}
	global.LOG.Debugf("AddApi,reqInfo: %v\n", reqInfo)

	// 参数校验
	err = param_check.VerifyParam(reqInfo)
	if err != nil {
		resp.Fail(ctx, err.Error(), "参数校验失败")
		return
	}

	//  判断Api是否存在
	var api models.SysApi
	// 获取数据库连接池
	db := global.DB

	err = db.Model(&models.SysApi{}).Where("name = ?", reqInfo.Name).First(&api).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		global.LOG.Errorln("判断Api是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断Api是否存在报错")
		return
	}

	if api.ID != 0 {
		resp.Fail(ctx, nil, "Api名已存在")
		return
	}

	if !strings.HasSuffix(reqInfo.Path, "/") {
		reqInfo.Path += "/"
	}

	api = models.SysApi{
		Name:     reqInfo.Name,
		Method:   strings.ToUpper(reqInfo.Method),
		Path:     reqInfo.Path,
		Category: reqInfo.Category,
		Common: models.Common{
			CreatedID: ctx.GetUint("userID"),
		},
	}
	global.LOG.Debugln("新增Api信息: ", utils.GetJson(api))

	//	创建Api
	err = db.Create(&api).Error
	if err != nil {
		global.LOG.Errorln("创建Api失败,err:", err)
		resp.Fail(ctx, nil, "创建Api失败")
		return
	}

	resp.Response(ctx, 201, nil, "创建Api成功")
}

// GetApiList
// @Tags     Api相关
// @Summary  获取所有的Api
// @Produce  application/json
// @Param    query  query  string     false  "过滤关键字"
// @Param    pagenum   query   int     false  "页码"
// @Param    pagesize   query    int     false  "返回数量"
// @Success  200   {object}  resp.Result{data=models.RespGetApiList}  "获取所有的Api"
// @Security ApiKeyAuth
// @Router   /system/apis [get]
func GetApiList(ctx *gin.Context) {
	db := global.DB.Model(&models.SysApi{})
	// 从请求的接口中获取相关数据
	queryStr := ctx.DefaultQuery("query", "")
	pagenum, _ := strconv.Atoi(ctx.DefaultQuery("pagenum", "1"))
	pagenum = pagenum - 1
	pagesize, _ := strconv.Atoi(ctx.DefaultQuery("pagesize", "10"))
	// 定义数据库接收的api结构体变量
	var apiList []models.RespGetApiByID
	var err error
	// Api总数量
	var total int64

	if pagesize > 100 || pagesize < 0 {
		pagesize = 20
	}

	if len(queryStr) != 0 {
		queryStr = "%" + queryStr + "%"
		db = db.Where("name LIKE ?", queryStr).
			Or("remark LIKE ?", queryStr)
	}

	err = db.Count(&total).Error
	if err != nil {
		global.LOG.Errorln("获取Api列表总数失败, err:", err)
		resp.Fail(ctx, nil, "获取Api列表总数失败")
		return
	}

	err = db.Scopes(models.HandleCommonUserPreload).
		Offset(pagenum * pagesize).
		Limit(pagesize).
		Order("id").
		Omit("deleted_at").
		Find(&apiList).Error
	if err != nil {
		global.LOG.Errorln("获取Api列表失败, err:", err)
		resp.Fail(ctx, nil, "获取Api列表失败")
		return
	}

	resp.Success(ctx, models.RespGetApiList{
		Apis:  apiList,
		Total: total,
	}, "获取Api列表成功")
}

// GetAllApiList
// @Tags     Api相关
// @Summary  获取所有的Api,提供给角色添加权限使用
// @Produce  application/json
// @Success  200   {object}  resp.Result{data=models.RespGetAllApiList}  "获取所有的Api"
// @Security ApiKeyAuth
// @Router   /system/all_apis [get]
func GetAllApiList(ctx *gin.Context) {
	db := global.DB
	// 定义数据库接收的user结构体变量
	var respApiList models.RespGetAllApiList
	var err error

	err = db.Model(&models.SysApi{}).Count(&respApiList.Total).Error
	if err != nil {
		global.LOG.Errorln("获取所有的Api列表总数失败, err:", err)
		resp.Fail(ctx, nil, "获取所有的Api列表总数失败")
		return
	}

	err = db.Model(&models.SysApi{}).Distinct("category").Find(&respApiList.Categories).Error
	if err != nil {
		global.LOG.Errorln("获取所有的Api 目录列表失败, err:", err)
		resp.Fail(ctx, nil, "获取所有的Api 目录列表失败")
		return
	}

	err = db.Model(&models.SysApi{}).Order("category").
		Find(&respApiList.Apis).Error
	if err != nil {
		global.LOG.Errorln("获取Api列表失败, err:", err)
		resp.Fail(ctx, nil, "获取Api列表失败")
		return
	}

	resp.Success(ctx, respApiList, "获取Api列表成功")
}

// GetApiByID
// @Tags     Api相关
// @Summary  通过id来获取Api信息
// @Produce  application/json
// @Param    id  path int true "ApiID"
// @Success  200   {object}  resp.Result{data=models.RespGetApiByID}  "通过id来获取Api信息"
// @Security ApiKeyAuth
// @Router   /system/api/{id} [get]
func GetApiByID(ctx *gin.Context) {
	//	获取ApiId
	id := ctx.Param("id")
	// 获取db
	db := global.DB.Model(&models.SysApi{})
	var api models.RespGetApiByID
	// 判断Api是否存在
	err := db.Scopes(models.HandleCommonUserPreload).
		Where("id = ?", id).
		Omit("deleted_at").First(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "Api不存在")
			return
		}
		global.LOG.Errorln("获取Api信息失败, err:", err)
		resp.Fail(ctx, nil, "获取Api信息失败")
		return
	}

	resp.Success(ctx, api, "获取Api信息成功")
}

// EditApiByID
// @Tags     Api相关
// @Summary  通过id来修改Api信息
// @Produce  application/json
// @Param    id  path int true "ApiID"
// @Param    data  body      models.ReqEditApiByID true "方法 路径 目录 备注"
// @Success  200   {object}  resp.Result{}  "通过id来修改Api信息"
// @Security ApiKeyAuth
// @Router   /system/api/{id} [put]
func EditApiByID(ctx *gin.Context) {
	//	获取接口信息
	var reqInfo models.ReqEditApiByID

	// 判断Api是否存在
	var api models.SysApi
	err := checkApiIsExist(ctx, &reqInfo, &api)
	if err != nil {
		return
	}

	// 访问路径 和 请求方式 唯一检验
	var apiList []models.SysApi
	global.DB.Model(&models.SysApi{}).
		Where("method = ? AND path = ?",
			reqInfo.Method, reqInfo.Path).
		Not("id = ?", api.ID).Find(&apiList)
	if len(apiList) > 0 {
		resp.Fail(ctx, nil, "请确认更新后的Api信息, 请求方式+请求路径联合唯一!!")
		return
	}

	if !strings.HasSuffix(reqInfo.Path, "/") {
		reqInfo.Path += "/"
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.SysApi{}).Where("id = ?", api.ID).
			Updates(map[string]interface{}{
				"name":      reqInfo.Name,
				"method":    strings.ToUpper(reqInfo.Method),
				"path":      reqInfo.Path,
				"category":  reqInfo.Category,
				"UpdatedID": ctx.GetUint("userID"),
			}).Error; err != nil {

			return err
		}

		if err = casbin.UpdateCasbinApi(tx, api.Path, reqInfo.Path, api.Method, reqInfo.Method); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		global.LOG.Errorln("更新Api失败, err:", err)
		resp.Fail(ctx, nil, "更新Api失败")
		return
	}

	err = casbin.FreshCasbin()
	if err != nil {
		global.LOG.Errorln("重新加载角色Api失败, err:", err)
		resp.Fail(ctx, nil, "重新加载角色Api失败")
		return
	}

	resp.Success(ctx, nil, "更新Api信息成功")
}

// DeleteApiByID
// @Tags     Api相关
// @Summary  通过id来删除Api
// @Produce  application/json
// @Param    id  path int true "ApiID"
// @Success  200   {object}  resp.Result{}  "通过id来删除Api"
// @Security ApiKeyAuth
// @Router   /system/api/{id} [delete]
func DeleteApiByID(ctx *gin.Context) {
	// 判断Api是否存在
	var api models.SysApi
	err := checkApiIsExist(ctx, nil, &api)
	if err != nil {
		return
	}

	err = global.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Unscoped().Delete(&api).Error
		if err != nil {
			return err
		}

		err = casbin.DeleteCasbinApi(tx, api.Path, api.Method)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		global.LOG.Errorln("Api删除失败, err:", err)
		resp.Fail(ctx, nil, "Api删除失败")
		return
	}

	err = casbin.FreshCasbin()
	if err != nil {
		global.LOG.Errorln("重新加载角色Api失败, err:", err)
		resp.Fail(ctx, nil, "重新加载角色Api失败")
		return
	}

	resp.Success(ctx, nil, "删除Api信息成功")
}

func checkApiIsExist(ctx *gin.Context, reqInfo any, api *models.SysApi) error {
	//获取接口参数
	id := ctx.Param("id")

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

	err := db.Model(&models.SysApi{}).Where("id = ?", id).First(api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Fail(ctx, nil, "Api不存在")
			return err
		}

		global.LOG.Errorln("判断Api是否存在报错, err:", err)
		resp.Fail(ctx, nil, "判断Api是否存在报错")
		return err
	}
	return nil
}
