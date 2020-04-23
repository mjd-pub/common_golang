package wechat

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/mjd-pub/common_golang/utils"
	"io"
	"net/http"
	"strings"
)

const (
	UNIFIED_ORDER     = "https://api.mch.weixin.qq.com/pay/unifiedorder"                      // 统一下单接口地址
	ORDER_QUERY       = "https://api.mch.weixin.qq.com/pay/orderquery"                        // 查询接口地址
	REFUND            = "https://api.mch.weixin.qq.com/secapi/pay/refund"                     // 退款接口地址
	REFEUN_QUERY      = "https://api.mch.weixin.qq.com/pay/refundquery"                       // 退款查询接口地址
	COMPANY_PAY       = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers" // 企业支付下单
	COMPANY_PAY_QUERY = "https://api.mch.weixin.qq.com/mmpaymkttransfers/gettransferinfo"     // 企业支付查询
)

const (
	DEFAULT        = 0
	PAY_SUCCESS    = 1
	REFUND_PROCESS = 11
	REFUND_SUCCESS = 12
	REFUND_FAIL    = 13
)

//　WechatPay 微信支付基础结构体
type wechatPay struct {
	apiclientCert string
	apiclientKey  string
	appid         string
	key           string
	mchid         string
}

/**
 * newWechatPay 微信支付初始化
 * @params appid 商户号绑定的appid
 * @params mchid 商户号
 * @params key   支付密钥
 * @params apiclientKey
 * @params apiclientCert
 */
func newWechatPay(appid, mchid, key, apiclientKey, apiclientCert string) *wechatPay {
	return &wechatPay{
		apiclientCert: apiclientCert,
		apiclientKey:  apiclientKey,
		appid:         appid,
		key:           key,
		mchid:         mchid,
	}
}

/**
 * Request 微信标准请求输出
 * @params uri 请求uri
 * @params requestData 请求参数为固定结构体
 *
 * @return resp 标准http请求 err 标准错误输出
 */
func (wechat *wechatPay) Request(uri string, requestData interface{}) (resp *http.Response, err error) {
	//2.签名和请求参数
	respXml, err := wechat.dealXmlRequest(requestData)
	if err != nil {
		return
	}
	//3. 证书和公钥文件
	cert := wechat.apiclientCert  //证书
	sslKey := wechat.apiclientKey //公钥
	//4. 带证书发送http请求
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(cert))
	//5.tls.X509KeyPair 直接读字符串
	cliCrt, err := tls.X509KeyPair([]byte(cert), []byte(sslKey))
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//RootCAs:      pool,	不能添加这一项
			Certificates: []tls.Certificate{cliCrt},
			ClientAuth:   tls.RequireAndVerifyClientCert,
		},
	}
	client := &http.Client{Transport: tr}
	//请求微信接口
	resp, err = client.Post(uri, "text/xml; charset=UTF8", strings.NewReader(respXml))
	return
}

// DealXmlRequest 处理初xml请求body
func (wechat *wechatPay) dealXmlRequest(params interface{}) (xmlRequest string, err error) {
	paramsMap := utils.Struct2Map(params)
	// 对数据进行签名
	sign, err := wechat.signData(paramsMap)
	if err != nil {
		return
	}
	paramsMap["sign"] = sign
	//4。转化为xml格式
	xmlRequest = utils.ToXml(paramsMap)
	return
}

// SignData 对数据进行签名
func (wechat *wechatPay) signData(data map[string]interface{}) (sign string, err error) {
	strs := utils.Ksort(data)
	//1.2 使用URL键值对的形式生成字符串
	str := utils.ToUrlParams(data, strs)
	//1.3 在str后加入KEY
	str = str + "&key=" + wechat.key
	//2. 将得到的数据做MD5结算得到signValue
	m := md5.New()
	_, err = io.WriteString(m, str)
	if err != nil {
		return "", errors.New("签名错误:" + err.Error())
	}
	//3。将签名转化为大写
	sign = strings.ToUpper(fmt.Sprintf("%x", m.Sum(nil)))
	return
}
