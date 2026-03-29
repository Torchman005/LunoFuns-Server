package services

import (
	"LunoFuns-Server/internal/models"

	"gorm.io/gorm"
)

type VideoService struct {
	db *gorm.DB
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{db: db}
}

func (s *VideoService) UploadVideo(req *models.UploadVideoRequest, userID uint64) (*models.Video, error) {
	video := &models.Video{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CoverURL:    req.CoverURL,
		VideoURL:    req.VideoURL,
		Duration:    req.Duration,
		Size:        req.Size,
		Format:      req.Format,
		Width:       req.Width,
		Height:      req.Height,
		CategoryID:  req.CategoryID,
		Tags:        req.Tags,
		Status:      1, // 默认已发布
	}

	if err := s.db.Create(video).Error; err != nil {
		return nil, err
	}

	return video, nil
}

func (s *VideoService) GetVideoList(req *models.VideoListRequest) (*models.VideoListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var videos []models.Video
	var total int64

	query := s.db.Model(&models.Video{}).Where("status = ?", 1)
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at desc").Offset(offset).Limit(req.PageSize).Find(&videos).Error; err != nil {
		return nil, err
	}

	return &models.VideoListResponse{
		Total: total,
		Items: videos,
	}, nil
}

func (s *VideoService) GetVideoDetail(id uint64) (*models.Video, error) {
	var video models.Video
	if err := s.db.Where("id = ? AND status = ?", id, 1).First(&video).Error; err != nil {
		return nil, err
	}
	// 增加播放量 (简单的并发非安全增加，要求高可考虑 redis 或者 gorm expr)
	s.db.Model(&video).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))
	video.ViewCount++
	return &video, nil
}
