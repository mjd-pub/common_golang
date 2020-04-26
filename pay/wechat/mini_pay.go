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

// miniPay 小程序支付
type MiniPay struct {
	wechatPay *wechatPay
}

// miniPayRequests 小程序支付请求参数
type MiniPayRequest struct {
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

// miniPayRespones 小程序支付请求返回参数
type MiniPayRespones struct {
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

// miniPayQueryRequests 小程序查询订单请求参数
type MiniPayQueryRequests struct {
	Appid         string `json:"appid" xml:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id"`
	TransactionId string `json:"transaction_id" xml:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type"`
}

// miniPayQueryRespones 小程序查询订单请求返回参数
type MiniPayQueryRespones struct {
	ReturnCode          string `json:"return_code,omitempty" xml:"return_code,omitempty"`
	ReturnMsg           string `json:"return_msg,omitempty" xml:"return_msg,omitempty"`
	Appid               string `json:"appid,omitempty" xml:"appid,omitempty"`
	MchId               string `json:"mch_id,omitempty" xml:"mch_id,omitempty"`
	NonceStr            string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty"`
	Sign                string `json:"sign,omitempty" xml:"sign,omitempty"`
	ResultCode          string `json:"result_code,omitempty" xml:"result_code,omitempty"`
	ErrCode             string `json:"err_code,omitempty" xml:"err_code,omitempty"`
	ErrCodeDes          string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty"`
	DeviceInfo          string `json:"device_info,omitempty" xml:"device_info,omitempty"`
	OppenId             string `json:"oppen_id,omitempty" xml:"oppen_id,omitempty"`
	IsSubscribe         string `json:"is_subscribe,omitempty" xml:"is_subscribe,omitempty"`
	TradeType           string `json:"trade_type,omitempty" xml:"trade_type,omitempty"`
	BankType            string `json:"bank_type,omitempty" xml:"bank_type,omitempty"`
	TotalFree           int    `json:"total_free,omitempty" xml:"total_free,omitempty"`
	SettlementTotalFree int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty"`
	FreeType            string `json:"free_type,omitempty" xml:"free_type,omitempty"`
	CashFee             int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty"`
	CashFeeType         string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty"`
	CouponFee           int    `xml:"coupon_fee,omitempty" json:"coupon_fee,omitempty"`
	CouponCount         int    `xml:"coupon_count,omitempty" json:"coupon_count,omitempty"`
	CouponType0         string `xml:"coupon_type_0,omitempty" json:"coupon_type_0,omitempty"`
	CouponType1         string `xml:"coupon_type_1,omitempty" json:"coupon_type_1,omitempty"`
	CouponId0           string `xml:"coupon_id_0,omitempty" json:"coupon_id_0,omitempty"`
	CouponId1           string `xml:"coupon_id_1,omitempty" json:"coupon_id_1,omitempty"`
	CouponFee0          int    `xml:"coupon_fee_0,omitempty" json:"coupon_fee_0,omitempty"`
	CouponFee1          int    `xml:"coupon_fee_1,omitempty" json:"coupon_fee_1,omitempty"`
	TransactionId       string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	OutTradeNo          string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
	Attach              string `xml:"attach,omitempty" json:"attach,omitempty"`
	TimeEnd             string `xml:"time_end,omitempty" json:"time_end,omitempty"`
	Trade               string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
}

// miniPayCloseRequests 小程序关闭订单请求参数
type MiniPayCloseRequests struct {
	Appid      string `json:"appid" xml:"appid"`
	MchId      string `json:"mch_id" xml:"mch_id"`
	OutTradeNo string `json:"out_trade_no" xml:"out_trade_no"`
	NonceStr   string `json:"nonce_str" xml:"nonce_str"`
	SignType   string `json:"sign_type" xml:"sign_type"`
}

// miniPayCloseRespones 小程序关闭订单请求返回参数
type MiniPayCloseRespones struct {
	ReturnCode string `json:"return_code,omitempty" xml:"return_code,omitempty"`
	ReturnMsg  string `json:"return_msg,omitempty" xml:"return_msg,omitempty"`
	Appid      string `json:"appid,omitempty" xml:"appid,omitempty"`
	MchId      string `json:"mch_id,omitempty" xml:"mch_id,omitempty"`
	NonceStr   string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty"`
	Sign       string `json:"sign,omitempty" xml:"sign,omitempty"`
	ResultCode string `json:"result_code,omitempty" xml:"result_code,omitempty"`
	ResultMsg  string `json:"result_msg,omitempty" xml:"result_msg,omitempty"`
	ErrCode    string `json:"err_code,omitempty" xml:"err_code,omitempty"`
	ErrCodeDes string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty"`
}

// miniPayRefundRequests 小程序申请退款请求参数
type MiniPayRefundRequests struct {
	Appid         string `json:"appid" xml:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id"`
	TransactionId string `json:"transaction_id" xml:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no"`
	TotalFee      int    `json:"total_fee" xml:"total_fee"`
	RefundFee     int    `json:"refund_fee" xml:"refund_fee"`
	RefundFeeType string `json:"refund_fee_type" xml:"refund_fee_type"`
	RefundDesc    string `json:"refund_desc" xml:"refund_desc"`
	RefundAccount string `json:"refund_account" xml:"refund_account"`
	NotifyUrl     string `json:"notify_url" xml:"notify_url"`
}

// miniPayRefundRespones 小程序申请退款请求返回参数
type MiniPayRefundRespones struct {
	ReturnCode          string `json:"return_code,omitempty" xml:"return_code,omitempty"`
	ReturnMsg           string `json:"return_msg,omitempty" xml:"return_msg,omitempty"`
	ResultCode          string `json:"result_code,omitempty" xml:"result_code,omitempty"`
	ErrCode             string `json:"err_code,omitempty" xml:"err_code,omitempty"`
	ErrCodeDes          string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty"`
	Appid               string `json:"appid,omitempty" xml:"appid,omitempty"`
	MchId               string `json:"mch_id,omitempty" xml:"mch_id,omitempty"`
	NonceStr            string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty"`
	Sign                string `json:"sign,omitempty" xml:"sign,omitempty"`
	TransactionId       string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	OutTradeNo          string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty"`
	OutRefundNo         string `xml:"out_refund_no,omitempty" json:"out_refund_no,omitempty"`
	RefundId            string `json:"refund_id,omitempty" xml:"refund_id,omitempty"`
	RefundFee           int    `json:"refund_fee,omitempty" xml:"refund_fee,omitempty"`
	SettlementTotalFree int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty"`
	FreeType            string `json:"free_type,omitempty" xml:"free_type,omitempty"`
	CashFee             int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty"`
	CashFeeType         string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty"`
	CashRefundFee       int    `json:"cash_refund_fee,omitempty" xml:"cash_refund_fee,omitempty"`
	CouponType0         string `json:"coupon_type_0,omitempty" xml:"coupon_type_0,omitempty"`
	CouponRefundFee     int    `json:"coupon_refund_fee,omitempty" xml:"coupon_refund_fee,omitempty"`
	CouponRefundFee0    int    `json:"coupon_refund_fee_0,omitempty" xml:"coupon_refund_fee_0,omitempty"`
	CouponRefundFee1    int    `json:"coupon_refund_fee_1,omitempty" xml:"coupon_refund_fee_1,omitempty"`
	ConponRefundCount   int    `json:"conpon_refund_count,omitempty" xml:"conpon_refund_count,omitempty"`
	ConponRefundId0     string `json:"conpon_refund_id_0,omitempty" xml:"conpon_refund_id_0,omitempty"`
	ConponRefundId1     string `json:"conpon_refund_id_1,omitempty" xml:"conpon_refund_id_1,omitempty"`
}

func NewMiniPayClient(appid, mchid, key, apiclientKey, apiclientCert string) *MiniPay {
	wechatPay := newWechatPay(appid, mchid, key, apiclientKey, apiclientCert)
	return &MiniPay{
		wechatPay: wechatPay,
	}
}

/**
 * NewMiniPayRequest 构造下单请求
 *
 * @params body
 * @params detail
 * @params orderId 订单id
 * @params userIp 用户ip
 * @params notifyUrl 异步回调url
 * @params openid
 * @params price 价格单位(元)
 *
 * @return NewMiniPayRequest
 */
func (h5Pay *H5Pay) NewMiniPayRequest(body, detail, orderId, userIp, notifyUrl, openid string, price float64) MiniPayRequest {
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
	return MiniPayRequest{
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
 * @params request MiniPayRequest
 * @return miniResp err
 */
func (miniPay *MiniPay) Pay(request MiniPayRequest) (miniResp *MiniPayRespones, err error) {
	// 向微信发送请求
	resp, err := miniPay.wechatPay.Request(UNIFIED_ORDER, request)
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
 * @params request MiniPayQueryRequests
 * @return MiniPayQueryRespones err
 */
func (miniPay *MiniPay) Query(request MiniPayQueryRequests) (queryResponse *MiniPayQueryRespones, err error) {
	// 向微信发送请求
	resp, err := miniPay.wechatPay.Request(ORDER_QUERY, request)
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
 * @params request MiniPayCloseRequests
 * @return MiniPayCloseRespones err
 */
func (miniPay *MiniPay) Close(request MiniPayCloseRequests) (queryResponse *MiniPayCloseRespones, err error) {
	// 向微信发送请求
	resp, err := miniPay.wechatPay.Request(ORDER_QUERY, request)
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
 * Refund 小程序申请退款
 *
 * @params request MiniPayRefundRequests
 * @return MiniPayRefundRespones err
 */
func (miniPay *MiniPay) Refund(request MiniPayRefundRequests) (queryResponse *MiniPayRefundRespones, err error) {
	// 向微信发送请求
	resp, err := miniPay.wechatPay.Request(REFUND, request)
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
