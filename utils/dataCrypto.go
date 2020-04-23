/**
* @Author: zhangjian@mioji.com
* @Date: 2019/8/21 上午11:43
 */
package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	syncAesMutex sync.Mutex
	commonAeskey []byte
)

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func EntryptDesECB(data, key []byte) (string, error) {
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", errors.New("EntryptDesECB newCipher error" + err.Error())
	}
	bs := block.BlockSize()
	data = PKCS5Padding(data, bs)
	if len(data)%bs != 0 {
		return "", errors.New("EntryptDesECB Need a multiple of the blocksize" + err.Error())
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Encrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	return base64.StdEncoding.EncodeToString(out), nil
}

func DecodeNotifyData(reqInfo string, paykey string) ([]byte, error) {

	//1.对加密串A做base64解码，得到加密串B
	b, err := base64.StdEncoding.DecodeString(reqInfo)
	if err != nil {
		return nil, err
	}

	//2.对商户key做md5，得到32位小写key
	err = SetAesKey(strings.ToLower(Md5(paykey)))
	if err != nil {
		return nil, err
	}

	//3.用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
	plaintext, err := AesECBDecrypt(b)

	return plaintext, nil
}

func DecryptDESECB(d string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(d)
	if err != nil {
		return "", errors.New("DecryptDES Decode base64 error" + err.Error())
	}
	if len(key) > 8 {
		key = key[:8]
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", errors.New("DecryptDES NewCipher error" + err.Error())
	}
	bs := block.BlockSize()
	if len(data)%bs != 0 {
		return "", errors.New("DecryptDES crypto/cipher: input not full blocks")
	}
	out := make([]byte, len(data))
	dst := out
	for len(data) > 0 {
		block.Decrypt(dst, data[:bs])
		data = data[bs:]
		dst = dst[bs:]
	}
	out = PKCS5UnPadding(out)
	return string(out), nil
}

func DesEncrypt(src string, key []byte) (string, error) {
	out_ori, err := EntryptDesECB([]byte(src), key)
	out_base64_decode, _ := base64.StdEncoding.DecodeString(out_ori)
	return strings.ToUpper(hex.EncodeToString(out_base64_decode)), err
}

func DesDecrypt(src string, key []byte) (string, error) {
	out_hex, _ := hex.DecodeString(src)
	out_base64 := base64.StdEncoding.EncodeToString(out_hex)
	return DecryptDESECB(out_base64, key)
}

func SetAesKey(key string) (err error) {
	syncAesMutex.Lock()
	defer syncAesMutex.Unlock()
	b := []byte(key)
	if len(b) == 16 || len(b) == 24 || len(b) == 32 {
		commonAeskey = b
		return nil
	}
	return fmt.Errorf("key size is not 16 or 24 or 32, but %d", len(b))

}

func Md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func AesECBDecrypt(ciphertext []byte, paddingType ...string) (plaintext []byte, err error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	// ECB mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	block, err := aes.NewCipher(commonAeskey)
	if err != nil {
		return nil, err
	}
	NewECBDecrypter(block).CryptBlocks(ciphertext, ciphertext)
	if len(paddingType) > 0 {
		switch paddingType[0] {
		case "ZeroUnPadding":
			plaintext = ZeroUnPadding(ciphertext)
		case "PKCS5UnPadding":
			plaintext = PKCS5UnPadding(ciphertext)
		}
	} else {
		plaintext = PKCS5UnPadding(ciphertext)
	}
	return plaintext, nil
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

type ecbDecrypter ecb

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//将data组织成输出xml字符
func ToXml(data map[string]interface{}) string {

	strs := Ksort(data)
	xml_temp := "<xml>"
	for _, key := range strs {
		if data[key] != nil {
			switch reflect.TypeOf(data[key]).Name() {
			case "int":
				if data[key].(int) != 0 {
					xml_temp += "<" + key + ">" + strconv.Itoa(data[key].(int)) + "</" + key + ">"
				}
			case "string":
				xml_temp += "<" + key + ">" + data[key].(string) + "</" + key + ">"
			case "float64":
				s := strconv.FormatFloat(data[key].(float64), 'f', -1, 64)
				xml_temp += "<" + key + ">" + s + "</" + key + ">"
			}
		}
	}
	xml_temp += "</xml>"
	return xml_temp
}

//map按字典排序
func Ksort(data map[string]interface{}) []string {
	//保存key值
	var strs []string
	for key, _ := range data {
		strs = append(strs, key)
	}
	//key值按字典排序
	sort.Strings(strs)

	return strs
}

//格式化参数格式化成url参数
func ToUrlParams(data map[string]interface{}, strs []string) string {
	buff := ""
	for _, str := range strs {
		if str != "sign" && data[str] != "" {
			switch reflect.TypeOf(data[str]).Name() {
			case "int":
				if data[str].(int) != 0 {
					buff += str + "=" + strconv.Itoa(data[str].(int)) + "&"
				}
				if str == "coupon_refund_count" || str == "coupon_refund_fee" {
					buff += str + "=" + strconv.Itoa(data[str].(int)) + "&"
				}
			case "float64":
				temp := data[str].(float64)
				buff += str + "=" + strconv.FormatFloat(temp, 'f', -1, 64) + "&"
			case "string":
				buff += str + "=" + data[str].(string) + "&"
			}
		}
	}
	buff = buff[:len(buff)-1]
	return buff
}