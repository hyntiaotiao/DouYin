package service

import (
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"os"
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
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	resumeUploader := storage.NewResumeUploaderV2(&cfg)
	ret := storage.PutRet{}
	recorder, err := storage.NewFileRecorder(os.TempDir())
	if err != nil {
		return err
	}
	putExtra := storage.RputV2Extra{
		Recorder: recorder,
	}

	dataLen := int64(len(data))

	err = resumeUploader.Put(context.Background(), ret, upToken, videoName, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return err
	}
	videoUrl := domain + "/" + videoName
	coverUrl := domain + "/" + coverName
	// 存入数据库
	err = videoDao.Addvideo(authorId, videoUrl, coverUrl, title)
	if err != nil {
		log.Println("发布视频失败，请稍后再试")
	}
	return err
}
