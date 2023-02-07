package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"fmt"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//type UserPublishRequest struct {
//	Token string `json:"token" binding:"required"`
//	Data  []byte `form:"data" json:"data" binding:"required"`
//	Title string `form:"title" json:"title" binding:"required"`
//}

var (
	// 视频图片格式判断
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

func Publish(c *gin.Context) {

	title := c.PostForm("title")
	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}
	files := form.File["data"]

	// 多个文件 文件信息计入数据库
	for _, file := range files {

		open, err := file.Open()
		if err != nil {
			PublishVideoError(c, "上传文件数据有误，无法读取")
			continue
		}
		defer open.Close()
		size := file.Size
		bytes := make([]byte, size)
		if _, err := open.Read(bytes); err != nil {
			PublishVideoError(c, "文件读取错误")
			continue
		}

		ext := filepath.Ext(file.Filename) // 得到后缀
		// 上传合法性判断
		if _, ok := videoIndexMap[ext]; !ok {
			PublishVideoError(c, "视频格式不支持")
			continue
		}

		index := strings.LastIndex(file.Filename, ".")
		newfileName := file.Filename[0:index]

		UserIDAny, _ := c.Get("UserID")
		UserID, _ := strconv.ParseInt(fmt.Sprintf("%v", UserIDAny), 0, 64)

		// 拼接视频文件名 用户id+视频title
		videoName := strconv.FormatInt(UserID, 10) + "-" + newfileName + ".mp4"
		coverName := strconv.FormatInt(UserID, 10) + "-" + newfileName + ".jpg"

		// 上传到数据库
		err = service.AddVideo(bytes, videoName, coverName, UserID, title)
		if err != nil {
			PublishVideoError(c, err.Error())
		}

		// 制作视频封面
		fileUrl := "http://rphysx900.hn-bkt.clouddn.com/" + videoName
		tmpCoverUrl := "tmpCover/" + coverName
		err = ffmpeg.Input(fileUrl, ffmpeg.KwArgs{"ss": "1"}).
			// "s": "320x240", "pix_fmt": "rgb24", "t": "3", "r": "3"
			Output(tmpCoverUrl, ffmpeg.KwArgs{"s": "368x208", "pix_fmt": "rgb24", "t": "3", "r": "3"}).
			OverWriteOutput().ErrorToStdOut().Run()
		if err != nil {
			PublishVideoError(c, err.Error())
		}

		// 文件转换为字节流文件
		openFile, err := os.Open(tmpCoverUrl)
		if err != nil {
			PublishVideoError(c, err.Error())
		}
		defer openFile.Close()
		var data []byte
		buf := make([]byte, 1024)
		for {
			// 将文件中读取的byte存储到buf中
			n, err := openFile.Read(buf)
			if err != nil && err != io.EOF {
				PublishVideoError(c, err.Error())
			}
			if n == 0 {
				break
			}
			// 将读取到的结果追加到data切片中
			data = append(data, buf[:n]...)
		}

		// 视频封面上传到oss
		service.UploadDataToOSS(data, coverName)
		err = os.Remove(tmpCoverUrl)
		if err != nil {
			PublishVideoError(c, err.Error())
		}
	}

	// 成功
	var response = &common.Response{}
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}

// 返回错误
func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: msg})
}
