/**
* @Author: zhangjian@mioji.com
* @Date: 2019/8/21 上午11:43
 */
package utils

import (
	"bytes"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
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
