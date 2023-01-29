package controller

import (
	"DouYIn/service"
	"DouYIn/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

//type UserPublishRequest struct {
//	Token string `json:"token" binding:"required"`
//	Data  byte   `form:"data" json:"data" binding:"required"`
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
		ext := filepath.Ext(file.Filename) // 得到后缀
		// 上传合法性判断
		if _, ok := videoIndexMap[ext]; !ok {
			PublishVideoError(c, "视频格式不支持")
			continue
		}
		// 存入本地，是否需要云端？
		err := c.SaveUploadedFile(file, "D:\\newproject\\go\\douyin\\video")
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		// 制作视频封面
		err = SaveImageFromVideo(file.Filename, true)
		if err != nil {
			PublishVideoError(c, err.Error())
			continue
		}
		// 数据存入数据库
		playUrl := file.Filename + ".mp4"
		coverUrl := file.Filename + ".jpg"
		authorId, exists := c.Get("UserId")
		if !exists {
			PublishVideoError(c, err.Error())
			continue
		}
		err = service.AddVideo(authorId.(int), playUrl, coverUrl, title)
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
	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: msg})
}

// SaveImageFromVideo 将视频切一帧保存到本地
// isDebug用于控制是否打印出执行的ffmepg命令
func SaveImageFromVideo(name string, isDebug bool) error {
	v2i := utils.NewVideo2Image()
	if isDebug {
		v2i.Debug()
	}
	v2i.InputPath = filepath.Join("D:\\newproject\\go\\douyin\\video", name+".mp4")
	v2i.OutputPath = filepath.Join("D:\\newproject\\go\\douyin\\video\\cover", name+".jpg")
	v2i.FrameCount = 1
	queryString, err := v2i.GetQueryString()
	if err != nil {
		return err
	}
	return v2i.ExecCommand(queryString)
}
