package service

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
)

var (
	accessKey = "XRPWoJnixwJaw_9Skz8VHUFMAb9tuiQqjiSNjsYl"
	secretKey = "QeWXT7tqjruO-xmPn1u5Ndu2EWlWMVftF_Smv-ki"
	bucket    = "douyin-video-1433"
	// 域名
	domain = "rphysx900.hn-bkt.clouddn.com"
)

func AddVideo(data []byte, videoName string, coverName string, authorId int64, title string) error {

	// 存入七牛云oss对象存储
	err := UploadDataToOSS(data, videoName)
	if err != nil {
		return err
	}
	videoUrl := domain + "/" + videoName
	coverUrl := domain + "/" + coverName

	//// 获取文件md5值,用于校验重复视频
	//sum := md5.Sum(data)
	//videoMD5 := fmt.Sprintf("%x", sum)
	// 存入数据库
	err = videoDao.Addvideo(authorId, videoUrl, coverUrl, title)
	if err != nil {
		log.Println("发布文件失败，请稍后再试")
	}
	return err

}

// 将本地文件上传到七牛云oss中
func UploadDataToOSS(data []byte, videoName string) error {
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	//recorder, err := storage.NewFileRecorder(os.TempDir())
	//if err != nil {
	//	return err
	//}
	putExtra := storage.PutExtra{}
	dataLen := int64(len(data))

	err := formUploader.Put(context.Background(), &ret, upToken, videoName, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return err
	}
	return nil
}
