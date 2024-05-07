package casbin

import (
	"dcss/global"
	"dcss/models"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
	"strconv"
)

func UpdateCasbins(AuthorityID uint, casbinInfo []models.CasbinInfo) error {
	authorityId := strconv.Itoa(int(AuthorityID))
	ClearCasbin(0, authorityId)
	rules := [][]string{}

	//做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfo {
		key := authorityId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
	}
	if len(casbinInfo) > 0 {
		success, _ := global.SyncedCachedEnforcer.AddPolicies(rules)
		if !success {
			return errors.New("存在相同api,添加失败,请联系管理员")
		}
	}
	return nil
}

func UpdateCasbinApi(tx *gorm.DB, oldPath string, newPath string, oldMethod string, newMethod string) error {
	return tx.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
}

func DeleteCasbinApi(tx *gorm.DB, Path, Method string) error {
	return tx.Unscoped().
		Where("v1 = ? AND v2 = ?", Path, Method).
		Delete(&gormadapter.CasbinRule{}).Error
}

func GetPolicyPathByAuthorityId(AuthorityID string) (pathMaps []models.CasbinInfo) {
	list := global.SyncedCachedEnforcer.GetFilteredPolicy(0, AuthorityID)
	for _, v := range list {
		pathMaps = append(pathMaps, models.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// ClearCasbin 清除匹配的权限
func ClearCasbin(v int, p ...string) bool {
	success, err := global.SyncedCachedEnforcer.RemoveFilteredPolicy(v, p...)
	if err != nil {
		global.LOG.Errorln("清除casbin失败, err: ", err)
	}
	return success
}

// RemoveFilteredPolicy  使用数据库方法清理筛选的politicy 此方法需要调用FreshCasbin方法才可以在系统中即刻生效
func RemoveFilteredPolicy(tx *gorm.DB, AuthorityID uint) error {
	authorityID := strconv.Itoa(int(AuthorityID))
	return tx.Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityID).Error

}

func FreshCasbin() (err error) {
	return global.SyncedCachedEnforcer.LoadPolicy()
}

func Init() {
	a, err := gormadapter.NewAdapterByDB(global.DB)
	if err != nil {
		global.LOG.Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", err)
		return
	}
	text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	m, err := model.NewModelFromString(text)
	if err != nil {
		global.LOG.Errorln("字符串加载模型失败!", err)
		return
	}
	global.SyncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
	global.SyncedCachedEnforcer.SetExpireTime(60 * 60)
	_ = global.SyncedCachedEnforcer.LoadPolicy()
}
