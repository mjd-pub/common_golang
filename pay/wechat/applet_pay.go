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

// AppletPay 小程序支付
type AppletPay struct {
	wechatPay *wechatPay
}

// AppletPayRequests 小程序支付请求参数
type AppletPayRequest struct {
	Appid          string `json:"appid" xml:"appid" structs:"appid"`
	MchId          string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	NonceStr       string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType       string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
	Body           string `json:"body" xml:"body" structs:"body"`
	Detail         string `json:"detail" xml:"detail" structs:"detail"`
	OutTradeNo     string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	FeeType        string `json:"fee_type" xml:"fee_type" structs:"fee_type"`
	TotalFee       int    `json:"total_fee" xml:"total_fee" structs:"total_fee"`
	SpbillCreateIp string `json:"spbill_create_ip" xml:"spbill_create_ip" structs:"spbill_create_ip"`
	TimeStart      string `json:"time_start" xml:"time_start" structs:"time_start"`
	TimeExpire     int64  `json:"time_expire" xml:"time_expire" structs:"time_expire"`
	NotifyUrl      string `json:"notify_url" xml:"notify_url" structs:"notify_url"`
	TradeType      string `json:"trade_type" xml:"trade_type" structs:"trade_type"`
	Openid         string `json:"openid" xml:"openid" structs:"openid"`
	SceneInfo      string `json:"scene_info" xml:"scene_info" structs:"scene_info"`
}

// AppletPayRespones 小程序支付请求返回参数
type AppletPayRespones struct {
	ReturnCode string `json:"return_code" xml:"return_code" structs:"return_code"`
	ReturnMsg  string `json:"return_msg" xml:"return_msg" structs:"apreturn_msgpid"`
	Appid      string `json:"appid" xml:"appid" structs:"appid"`
	MchId      string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	DeviceInfo string `json:"device_info" xml:"device_info" structs:"device_info"`
	NonceStr   string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	Sign       string `json:"sign" xml:"sign" structs:"sign"`
	ResultCode string `json:"result_code" xml:"result_code" structs:"result_code"`
	ErrCode    string `json:"err_code" xml:"err_code" structs:"err_code"`
	ErrCodeDes string `json:"err_code_des" xml:"err_code_des" structs:"err_code_des"`
	TradeType  string `json:"trade_type" xml:"trade_type" structs:"trade_type"`
	PrepayId   string `json:"prepay_id" xml:"prepay_id" structs:"prepay_id"`
	MwebUrl    string `json:"mweb_url" xml:"mweb_url" structs:"mweb_url"`
}

// AppletPayQueryRequests 小程序查询订单请求参数
type AppletPayQueryRequests struct {
	Appid         string `json:"appid" xml:"appid" structs:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	TransactionId string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
}

// AppletPayQueryRespones 小程序查询订单请求返回参数
type AppletPayQueryRespones struct {
	ReturnCode          string `json:"return_code,omitempty" xml:"return_code,omitempty" structs:"return_code"`
	ReturnMsg           string `json:"return_msg,omitempty" xml:"return_msg,omitempty" structs:"return_msg"`
	Appid               string `json:"appid,omitempty" xml:"appid,omitempty" structs:"appid"`
	MchId               string `json:"mch_id,omitempty" xml:"mch_id,omitempty" structs:"mch_id"`
	NonceStr            string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty" structs:"nonce_str"`
	Sign                string `json:"sign,omitempty" xml:"sign,omitempty" structs:"sign"`
	ResultCode          string `json:"result_code,omitempty" xml:"result_code,omitempty" structs:"result_code"`
	ErrCode             string `json:"err_code,omitempty" xml:"err_code,omitempty" structs:"err_code"`
	ErrCodeDes          string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty" structs:"err_code_des"`
	DeviceInfo          string `json:"device_info,omitempty" xml:"device_info,omitempty" structs:"device_info"`
	OppenId             string `json:"oppen_id,omitempty" xml:"oppen_id,omitempty" structs:"oppen_id"`
	IsSubscribe         string `json:"is_subscribe,omitempty" xml:"is_subscribe,omitempty" structs:"is_subscribe"`
	TradeType           string `json:"trade_type,omitempty" xml:"trade_type,omitempty" structs:"trade_type"`
	BankType            string `json:"bank_type,omitempty" xml:"bank_type,omitempty" structs:"bank_type"`
	TotalFree           int    `json:"total_free,omitempty" xml:"total_free,omitempty" structs:"total_free"`
	SettlementTotalFree int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty" structs:"settlement_total_free"`
	FreeType            string `json:"free_type,omitempty" xml:"free_type,omitempty" structs:"free_type"`
	CashFee             int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty" structs:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty" structs:"cash_fee_type"`
	CouponFee           int    `xml:"coupon_fee,omitempty" json:"coupon_fee,omitempty" structs:"coupon_fee"`
	CouponCount         int    `xml:"coupon_count,omitempty" json:"coupon_count,omitempty" structs:"coupon_count"`
	CouponType0         string `xml:"coupon_type_0,omitempty" json:"coupon_type_0,omitempty" structs:"coupon_type_0"`
	CouponId0           string `xml:"coupon_id_0,omitempty" json:"coupon_id_0,omitempty" structs:"coupon_id_0"`
	CouponFee0          int    `xml:"coupon_fee_0,omitempty" json:"coupon_fee_0,omitempty" structs:"coupon_fee_0"`
	TransactionId       string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty" structs:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty" structs:"out_trade_no"`
	Attach              string `xml:"attach,omitempty" json:"attach,omitempty" structs:"attach"`
	TimeEnd             string `xml:"time_end,omitempty" json:"time_end,omitempty" structs:"time_end"`
	Trade               string `xml:"trade,omitempty" json:"trade,omitempty" structs:"trade"`
}

