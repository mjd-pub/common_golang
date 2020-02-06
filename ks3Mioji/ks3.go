/**
* @Author: zhangjian@mioji.com
* @Date: 2019/7/26 上午10:47
 */
package ks3Mioji

import (
	"bytes"
	"github.com/ks3sdklib/aws-sdk-go/aws"
	"github.com/ks3sdklib/aws-sdk-go/aws/credentials"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"
	"os"
)

const (
	Ks3EndPoint     = "kss.ksyun.com"
	CdnAddr         = "cdn.magicdmc.com"
)

type Ks3 struct {
	ClientConn      *s3.S3
	AccessKeyID     string
	AccessKeySecret string
}

//初始化客户端
func (k *Ks3) InitClient() *s3.S3 {
	mCredentials := credentials.NewStaticCredentials(k.AccessKeyID, k.AccessKeySecret, "")
	client := s3.New(&aws.Config{
		Region:           "BEIJING",
		Credentials:      mCredentials,
		Endpoint:         Ks3EndPoint, //ks3地址
		DisableSSL:       false,       //是否禁用https
		LogLevel:         0,           //是否开启日志,0为关闭日志，1为开启日志
		S3ForcePathStyle: false,       //是否强制使用path style方式访问
		LogHTTPBody:      true,        //是否把HTTP请求body打入日志
		Logger:           os.Stdout,   //打日志的位置
	})
	k.ClientConn = client
	return client
}

//上传文件
func (k *Ks3) UploadFileByByte(bucket string, uploadPath string, resource []byte, contentType string) (string, error) {
	params := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),        // bucket名称
		Key:         aws.String(uploadPath),    // object key 文件夹+文件名
		ACL:         aws.String("public-read"), //权限，支持private(私有)，public-read(公开读)
		Body:        bytes.NewReader(resource), //要上传的内容
		ContentType: aws.String(contentType),   //设置content-type
		Metadata:    map[string]*string{
			//"Key": aws.String("MetadataValue"), // 设置用户元数据
			// More values...
		},
	}
	_, err := k.ClientConn.PutObject(params)
	if err != nil {
		return "", err
	}
	url := "https://" + bucket + "." + CdnAddr + "/" + uploadPath
	return url, nil
}
