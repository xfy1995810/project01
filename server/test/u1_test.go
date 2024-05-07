package test

import (
	"dcss/global"
	"dcss/global/db"
	"dcss/models"
	"fmt"
	"testing"
)

type GetAllName struct {
	AllName string `gorm:"column:allname"`
}

func TestM1(t *testing.T) {
	err := db.InitDB()
	if err != nil {
		global.LOG.Errorln(err)
		panic(err)
	}
	db := global.DB
	var gn GetAllName
	err = db.Model(&models.Muban{}).Where("allname", "mobanNone111").Scan(&gn).Error
	fmt.Println("||||", gn.AllName, err)
	err = db.Debug().Where("allname = ?", "mobanNone111").Unscoped().Delete(&models.Muban{}).Error
	fmt.Println("||||delete ", err)

}