// AppletPayCloseRequests 小程序关闭订单请求参数
type AppletPayCloseRequests struct {
	Appid      string `json:"appid" xml:"appid" structs:"appid"`
	MchId      string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	OutTradeNo string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	NonceStr   string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType   string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
}

// AppletPayCloseRespones 小程序关闭订单请求返回参数
type AppletPayCloseRespones struct {
	ReturnCode string `json:"return_code,omitempty" xml:"return_code,omitempty" structs:"return_code"`
	ReturnMsg  string `json:"return_msg,omitempty" xml:"return_msg,omitempty" structs:"return_msg"`
	Appid      string `json:"appid,omitempty" xml:"appid,omitempty" structs:"appid"`
	MchId      string `json:"mch_id,omitempty" xml:"mch_id,omitempty" structs:"mch_id"`
	NonceStr   string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty" structs:"nonce_str"`
	Sign       string `json:"sign,omitempty" xml:"sign,omitempty" structs:"sign"`
	ResultCode string `json:"result_code,omitempty" xml:"result_code,omitempty" structs:"result_code"`
	ResultMsg  string `json:"result_msg,omitempty" xml:"result_msg,omitempty" structs:"result_msg"`
	ErrCode    string `json:"err_code,omitempty" xml:"err_code,omitempty" structs:"err_code"`
	ErrCodeDes string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty" structs:"err_code_des"`
}

func NewAppletPayClient(appid, mchid, key, apiclientKey, apiclientCert string) *AppletPay {
	wechatPay := NewWechatPay(appid, mchid, key, apiclientKey, apiclientCert)
	return &AppletPay{
		wechatPay: wechatPay,
	}
}

/**
 * NewAppletPayRequest 构造下单请求
 *
 * @params body
 * @params detail
 * @params orderId 订单id
 * @params userIp 用户ip
 * @params notifyUrl 异步回调url
 * @params openid
 * @params price 价格单位(元)
 *
 * @return NewAppletPayRequest
 */
func (h5Pay *H5Pay) NewAppletPayRequest(body, detail, orderId, userIp, notifyUrl, openid string, price float64) AppletPayRequest {
	//scene_info
	scene_info := map[string]interface{}{
		"store_info": map[string]interface{}{
			"id":        "门店id",
			"name":      "门店名称",
			"area_code": "门店行政区划码",
			"address":   "门店详细地址",
		},
	}
	sceneInfo, _ := json.Marshal(scene_info)
	return AppletPayRequest{
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
		TradeType:      "JSAPI",
		Openid:         openid,
		SceneInfo:      string(sceneInfo),
	}
}

/**
 * Pay 发起支付
 *
 * @params request AppletPayRequest
 * @return AppletPayRespones error
 */
func (AppletPay *AppletPay) Pay(request AppletPayRequest) (miniResp *AppletPayRespones, err error) {
	// 向微信发送请求
	resp, err := AppletPay.wechatPay.Request(UNIFIED_ORDER, request)
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
	err = xml.Unmarshal(respData, miniResp)
	if err != nil {
		return nil, err
	}
	return
}

/**
 * Query 小程序支付查询
 *
 * @params request AppletPayQueryRequests
 * @return AppletPayQueryRespones err
 */
func (AppletPay *AppletPay) Query(request AppletPayQueryRequests) (queryResponse *AppletPayQueryRespones, err error) {
	// 向微信发送请求
	resp, err := AppletPay.wechatPay.Request(ORDER_QUERY, request)
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
	err = xml.Unmarshal(respData, queryResponse)
	if err != nil {
		return nil, err
	}
	return
}

/**
 * Close 小程序支付关闭
 *
 * @params request AppletPayCloseRequests
 * @return AppletPayCloseRespones err
 */
func (AppletPay *AppletPay) Close(request AppletPayCloseRequests) (queryResponse *AppletPayCloseRespones, err error) {
	// 向微信发送请求
	resp, err := AppletPay.wechatPay.Request(ORDER_QUERY, request)
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
	err = xml.Unmarshal(respData, queryResponse)
	if err != nil {
		return nil, err
	}
	return
}
