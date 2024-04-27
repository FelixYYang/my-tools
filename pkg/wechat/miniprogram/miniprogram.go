package miniprogram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	appId     string
	appSecret string
	client    = &http.Client{}
)

func Init(id, secret string) error {
	appId = id
	appSecret = secret
	return nil
}

type Code2SessionS struct {
	OpenId string `json:"openid"`
}

// Code2Session 通过code交换用户登录信息
func Code2Session(code string) (res Code2SessionS, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, appSecret, code)
	err = httpGet(url, &res, false)
	return
}

type QRCodeForm struct {
	Page       string `json:"page"`
	Scene      string `json:"scene"`
	CheckPath  bool   `json:"check_path"`
	EnvVersion string `json:"env_version"`
}

// GetUnlimitedQRCode 获取无限小程序码
func GetUnlimitedQRCode(form QRCodeForm, token string) (res []byte, err error) {
	url := "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=" + token
	err = httpPost(url, form, &res, true)
	return
}

type AccessToken struct {
	ErrCode     string `json:"errcode"`
	AccessToken string `json:"access_token"`
	ErrMsg      string `json:"errmsg"`
}

// GetAccessToken 获取小程序全局唯一后台接口调用凭据
func GetAccessToken() (token AccessToken, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appId, appSecret)
	err = httpGet(url, &token, false)
	return
}

func httpGet(url string, res any, raw bool) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if raw {
		res = dataBytes
		return nil
	}
	return json.Unmarshal(dataBytes, res)
}

func httpPost(url string, form any, res any, raw bool) error {
	formData, err := json.Marshal(form)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(formData))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json;charset='utf-8'")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if raw {
		res = dataBytes
		return nil
	}
	return json.Unmarshal(dataBytes, res)
}
