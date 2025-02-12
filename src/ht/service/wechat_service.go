package service

import (
	"context"
	"encoding/json"
	"ht-crm/logger"
	"ht-crm/src/ht/config/db"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	appid     = "wx963823904279a30e"
	appSecret = "e55a06f73e1fa8f41d0ff806f56ca1f3"
)

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type Code2SessionResp struct {
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrMsg     string `json:"errmsg"`
	OpenId     string `json:"openid"`
	ErrCode    int    `json:"errcode"`
}

func getAccToken() string {
	key := "crm:wechat:accessToken"
	accessToken, err := db.Redis.Get(context.Background(), key).Result()
	if err != nil {
		logger.Error(err.Error())
	}
	if accessToken == "" {
		getTokenUrl := "https://api.weixin.qq.com/cgi-bin/token"
		param := url.Values{}
		param.Set("appid", appid)
		param.Set("secret", appSecret)
		param.Set("grant_type", "client_credential")

		req, _ := http.NewRequest(http.MethodGet, getTokenUrl+"?"+param.Encode(), nil)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Info("获取微信AppToken失败", err.Error())
		}
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return ""
		}
		accessToken := &accessTokenResp{}
		err = json.Unmarshal(all, accessToken)
		if err != nil {
			logger.Error(err.Error())
		}
		db.Redis.Set(context.Background(), key, accessToken.AccessToken, time.Duration(accessToken.ExpiresIn-60)*time.Second)
	}

	return accessToken
}

func Code2Session(code string) (*Code2SessionResp, error) {
	reqUrl := "https://api.weixin.qq.com/sns/jscode2session"
	param := url.Values{}
	param.Set("appid", appid)
	param.Set("secret", appSecret)
	param.Set("js_code", code)
	param.Set("grant_type", "authorization_code")
	req, _ := http.NewRequest(http.MethodGet, reqUrl+"?"+param.Encode(), nil)
	client := &http.Client{}
	do, err := client.Do(req)
	if err != nil {
		logger.Error("微信小程序登录错误", err.Error())
		return nil, err
	}
	all, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	resp := Code2SessionResp{}
	json.Unmarshal(all, &resp)
	return &resp, nil
}
