package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"DouYIn/utils"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

type UserPublishRequest struct {
	Token string `json:"token" binding:"required"`
	Data  []byte `form:"data" json:"data" binding:"required"`
	Title string `form:"title" json:"title" binding:"required"`
}

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

	var request UserPublishRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Println("request参数绑定失败")
		return
	}

	title := c.PostForm("title")
	form, err := c.MultipartForm()
	if err != nil {
		PublishVideoError(c, err.Error())
		return
	}
	files := form.File["data"]

	// 多个文件 文件信息计入数据库
	for _, file := range files {

		ext := filepath.Ext(file.Filename) // 得到后缀
		// 上传合法性判断
		if _, ok := videoIndexMap[ext]; !ok {
			PublishVideoError(c, "视频格式不支持")
			continue
		}
		// 存入本地，是否需要云端？
		//err := c.SaveUploadedFile(file, "D:\\newproject\\go\\douyin\\video")
		//if err != nil {
		//	PublishVideoError(c, err.Error())
		//	continue
		//}
		// 制作视频封面
		coverName := ""
		err, coverName = SaveImageFromVideo(file.Filename, true)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		// 数据存入数据库

		authorId, exists := c.Get("UserId")
		if !exists {
			PublishVideoError(c, err.Error())
			continue
		}
		// 获取文件md5值
		sum := md5.Sum(request.Data)
		s := fmt.Sprintf("%x", sum)

		// 拼接视频文件名
		videoName := authorId.(string) + "-" + s

		err = service.AddVideo(request.Data, videoName, coverName, authorId.(int64), title)
		if err != nil {
			PublishVideoError(c, err.Error())
		}
	}

	// 成功
	var response = &UserLoginResponse{}
	response.StatusCode = 0
	response.StatusMsg = "success"
	c.JSON(200, response)
}

// 返回错误
func PublishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: msg})
}

// SaveImageFromVideo 将视频切一帧保存到本地
// isDebug用于控制是否打印出执行的ffmepg命令
func SaveImageFromVideo(flieName string, isDebug bool) (error, string) {
	v2i := utils.NewVideo2Image()
	if isDebug {
		v2i.Debug()
	}
	coverName := flieName + ".jpg"
	v2i.InputPath = filepath.Join("rphysx900.hn-bkt.clouddn.com", flieName+".mp4")
	v2i.OutputPath = filepath.Join("rphysx900.hn-bkt.clouddn.com", coverName)
	v2i.FrameCount = 1
	queryString, err := v2i.GetQueryString()
	if err != nil {
		return err, ""
	}
	return v2i.ExecCommand(queryString), coverName
}
