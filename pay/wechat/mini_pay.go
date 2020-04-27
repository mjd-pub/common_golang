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

// miniPayRespones 小程序支付请求返回参数
type MiniPayRespones struct {
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

// miniPayQueryRequests 小程序查询订单请求参数
type MiniPayQueryRequests struct {
	Appid         string `json:"appid" xml:"appid" structs:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	TransactionId string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
}

// miniPayQueryRespones 小程序查询订单请求返回参数
type MiniPayQueryRespones struct {
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

// miniPayCloseRequests 小程序关闭订单请求参数
type MiniPayCloseRequests struct {
	Appid      string `json:"appid" xml:"appid" structs:"appid"`
	MchId      string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	OutTradeNo string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	NonceStr   string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType   string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
}

// miniPayCloseRespones 小程序关闭订单请求返回参数
type MiniPayCloseRespones struct {
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

// miniPayRefundRequests 小程序申请退款请求参数
type MiniPayRefundRequests struct {
	Appid         string `json:"appid" xml:"appid" structs:"appid"`
	MchId         string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	TransactionId string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no" structs:"out_refund_no"`
	TotalFee      int    `json:"total_fee" xml:"total_fee" structs:"total_fee"`
	RefundFee     int    `json:"refund_fee" xml:"refund_fee" structs:"refund_fee"`
	RefundFeeType string `json:"refund_fee_type" xml:"refund_fee_type" structs:"refund_fee_type"`
	RefundDesc    string `json:"refund_desc" xml:"refund_desc" structs:"refund_desc"`
	RefundAccount string `json:"refund_account" xml:"refund_account" structs:"refund_account"`
	NotifyUrl     string `json:"notify_url" xml:"notify_url" structs:"notify_url"`
}

// miniPayRefundRespones 小程序申请退款请求返回参数
type MiniPayRefundRespones struct {
	ReturnCode          string `json:"return_code,omitempty" xml:"return_code,omitempty" structs:"return_code"`
	ReturnMsg           string `json:"return_msg,omitempty" xml:"return_msg,omitempty" structs:"return_msg"`
	ResultCode          string `json:"result_code,omitempty" xml:"result_code,omitempty" structs:"result_code"`
	ErrCode             string `json:"err_code,omitempty" xml:"err_code,omitempty" structs:"err_code"`
	ErrCodeDes          string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty" structs:"err_code_des"`
	Appid               string `json:"appid,omitempty" xml:"appid,omitempty" structs:"appid"`
	MchId               string `json:"mch_id,omitempty" xml:"mch_id,omitempty" structs:"mch_id"`
	NonceStr            string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty" structs:"nonce_str"`
	Sign                string `json:"sign,omitempty" xml:"sign,omitempty" structs:"sign"`
	TransactionId       string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty" structs:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty" structs:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no,omitempty" json:"out_refund_no,omitempty" structs:"out_refund_no"`
	RefundId            string `json:"refund_id,omitempty" xml:"refund_id,omitempty" structs:"refund_id"`
	RefundFee           int    `json:"refund_fee,omitempty" xml:"refund_fee,omitempty" structs:"refund_fee"`
	SettlementTotalFree int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty" structs:"settlement_total_free"`
	FreeType            string `json:"free_type,omitempty" xml:"free_type,omitempty" structs:"free_type"`
	CashFee             int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty" structs:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type,omitempty" json:"cash_fee_type,omitempty" structs:"cash_fee_type"`
	CashRefundFee       int    `json:"cash_refund_fee,omitempty" xml:"cash_refund_fee,omitempty" structs:"cash_refund_fee"`
	CouponType0         string `json:"coupon_type_0,omitempty" xml:"coupon_type_0,omitempty" structs:"coupon_type_0"`
	CouponRefundFee     int    `json:"coupon_refund_fee,omitempty" xml:"coupon_refund_fee,omitempty" structs:"coupon_refund_fee"`
	CouponRefundFee0    int    `json:"coupon_refund_fee_0,omitempty" xml:"coupon_refund_fee_0,omitempty" structs:"coupon_refund_fee_0"`
	ConponRefundCount   int    `json:"conpon_refund_count,omitempty" xml:"conpon_refund_count,omitempty" structs:"conpon_refund_count"`
	ConponRefundId0     string `json:"conpon_refund_id_0,omitempty" xml:"conpon_refund_id_0,omitempty" structs:"conpon_refund_id_0"`
}

type MiniPayRefundQueryRequests struct {
	AppId         string `json:"app_id" xml:"app_id" structs:"app_id"`
	MchId         string `json:"mch_id" xml:"mch_id" structs:"mch_id"`
	NonceStr      string `json:"nonce_str" xml:"nonce_str" structs:"nonce_str"`
	SignType      string `json:"sign_type" xml:"sign_type" structs:"sign_type"`
	TransactionId string `json:"transaction_id" xml:"transaction_id" structs:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no" xml:"out_trade_no" structs:"out_trade_no"`
	OutRefundNo   string `json:"out_refund_no" xml:"out_refund_no" structs:"out_refund_no"`
	RefundId      string `json:"refund_id" xml:"refund_id" structs:"refund_id"`
	Offset        int    `json:"offset" xml:"offset" structs:"offset"`
}

type MiniPayRefundQueryRespones struct {
	ReturnCode           string `json:"return_code,omitempty" xml:"return_code,omitempty" structs:"return_code"`
	ReturnMsg            string `json:"return_msg,omitempty" xml:"return_msg,omitempty" structs:"return_msg"`
	ResultCode           string `json:"result_code,omitempty" xml:"result_code,omitempty" structs:"result_code"`
	ErrCode              string `json:"err_code,omitempty" xml:"err_code,omitempty" structs:"err_code"`
	ErrCodeDes           string `json:"err_code_des,omitempty" xml:"err_code_des,omitempty" structs:"err_code_des"`
	Appid                string `json:"appid,omitempty" xml:"appid,omitempty" structs:"appid"`
	MchId                string `json:"mch_id,omitempty" xml:"mch_id,omitempty" structs:"mch_id"`
	NonceStr             string `json:"nonce_str,omitempty" xml:"nonce_str,omitempty" structs:"nonce_str"`
	Sign                 string `json:"sign,omitempty" xml:"sign,omitempty" structs:"sign"`
	TotalRefundCount     int    `json:"total_refund_count,omitempty" xml:"total_refund_count,omitempty" structs:"total_refund_count"`
	TransactionId        string `xml:"transaction_id,omitempty" json:"transaction_id,omitempty" structs:"transaction_id"`
	OutTradeNo           string `xml:"out_trade_no,omitempty" json:"out_trade_no,omitempty" structs:"out_trade_no"`
	TotalFee             int    `json:"total_fee,omitempty" xml:"total_fee,omitempty" structs:"total_fee"`
	SettlementTotalFree  int    `json:"settlement_total_free,omitempty" xml:"settlement_total_free,omitempty" structs:"settlement_total_free"`
	FreeType             string `json:"free_type,omitempty" xml:"free_type,omitempty" structs:"free_type"`
	CashFee              int    `xml:"cash_fee,omitempty" json:"cash_fee,omitempty" structs:"cash_fee"`
	RefundCount          int    `json:"refund_count,omitempty" xml:"refund_count,omitempty" structs:"refund_count"`
	OutRefundNo0         string `xml:"out_refund_no_0,omitempty" json:"out_refund_no_0,omitempty" structs:"out_refund_no_0"`
	RefundId0            string `json:"refund_id_0,omitempty" xml:"refund_id_0,omitempty" structs:"refund_id_0"`
	RefundFee0           int    `json:"refund_fee_0,omitempty" xml:"refund_fee_0,omitempty" structs:"refund_fee_0"`
	SettleMentRefundFee0 int    `json:"settle_ment_refund_fee_0" xml:"settle_ment_refund_fee_0" structs:"settle_ment_refund_fee_0"`
	CouponType00         string `json:"coupon_type_0_0" xml:"coupon_type_0_0" structs:"coupon_type_0_0"`
	ConponRefundFee0     int    `json:"conpon_refund_fee_0" xml:"conpon_refund_fee_0" structs:"conpon_refund_fee_0"`
	ConponRefundCount0   int    `json:"conpon_refund_count_0" xml:"conpon_refund_count_0" structs:"conpon_refund_count_0"`
	ConponRefundId00     string `json:"conpon_refund_id_0_0" xml:"conpon_refund_id_0_0" structs:"conpon_refund_id_0_0"`
	ConponRefundFee00     string `json:"conpon_refund_fee_0_0" xml:"conpon_refund_fee_0_0" structs:"conpon_refund_fee_0_0"`
	RefundStatus0        string `json:"refund_status_0" xml:"refund_status_0" structs:"refund_status_0"`
	RefundAccount0       string `json:"refund_account_0" xml:"refund_account_0" structs:"refund_account_0"`
	RefundRecvAccount0   string  `json:"refund_recv_account_0" xml:"refund_recv_account_0" structs:"refund_recv_account_0"`
	RefundSuccessTime0    string  `json:"refund_success_time_0" xml:"refund_success_time_0" structs:"refund_success_time_0"`
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
 * @return MiniPayRespones error
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

/**
 * RefundQuery 小程序退款查询
 *
 * @params request MiniPayRefundQueryRequests
 * @return MiniPayRefundQueryRespones err
 */
func (miniPay *MiniPay) RefundQuery(request MiniPayRefundQueryRequests) (queryResponse *MiniPayRefundQueryRespones, err error) {
	// 向微信发送请求
	resp, err := miniPay.wechatPay.Request(REFEUN_QUERY, request)
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