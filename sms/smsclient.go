package sms

import (
	"encoding/hex"
	"crypto/md5"
	"time"
	"fmt"
	"encoding/base64"
	"net/url"
	"net/http"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"errors"
)

const (
	API_VERSION = "2013-12-26"
)

type SMSClient interface {
	/**
	获取短信模板
	 */
	GetSMSTemplates() (string, error)
	/**
	获取子账户
	 */
	GetSubAccount() (string, error)
	/**
	发送短信
	 */
	SendSMS(templateId string, to string, data ...string) (error)
}

type smsClientImpl struct {
	Host	string
	Account string
	Token   string
	AppId   string
}

func New(host string, account string, token string, appId string) (SMSClient) {
	return &smsClientImpl{Host:host, Account:account, Token:token, AppId:appId}
}

func (client*smsClientImpl) GetSigParamater() (string, string) {
	date := time.Now()
	sig := getMd5String([]byte(fmt.Sprintf("%s%s%s", client.Account, client.Token, date.Format("20060102150405"))))
	auth := getBase64String([]byte(fmt.Sprintf("%s:%s", client.Account, date.Format("20060102150405"))))
	return sig, auth
}

func (client*smsClientImpl) GetSMSTemplates() (string, error) {
	sig, auth := client.GetSigParamater();

	values := url.Values{}
	values.Add("sig", sig)

	u := &url.URL{Scheme:"https", Host:client.Host,
		Path:fmt.Sprintf("/%s/Accounts/%s/SMS/QuerySMSTemplate",	API_VERSION, client.Account),
		RawQuery:values.Encode()}
	request := GetSMSTemplateRequest{AppId:client.AppId}
	body, error := json.Marshal(request)
	if error != nil {
		return "", error
	}
	req, error := http.NewRequest("POST", u.String(), bytes.NewReader(body))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json;charset=utf-8;")
	req.Header.Add("Accept", "application/json;")
	if error != nil {
		return "", error
	}
	httpClient := http.Client{}
	res, error := httpClient.Do(req)
	if error != nil {
		return "", error
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	data, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return "", error
	}
	return string(data), error
}

type GetSMSTemplateRequest struct {
	AppId string `json:"appId"`
}

func (client*smsClientImpl) GetSubAccount() (string, error) {
	sig, auth := client.GetSigParamater();

	values := url.Values{}
	values.Add("sig", sig)

	u := &url.URL{Scheme:"https", Host:client.Host,
		Path:fmt.Sprintf("/%s/Accounts/%s/GetSubAccounts",	API_VERSION, client.Account),
		RawQuery:values.Encode()}
	request := GetSubAccountRequest{AppId:client.AppId}
	body, error := json.Marshal(request)
	if error != nil {
		return "", error
	}
	req, error := http.NewRequest("POST", u.String(), bytes.NewReader(body))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json;charset=utf-8;")
	req.Header.Add("Accept", "application/json;")
	if error != nil {
		return "", error
	}
	httpClient := http.Client{}
	res, error := httpClient.Do(req)
	if error != nil {
		return "", error
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	data, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return "", error
	}
	return string(data), error
}

type GetSubAccountRequest struct {
	AppId string `json:"appId"`
}

func (client*smsClientImpl) SendSMS(templateId string, to string, args ...string) (error){
	sig, auth := client.GetSigParamater();

	values := url.Values{}
	values.Add("sig", sig)

	u := &url.URL{Scheme:"https", Host:client.Host,
		Path:fmt.Sprintf("/%s/Accounts/%s/SMS/TemplateSMS",	API_VERSION, client.Account),
		RawQuery:values.Encode()}
	request := SendSMSRequest{AppId:client.AppId, TemplateId:templateId, To:to, Datas:args}
	body, error := json.Marshal(request)
	if error != nil {
		return error
	}
	req, error := http.NewRequest("POST", u.String(), bytes.NewReader(body))
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json;charset=utf-8;")
	req.Header.Add("Accept", "application/json;")
	if error != nil {
		return error
	}
	httpClient := http.Client{}
	res, error := httpClient.Do(req)
	if error != nil {
		return error
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	data, error := ioutil.ReadAll(res.Body)
	if error != nil {
		return error
	}
	response := &SendSMSResponse{}
	error = json.Unmarshal(data, response)
	if error != nil {
		return error
	}
	if response.StatusCode != "000000" {
		return errors.New(response.StatusMsg)
	}
	return nil
}

type SendSMSRequest struct {
	AppId string `json:"appId"`
	To string `json:"to"`
	TemplateId string `json:"templateId"`
	Datas []string `json:"datas"`
}

type SendSMSResponse struct {
	StatusCode string `json:"statusCode"`
	StatusMsg string `json:"statusMsg"`
}

func getMd5String(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func getBase64String(data []byte) string {
	h := base64.StdEncoding
	return h.EncodeToString(data)
}
