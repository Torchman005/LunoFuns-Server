package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"LunoFuns-Server/internal/models"
	"LunoFuns-Server/internal/services"
)

type VideoController struct {
	videoService *services.VideoService
}

func NewVideoController(videoService *services.VideoService) *VideoController {
	return &VideoController{
		videoService: videoService,
	}
}

// UploadVideo 提交视频信息
func (ctrl *VideoController) UploadVideo(c *gin.Context) {
	var req models.UploadVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	video, err := ctrl.videoService.UploadVideo(&req, userID.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传视频失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "视频上传成功",
		"data":    video,
	})
}

// GetVideoList 获取视频列表
func (ctrl *VideoController) GetVideoList(c *gin.Context) {
	var req models.VideoListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := ctrl.videoService.GetVideoList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取视频列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resp,
	})
}

// GetVideoDetail 获取视频详情
func (ctrl *VideoController) GetVideoDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的视频ID"})
		return
	}

	video, err := ctrl.videoService.GetVideoDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "视频不存在或已下架"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    video,
	})
}
