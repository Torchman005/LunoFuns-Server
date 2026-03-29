package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	minioGo "github.com/minio/minio-go/v7"

	"LunoFuns-Server/internal/services"
)

type UploadController struct {
	uploadService *services.UploadService
}

func NewUploadController(uploadService *services.UploadService) *UploadController {
	return &UploadController{
		uploadService: uploadService,
	}
}

// GetUploadToken 获取简单上传凭证（例如封面）
func (ctrl *UploadController) GetUploadToken(c *gin.Context) {
	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadURL, fileURL, err := ctrl.uploadService.GeneratePresignedURL(req.Filename, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成上传凭证失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_url": uploadURL,
		"file_url":   fileURL,
	})
}

// InitMultipartUpload 初始化分片上传
func (ctrl *UploadController) InitMultipartUpload(c *gin.Context) {
	var req struct {
		Filename    string `json:"filename" binding:"required"`
		ContentType string `json:"content_type"`
		PartCount   int    `json:"part_count" binding:"required,min=1,max=10000"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadID, objectName, err := ctrl.uploadService.InitMultipartUpload(req.Filename, req.ContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "初始化分片上传失败: " + err.Error()})
		return
	}

	urls, err := ctrl.uploadService.GetMultipartUploadURLs(objectName, uploadID, req.PartCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分片上传链接失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_id":   uploadID,
		"object_name": objectName,
		"upload_urls": urls,
	})
}

// CompleteMultipartUpload 完成分片上传
func (ctrl *UploadController) CompleteMultipartUpload(c *gin.Context) {
	var req struct {
		UploadID   string                 `json:"upload_id" binding:"required"`
		ObjectName string                 `json:"object_name" binding:"required"`
		Parts      []minioGo.CompletePart `json:"parts" binding:"required"` // 需要包含 PartNumber 和 ETag
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileURL, err := ctrl.uploadService.CompleteMultipartUpload(req.ObjectName, req.UploadID, req.Parts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "合并分片失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "上传成功",
		"file_url": fileURL,
	})
}
