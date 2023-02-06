package controller

import (
	"DouYIn/common"
	"DouYIn/service"
	"DouYIn/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

	//var request UserPublishRequest
	//if err := c.Bind(&request); err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	log.Println("request参数绑定失败")
	//	return
	//}

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
		coverName := newfileName + ".jpg"

		UserIDAny, _ := c.Get("UserID")
		UserID, _ := strconv.ParseInt(fmt.Sprintf("%v", UserIDAny), 0, 64)

		// 拼接视频文件名
		videoName := strconv.FormatInt(UserID, 10) + "-" + file.Filename

		// 上传到数据库
		err = service.AddVideo(bytes, videoName, coverName, UserID, title)
		if err != nil {
			PublishVideoError(c, err.Error())
		}

		// 制作视频封面（有问题，读不到文件）
		//err = SaveImageFromVideo(videoName, coverName, true)
		//if err != nil {
		//	PublishVideoError(c, err.Error())
		//	continue
		//}
		//
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

// SaveImageFromVideo 将视频切一帧保存到本地
// isDebug用于控制是否打印出执行的ffmepg命令
func SaveImageFromVideo(fileName string, coverName string, isDebug bool) error {
	v2i := utils.NewVideo2Image()
	if isDebug {
		v2i.Debug()
	}

	v2i.InputPath = filepath.Join("rphysx900.hn-bkt.clouddn.com", fileName)
	v2i.OutputPath = filepath.Join("rphysx900.hn-bkt.clouddn.com", coverName)
	v2i.FrameCount = 1
	queryString, err := v2i.GetQueryString()
	if err != nil {
		return err
	}
	return v2i.ExecCommand(queryString)
}
