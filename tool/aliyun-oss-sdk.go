package tool

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func SignURLSample(object string) (signedURL string, err error) {
	accessKeyId := "a5etodit71tlznjt3pdx7lch"  //AccessKeyId，需要使用用户自己的
	accessKeySecret := "secret_key"            //AccessKeySecret，需要用用户自己的
	endpoint := "oss-cn-hangzhou.aliyuncs.com" //Endpoint，根据Bucket创建的区域来选择，本文中是杭州
	bucketname := "referer-test"               //Bucket，需要用用户自己的
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		// HandleError(err)
	}
	bucket, err := client.Bucket(bucketname)

	signedURL, err = bucket.SignURL(object, oss.HTTPGet, 600)
	if err != nil {

	}
	return
}
