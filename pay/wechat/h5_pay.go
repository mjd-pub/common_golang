package wechat

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/mjd-pub/common_golang/utils"
	"io/ioutil"
	"strconv"
	"time"
)

// H5Pay h5支付
type H5Pay struct {
	wechatPay *wechatPay
}

// H5PayRequests h5支付请求参数
type H5PayRequest struct {
	Appid          string `json:"appid" xml:"appid"`
	MchId          string `json:"mch_id" xml:"mch_id"`
	NonceStr       string `json:"nonce_str" xml:"nonce_str"`
	SignType       string `json:"sign_type" xml:"sign_type"`
	Body           string `json:"body" xml:"body"`
	Detail         string `json:"detail" xml:"detail"`
	OutTradeNo     string `json:"out_trade_no" xml:"out_trade_no"`
	FeeType        string `json:"fee_type" xml:"fee_type"`
	TotalFee       int    `json:"total_fee" xml:"total_fee"`
	SpbillCreateIp string `json:"spbill_create_ip" xml:"spbill_create_ip"`
	TimeStart      string `json:"time_start" xml:"time_start"`
	TimeExpire     int64  `json:"time_expire" xml:"time_expire"`
	NotifyUrl      string `json:"notify_url" xml:"notify_url"`
	TradeType      string `json:"trade_type" xml:"trade_type"`
	Openid         string `json:"openid" xml:"openid"`
	SceneInfo      string `json:"scene_info" xml:"scene_info"`
}

// H5PayRespones h5支付请求返回参数
type H5PayRespones struct {
	ReturnCode string `json:"return_code" xml:"return_code"`
	ReturnMsg  string `json:"return_msg" xml:"return_msg"`
	Appid      string `json:"appid" xml:"appid"`
	MchId      string `json:"mch_id" xml:"mch_id"`
	DeviceInfo string `json:"device_info" xml:"device_info"`
	NonceStr   string `json:"nonce_str" xml:"nonce_str"`
	Sign       string `json:"sign" xml:"sign"`
	ResultCode string `json:"result_code" xml:"result_code"`
	ErrCode    string `json:"err_code" xml:"err_code"`
	ErrCodeDes string `json:"err_code_des" xml:"err_code_des"`
	TradeType  string `json:"trade_type" xml:"trade_type"`
	PrepayId   string `json:"prepay_id" xml:"prepay_id"`
	MwebUrl    string `json:"mweb_url" xml:"mweb_url"`
}

func NewH5PayClient(appid, mchid, key, apiclientKey, apiclientCert string) *H5Pay {
	wechatPay := newWechatPay(appid, mchid, key, apiclientKey, apiclientCert)
	return &H5Pay{
		wechatPay: wechatPay,
	}
}

/**
 * NewH5PayRequest 构造下单请求
 *
 * @params body
 * @params detail
 * @params orderId 订单id
 * @params userIp 用户ip
 * @params notifyUrl 异步回掉url
 * @params openid
 * @params price 价格单位(元)
 *
 * @return CompanyPayRequest
 */
func (h5Pay *H5Pay) NewH5PayRequest(body, detail, orderId, userIp, notifyUrl, openid string, price float64) H5PayRequest {
	//3.scene_info
	scene_info := map[string]interface{}{
		"h5_info": map[string]interface{}{
			"type":     "Wap",
			"wap_url":  "",
			"wap_name": "MagicData",
		},
	}
	sceneInfo, _ := json.Marshal(scene_info)
	return H5PayRequest{
		Appid:          h5Pay.wechatPay.appid,
		MchId:          h5Pay.wechatPay.mchid,
		NonceStr:       utils.GetNonceStr(),
		SignType:       "MD5",
		Body:           body,
		Detail:         detail,
		OutTradeNo:     orderId,
		FeeType:        "CNY",
		TotalFee:       int(price * 100), // 单位为分
		SpbillCreateIp: userIp,
		TimeStart:      time.Now().Format("20060102150405"),
		TimeExpire:     time.Now().Unix() + 2*3600,
		NotifyUrl:      notifyUrl,
		TradeType:      "MWEB",
		Openid:         openid,
		SceneInfo:      string(sceneInfo),
	}
}

/**
 * Pay 发起支付
 *
 * @params request H5PayRequest
 * @return h5Resp err
 */
func (h5Pay *H5Pay) Pay(request H5PayRequest) (h5Resp *H5PayRespones, err error) {
	// 向微信发送请求
	resp, err := h5Pay.wechatPay.Request(COMPANY_PAY, request)
	if err != nil {
		return nil, errors.New("请求异常:" + err.Error())
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("httpCode Err:" + strconv.Itoa(resp.StatusCode))
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//xml解码
	err = xml.Unmarshal(respData, &h5Resp)
	if err != nil {
		return nil, err
	}
	return
}
