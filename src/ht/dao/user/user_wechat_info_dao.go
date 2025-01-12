package user

import (
	"ht-crm/src/ht/config/db"
	"ht-crm/src/ht/model"
)

func QueryUserWechatInfo(openId string) *model.WechatInfo {
	wechatInfo := model.WechatInfo{}
	tx := db.Mysql.Where("open_id=?", openId).Take(&wechatInfo)
	if tx.Error != nil {
		return nil
	}
	return &wechatInfo
}
