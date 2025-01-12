package user

import (
	"context"
	error2 "ht-crm/src/ht/common/error"
	"ht-crm/src/ht/common/result"
	"ht-crm/src/ht/common/utils"
	"ht-crm/src/ht/config/db"
	"ht-crm/src/ht/dao/user"
	"ht-crm/src/ht/dto/req"
	"ht-crm/src/ht/model"
	"ht-crm/src/ht/service"
)

func WechatUserLogin(ctx context.Context, code string) (string, error) {
	session, _ := service.Code2Session(code)

	userWechatInfo := user.QueryUserWechatInfo(session.OpenId)
	if userWechatInfo == nil {
		return "", error2.NewHtErrorFromTemplate(ctx, result.NoReg)
	}

	return utils.JwtGen(userWechatInfo.UserId), nil
}

func WechatUserReg(ctx context.Context, req req.WechatUserRegReq) error {
	session, _ := service.Code2Session(req.Code)
	userWechatInfo := user.QueryUserWechatInfo(session.OpenId)
	if userWechatInfo != nil {
		return error2.NewHtErrorFromTemplate(ctx, result.RepeatReg)
	}
	baseInfo := model.BaseInfo{AvatarUrl: req.AvatarUrl, Nickname: req.Nickname}
	db.Mysql.Create(&baseInfo)

	wechatInfo := model.WechatInfo{UserId: baseInfo.ID, AvatarUrl: baseInfo.AvatarUrl, NickName: baseInfo.Nickname, OpenId: session.OpenId}
	db.Mysql.Create(&wechatInfo)
	return nil
}
